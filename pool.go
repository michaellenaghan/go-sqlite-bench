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
)

// poolObject represents an object in the pool along with its last usage time.
type poolObject[T any] struct {
	object   T
	lastUsed time.Time
}

// Pool is a generic object pool that manages a collection of objects of type T.
// It maintains a minimum and maximum number of objects, handles object creation
// and destruction, and manages object availability and idle time.
type Pool[T any] struct {
	min         int               // must be >= 0
	max         int               // must be >= min
	idleTime    time.Duration     // must be >= 0; 0 == "never idle out"
	newFunc     func() (T, error) // required
	destroyFunc func(T)           // optional

	mu   sync.Mutex
	cond *sync.Cond // signal "available" (or, occassionally: ctx.Done())

	objectCount int             // min <= objectCount <= max
	objects     []poolObject[T] // available objects, ordered by lastUsed

	stoppingOrStopped bool
	stopping          chan struct{}

	createdCount int // didn't have available, had capacity; created new pool object
	blockedCount int // didn't have available, didn't have capacity; blocked waiting for old pool object
}

// Stats represents statistical information about the pool's performance.
type Stats struct {
	// Didn't have available, had capacity; created new pool object.
	Created int
	// Didn't have available, didn't have capacity; blocked waiting for old pool object.
	Blocked int
}

// NewPool creates a new pool of objects of type T. The pool will maintain
// somewhere between a minimum number of objects (min) and a maximum number
// of objects (max). Objects will be reused up to the idle time (idleTime)
// before being destroyed and recreated. The newFunc function is required
// and used to create new objects, and the destroyFunc function is optional
// and used to destroy objects when they're removed from the pool.
func NewPool[T any](min, max int, idleTime time.Duration, newFunc func() (T, error), destroyFunc func(T)) (*Pool[T], error) {
	if min < 0 {
		return nil, errors.New("min must be greater than or equal to zero")
	}
	if min > max {
		return nil, errors.New("min must be less than or equal to max")
	}
	if idleTime < 0 {
		return nil, errors.New("idle time must be greater than or equal to zero")
	}
	if newFunc == nil {
		return nil, errors.New("newFunc is required")
	}
	p := &Pool[T]{
		min:         min,
		max:         max,
		idleTime:    idleTime,
		newFunc:     newFunc,
		destroyFunc: destroyFunc,
		stopping:    make(chan struct{}),
	}
	p.cond = sync.NewCond(&p.mu)
	return p, nil
}

// Start initializes the pool and prepares it for use. If immediately is true,
// it creates the minimum number of pool objects right away. If immediately is
// false, it doesn't create any pool objects right away; rather, it creates
// them lazily, on demand.
//
// If immediately is true and there's an error creating an initial object, the
// pool will destroy any objects it created and return the error.
//
// This function should be called after Init to ensure the pool is ready to use.
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
		for i := 0; i < p.min; i++ {
			object, err := p.newFunc()
			if err != nil {
				for _, pObject := range p.objects {
					if p.destroyFunc != nil {
						p.destroyFunc(pObject.object)
					}
					p.objectCount--
				}
				p.objects = nil

				return errors.Join(ErrNew, err)
			}
			p.objects = append(p.objects, poolObject[T]{object: object, lastUsed: time.Now()})
			p.objectCount++
		}
	}

	go p.cleanupTick()

	return nil
}

// Stop stops the pool. It marks the pool as stopping or stopped, closes the
// stopping channel, broadcasts the condition variable to wake up any waiting
// goroutines, and destroys all objects in the pool. If the pool is already
// stopping or stopped, this function does nothing.
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

	p.cond.Broadcast()

	for _, pObject := range p.objects {
		if p.destroyFunc != nil {
			p.destroyFunc(pObject.object)
		}
		p.objectCount--
	}
	p.objects = nil
}

// Get retrieves an object from the pool. If an object is available, it is
// returned immediately. If no objects are available and the pool has not
// reached its maximum capacity, a new object is created. If the pool is at
// capacity, the call blocks until an object becomes available or the provided
// context is cancelled.
//
// If the pool is stopping or stopped, an error is returned.
func (p *Pool[T]) Get(ctx context.Context) (T, error) {
	// log.Println("> Get")
	// defer log.Println("< Get")
	p.mu.Lock()
	defer p.mu.Unlock()

	for {
		// 1. pool is stopping or stopped
		if p.stoppingOrStopped {
			var zero T
			return zero, ErrStoppingOrStopped
		}

		// 2. pool has available objects
		if len(p.objects) > 0 {
			// Pop from the front to preserve lastUsed ordering.
			pObject := p.objects[0]
			p.objects = p.objects[1:]
			return pObject.object, nil
		}

		// 3. pool has available capacity
		if p.objectCount < p.max {
			p.createdCount++

			object, err := p.newFunc()
			if err != nil {
				return object, errors.Join(ErrNew, err)
			}
			p.objectCount++

			return object, nil
		}

		// 4. pool is full, wait for an object
		p.blockedCount++

		waitCh := make(chan struct{})
		go func() {
			p.mu.Lock()
			p.cond.Wait() // signalled by `Put()`, `Stop()` -- and `ctx.Done()`
			close(waitCh)
			p.mu.Unlock()
		}()

		p.mu.Unlock()
		select {
		case <-waitCh:
			p.mu.Lock()
		case <-ctx.Done():
			p.mu.Lock()

			// Broadcast the condition to ensure that the `Wait()` in the
			// goroutine is triggered; otherwise the goroutine will remain
			// blocked until `Stop()` is eventually called. `Broadcast()`
			// may result in spurious wakeups, but it's our only choice;
			// `Signal()` will only signal *one* waiter, and it may not be
			// *our* waiter.
			p.cond.Broadcast()

			var zero T
			return zero, ctx.Err()
		}
	}
}

// Put returns an object to the pool. If the pool is stopping or stopped,
// the object is destroyed immediately instead of being returned to the
// pool.
func (p *Pool[T]) Put(object T) {
	// log.Println("> Put")
	// defer log.Println("< Put")
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.stoppingOrStopped {
		if p.destroyFunc != nil {
			p.destroyFunc(object)
		}
		p.objectCount--
		return
	}

	p.objects = append(p.objects, poolObject[T]{object: object, lastUsed: time.Now()})
	p.cond.Signal()
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

// cleanupTock is an internal method that performs the actual cleanup of idle
// objects in the pool.
func (p *Pool[T]) cleanupTock() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.stoppingOrStopped {
		// If stopping or stopped, don't clean up.
		return
	}

	if p.objectCount <= p.min {
		// If there are too few total objects, don't clean up.
		return
	}
	if len(p.objects) <= 0 {
		// If there are too few available objects, don't clean up.
		return
	}

	if time.Since(p.objects[0].lastUsed) < p.idleTime {
		// If the first available object hasn't idled out, don't clean up.
		return
	}

	// Now we know at least one available object has idled out;
	// filter the queue. (Since the queue is ordered by lastUsed
	// the objects that have idled out will all be at the front.)

	for i, pObject := range p.objects {
		if p.objectCount > p.min && time.Since(pObject.lastUsed) >= p.idleTime {
			if p.destroyFunc != nil {
				p.destroyFunc(pObject.object)
			}
			p.objectCount--
		} else {
			p.objects = p.objects[i:]
			return
		}
	}

	// If we get here, all objects were idle and destroyed. Clear
	// the slice to avoid keeping references to destroyed objects.
	p.objects = p.objects[:0]
}

// Stats returns statistical information about the pool's performance.
func (p *Pool[T]) Stats() Stats {
	p.mu.Lock()
	defer p.mu.Unlock()

	return Stats{
		Created: p.createdCount,
		Blocked: p.blockedCount,
	}
}
