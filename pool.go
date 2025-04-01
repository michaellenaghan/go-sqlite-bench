// Package pool provides a generic object pool implementation. It allows
// for efficient reuse of expensive-to-create objects, managing their
// lifecycle and availability based on configurable parameters.
package go_sqlite_bench

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrNew               = errors.New("failed to make new pool object")
	ErrStoppingOrStopped = errors.New("pool is stopping or stopped")
	ErrRingIsEmpty       = errors.New("ring is empty")
	ErrRingIsFull        = errors.New("ring is full")
)

// Pool is a generic object pool that manages a collection of objects of type T.
// It maintains a minimum and maximum number of objects, handles object creation
// and destruction, and manages object availability and idle time.
type Pool[T any] struct {
	min         int               // must be >= 0
	max         int               // must be >= min
	idleTime    time.Duration     // must be >= 0; 0 == "never idle out"
	newFunc     func() (T, error) // required
	checkFunc   func(T) error     // optional
	destroyFunc func(T)           // optional

	mu sync.Mutex

	count int // min <= count <= max

	idle    ring[T] // idle object ring buffer, ordered by lastUsed
	waiting chan T  // "Get"ers waiting for "Put"ers

	stoppingOrStopped bool
	stopping          chan struct{}

	createdTotal   uint // created an object (stats only)
	waitedTotal    uint // waited for an object (stats only)
	destroyedTotal uint // destroyed an object (stats only)
}

// Stats represents statistical information about the pool's performance.
type Stats struct {
	// Created an object.
	CreatedTotal uint
	// Waited for an object.
	WaitedTotal uint
	// Destroyed an object.
	DestroyedTotal uint

	// Number of objects (count == busy + idle)
	CountNow int
	// Number of busy objects (busy == count - idle)
	BusyNow int
	// Number of idle objects (idle == count - busy)
	IdleNow int
	// Number of waiting "Get"ers
	WaitingNow int
}

// NewPool creates a new pool of objects of type T.
//
// The pool will maintain somewhere between a minimum number of objects (min)
// and a maximum number of objects (max).
//
// Objects will be reused up to the idle time (idleTime) before being destroyed
// and recreated.
//
// An idle time of zero means that objects never idle out.
//
// (In that case, min and max should be the same.)
//
// The newFunc function is required and used to create new objects.
//
// The checkFunc function is optional and used to check objects when they're
// returned to the pool. If check returns an error the object is destroyed
// (via destroyFunc) and not returned.
//
// The destroyFunc function is optional and used to destroy objects when
// they're no longer needed.
func NewPool[T any](min, max int, idleTime time.Duration, newFunc func() (T, error), checkFunc func(T) error, destroyFunc func(T)) (*Pool[T], error) {
	if min < 0 {
		return nil, errors.New("min must be greater than or equal to zero")
	}
	if min > max {
		return nil, errors.New("min must be less than or equal to max")
	}
	if idleTime < 0 {
		return nil, errors.New("idle time must be greater than or equal to zero")
	}
	// This isn't *really* an error, it's just an indication that someone's
	// mental model isn't quite right. If idle time is zero, the only
	// difference between min and max is that min objects might be created
	// immediately during `Start()`. Other than that, once more than min
	// objects exist, they'll continue to exist until `Stop()`. Again, not
	// *really* an error, but probably not what was intended?
	if idleTime == 0 && min != max {
		return nil, errors.New("when idle time equals zero min should equal max")
	}
	if newFunc == nil {
		return nil, errors.New("newFunc is required")
	}
	p := &Pool[T]{
		min:         min,
		max:         max,
		idleTime:    idleTime,
		newFunc:     newFunc,
		checkFunc:   checkFunc,
		destroyFunc: destroyFunc,
		idle:        newRing[T](max, idleTime),
		waiting:     make(chan T),
		stopping:    make(chan struct{}),
	}
	return p, nil
}

// Start initializes the pool and prepares it for use.
//
// If the pool is already stopping or stopped, Start returns an error.
//
// If immediately is true, Start creates the minimum number of pool objects
// right away. If there's an error creating one of the objects Start destroys
// any objects it created and returns an error.
//
// If immediately is false, Start doesn't create any pool objects right away.
//
// Start should be called after NewPool to ensure the pool is ready to use.
func (p *Pool[T]) Start(immediately bool) error {
	// log.Println("> Start")
	// defer log.Println("< Start")
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.stoppingOrStopped {
		return ErrStoppingOrStopped
	}

	// Immediately create the minimum number of objects?
	if immediately {
		for range p.min {
			object, err := p.newFunc()
			if err != nil {
				for p.idle.count > 0 {
					object := p.idle.popOldest()
					if p.destroyFunc != nil {
						p.destroyFunc(object)
					}
					p.destroyedTotal++
					p.count--
				}

				return errors.Join(ErrNew, err)
			}
			p.idle.pushNewest(object)
			p.createdTotal++
			p.count++
		}
	}

	go p.cleanupTick()

	return nil
}

// Stop stops the pool.
//
// Stop marks the pool as stopping or stopped, closes the waiting and stopping
// channels, and destroys all idle objects.
//
// If the pool is already stopping or stopped, Stop does nothing.
func (p *Pool[T]) Stop() {
	// log.Println("> Stop")
	// defer log.Println("< Stop")
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.stoppingOrStopped {
		return
	}

	p.stoppingOrStopped = true
	close(p.stopping)

	for p.idle.count > 0 {
		object := p.idle.popOldest()
		if p.destroyFunc != nil {
			p.destroyFunc(object)
		}
		p.destroyedTotal++
		p.count--
	}
}

// Get returns an object from the pool.
//
// If the pool is stopping or stopped, Get returns an error.
//
// If there are idle objects, Get returns an idle object.
//
// If there are no idle objects and the pool has capacity, Get returns a newly
// created object.
//
// If there are no idle objects and the pool has no capacity, Get waits for an
// object to be returned to the pool by Put. Get stops waiting when Stop is
// called or when the provided context is cancelled.
func (p *Pool[T]) Get(ctx context.Context) (T, error) {
	// log.Println("> Get")
	// defer log.Println("< Get")
	p.mu.Lock()

	// 1. pool is stopping or stopped
	if p.stoppingOrStopped {
		p.mu.Unlock()

		var zero T
		return zero, ErrStoppingOrStopped
	}

	// 2. pool has idle objects
	if p.idle.count > 0 {
		// popOldest() makes the oldest connections less likely to idle out
		// since they get used more frequently.
		//
		// popNewest() makes the oldest connections more likely to idle out
		// since they get used less frequently.
		object := p.idle.popNewest()
		p.mu.Unlock()

		return object, nil
	}

	// 3. pool has capacity
	if p.count < p.max {
		// Don't hold the lock while we call newFunc().
		p.createdTotal++
		p.count++
		p.mu.Unlock()

		object, err := p.newFunc()
		if err != nil {
			p.mu.Lock()
			p.count--
			p.mu.Unlock()

			var zero T
			return zero, errors.Join(ErrNew, err)
		}

		return object, nil
	}

	// 4. pool has no idle objects and no capacity; wait for an object
	p.waitedTotal++
	p.mu.Unlock()

	select {
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	case object := <-p.waiting:
		return object, nil
	case <-p.stopping:
		var zero T
		return zero, ErrStoppingOrStopped
	}
}

// Put returns an object to the pool.
//
// If the pool is stopping or stopped, Put destroys the object rather than
// return it.
func (p *Pool[T]) Put(object T) {
	// log.Println("> Put")
	// defer log.Println("< Put")
	if p.checkFunc != nil {
		err := p.checkFunc(object)
		if err != nil {
			// Don't hold the lock while we call destroyFunc().
			p.mu.Lock()
			p.destroyedTotal++
			p.count--
			p.mu.Unlock()

			if p.destroyFunc != nil {
				p.destroyFunc(object)
			}

			return
		}
	}

	p.mu.Lock()

	if p.stoppingOrStopped {
		// Don't hold the lock while we call destroyFunc().
		p.destroyedTotal++
		p.count--
		p.mu.Unlock()

		if p.destroyFunc != nil {
			p.destroyFunc(object)
		}

		return
	}

	select {
	case p.waiting <- object:
	default:
		p.idle.pushNewest(object)
	}

	p.mu.Unlock()
}

// cleanupTick is an internal method that periodically checks for and removes
// idle objects from the pool.
func (p *Pool[T]) cleanupTick() {
	if p.max > p.min && p.idleTime > 0 {
		ticker := time.NewTicker(p.idleTime / 2)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				p.cleanupTock()
			case <-p.stopping:
				return
			}
		}
	}
}

// cleanupTock is an internal method that periodically checks for and removes
// idle objects from the pool.
func (p *Pool[T]) cleanupTock() {
	p.mu.Lock()
	for !p.stoppingOrStopped && p.count > p.min && p.idle.count > 0 && p.idle.oldestIdleTooLong() {
		// Don't hold the lock while we call destroyFunc().
		object := p.idle.popOldest()
		p.destroyedTotal++
		p.count--
		p.mu.Unlock()

		if p.destroyFunc != nil {
			p.destroyFunc(object)
		}

		p.mu.Lock()
	}
	p.mu.Unlock()
}

// Stats returns statistics about the pool.
func (p *Pool[T]) Stats() Stats {
	p.mu.Lock()
	defer p.mu.Unlock()

	return Stats{
		CreatedTotal:   p.createdTotal,
		WaitedTotal:    p.waitedTotal,
		DestroyedTotal: p.destroyedTotal,

		CountNow:   p.count,
		BusyNow:    p.count - p.idle.count,
		IdleNow:    p.idle.count,
		WaitingNow: len(p.waiting),
	}
}

// ===

// ring is a generic ring buffer that stores objects along with their last
// used time.
type ring[T any] struct {
	buffer   []ringObject[T] // The actual ring buffer storage
	head     int             // Index of the first object (0 <= head < cap(buffer))
	tail     int             // Index of the next object (0 <= tail < cap(buffer))
	count    int             // Current number of items in the buffer (0 <= count <= cap(buffer))
	idleTime time.Duration   // Maximum time an object should be idle (>= 0; 0 == "never idle out")
}

// ringObject represents a generic object in the ring buffer along with its
// last used time.
type ringObject[T any] struct {
	object   T
	lastUsed time.Time
}

// newRing creates a new generic ring buffer with the given capacity.
func newRing[T any](max int, idleTime time.Duration) ring[T] {
	return ring[T]{
		buffer:   make([]ringObject[T], max),
		idleTime: idleTime,
	}
}

// oldestIdleTooLong returns true if the oldest object has been idle too long.
func (r *ring[T]) oldestIdleTooLong() bool {
	if r.count > 0 {
		return r.idleTime > 0 && time.Since(r.buffer[r.head].lastUsed) >= r.idleTime
	} else {
		panic(ErrRingIsEmpty)
	}
}

// popOldest removes and returns an object from the head of the ring buffer.
// This implements FIFO (First-In-First-Out) behavior.
func (r *ring[T]) popOldest() T {
	if r.count > 0 {
		var zero T
		object := r.buffer[r.head].object
		r.buffer[r.head].object = zero

		if r.head < cap(r.buffer)-1 {
			r.head++
		} else {
			r.head = 0
		}
		r.count--

		return object
	} else {
		panic(ErrRingIsEmpty)
	}
}

// popNewest removes and returns an object from the tail of the ring buffer.
// This implements LIFO (Last-In-First-Out) behavior.
func (r *ring[T]) popNewest() T {
	if r.count > 0 {
		if r.tail > 0 {
			r.tail--
		} else {
			r.tail = cap(r.buffer) - 1
		}
		r.count--

		var zero T
		object := r.buffer[r.tail].object
		r.buffer[r.tail].object = zero

		return object
	} else {
		panic(ErrRingIsEmpty)
	}
}

// pushNewest adds an object to the tail of the ring buffer.
func (r *ring[T]) pushNewest(object T) {
	if r.count < cap(r.buffer) {
		r.buffer[r.tail].object = object
		// Don't record lastUsed time if we don't need it.
		// (Recording time takes time.)
		if r.idleTime > 0 {
			r.buffer[r.tail].lastUsed = time.Now()
		}

		if r.tail < cap(r.buffer)-1 {
			r.tail++
		} else {
			r.tail = 0
		}
		r.count++
	} else {
		panic(ErrRingIsFull)
	}
}
