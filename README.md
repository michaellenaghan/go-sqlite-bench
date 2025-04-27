This project, originally inspired by [Benchmarking SQLite Performance in Go](https://www.golang.dk/articles/benchmarking-sqlite-performance-in-go), bechmarks various SQLite implementations:

* [github.com/eatonphil/gosqlite](https://github.com/eatonphil/gosqlite) (aka "eatonphil")
* [github.com/glebarez/go-sqlite](https://github.com/glebarez/go-sqlite) (aka "glebarez")
* [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) (aka "mattn")
* [github.com/ncruces/go-sqlite3](https://github.com/ncruces/go-sqlite3) (aka "ncruces")
* [github.com/tailscale/sqlite](https://github.com/tailscale/sqlite) (aka "tailscale")
*	[github.com/zombiezen/go-sqlite](https://github.com/zombiezen/go-sqlite) (aka "zombiezen")
*	[gitlab.com/cznic/sqlite](https://gitlab.com/cznic/sqlite) (aka "modernc")

Here are some quick descriptions:

* **eatonphil**

  eatonphil is a CGO-based implementation. It offers a direct interface.

* **glebarez**

  glebarez is a non-CGO transpilation-based implementation, based on modernc. It offers a `database/sql` interface.

* **mattn**

  mattn is a CGO-based implementation. It offers a `database/sql` interface.

* **modernc**

  modernc is a non-CGO transpilation-based implementation. It offers a `database/sql` interface.

* **ncruces**

  ncruces is a non-CGO WASM-based implementation. It offers both direct and `database/sql` interfaces.

* **tailscale**

  tailscale is a CGO-based implementation. It offers both direct and `database/sql` interfaces.

* **zombiezen**

  zombiezen is a non-CGO transpilation-based implementation, based on modernc. It offers a direct interface.

The benchmark results are [below](#reviewing-the-benchmarks).

# Running the Benchmarks

The benchmark results in the [quick](./quick/) directory were generated using:

```sh
make benchstat-by-category
```

It takes ~20m to run the "quick" benchmarks on my laptop.

The benchmark results in the [slow](./slow/) directory were generated using:

```sh
make benchstat-by-category BENCH_BIG=1 BENCH_SLOW=1
```

It takes ~1h10m  to run the "slow" benchmarks on my laptop.

The tests in the [tests](./tests/) directory were generated using:

```sh
make test-all
```

Among other things, the tests capture the [compile-time options, pragmas and SQLite version](#reviewing-the-implementations) used by each implementation.

(You can diff one test file against another to get a sense of how and where the implementations... well, differ.)

There are lots of ways to play the benchmark game; take a look at the examples below.

<!--COMMAND:make help-->
```
$ make help
Targets:

  bench-all                                          - Run all benchmarks
  bench-by-category                                  - Run all benchmark categories
  bench-category-baseline                            - Run baseline benchmarks
  bench-category-baseline-parallel                   - Run baseline benchmarks in parallel
  bench-category-populate                            - Run populate benchmarks
  bench-category-readwrite                           - Run readwrite benchmarks
  bench-category-readwrite-parallel                  - Run readwrite benchmarks in parallel
  bench-category-query-correlated                    - Run correlated query benchmarks
  bench-category-query-correlated-parallel           - Run correlated query benchmarks in parallel
  bench-category-query-groupby                       - Run groupby query benchmarks
  bench-category-query-groupby-parallel              - Run groupby query benchmarks in parallel
  bench-category-query-json                          - Run json query benchmarks
  bench-category-query-json-parallel                 - Run json query benchmarks in parallel
  bench-category-query-orderby                       - Run orderby query benchmarks
  bench-category-query-orderby-parallel              - Run orderby query benchmarks in parallel
  bench-category-query-recursivecte                  - Run recursivecte query benchmarks
  bench-category-query-recursivecte-parallel         - Run recursivecte query benchmarks in parallel
  bench-category-query-window                        - Run window query benchmarks
  bench-category-query-window-parallel               - Run window query benchmarks in parallel
  benchstat-all                                      - Compare all benchmarks
  benchstat-by-category                              - Run and compare all benchmark categories
  benchstat-category-baseline                        - Run and compare baseline benchmarks
  benchstat-category-baseline-parallel               - Run and compare baseline benchmarks in parallel
  benchstat-category-populate                        - Run and compare populate benchmarks
  benchstat-category-readwrite                       - Run and compare readwrite benchmarks
  benchstat-category-readwrite-parallel              - Run and compare readwrite benchmarks in parallel
  benchstat-category-query-correlated                - Run and compare correlated query benchmarks
  benchstat-category-query-correlated-parallel       - Run and compare correlated query benchmarks in parallel
  benchstat-category-query-groupby                   - Run and compare groupby query benchmarks
  benchstat-category-query-groupby-parallel          - Run and compare groupby query benchmarks in parallel
  benchstat-category-query-json                      - Run and compare json query benchmarks
  benchstat-category-query-json-parallel             - Run and compare json query benchmarks in parallel
  benchstat-category-query-orderby                   - Run and compare orderby query benchmarks
  benchstat-category-query-orderby-parallel          - Run and compare orderby query benchmarks in parallel
  benchstat-category-query-recursivecte              - Run and compare recursivecte query benchmarks
  benchstat-category-query-recursivecte-parallel     - Run and compare recursivecte query benchmarks in parallel
  benchstat-category-query-window                    - Run and compare window query benchmarks
  benchstat-category-query-window-parallel           - Run and compare window query benchmarks in parallel
  clean                                              - Remove all benchmark, benchstat and test files
  test-all                                           - Run all tests
  update-all                                         - Update the quick/, slow/, and tests/ directories and README
  update-quick                                       - Update the quick/ directory
  update-slow                                        - Update the slow/ directory
  update-tests                                       - Update the tests/ directory
  update-readme                                      - Update the README

Variables:

  TAGS="ncruces_direct ncruces_driver eatonphil_direct glebarez_driver mattn_driver modernc_driver tailscale_driver zombiezen_direct"

    The first TAG listed becomes the baseline for benchstat comparisons.

  BENCH_COUNT=1
  BENCH_CPU=''
  BENCH_CPU_PARALLEL=1,2,4,8,16
  BENCH_OPTS="-benchmem -short"
  BENCH_PATTERN=.
  BENCH_SKIP=''
  BENCH_TIME=1s
  BENCH_TIMEOUT=15m

  BENCH_SLOW=

    BENCH_SLOW changes some of the values shown above. To see the BENCH_SLOW values,
    try "make help BENCH_SLOW=1"

  BENCH_MAX_READ_CONNECTIONS=0
  BENCH_MAX_WRITE_CONNECTIONS=32
  BENCH_POSTS=200
  BENCH_POST_PARAGRAPHS=10
  BENCH_COMMENTS=10
  BENCH_COMMENT_PARAGRAPHS=1

  BENCH_BIG=

    BENCH_BIG changes some of the values shown above. To see the BENCH_BIG values,
    try "make help BENCH_BIG=1"

  TEST_COUNT=1
  TEST_CPU=''
  TEST_OPTS="-short -v"
  TEST_PATTERN=.
  TEST_SKIP=''

Examples:

  make bench-all
  make benchstat-all

  make bench-all BENCH_COUNT=3
  make benchstat-all

  make bench-all BENCH_SLOW=1
  make benchstat-all

  make bench-all BENCH_BIG=1
  make benchstat-all

  make bench-all BENCH_SLOW=1 BENCH_BIG=1
  make benchstat-all

  make bench-all TAGS="ncruces_direct ncruces_driver"
  make benchstat-all TAGS="ncruces_direct ncruces_driver"

  make benchstat-by-category

  make benchstat-by-category BENCH_COUNT=3

  make benchstat-by-category BENCH_SLOW=1

  make benchstat-by-category BENCH_BIG=1

  make benchstat-by-category BENCH_SLOW=1 BENCH_BIG=1

  make benchstat-by-category TAGS="ncruces_direct ncruces_driver"

  make test-all

  make test-all TAGS="ncruces_direct ncruces_driver"

  make clean

```
<!--END_COMMAND-->

# Reviewing the Benchmarks

### What to Know

* **Configuration**

  All implementations use the following configuration. Some implementations use additional configuration.

      PRAGMA busy_timeout(5000);
      PRAGMA foreign_keys(true);
      PRAGMA journal_mode(wal);
      PRAGMA synchronous(normal);

* **Sequential vs. Parallel**

  Benchmarks that don't end with "Parallel" are sequential benchmarks; they run on a single goroutine, in a loop:

      for b.Loop() {
        ...
      }

  Benchmarks that end with "Parallel" are parallel benchmarks; they run on `GOMAXPROCS` goroutines, in a loop:

      b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
          ...
        }
      })

  "Sequential & Parallel Performance" runs have lines like these:

      Baseline/Conn-12
      Baseline/ConnParallel-12

  They show the same benchmark, "Conn", being run both sequentially and in parallel (with `GOMAXPROCS` set to 12 in this case, the default value on my laptop).

  "Parallel Performance" runs have lines like these:

      Baseline/ConnParallel
      Baseline/ConnParallel-2
      Baseline/ConnParallel-4
      Baseline/ConnParallel-8
      Baseline/ConnParallel-16

  They show the same benchmark, "ConnParallel", being run in parallel (with `GOMAXPROCS` set to 1, 2, 4, 8, and 16 in this case).

* **Queries**

  The query benchmarks were really meant to push the non-CGO-based implementations. They do that — but they also revealed unexpected differences between the CGO-based implementations.

### What to Look For

* **Sequential Performance**

  Sequential benchmarks answer the question: if I'm doing this *one* thing, how fast (or slow) will it be?

* **Parallel Performance**

  Parallel benchmarks answer the question: if I'm doing more than one of these things *at the same time*, how fast (or slow) will they be?

  Depending on the nature of the task, the answer might be "faster than sequential;" "the same as sequential;" or "slower than sequential."

* **Memory Usage**

  How much memory does one implementation use compared to another? Is the difference big enough to matter in your environment?

  (Just keep in mind that the benchmarks only measure memory usage that's visible to Go.)

## Baseline

* **Conn**

  Get one connection from the pool and then return it to the pool. (Every benchmark does that, but this benchmark does *only* that.)

* **Select1**

  Execute a "SELECT 1" query, preparing it each time.

* **Select1PrePrepared**

  Execute a "SELECT 1" query, preparing it an advance.

### - Sequential & Parallel Performance

<!--BENCHMARK:slow/benchstat_baseline.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                                       │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │         bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt        │
                                       │          sec/op          │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                  │
Baseline/Conn-12                                     189.2n ± ∞ ¹    224.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   187.9n ± ∞ ¹        ~ (p=0.100 n=3) ²    225.4n ± ∞ ¹        ~ (p=0.100 n=3) ²    222.5n ± ∞ ¹        ~ (p=0.100 n=3) ²    226.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    227.7n ± ∞ ¹       ~ (p=0.100 n=3) ²   1119.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-12                             233.1n ± ∞ ¹    534.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   223.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    542.5n ± ∞ ¹        ~ (p=0.100 n=3) ²    538.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    537.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    541.0n ± ∞ ¹       ~ (p=0.100 n=3) ²   1314.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1-12                                 1902.0n ± ∞ ¹   1947.0n ± ∞ ¹        ~ (p=0.200 n=3) ²   924.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   2001.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2635.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2024.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1118.0n ± ∞ ¹       ~ (p=0.100 n=3) ²   3110.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-12                         1239.0n ± ∞ ¹   1159.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   751.2n ± ∞ ¹        ~ (p=0.100 n=3) ²   2292.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1513.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2305.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    914.9n ± ∞ ¹       ~ (p=0.100 n=3) ²   2652.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePrepared-12                       504.4n ± ∞ ¹    467.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   487.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   1976.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1304.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2002.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    409.9n ± ∞ ¹       ~ (p=0.100 n=3) ²   1507.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-12               654.0n ± ∞ ¹    897.4n ± ∞ ¹        ~ (p=0.100 n=3) ²   477.8n ± ∞ ¹        ~ (p=0.100 n=3) ²   1787.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1108.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1784.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    774.5n ± ∞ ¹       ~ (p=0.100 n=3) ²   1533.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                              570.0n          695.6n        +22.04%                   435.1n        -23.67%                    1.121µ        +96.61%                    939.9n        +64.91%                    1.125µ        +97.31%                    584.8n        +2.60%                    1.743µ        +205.76%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                       │ bench_ncruces_direct.txt │    bench_ncruces_driver.txt    │      bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt    │     bench_mattn_driver.txt      │    bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt    │
                                       │           B/op           │    B/op      vs base           │    B/op      vs base                 │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │
Baseline/Conn-12                                       0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   353.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-12                               0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1-12                                   48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   705.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-12                           48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePrepared-12                         0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-12                 0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   353.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                         ⁴   46.15        ?                                +25.99%               ⁴    159.6        ?                    164.2        ?                    159.6        ?                    90.34        ?                    444.0        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                                       │ bench_ncruces_direct.txt │    bench_ncruces_driver.txt    │      bench_eatonphil_direct.txt      │   bench_glebarez_driver.txt    │     bench_mattn_driver.txt     │    bench_modernc_driver.txt    │   bench_tailscale_driver.txt   │   bench_zombiezen_direct.txt   │
                                       │        allocs/op         │  allocs/op   vs base           │  allocs/op   vs base                 │  allocs/op   vs base           │  allocs/op   vs base           │  allocs/op   vs base           │  allocs/op   vs base           │  allocs/op   vs base           │
Baseline/Conn-12                                      0.000 ± ∞ ¹   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹        ~ (p=1.000 n=3) ³   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-12                              0.000 ± ∞ ¹   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹        ~ (p=1.000 n=3) ³   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1-12                                   1.000 ± ∞ ¹   1.000 ± ∞ ¹  ~ (p=1.000 n=3) ³   4.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   9.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   8.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-12                           1.000 ± ∞ ¹   1.000 ± ∞ ¹  ~ (p=1.000 n=3) ³   4.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   9.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   8.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePrepared-12                        0.000 ± ∞ ¹   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹        ~ (p=1.000 n=3) ³   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-12                0.000 ± ∞ ¹   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹        ~ (p=1.000 n=3) ³   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                         ⁴   1.260        ?                                +58.74%               ⁴   3.659        ?                   3.979        ?                   3.659        ?                   2.154        ?                   5.848        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

### - Parallel Performance

<!--BENCHMARK:slow/benchstat_baseline_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                                       │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │         bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt        │
                                       │          sec/op          │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                  │
Baseline/ConnParallel                                186.4n ± ∞ ¹    250.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   188.1n ± ∞ ¹        ~ (p=0.200 n=3) ²    245.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    233.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    243.9n ± ∞ ¹        ~ (p=0.100 n=3) ²    233.6n ± ∞ ¹       ~ (p=0.100 n=3) ²   1351.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-2                              160.4n ± ∞ ¹    216.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   155.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    221.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    218.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    224.6n ± ∞ ¹        ~ (p=0.100 n=3) ²    218.9n ± ∞ ¹       ~ (p=0.100 n=3) ²    693.6n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-4                              217.1n ± ∞ ¹    266.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   205.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    269.3n ± ∞ ¹        ~ (p=0.100 n=3) ²    271.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    265.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    268.1n ± ∞ ¹       ~ (p=0.100 n=3) ²    685.7n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-8                              229.7n ± ∞ ¹    509.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   216.6n ± ∞ ¹        ~ (p=0.100 n=3) ²    509.9n ± ∞ ¹        ~ (p=0.100 n=3) ²    514.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    518.9n ± ∞ ¹        ~ (p=0.100 n=3) ²    516.0n ± ∞ ¹       ~ (p=0.100 n=3) ²   1247.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-16                             227.0n ± ∞ ¹    542.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   219.6n ± ∞ ¹        ~ (p=0.100 n=3) ²    538.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    547.6n ± ∞ ¹        ~ (p=0.100 n=3) ²    551.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    547.5n ± ∞ ¹       ~ (p=0.100 n=3) ²   1378.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel                            1933.0n ± ∞ ¹   1927.0n ± ∞ ¹        ~ (p=1.000 n=3) ²   935.6n ± ∞ ¹        ~ (p=0.100 n=3) ²   1898.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1647.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1902.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1156.0n ± ∞ ¹       ~ (p=0.100 n=3) ²   2858.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-2                          1057.0n ± ∞ ¹   1019.0n ± ∞ ¹        ~ (p=0.200 n=3) ²   550.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   1098.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1005.0n ± ∞ ¹        ~ (p=0.400 n=3) ²   1085.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    649.8n ± ∞ ¹       ~ (p=0.100 n=3) ²   1529.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-4                           622.6n ± ∞ ¹    563.6n ± ∞ ¹        ~ (p=0.200 n=3) ²   349.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1004.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    631.5n ± ∞ ¹        ~ (p=0.700 n=3) ²    983.3n ± ∞ ¹        ~ (p=0.100 n=3) ²    399.0n ± ∞ ¹       ~ (p=0.100 n=3) ²   1110.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-8                          1340.0n ± ∞ ¹    970.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   760.4n ± ∞ ¹        ~ (p=0.100 n=3) ²   2214.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1473.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2234.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    879.1n ± ∞ ¹       ~ (p=0.100 n=3) ²   2481.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-16                         1369.0n ± ∞ ¹   1024.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   766.5n ± ∞ ¹        ~ (p=0.100 n=3) ²   2269.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1493.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2295.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    926.8n ± ∞ ¹       ~ (p=0.100 n=3) ²   2703.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel                  498.8n ± ∞ ¹    468.5n ± ∞ ¹        ~ (p=0.100 n=3) ²   488.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   1868.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    858.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   1885.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    429.8n ± ∞ ¹       ~ (p=0.100 n=3) ²   1616.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-2                326.0n ± ∞ ¹    425.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   339.5n ± ∞ ¹        ~ (p=0.100 n=3) ²   1053.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    634.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   1043.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    354.4n ± ∞ ¹       ~ (p=0.100 n=3) ²    853.2n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-4                440.8n ± ∞ ¹    304.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   313.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    759.3n ± ∞ ¹        ~ (p=0.100 n=3) ²    524.7n ± ∞ ¹        ~ (p=0.400 n=3) ²    774.9n ± ∞ ¹        ~ (p=0.100 n=3) ²    288.0n ± ∞ ¹       ~ (p=0.100 n=3) ²    719.7n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-8                660.7n ± ∞ ¹    817.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   488.8n ± ∞ ¹        ~ (p=0.100 n=3) ²   1741.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1049.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1742.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    738.6n ± ∞ ¹       ~ (p=0.100 n=3) ²   1467.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-16               665.4n ± ∞ ¹    853.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   491.8n ± ∞ ¹        ~ (p=0.100 n=3) ²   1779.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1097.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1784.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    778.3n ± ∞ ¹       ~ (p=0.100 n=3) ²   1595.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                              493.2n          564.1n        +14.36%                   372.8n        -24.42%                    896.9n        +81.83%                    678.2n        +37.49%                    899.6n        +82.39%                    491.9n        -0.28%                    1.343µ        +172.28%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                       │ bench_ncruces_direct.txt │    bench_ncruces_driver.txt    │      bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt    │     bench_mattn_driver.txt      │    bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt    │
                                       │           B/op           │    B/op      vs base           │    B/op      vs base                 │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │
Baseline/ConnParallel                                  0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   381.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-2                                0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-4                                0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-8                                0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-16                               0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel                              48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   727.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-2                            48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-4                            48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-8                            48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-16                           48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel                    0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   373.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-2                  0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-4                  0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-8                  0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-16                 0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                         ⁴   46.15        ?                                +25.99%               ⁴    159.6        ?                    164.2        ?                    159.6        ?                    90.34        ?                    448.5        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                                       │ bench_ncruces_direct.txt │    bench_ncruces_driver.txt    │      bench_eatonphil_direct.txt      │   bench_glebarez_driver.txt    │     bench_mattn_driver.txt     │    bench_modernc_driver.txt    │   bench_tailscale_driver.txt   │   bench_zombiezen_direct.txt   │
                                       │        allocs/op         │  allocs/op   vs base           │  allocs/op   vs base                 │  allocs/op   vs base           │  allocs/op   vs base           │  allocs/op   vs base           │  allocs/op   vs base           │  allocs/op   vs base           │
Baseline/ConnParallel                                 0.000 ± ∞ ¹   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹        ~ (p=1.000 n=3) ³   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-2                               0.000 ± ∞ ¹   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹        ~ (p=1.000 n=3) ³   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-4                               0.000 ± ∞ ¹   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹        ~ (p=1.000 n=3) ³   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-8                               0.000 ± ∞ ¹   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹        ~ (p=1.000 n=3) ³   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-16                              0.000 ± ∞ ¹   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹        ~ (p=1.000 n=3) ³   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel                              1.000 ± ∞ ¹   1.000 ± ∞ ¹  ~ (p=1.000 n=3) ³   4.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   9.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   8.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-2                            1.000 ± ∞ ¹   1.000 ± ∞ ¹  ~ (p=1.000 n=3) ³   4.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   9.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   8.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-4                            1.000 ± ∞ ¹   1.000 ± ∞ ¹  ~ (p=1.000 n=3) ³   4.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   9.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   8.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-8                            1.000 ± ∞ ¹   1.000 ± ∞ ¹  ~ (p=1.000 n=3) ³   4.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   9.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   8.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-16                           1.000 ± ∞ ¹   1.000 ± ∞ ¹  ~ (p=1.000 n=3) ³   4.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   9.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   8.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel                   0.000 ± ∞ ¹   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹        ~ (p=1.000 n=3) ³   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-2                 0.000 ± ∞ ¹   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹        ~ (p=1.000 n=3) ³   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-4                 0.000 ± ∞ ¹   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹        ~ (p=1.000 n=3) ³   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-8                 0.000 ± ∞ ¹   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹        ~ (p=1.000 n=3) ³   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-16                0.000 ± ∞ ¹   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹        ~ (p=1.000 n=3) ³   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                         ⁴   1.260        ?                                +58.74%               ⁴   3.659        ?                   3.979        ?                   3.659        ?                   2.154        ?                   5.848        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

## Populate

* **PopulateDB**

  Populate the database without using any explicit transactions.

* **PopulateDBWithTx**

  Populate the database using one (big) explicit write transaction.

* **PopulateDBWithTxs**

  Populate the database using one explicit write transaction per post + comments.

### - Sequential Performance

<!--BENCHMARK:slow/benchstat_populate.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                              │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │      bench_zombiezen_direct.txt       │
                              │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                │
Populate/PopulateDB-12                       2.490 ± ∞ ¹     2.524 ± ∞ ¹       ~ (p=0.700 n=3) ²    2.050 ± ∞ ¹        ~ (p=0.100 n=3) ²     2.736 ± ∞ ¹        ~ (p=0.100 n=3) ²     2.219 ± ∞ ¹       ~ (p=0.100 n=3) ²     2.751 ± ∞ ¹        ~ (p=0.100 n=3) ²     2.038 ± ∞ ¹       ~ (p=0.100 n=3) ²     2.388 ± ∞ ¹       ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTx-12               1166.8m ± ∞ ¹   1170.1m ± ∞ ¹       ~ (p=1.000 n=3) ²   913.3m ± ∞ ¹        ~ (p=0.100 n=3) ²   1478.1m ± ∞ ¹        ~ (p=0.100 n=3) ²   1129.6m ± ∞ ¹       ~ (p=1.000 n=3) ²   1409.9m ± ∞ ¹        ~ (p=0.100 n=3) ²   1024.0m ± ∞ ¹       ~ (p=0.100 n=3) ²   1130.9m ± ∞ ¹       ~ (p=0.400 n=3) ²
Populate/PopulateDBWithTxs-12                1.115 ± ∞ ¹     1.381 ± ∞ ¹       ~ (p=0.100 n=3) ²    1.088 ± ∞ ¹        ~ (p=0.200 n=3) ²     1.398 ± ∞ ¹        ~ (p=0.100 n=3) ²     1.188 ± ∞ ¹       ~ (p=0.100 n=3) ²     1.391 ± ∞ ¹        ~ (p=0.100 n=3) ²     1.137 ± ∞ ¹       ~ (p=1.000 n=3) ²     1.194 ± ∞ ¹       ~ (p=0.100 n=3) ²
geomean                                      1.480           1.598        +7.99%                    1.268        -14.32%                     1.781        +20.41%                     1.439        -2.75%                     1.754        +18.56%                     1.334        -9.86%                     1.477        -0.15%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │         bench_ncruces_driver.txt          │       bench_eatonphil_direct.txt        │         bench_glebarez_driver.txt         │           bench_mattn_driver.txt            │         bench_modernc_driver.txt          │        bench_tailscale_driver.txt         │       bench_zombiezen_direct.txt        │
                              │           B/op           │      B/op       vs base                   │     B/op       vs base                  │      B/op       vs base                   │      B/op        vs base                    │      B/op       vs base                   │      B/op       vs base                   │     B/op       vs base                  │
Populate/PopulateDB-12                     2.362Mi ± ∞ ¹   18.968Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   5.771Mi ± ∞ ¹         ~ (p=0.100 n=3) ²   31.505Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   273.296Mi ± ∞ ¹           ~ (p=0.100 n=3) ²   31.492Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   18.963Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   6.167Mi ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTx-12               2.541Mi ± ∞ ¹   32.394Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   5.774Mi ± ∞ ¹         ~ (p=0.100 n=3) ²   44.496Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   286.530Mi ± ∞ ¹           ~ (p=0.100 n=3) ²   44.762Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   32.553Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   6.167Mi ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTxs-12              2.373Mi ± ∞ ¹   31.219Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   5.772Mi ± ∞ ¹         ~ (p=0.100 n=3) ²   43.667Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   286.004Mi ± ∞ ¹           ~ (p=0.100 n=3) ²   43.848Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   31.256Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   6.213Mi ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                    2.424Mi          26.77Mi        +1004.38%                   5.772Mi        +138.13%                    39.41Mi        +1525.94%                     281.9Mi        +11529.03%                    39.54Mi        +1531.19%                    26.82Mi        +1006.52%                   6.182Mi        +155.05%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │         bench_mattn_driver.txt          │        bench_modernc_driver.txt         │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                              │        allocs/op         │  allocs/op    vs base                  │  allocs/op    vs base                  │   allocs/op    vs base                  │   allocs/op    vs base                  │   allocs/op    vs base                  │  allocs/op    vs base                  │  allocs/op    vs base                  │
Populate/PopulateDB-12                      140.0k ± ∞ ¹   598.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   394.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1209.5k ± ∞ ¹         ~ (p=0.100 n=3) ²   1057.5k ± ∞ ¹         ~ (p=0.100 n=3) ²   1209.4k ± ∞ ¹         ~ (p=0.100 n=3) ²   598.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   445.1k ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTx-12                140.1k ± ∞ ¹   751.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   394.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1362.4k ± ∞ ¹         ~ (p=0.100 n=3) ²   1210.5k ± ∞ ¹         ~ (p=0.100 n=3) ²   1362.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   751.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   445.1k ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTxs-12               140.0k ± ∞ ¹   765.2k ± ∞ ¹         ~ (p=0.100 n=3) ²   394.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1376.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   1239.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1380.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   767.2k ± ∞ ¹         ~ (p=0.100 n=3) ²   448.1k ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                     140.1k         700.5k        +400.17%                   394.0k        +181.34%                    1.314M        +838.08%                    1.166M        +732.65%                    1.315M        +838.94%                   701.1k        +400.60%                   446.1k        +218.50%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
```
<!--END_BENCHMARK-->

## ReadWrite

* **ReadPost**

  Read one post, without a transaction.

* **ReadPostWithTx**

  Read one post, in a read transaction.

* **ReadPostAndComments**

  Read one post + comments, without a transaction.

* **ReadPostAndCommentsWithTx**

  Read one post + comments, in a read transaction.

* **WritePost**

  Write one post, without a transaction.

* **WritePostWithTx**

  Write one post, in a write transaction.

* **WritePostAndComments**

  Write one post + comments, without a transaction.

* **WritePostAndCommentsWithTx**

  Write one post + comments, in a write transaction.

* **ReadOrWritePostAndComments**

  Read or write one post + comments, without a transaction.

  ("Read or write" is determined by a "write rate" — either 10% or 90% in the current benchmarks.)

* **ReadOrWritePostAndCommentsWithTx**

  Read or write one post + comments, in a read or write transaction.

  ("Read or write" is determined by a "write rate" — either 10% or 90% in the current benchmarks.)

### - Sequential & Parallel Performance

<!--BENCHMARK:slow/benchstat_readwrite.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                                                                    │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │         bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                                                                    │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                 │
ReadWrite/ReadPost-12                                                            13.405µ ± ∞ ¹   21.673µ ± ∞ ¹       ~ (p=0.100 n=3) ²   9.511µ ± ∞ ¹        ~ (p=0.100 n=3) ²   36.905µ ± ∞ ¹        ~ (p=0.100 n=3) ²   17.720µ ± ∞ ¹        ~ (p=0.100 n=3) ²   36.533µ ± ∞ ¹        ~ (p=0.100 n=3) ²   16.470µ ± ∞ ¹       ~ (p=0.100 n=3) ²   25.550µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadPostWithTx-12                                                      15.449µ ± ∞ ¹   37.716µ ± ∞ ¹       ~ (p=0.100 n=3) ²   9.971µ ± ∞ ¹        ~ (p=0.100 n=3) ²   39.974µ ± ∞ ¹        ~ (p=0.100 n=3) ²   22.053µ ± ∞ ¹        ~ (p=0.100 n=3) ²   40.170µ ± ∞ ¹        ~ (p=0.100 n=3) ²   24.476µ ± ∞ ¹       ~ (p=0.100 n=3) ²   25.849µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndComments-12                                                  92.93µ ± ∞ ¹   113.07µ ± ∞ ¹       ~ (p=0.100 n=3) ²   72.39µ ± ∞ ¹        ~ (p=0.100 n=3) ²   177.97µ ± ∞ ¹        ~ (p=0.100 n=3) ²   173.02µ ± ∞ ¹        ~ (p=0.100 n=3) ²   183.85µ ± ∞ ¹        ~ (p=0.100 n=3) ²    98.54µ ± ∞ ¹       ~ (p=0.100 n=3) ²   102.35µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndCommentsWithTx-12                                            92.03µ ± ∞ ¹   141.16µ ± ∞ ¹       ~ (p=0.100 n=3) ²   70.75µ ± ∞ ¹        ~ (p=0.100 n=3) ²   179.73µ ± ∞ ¹        ~ (p=0.100 n=3) ²   176.75µ ± ∞ ¹        ~ (p=0.100 n=3) ²   185.43µ ± ∞ ¹        ~ (p=0.100 n=3) ²   110.15µ ± ∞ ¹       ~ (p=0.100 n=3) ²   101.50µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePost-12                                                            134.2µ ± ∞ ¹    119.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²   124.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    158.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    126.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    153.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    122.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²    176.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePostWithTx-12                                                      139.8µ ± ∞ ¹    128.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²   129.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    157.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    134.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    159.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    124.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²    176.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndComments-12                                                 2.618m ± ∞ ¹    2.542m ± ∞ ¹       ~ (p=0.100 n=3) ²   2.160m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.899m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.341m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.820m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.146m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.607m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndCommentsWithTx-12                                          1005.0µ ± ∞ ¹    968.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²   822.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1241.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    997.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1233.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    841.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²   1102.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=10-12                             430.2µ ± ∞ ¹    350.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²   311.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    461.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    400.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    451.4µ ± ∞ ¹        ~ (p=0.200 n=3) ²    303.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²    434.8µ ± ∞ ¹        ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=90-12                             2.409m ± ∞ ¹    2.315m ± ∞ ¹       ~ (p=0.100 n=3) ²   1.990m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.639m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.088m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.586m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.002m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.380m ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-12                     368.8µ ± ∞ ¹    375.6µ ± ∞ ¹       ~ (p=0.400 n=3) ²   299.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    414.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    336.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    429.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    300.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²    415.5µ ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-12                     3.158m ± ∞ ¹    3.022m ± ∞ ¹       ~ (p=0.700 n=3) ²   2.462m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.758m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.455m ± ∞ ¹        ~ (p=0.200 n=3) ²    2.880m ± ∞ ¹        ~ (p=0.200 n=3) ²    2.481m ± ∞ ¹       ~ (p=0.100 n=3) ²    3.386m ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=10-12                       258.1µ ± ∞ ¹    218.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²   195.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    287.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    270.0µ ± ∞ ¹        ~ (p=0.700 n=3) ²    309.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    182.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²    270.2µ ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=90-12                       970.3µ ± ∞ ¹    905.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²   760.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1162.7µ ± ∞ ¹        ~ (p=0.700 n=3) ²    937.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1177.9µ ± ∞ ¹        ~ (p=0.700 n=3) ²    770.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²   1049.1µ ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-12               111.8µ ± ∞ ¹    119.7µ ± ∞ ¹       ~ (p=0.200 n=3) ²   129.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    238.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    185.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    230.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    134.1µ ± ∞ ¹       ~ (p=0.100 n=3) ²    186.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-12               869.3µ ± ∞ ¹    887.8µ ± ∞ ¹       ~ (p=0.400 n=3) ²   940.8µ ± ∞ ¹        ~ (p=0.400 n=3) ²   3205.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1259.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   5038.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    926.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²   2418.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                                                           290.5µ          315.3µ        +8.57%                   239.2µ        -17.65%                    432.6µ        +48.93%                    329.5µ        +13.45%                    447.7µ        +54.13%                    273.0µ        -6.01%                    368.3µ        +26.80%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                                                    │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │    bench_eatonphil_direct.txt    │     bench_glebarez_driver.txt     │       bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                                                                    │           B/op           │      B/op       vs base           │     B/op       vs base           │      B/op       vs base           │      B/op        vs base           │      B/op       vs base           │      B/op       vs base           │     B/op       vs base           │
ReadWrite/ReadPost-12                                                            40.17Ki ± ∞ ¹    41.23Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.18Ki ± ∞ ¹  ~ (p=0.400 n=3) ²    81.35Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     41.39Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    81.30Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    41.19Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.56Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostWithTx-12                                                      40.17Ki ± ∞ ¹    42.29Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.18Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    82.22Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     42.60Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    82.36Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    42.13Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.61Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndComments-12                                                 250.2Ki ± ∞ ¹    260.3Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.2Ki ± ∞ ¹  ~ (p=0.600 n=3) ²    501.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     263.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    501.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    257.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndCommentsWithTx-12                                           250.1Ki ± ∞ ¹    261.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    503.1Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     265.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    503.3Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    258.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePost-12                                                              0.00 ± ∞ ¹     240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     48.00 ± ∞ ¹  ~ (p=0.400 n=3) ²     496.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    42288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     496.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    410.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostWithTx-12                                                        0.00 ± ∞ ¹     857.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    1104.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    43255.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    1290.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     880.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    458.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndComments-12                                                0.000Ki ± ∞ ¹   13.907Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   2.781Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.663Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   308.876Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.669Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   13.906Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   3.535Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndCommentsWithTx-12                                          0.000Ki ± ∞ ¹   26.448Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   2.781Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   39.193Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   321.743Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   39.380Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.465Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   3.579Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=10-12                            224.2Ki ± ∞ ¹    236.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   226.6Ki ± ∞ ¹  ~ (p=0.400 n=3) ²    453.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     268.4Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    455.5Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    233.0Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   224.7Ki ± ∞ ¹  ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=90-12                            22.70Ki ± ∞ ¹    38.12Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   27.95Ki ± ∞ ¹  ~ (p=1.000 n=3) ²    66.26Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    303.61Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    76.57Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    34.67Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   30.29Ki ± ∞ ¹  ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-12                    225.9Ki ± ∞ ¹    236.1Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.2Ki ± ∞ ¹  ~ (p=0.200 n=3) ²    456.5Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     268.4Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    453.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    232.4Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.4Ki ± ∞ ¹  ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-12                    30.05Ki ± ∞ ¹    40.05Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   27.09Ki ± ∞ ¹  ~ (p=0.700 n=3) ²    69.96Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    303.92Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    74.44Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    40.76Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.83Ki ± ∞ ¹  ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=10-12                      225.1Ki ± ∞ ¹    238.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.3Ki ± ∞ ¹  ~ (p=0.700 n=3) ²    457.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     271.4Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    452.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    234.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.4Ki ± ∞ ¹  ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=90-12                      23.22Ki ± ∞ ¹    50.64Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   30.00Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    81.15Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    315.79Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    90.18Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    51.17Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   28.61Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-12              225.2Ki ± ∞ ¹    237.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   226.0Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    455.3Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     271.4Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    457.3Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    235.1Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   226.5Ki ± ∞ ¹  ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-12              25.64Ki ± ∞ ¹    50.02Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   27.23Ki ± ∞ ¹  ~ (p=0.400 n=3) ²    82.79Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    315.55Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    83.97Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    48.76Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   32.62Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                                                      ³    42.31Ki        ?                   21.88Ki        ?                    75.99Ki        ?                     178.3Ki        ?                    78.31Ki        ?                    41.90Ki        ?                   30.12Ki        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean

                                                                    │ bench_ncruces_direct.txt │    bench_ncruces_driver.txt     │   bench_eatonphil_direct.txt    │    bench_glebarez_driver.txt     │     bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt    │
                                                                    │        allocs/op         │  allocs/op    vs base           │  allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op    vs base           │
ReadWrite/ReadPost-12                                                              5.000 ± ∞ ¹   39.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    43.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   42.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    42.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   31.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   12.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostWithTx-12                                                        5.000 ± ∞ ¹   60.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    57.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   70.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    60.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   50.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   15.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndComments-12                                                   262.0 ± ∞ ¹    773.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    264.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     976.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    684.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    562.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    272.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndCommentsWithTx-12                                             262.0 ± ∞ ¹    801.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    264.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     997.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    719.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1000.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    587.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    275.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePost-12                                                             0.000 ± ∞ ¹    7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    3.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    17.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   15.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    17.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    9.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostWithTx-12                                                       0.000 ± ∞ ¹   18.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    3.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    28.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   40.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    32.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   20.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   12.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndComments-12                                                    0.0 ± ∞ ¹    407.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    203.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     967.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    815.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     967.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    407.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    259.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndCommentsWithTx-12                                              0.0 ± ∞ ¹    574.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    203.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1134.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    996.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1138.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    576.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    262.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=10-12                              234.0 ± ∞ ¹    737.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    697.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=90-12                              23.00 ± ∞ ¹   443.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   799.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   420.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-12                      236.0 ± ∞ ¹    737.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    697.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-12                      30.00 ± ∞ ¹   445.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   800.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   424.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=10-12                        235.0 ± ∞ ¹    778.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1010.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    747.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1015.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    274.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=90-12                        24.00 ± ∞ ¹   597.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1121.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   966.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1123.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-12                235.0 ± ∞ ¹    777.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1011.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    747.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    274.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-12                26.00 ± ∞ ¹   596.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1121.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   965.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1125.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                                                      ³    273.5        ?                    85.27        ?                     431.1        ?                    368.1        ?                     436.0        ?                    237.4        ?                    122.4        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

### - Parallel Performance

<!--BENCHMARK:slow/benchstat_readwrite_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                                                                    │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt         │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt        │
                                                                    │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │     sec/op      vs base                 │    sec/op      vs base                │     sec/op      vs base                 │    sec/op      vs base                 │     sec/op      vs base                 │
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10                        451.0µ ± ∞ ¹    346.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²   319.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²     428.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    347.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²     465.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    287.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²     442.1µ ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-2                      374.1µ ± ∞ ¹    336.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²   280.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²     389.2µ ± ∞ ¹        ~ (p=0.400 n=3) ²    310.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²     448.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    256.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²     378.2µ ± ∞ ¹        ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-4                      334.0µ ± ∞ ¹    352.5µ ± ∞ ¹       ~ (p=0.700 n=3) ²   286.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²     365.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    297.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²     389.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    271.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²     364.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-8                      373.3µ ± ∞ ¹    350.2µ ± ∞ ¹       ~ (p=0.400 n=3) ²   288.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²     408.8µ ± ∞ ¹        ~ (p=0.200 n=3) ²    326.1µ ± ∞ ¹       ~ (p=0.100 n=3) ²     418.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    268.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²     356.8µ ± ∞ ¹        ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-16                     427.3µ ± ∞ ¹    430.7µ ± ∞ ¹       ~ (p=0.700 n=3) ²   331.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²     441.8µ ± ∞ ¹        ~ (p=0.400 n=3) ²    339.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²     481.7µ ± ∞ ¹        ~ (p=0.400 n=3) ²    307.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²     387.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90                        2.413m ± ∞ ¹    2.294m ± ∞ ¹       ~ (p=0.100 n=3) ²   1.946m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.446m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.042m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.638m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.022m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.386m ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-2                      2.562m ± ∞ ¹    2.446m ± ∞ ¹       ~ (p=0.100 n=3) ²   1.935m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.607m ± ∞ ¹        ~ (p=1.000 n=3) ²    2.063m ± ∞ ¹       ~ (p=0.100 n=3) ²     3.033m ± ∞ ¹        ~ (p=0.100 n=3) ²    1.989m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.418m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-4                      2.705m ± ∞ ¹    2.685m ± ∞ ¹       ~ (p=0.700 n=3) ²   2.086m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.549m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.105m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.899m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.091m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.378m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-8                      2.836m ± ∞ ¹    2.942m ± ∞ ¹       ~ (p=0.700 n=3) ²   2.301m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.912m ± ∞ ¹        ~ (p=0.700 n=3) ²    2.344m ± ∞ ¹       ~ (p=0.100 n=3) ²     3.115m ± ∞ ¹        ~ (p=0.700 n=3) ²    2.222m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.522m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-16                     3.367m ± ∞ ¹    3.339m ± ∞ ¹       ~ (p=0.400 n=3) ²   2.855m ± ∞ ¹        ~ (p=0.100 n=3) ²     3.425m ± ∞ ¹        ~ (p=0.700 n=3) ²    2.818m ± ∞ ¹       ~ (p=0.100 n=3) ²     3.958m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.734m ± ∞ ¹        ~ (p=0.100 n=3) ²     8.487m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10                  275.0µ ± ∞ ¹    218.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²   190.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²     288.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    230.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²     307.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    181.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²     288.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-2                163.4µ ± ∞ ¹    138.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²   144.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²     222.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    171.8µ ± ∞ ¹       ~ (p=0.400 n=3) ²     220.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    130.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²     199.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-4                119.8µ ± ∞ ¹    119.2µ ± ∞ ¹       ~ (p=1.000 n=3) ²   121.9µ ± ∞ ¹        ~ (p=0.700 n=3) ²     199.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    142.9µ ± ∞ ¹       ~ (p=0.100 n=3) ²     195.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    115.4µ ± ∞ ¹        ~ (p=0.700 n=3) ²     175.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-8                110.3µ ± ∞ ¹    110.3µ ± ∞ ¹       ~ (p=1.000 n=3) ²   121.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²     209.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    180.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²     223.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    126.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²     188.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-16               122.0µ ± ∞ ¹    129.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²   131.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²     240.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    222.1µ ± ∞ ¹       ~ (p=0.100 n=3) ²     234.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    136.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²     204.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90                  953.2µ ± ∞ ¹    922.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²   781.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1128.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    884.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²    1158.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    774.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1060.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-2                856.6µ ± ∞ ¹    776.1µ ± ∞ ¹       ~ (p=0.200 n=3) ²   764.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1143.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    891.3µ ± ∞ ¹       ~ (p=0.400 n=3) ²    1222.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    752.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1005.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-4                774.0µ ± ∞ ¹    788.0µ ± ∞ ¹       ~ (p=0.200 n=3) ²   760.0µ ± ∞ ¹        ~ (p=1.000 n=3) ²    1205.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    956.6µ ± ∞ ¹       ~ (p=0.100 n=3) ²    1246.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    791.6µ ± ∞ ¹        ~ (p=0.200 n=3) ²    1067.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-8                816.5µ ± ∞ ¹    829.4µ ± ∞ ¹       ~ (p=0.400 n=3) ²   810.5µ ± ∞ ¹        ~ (p=1.000 n=3) ²    1264.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1040.9µ ± ∞ ¹       ~ (p=0.100 n=3) ²    1580.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    846.2µ ± ∞ ¹        ~ (p=0.200 n=3) ²    1120.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-16              1012.0µ ± ∞ ¹   1038.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²   878.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   10245.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1545.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²   12614.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1029.6µ ± ∞ ¹        ~ (p=0.700 n=3) ²   10587.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                                                           612.0µ          586.0µ        -4.24%                   521.4µ        -14.80%                     828.6µ        +35.40%                    613.5µ        +0.25%                     897.3µ        +46.62%                    513.1µ        -16.16%                     804.9µ        +31.53%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                                                    │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │          bench_mattn_driver.txt          │        bench_modernc_driver.txt         │       bench_tailscale_driver.txt       │      bench_zombiezen_direct.txt       │
                                                                    │           B/op           │     B/op       vs base                 │     B/op       vs base                │     B/op       vs base                  │      B/op       vs base                  │     B/op       vs base                  │     B/op       vs base                 │     B/op       vs base                │
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10                       224.2Ki ± ∞ ¹   235.4Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.1Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   456.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   453.7Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   232.9Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.0Ki ± ∞ ¹       ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-2                     224.9Ki ± ∞ ¹   235.2Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.6Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   456.9Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   452.8Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   233.6Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.1Ki ± ∞ ¹       ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-4                     225.6Ki ± ∞ ¹   235.6Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.1Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   454.5Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   453.9Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   232.1Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.4Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-8                     224.4Ki ± ∞ ¹   236.6Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.1Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   455.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.5Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   452.7Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   233.1Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   227.4Ki ± ∞ ¹       ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-16                    226.0Ki ± ∞ ¹   235.7Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.0Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   453.5Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   451.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   232.1Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.0Ki ± ∞ ¹       ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90                       26.59Ki ± ∞ ¹   38.10Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.00Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   75.05Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   303.58Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   73.14Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   37.50Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   25.94Ki ± ∞ ¹       ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-2                     21.52Ki ± ∞ ¹   40.47Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.25Ki ± ∞ ¹       ~ (p=0.200 n=3) ²   69.50Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.60Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   73.55Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   37.43Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   26.62Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-4                     27.12Ki ± ∞ ¹   34.83Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   29.78Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   76.10Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.21Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   65.67Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   39.57Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   30.30Ki ± ∞ ¹       ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-8                     28.49Ki ± ∞ ¹   35.98Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.95Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   71.65Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.48Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   80.66Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   40.79Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   25.03Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-16                    23.55Ki ± ∞ ¹   40.41Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   29.02Ki ± ∞ ¹       ~ (p=0.200 n=3) ²   75.05Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.03Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   71.53Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   40.43Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   23.36Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10                 224.4Ki ± ∞ ¹   237.7Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.5Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   459.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.0Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   456.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   235.8Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.0Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-2               224.5Ki ± ∞ ¹   238.3Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.2Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   454.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   456.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   234.8Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.8Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-4               225.5Ki ± ∞ ¹   237.8Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.9Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   456.6Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   455.8Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   235.6Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.2Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-8               224.5Ki ± ∞ ¹   239.3Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.7Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   457.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   455.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   235.1Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.9Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-16              226.0Ki ± ∞ ¹   238.5Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.7Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   453.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   457.8Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   234.9Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.2Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90                 25.66Ki ± ∞ ¹   50.97Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   29.67Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   89.82Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.24Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   81.22Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   48.96Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.44Ki ± ∞ ¹       ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-2               23.68Ki ± ∞ ¹   50.13Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.57Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   88.58Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   315.94Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   91.25Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   50.02Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   25.19Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-4               25.52Ki ± ∞ ¹   49.48Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.83Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   84.91Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.09Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   84.41Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   48.91Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.50Ki ± ∞ ¹       ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-8               24.92Ki ± ∞ ¹   50.41Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   25.55Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   88.20Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   315.39Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   80.43Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   51.00Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.36Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-16              24.14Ki ± ∞ ¹   49.54Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   26.50Ki ± ∞ ¹       ~ (p=0.200 n=3) ²   86.19Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.02Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   81.17Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   51.52Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   23.50Ki ± ∞ ¹       ~ (p=1.000 n=3) ²
geomean                                                                          75.07Ki         101.6Ki        +35.38%                   79.26Ki        +5.58%                   191.1Ki        +154.61%                    289.2Ki        +285.24%                   188.3Ki        +150.83%                   101.8Ki        +35.57%                   76.98Ki        +2.55%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                                                    │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt         │         bench_mattn_driver.txt         │         bench_modernc_driver.txt         │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                                                                    │        allocs/op         │  allocs/op    vs base                  │  allocs/op    vs base                  │   allocs/op    vs base                   │  allocs/op    vs base                  │   allocs/op    vs base                   │  allocs/op    vs base                  │  allocs/op    vs base                  │
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10                         234.0 ± ∞ ¹    734.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     973.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    695.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     973.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    544.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    270.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-2                       235.0 ± ∞ ¹    734.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    696.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    270.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-4                       236.0 ± ∞ ¹    736.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    696.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-8                       235.0 ± ∞ ¹    737.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    697.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-16                      236.0 ± ∞ ¹    736.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    696.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90                         27.00 ± ∞ ¹   442.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   799.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   421.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-2                       22.00 ± ∞ ¹   446.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   802.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   422.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-4                       28.00 ± ∞ ¹   438.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   801.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   423.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-8                       29.00 ± ∞ ¹   439.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   802.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   424.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-16                      24.00 ± ∞ ¹   446.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   801.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   424.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10                   235.0 ± ∞ ¹    777.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1009.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    745.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    585.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-2                 235.0 ± ∞ ¹    777.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1011.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    585.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-4                 235.0 ± ∞ ¹    777.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1010.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    745.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-8                 235.0 ± ∞ ¹    779.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1011.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1015.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-16                236.0 ± ∞ ¹    778.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1012.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    747.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90                   26.00 ± ∞ ¹   597.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1119.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   968.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1125.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-2                 24.00 ± ∞ ¹   596.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1119.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   967.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1122.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-4                 26.00 ± ∞ ¹   596.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1120.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   968.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1124.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-8                 26.00 ± ∞ ¹   597.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1119.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   964.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1126.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-16                25.00 ± ∞ ¹   596.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1121.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   967.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1126.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   264.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                                                            77.63          623.2        +702.71%                    231.7        +198.49%                    1.016k        +1209.09%                    796.3        +925.65%                    1.018k        +1211.58%                    528.4        +580.61%                    266.6        +243.47%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
```
<!--END_BENCHMARK-->

## Query - Correlated

<!--SQL:gsb_common_sql_query_correlated.sql-->
```sql
SELECT
  id,
  title,
  (SELECT COUNT(*) FROM comments WHERE post_id = posts.id) as comment_count,
  (SELECT AVG(LENGTH(content)) FROM comments WHERE post_id = posts.id) AS avg_comment_length,
  (SELECT MAX(LENGTH(content)) FROM comments WHERE post_id = posts.id) AS max_comment_length
FROM posts
```
<!--END_SQL-->

### - Sequential & Parallel Performance

<!--BENCHMARK:slow/benchstat_query_correlated.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                            │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                            │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │
Query/Correlated-12                       635.3m ± ∞ ¹    646.0m ± ∞ ¹       ~ (p=0.100 n=3) ²   288.5m ± ∞ ¹        ~ (p=0.100 n=3) ²    304.6m ± ∞ ¹        ~ (p=0.100 n=3) ²   176.5m ± ∞ ¹        ~ (p=0.100 n=3) ²    303.8m ± ∞ ¹        ~ (p=0.100 n=3) ²   173.6m ± ∞ ¹        ~ (p=0.100 n=3) ²    323.4m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-12              120.94m ± ∞ ¹   120.86m ± ∞ ¹       ~ (p=1.000 n=3) ²   50.85m ± ∞ ¹        ~ (p=0.100 n=3) ²   123.25m ± ∞ ¹        ~ (p=0.400 n=3) ²   46.98m ± ∞ ¹        ~ (p=0.100 n=3) ²   120.70m ± ∞ ¹        ~ (p=1.000 n=3) ²   44.85m ± ∞ ¹        ~ (p=0.100 n=3) ²   186.00m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                   277.2m          279.4m        +0.81%                   121.1m        -56.30%                    193.8m        -30.10%                   91.05m        -67.15%                    191.5m        -30.92%                   88.24m        -68.17%                    245.3m        -11.51%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │         bench_ncruces_driver.txt         │      bench_eatonphil_direct.txt      │        bench_glebarez_driver.txt         │          bench_mattn_driver.txt           │         bench_modernc_driver.txt         │       bench_tailscale_driver.txt        │      bench_zombiezen_direct.txt       │
                            │           B/op           │     B/op       vs base                   │    B/op      vs base                 │     B/op       vs base                   │      B/op       vs base                   │     B/op       vs base                   │     B/op       vs base                  │     B/op       vs base                │
Query/Correlated-12                      57168.0 ± ∞ ¹   71264.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   732.0 ± ∞ ¹        ~ (p=0.100 n=3) ²   63072.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   207197.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   63160.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   47021.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    9280.0 ± ∞ ¹       ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-12                335.0 ± ∞ ¹   71391.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   157.0 ± ∞ ¹        ~ (p=0.700 n=3) ²   64257.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   207526.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   63265.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   47773.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1788.0 ± ∞ ¹       ~ (p=0.700 n=3) ²
geomean                                  4.274Ki         69.66Ki        +1529.89%                   339.0        -92.25%                   62.17Ki        +1354.72%                    202.5Ki        +4638.37%                   61.73Ki        +1344.45%                   46.28Ki        +983.02%                   3.978Ki        -6.92%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt           │      bench_eatonphil_direct.txt      │          bench_glebarez_driver.txt          │           bench_mattn_driver.txt            │          bench_modernc_driver.txt           │         bench_tailscale_driver.txt         │       bench_zombiezen_direct.txt       │
                            │        allocs/op         │   allocs/op     vs base                     │  allocs/op   vs base                 │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                    │  allocs/op    vs base                  │
Query/Correlated-12                       12.000 ± ∞ ¹   5770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹           ~ (p=0.100 n=3) ²   30.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-12                2.000 ± ∞ ¹   5770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.300 n=3) ²   6788.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6773.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3768.000 ± ∞ ¹           ~ (p=0.100 n=3) ²   15.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                    4.899           5.770k        +117679.63%                   2.449        -50.00%                     6.778k        +138265.42%                     6.772k        +138132.87%                     6.771k        +138112.46%                     3.767k        +76793.56%                    21.21        +333.01%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
```
<!--END_BENCHMARK-->

### - Parallel Performance

<!--BENCHMARK:slow/benchstat_query_correlated_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                            │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                            │          sec/op          │    sec/op      vs base                │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │
Query/CorrelatedParallel                  645.9m ± ∞ ¹    647.5m ± ∞ ¹       ~ (p=1.000 n=3) ²    289.0m ± ∞ ¹        ~ (p=0.100 n=3) ²    302.7m ± ∞ ¹        ~ (p=0.100 n=3) ²   174.8m ± ∞ ¹        ~ (p=0.100 n=3) ²    302.6m ± ∞ ¹        ~ (p=0.100 n=3) ²   172.6m ± ∞ ¹        ~ (p=0.100 n=3) ²    327.1m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-2               443.68m ± ∞ ¹   436.45m ± ∞ ¹       ~ (p=0.100 n=3) ²   170.37m ± ∞ ¹        ~ (p=0.100 n=3) ²   184.01m ± ∞ ¹        ~ (p=0.100 n=3) ²   96.78m ± ∞ ¹        ~ (p=0.100 n=3) ²   184.50m ± ∞ ¹        ~ (p=0.100 n=3) ²   95.84m ± ∞ ¹        ~ (p=0.100 n=3) ²   189.30m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-4               227.99m ± ∞ ¹   225.66m ± ∞ ¹       ~ (p=0.200 n=3) ²    84.66m ± ∞ ¹        ~ (p=0.100 n=3) ²   117.50m ± ∞ ¹        ~ (p=0.100 n=3) ²   58.08m ± ∞ ¹        ~ (p=0.100 n=3) ²   120.16m ± ∞ ¹        ~ (p=0.100 n=3) ²   56.86m ± ∞ ¹        ~ (p=0.100 n=3) ²   123.51m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-8               153.55m ± ∞ ¹   154.95m ± ∞ ¹       ~ (p=0.400 n=3) ²    53.84m ± ∞ ¹        ~ (p=0.100 n=3) ²   138.45m ± ∞ ¹        ~ (p=0.100 n=3) ²   43.37m ± ∞ ¹        ~ (p=0.100 n=3) ²   138.67m ± ∞ ¹        ~ (p=0.100 n=3) ²   42.05m ± ∞ ¹        ~ (p=0.100 n=3) ²   187.21m ± ∞ ¹        ~ (p=0.700 n=3) ²
Query/CorrelatedParallel-16               78.47m ± ∞ ¹    78.53m ± ∞ ¹       ~ (p=1.000 n=3) ²    53.68m ± ∞ ¹        ~ (p=0.100 n=3) ²   123.65m ± ∞ ¹        ~ (p=0.100 n=3) ²   49.66m ± ∞ ¹        ~ (p=0.100 n=3) ²   117.56m ± ∞ ¹        ~ (p=0.100 n=3) ²   44.19m ± ∞ ¹        ~ (p=0.100 n=3) ²   184.86m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                   239.4m          238.8m        -0.29%                    103.8m        -56.65%                    162.1m        -32.29%                   73.31m        -69.38%                    161.4m        -32.61%                   70.55m        -70.53%                    192.5m        -19.59%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │         bench_ncruces_driver.txt          │      bench_eatonphil_direct.txt      │         bench_glebarez_driver.txt         │           bench_mattn_driver.txt           │         bench_modernc_driver.txt          │        bench_tailscale_driver.txt        │      bench_zombiezen_direct.txt       │
                            │           B/op           │      B/op       vs base                   │    B/op      vs base                 │      B/op       vs base                   │      B/op        vs base                   │      B/op       vs base                   │      B/op       vs base                  │     B/op      vs base                 │
Query/CorrelatedParallel                 57140.0 ± ∞ ¹    71132.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   718.0 ± ∞ ¹        ~ (p=0.100 n=3) ²    63050.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    207185.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    63066.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    47009.0 ± ∞ ¹         ~ (p=0.100 n=3) ²   2178.0 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-2              57157.00 ± ∞ ¹   71082.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   49.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   63013.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   207154.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   63029.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   47002.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   492.00 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-4              11745.00 ± ∞ ¹   71028.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   48.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   63078.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   207175.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   63080.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   46991.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   464.00 ± ∞ ¹        ~ (p=0.700 n=3) ²
Query/CorrelatedParallel-8                 388.0 ± ∞ ¹    71304.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   208.0 ± ∞ ¹        ~ (p=0.400 n=3) ²    63190.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    207386.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    63250.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    47204.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    858.0 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-16                118.0 ± ∞ ¹    72441.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   625.0 ± ∞ ¹        ~ (p=0.300 n=3) ²    63888.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    207687.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    63420.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    47213.0 ± ∞ ¹         ~ (p=0.100 n=3) ²   1782.0 ± ∞ ¹        ~ (p=0.200 n=3) ²
geomean                                  4.351Ki          69.72Ki        +1502.34%                   185.5        -95.84%                    61.76Ki        +1319.37%                     202.5Ki        +4552.86%                    61.69Ki        +1317.71%                    45.98Ki        +956.71%                    946.6        -78.75%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt           │      bench_eatonphil_direct.txt      │          bench_glebarez_driver.txt          │           bench_mattn_driver.txt            │          bench_modernc_driver.txt           │         bench_tailscale_driver.txt         │       bench_zombiezen_direct.txt       │
                            │        allocs/op         │   allocs/op     vs base                     │  allocs/op   vs base                 │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                    │  allocs/op    vs base                  │
Query/CorrelatedParallel                  14.000 ± ∞ ¹   5771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹           ~ (p=0.100 n=3) ²   30.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-2                14.000 ± ∞ ¹   5770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹           ~ (p=0.100 n=3) ²    8.000 ± ∞ ¹         ~ (p=0.700 n=3) ²
Query/CorrelatedParallel-4                 4.000 ± ∞ ¹   5769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹           ~ (p=0.100 n=3) ²    7.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-8                 3.000 ± ∞ ¹   5770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹           ~ (p=0.100 n=3) ²   10.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-16                2.000 ± ∞ ¹   5773.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹        ~ (p=1.000 n=3) ²   6782.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6774.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6774.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3769.000 ± ∞ ¹           ~ (p=0.100 n=3) ²   15.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                    5.426           5.771k        +106247.55%                   1.695        -68.76%                     6.772k        +124709.90%                     6.772k        +124706.24%                     6.771k        +124684.12%                     3.767k        +69330.18%                    12.03        +121.71%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
```
<!--END_BENCHMARK-->

## Query - GroupBy

<!--SQL:gsb_common_sql_query_groupby.sql-->
```sql
SELECT
  strftime('%Y-%m', created) AS month,
  COUNT(*) as month_total
FROM posts
GROUP BY month
ORDER BY month
```
<!--END_SQL-->

### - Sequential & Parallel Performance

<!--BENCHMARK:slow/benchstat_query_groupby.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                         │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │      bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │         bench_mattn_driver.txt         │        bench_modernc_driver.txt         │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt        │
                         │          sec/op          │    sec/op     vs base                │    sec/op     vs base                 │    sec/op      vs base                  │    sec/op      vs base                 │    sec/op      vs base                  │    sec/op     vs base                 │    sec/op      vs base                  │
Query/GroupBy-12                       569.5µ ± ∞ ¹   582.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²   277.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    489.2µ ± ∞ ¹         ~ (p=0.100 n=3) ²    333.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    498.1µ ± ∞ ¹         ~ (p=0.100 n=3) ²   261.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    479.1µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-12               62.25µ ± ∞ ¹   66.33µ ± ∞ ¹       ~ (p=0.100 n=3) ²   38.59µ ± ∞ ¹        ~ (p=0.100 n=3) ²   687.89µ ± ∞ ¹         ~ (p=0.100 n=3) ²   404.86µ ± ∞ ¹        ~ (p=0.100 n=3) ²   691.06µ ± ∞ ¹         ~ (p=0.100 n=3) ²   28.07µ ± ∞ ¹        ~ (p=0.100 n=3) ²   689.54µ ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                188.3µ         196.6µ        +4.42%                   103.4µ        -45.09%                    580.1µ        +208.10%                    367.3µ        +95.07%                    586.7µ        +211.61%                   85.64µ        -54.51%                    574.8µ        +205.28%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │    bench_eatonphil_direct.txt     │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │    bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                         │           B/op           │     B/op       vs base           │   B/op     vs base                │     B/op       vs base           │     B/op       vs base           │     B/op       vs base           │     B/op       vs base           │    B/op      vs base           │
Query/GroupBy-12                          0.0 ± ∞ ¹    2340.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ²    1932.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    6989.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1944.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1558.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   361.0 ± ∞ ¹  ~ (p=0.600 n=3) ²
Query/GroupByParallel-12                  0.0 ± ∞ ¹    2328.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1935.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    6995.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1948.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1503.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   370.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ⁴   2.279Ki        ?                              +0.00%               ⁴   1.888Ki        ?                   6.828Ki        ?                   1.900Ki        ?                   1.494Ki        ?                   365.5        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                         │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/GroupBy-12                        0.000 ± ∞ ¹   116.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   119.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   151.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   119.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-12                0.000 ± ∞ ¹   115.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   119.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   151.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   119.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   50.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ⁴     115.5        ?                                +0.00%               ⁴     119.0        ?                     151.0        ?                     119.0        ?                    50.50        ?                   6.000        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

### - Parallel Performance

<!--BENCHMARK:slow/benchstat_query_groupby_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                         │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │         bench_mattn_driver.txt          │        bench_modernc_driver.txt         │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt        │
                         │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op     vs base                 │    sec/op      vs base                  │
Query/GroupByParallel                  566.0µ ± ∞ ¹    562.1µ ± ∞ ¹       ~ (p=0.200 n=3) ²   287.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    484.5µ ± ∞ ¹         ~ (p=0.100 n=3) ²    301.7µ ± ∞ ¹         ~ (p=0.100 n=3) ²    483.4µ ± ∞ ¹         ~ (p=0.100 n=3) ²   264.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    467.3µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-2                285.1µ ± ∞ ¹    287.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²   146.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    384.0µ ± ∞ ¹         ~ (p=0.100 n=3) ²    370.9µ ± ∞ ¹         ~ (p=0.100 n=3) ²    384.9µ ± ∞ ¹         ~ (p=0.100 n=3) ²   130.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    378.6µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-4               145.70µ ± ∞ ¹   147.83µ ± ∞ ¹       ~ (p=0.100 n=3) ²   79.85µ ± ∞ ¹        ~ (p=0.100 n=3) ²   398.64µ ± ∞ ¹         ~ (p=0.100 n=3) ²   448.10µ ± ∞ ¹         ~ (p=0.100 n=3) ²   390.96µ ± ∞ ¹         ~ (p=0.100 n=3) ²   65.92µ ± ∞ ¹        ~ (p=0.100 n=3) ²   393.71µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-8                73.51µ ± ∞ ¹    75.90µ ± ∞ ¹       ~ (p=0.100 n=3) ²   42.98µ ± ∞ ¹        ~ (p=0.100 n=3) ²   651.60µ ± ∞ ¹         ~ (p=0.100 n=3) ²   453.48µ ± ∞ ¹         ~ (p=0.100 n=3) ²   656.66µ ± ∞ ¹         ~ (p=0.100 n=3) ²   34.31µ ± ∞ ¹        ~ (p=0.100 n=3) ²   651.81µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-16               61.75µ ± ∞ ¹    64.71µ ± ∞ ¹       ~ (p=0.100 n=3) ²   37.65µ ± ∞ ¹        ~ (p=0.100 n=3) ²   687.56µ ± ∞ ¹         ~ (p=0.100 n=3) ²   414.56µ ± ∞ ¹         ~ (p=0.100 n=3) ²   688.05µ ± ∞ ¹         ~ (p=0.100 n=3) ²   28.35µ ± ∞ ¹        ~ (p=0.100 n=3) ²   692.11µ ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                160.6µ          163.6µ        +1.91%                   88.47µ        -44.90%                    506.2µ        +215.25%                    393.4µ        +145.04%                    505.1µ        +214.57%                   73.98µ        -53.92%                    500.6µ        +211.76%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │    bench_eatonphil_direct.txt     │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │    bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                         │           B/op           │     B/op       vs base           │   B/op     vs base                │     B/op       vs base           │     B/op       vs base           │     B/op       vs base           │     B/op       vs base           │    B/op      vs base           │
Query/GroupByParallel                     0.0 ± ∞ ¹    2324.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1817.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    6984.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1832.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1557.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   362.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-2                   0.0 ± ∞ ¹    2324.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1928.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    6979.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1874.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1557.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   360.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-4                   0.0 ± ∞ ¹    2324.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1930.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    6951.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1923.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1516.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   361.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-8                   0.0 ± ∞ ¹    2326.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1932.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    6982.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1948.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1547.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   362.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-16                  0.0 ± ∞ ¹    2326.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ²    1938.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    6988.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1953.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1490.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   368.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ⁴   2.270Ki        ?                              +0.00%               ⁴   1.864Ki        ?                   6.813Ki        ?                   1.861Ki        ?                   1.497Ki        ?                   362.6        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                         │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/GroupByParallel                   0.000 ± ∞ ¹   115.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   117.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   151.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   117.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-2                 0.000 ± ∞ ¹   115.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   119.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   150.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-4                 0.000 ± ∞ ¹   115.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   119.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   150.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   50.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-8                 0.000 ± ∞ ¹   115.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   119.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   150.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   119.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   50.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-16                0.000 ± ∞ ¹   115.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   119.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   151.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   119.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   50.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ⁴     115.0        ?                                +0.00%               ⁴     118.6        ?                     150.4        ?                     118.2        ?                    50.40        ?                   6.000        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

## Query - JSON

<!--SQL:gsb_common_sql_query_json.sql-->
```sql
SELECT
  date(created) as day,
  SUM(json_extract(stats, '$.lorem')) as sum_lorem,
  AVG(json_extract(stats, '$.ipsum.dolor')) as avg_dolor,
  MAX(json_extract(stats, '$.lorem.sit')) as max_sit
FROM posts
GROUP BY day
ORDER BY day
```
<!--END_SQL-->

### - Sequential & Parallel Performance

<!--BENCHMARK:slow/benchstat_query_json.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                      │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt        │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                      │          sec/op          │    sec/op     vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │
Query/JSON-12                       9.797m ± ∞ ¹   9.787m ± ∞ ¹       ~ (p=1.000 n=3) ²   7.264m ± ∞ ¹        ~ (p=0.100 n=3) ²    9.312m ± ∞ ¹        ~ (p=0.100 n=3) ²   8.260m ± ∞ ¹       ~ (p=0.100 n=3) ²    9.375m ± ∞ ¹        ~ (p=0.100 n=3) ²   7.055m ± ∞ ¹        ~ (p=0.100 n=3) ²   10.277m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-12               8.500m ± ∞ ¹   9.483m ± ∞ ¹       ~ (p=0.100 n=3) ²   8.571m ± ∞ ¹        ~ (p=1.000 n=3) ²   16.088m ± ∞ ¹        ~ (p=0.100 n=3) ²   8.649m ± ∞ ¹       ~ (p=0.700 n=3) ²   15.201m ± ∞ ¹        ~ (p=0.100 n=3) ²   8.718m ± ∞ ¹        ~ (p=0.400 n=3) ²   13.913m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                             9.126m         9.634m        +5.57%                   7.891m        -13.53%                    12.24m        +34.13%                   8.452m        -7.38%                    11.94m        +30.82%                   7.842m        -14.06%                    11.96m        +31.04%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                      │ bench_ncruces_direct.txt │           bench_ncruces_driver.txt           │      bench_eatonphil_direct.txt      │          bench_glebarez_driver.txt           │             bench_mattn_driver.txt             │           bench_modernc_driver.txt           │          bench_tailscale_driver.txt          │        bench_zombiezen_direct.txt        │
                      │           B/op           │      B/op        vs base                     │    B/op      vs base                 │      B/op        vs base                     │       B/op        vs base                      │      B/op        vs base                     │      B/op        vs base                     │     B/op       vs base                   │
Query/JSON-12                        2.000 ± ∞ ¹   64953.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.300 n=3) ²   56949.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   201026.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56966.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   32891.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   403.000 ± ∞ ¹          ~ (p=0.600 n=3) ²
Query/JSONParallel-12                40.00 ± ∞ ¹    64887.00 ± ∞ ¹            ~ (p=0.100 n=3) ²   43.00 ± ∞ ¹        ~ (p=0.200 n=3) ²    57450.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    201235.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    57471.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    32945.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    469.00 ± ∞ ¹          ~ (p=0.100 n=3) ²
geomean                              8.944           63.40Ki        +725727.57%                   6.557        -26.69%                     55.86Ki        +639403.72%                      196.4Ki        +2248607.05%                     55.88Ki        +639616.05%                     32.15Ki        +367934.30%                     434.7        +4760.65%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                      │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │   bench_zombiezen_direct.txt   │
                      │        allocs/op         │   allocs/op     vs base           │  allocs/op   vs base                │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │  allocs/op   vs base           │
Query/JSON-12                        0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-12                0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                        ⁴     4.019k        ?                                +0.00%               ⁴     4.022k        ?                     5.021k        ?                     4.022k        ?                     2.020k        ?                   6.000        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

### - Parallel Performance

<!--BENCHMARK:slow/benchstat_query_json_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                      │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt        │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                      │          sec/op          │    sec/op     vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │
Query/JSONParallel                 10.041m ± ∞ ¹   9.792m ± ∞ ¹       ~ (p=0.700 n=3) ²   7.274m ± ∞ ¹        ~ (p=0.100 n=3) ²    9.453m ± ∞ ¹        ~ (p=0.100 n=3) ²   7.732m ± ∞ ¹       ~ (p=0.100 n=3) ²    9.413m ± ∞ ¹        ~ (p=0.100 n=3) ²   7.072m ± ∞ ¹        ~ (p=0.100 n=3) ²   10.125m ± ∞ ¹        ~ (p=0.700 n=3) ²
Query/JSONParallel-2                7.359m ± ∞ ¹   6.987m ± ∞ ¹       ~ (p=0.100 n=3) ²   5.792m ± ∞ ¹        ~ (p=0.400 n=3) ²    7.459m ± ∞ ¹        ~ (p=0.700 n=3) ²   7.342m ± ∞ ¹       ~ (p=1.000 n=3) ²    7.637m ± ∞ ¹        ~ (p=0.700 n=3) ²   5.789m ± ∞ ¹        ~ (p=0.100 n=3) ²    8.633m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-4                4.882m ± ∞ ¹   4.816m ± ∞ ¹       ~ (p=1.000 n=3) ²   4.619m ± ∞ ¹        ~ (p=0.400 n=3) ²    5.927m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.256m ± ∞ ¹       ~ (p=0.200 n=3) ²    5.688m ± ∞ ¹        ~ (p=0.100 n=3) ²   4.441m ± ∞ ¹        ~ (p=0.100 n=3) ²    8.151m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-8                7.702m ± ∞ ¹   7.086m ± ∞ ¹       ~ (p=0.100 n=3) ²   7.254m ± ∞ ¹        ~ (p=0.100 n=3) ²   13.831m ± ∞ ¹        ~ (p=0.100 n=3) ²   7.696m ± ∞ ¹       ~ (p=1.000 n=3) ²   13.719m ± ∞ ¹        ~ (p=0.100 n=3) ²   6.870m ± ∞ ¹        ~ (p=0.100 n=3) ²   15.200m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-16               9.570m ± ∞ ¹   8.875m ± ∞ ¹       ~ (p=0.100 n=3) ²   9.005m ± ∞ ¹        ~ (p=0.100 n=3) ²   16.323m ± ∞ ¹        ~ (p=0.100 n=3) ²   9.646m ± ∞ ¹       ~ (p=0.700 n=3) ²   16.343m ± ∞ ¹        ~ (p=0.100 n=3) ²   8.795m ± ∞ ¹        ~ (p=0.100 n=3) ²   13.976m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                             7.672m         7.299m        -4.86%                   6.620m        -13.72%                    9.884m        +28.83%                   7.397m        -3.59%                    9.828m        +28.09%                   6.429m        -16.20%                    10.86m        +41.60%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                      │ bench_ncruces_direct.txt │           bench_ncruces_driver.txt           │      bench_eatonphil_direct.txt      │          bench_glebarez_driver.txt           │             bench_mattn_driver.txt             │           bench_modernc_driver.txt           │          bench_tailscale_driver.txt          │        bench_zombiezen_direct.txt        │
                      │           B/op           │      B/op        vs base                     │    B/op      vs base                 │      B/op        vs base                     │       B/op        vs base                      │      B/op        vs base                     │      B/op        vs base                     │     B/op       vs base                   │
Query/JSONParallel                   2.000 ± ∞ ¹   64845.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   56877.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   201019.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56861.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   32891.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   367.000 ± ∞ ¹          ~ (p=0.600 n=3) ²
Query/JSONParallel-2                 4.000 ± ∞ ¹   64849.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹        ~ (p=0.400 n=3) ²   56949.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   201022.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56905.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   32885.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   365.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-4                 6.000 ± ∞ ¹   64846.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹        ~ (p=1.000 n=3) ²   56955.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   201019.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56948.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   32892.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   368.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-8                20.000 ± ∞ ¹   64847.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹        ~ (p=1.000 n=3) ²   56995.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   201070.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   57002.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   32907.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   426.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-16                48.00 ± ∞ ¹    64944.00 ± ∞ ¹            ~ (p=0.100 n=3) ²   70.00 ± ∞ ¹        ~ (p=0.700 n=3) ²    57048.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    201112.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    57075.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    33031.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    538.00 ± ∞ ¹          ~ (p=0.100 n=3) ²
geomean                              8.565           63.35Ki        +757283.58%                   5.305        -38.06%                     55.63Ki        +665025.92%                      196.3Ki        +2347359.56%                     55.62Ki        +664948.61%                     32.15Ki        +384290.41%                     407.9        +4663.18%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                      │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │   bench_zombiezen_direct.txt   │
                      │        allocs/op         │   allocs/op     vs base           │  allocs/op   vs base                │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │  allocs/op   vs base           │
Query/JSONParallel                   0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-2                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-4                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-8                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-16                0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4023.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4023.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                        ⁴     4.019k        ?                                +0.00%               ⁴     4.022k        ?                     5.021k        ?                     4.022k        ?                     2.020k        ?                   6.000        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

## Query - OrderBy

<!--SQL:gsb_common_sql_query_orderby.sql-->
```sql
SELECT
  name, created, id
FROM comments
ORDER BY name, created, id
```
<!--END_SQL-->

### - Sequential & Parallel Performance

<!--BENCHMARK:slow/benchstat_query_orderby.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                         │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt       │         bench_mattn_driver.txt         │       bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │      bench_zombiezen_direct.txt       │
                         │          sec/op          │    sec/op     vs base                │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │
Query/OrderBy-12                       83.46m ± ∞ ¹   90.57m ± ∞ ¹       ~ (p=0.100 n=3) ²   47.05m ± ∞ ¹        ~ (p=0.100 n=3) ²   81.18m ± ∞ ¹        ~ (p=0.700 n=3) ²   109.21m ± ∞ ¹        ~ (p=0.100 n=3) ²   83.44m ± ∞ ¹        ~ (p=1.000 n=3) ²   61.86m ± ∞ ¹        ~ (p=0.100 n=3) ²   84.56m ± ∞ ¹        ~ (p=0.700 n=3) ²
Query/OrderByParallel-12               48.06m ± ∞ ¹   51.52m ± ∞ ¹       ~ (p=0.100 n=3) ²   48.50m ± ∞ ¹        ~ (p=0.400 n=3) ²   75.27m ± ∞ ¹        ~ (p=0.100 n=3) ²    91.74m ± ∞ ¹        ~ (p=0.100 n=3) ²   71.91m ± ∞ ¹        ~ (p=0.100 n=3) ²   51.47m ± ∞ ¹        ~ (p=0.100 n=3) ²   79.94m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                63.33m         68.31m        +7.85%                   47.77m        -24.58%                   78.17m        +23.43%                    100.1m        +58.04%                   77.46m        +22.31%                   56.43m        -10.91%                   82.22m        +29.82%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │           bench_ncruces_driver.txt           │      bench_eatonphil_direct.txt      │          bench_glebarez_driver.txt           │            bench_mattn_driver.txt             │           bench_modernc_driver.txt           │         bench_tailscale_driver.txt          │       bench_zombiezen_direct.txt       │
                         │           B/op           │      B/op        vs base                     │    B/op      vs base                 │      B/op        vs base                     │       B/op        vs base                     │      B/op        vs base                     │      B/op        vs base                    │     B/op       vs base                 │
Query/OrderBy-12                      57461.0 ± ∞ ¹   6399308.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   118.0 ± ∞ ¹        ~ (p=0.100 n=3) ²   6397822.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   11999108.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   6397808.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   2798877.0 ± ∞ ¹           ~ (p=0.100 n=3) ²    7222.0 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-12                695.0 ± ∞ ¹   6399480.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   532.0 ± ∞ ¹        ~ (p=0.700 n=3) ²   6398311.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   12000026.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   6398237.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   2799192.0 ± ∞ ¹           ~ (p=0.100 n=3) ²    1448.0 ± ∞ ¹        ~ (p=0.700 n=3) ²
geomean                               6.171Ki           6.103Mi        +101165.11%                   250.6        -96.04%                     6.102Mi        +101144.11%                      11.44Mi        +189783.22%                     6.102Mi        +101143.41%                     2.669Mi        +44192.40%                   3.158Ki        -48.83%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt            │      bench_eatonphil_direct.txt      │           bench_glebarez_driver.txt            │             bench_mattn_driver.txt             │            bench_modernc_driver.txt            │           bench_tailscale_driver.txt           │      bench_zombiezen_direct.txt       │
                         │        allocs/op         │    allocs/op      vs base                      │  allocs/op   vs base                 │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │  allocs/op    vs base                 │
Query/OrderBy-12                       20.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   449782.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349771.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449781.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   37.000 ± ∞ ¹        ~ (p=0.300 n=3) ²
Query/OrderByParallel-12                9.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   449784.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449785.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149767.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   19.000 ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                 13.42             349.8k        +2606995.01%                   1.000        -92.55%                       449.8k        +3352384.54%                       349.8k        +2606972.65%                       449.8k        +3352384.54%                       149.8k        +1116193.58%                    26.51        +97.62%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
```
<!--END_BENCHMARK-->

### - Parallel Performance

<!--BENCHMARK:slow/benchstat_query_orderby_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                         │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt       │        bench_mattn_driver.txt         │       bench_modernc_driver.txt        │      bench_tailscale_driver.txt      │      bench_zombiezen_direct.txt       │
                         │          sec/op          │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                │    sec/op     vs base                 │
Query/OrderByParallel                  79.20m ± ∞ ¹   90.64m ± ∞ ¹        ~ (p=0.100 n=3) ²   46.44m ± ∞ ¹        ~ (p=0.100 n=3) ²   83.80m ± ∞ ¹        ~ (p=0.100 n=3) ²   87.88m ± ∞ ¹        ~ (p=0.100 n=3) ²   82.67m ± ∞ ¹        ~ (p=0.100 n=3) ²   61.85m ± ∞ ¹       ~ (p=0.100 n=3) ²   87.37m ± ∞ ¹        ~ (p=0.700 n=3) ²
Query/OrderByParallel-2                49.64m ± ∞ ¹   59.12m ± ∞ ¹        ~ (p=0.100 n=3) ²   34.91m ± ∞ ¹        ~ (p=0.100 n=3) ²   56.34m ± ∞ ¹        ~ (p=0.100 n=3) ²   63.23m ± ∞ ¹        ~ (p=0.100 n=3) ²   56.53m ± ∞ ¹        ~ (p=0.100 n=3) ²   42.48m ± ∞ ¹       ~ (p=0.100 n=3) ²   51.81m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-4                33.24m ± ∞ ¹   37.22m ± ∞ ¹        ~ (p=0.100 n=3) ²   24.03m ± ∞ ¹        ~ (p=0.200 n=3) ²   38.87m ± ∞ ¹        ~ (p=0.100 n=3) ²   51.52m ± ∞ ¹        ~ (p=0.100 n=3) ²   37.12m ± ∞ ¹        ~ (p=0.100 n=3) ²   29.48m ± ∞ ¹       ~ (p=0.100 n=3) ²   43.71m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-8                40.55m ± ∞ ¹   44.51m ± ∞ ¹        ~ (p=0.100 n=3) ²   38.15m ± ∞ ¹        ~ (p=0.700 n=3) ²   66.22m ± ∞ ¹        ~ (p=0.100 n=3) ²   74.20m ± ∞ ¹        ~ (p=0.100 n=3) ²   65.41m ± ∞ ¹        ~ (p=0.100 n=3) ²   41.51m ± ∞ ¹       ~ (p=0.100 n=3) ²   76.77m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-16               51.16m ± ∞ ¹   52.81m ± ∞ ¹        ~ (p=0.700 n=3) ²   48.87m ± ∞ ¹        ~ (p=0.200 n=3) ²   71.58m ± ∞ ¹        ~ (p=0.100 n=3) ²   97.17m ± ∞ ¹        ~ (p=0.100 n=3) ²   68.18m ± ∞ ¹        ~ (p=0.100 n=3) ²   51.36m ± ∞ ¹       ~ (p=0.700 n=3) ²   78.76m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                48.60m         54.22m        +11.58%                   37.34m        -23.15%                   61.36m        +26.27%                   72.94m        +50.09%                   59.94m        +23.33%                   44.01m        -9.44%                   65.40m        +34.57%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt            │      bench_eatonphil_direct.txt      │           bench_glebarez_driver.txt            │             bench_mattn_driver.txt              │            bench_modernc_driver.txt            │           bench_tailscale_driver.txt           │       bench_zombiezen_direct.txt        │
                         │           B/op           │       B/op         vs base                     │    B/op      vs base                 │       B/op         vs base                     │        B/op         vs base                     │       B/op         vs base                     │       B/op         vs base                     │      B/op       vs base                 │
Query/OrderByParallel                55903.00 ± ∞ ¹    6399279.00 ± ∞ ¹            ~ (p=0.100 n=3) ²   83.00 ± ∞ ¹        ~ (p=0.100 n=3) ²    6397704.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    11999009.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    6397711.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    2798848.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    6477.00 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-2               392.000 ± ∞ ¹   6399274.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   9.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6397723.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   11999045.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6397741.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   2798844.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1131.000 ± ∞ ¹        ~ (p=0.600 n=3) ²
Query/OrderByParallel-4                413.00 ± ∞ ¹    6399273.00 ± ∞ ¹            ~ (p=0.100 n=3) ²   14.00 ± ∞ ¹        ~ (p=0.100 n=3) ²    6397889.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    11999095.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    6397862.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    2798861.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    1128.00 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-8                 493.0 ± ∞ ¹     6399360.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   233.0 ± ∞ ¹        ~ (p=0.200 n=3) ²     6398016.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     12000246.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     6398052.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     2798894.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     1480.0 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-16                715.0 ± ∞ ¹     6399513.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   312.0 ± ∞ ¹        ~ (p=0.700 n=3) ²     6398773.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     12001562.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     6398608.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     2799260.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     1661.0 ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                               1.232Ki             6.103Mi        +507323.21%                   59.73        -95.26%                       6.102Mi        +507218.64%                        11.44Mi        +951400.14%                       6.102Mi        +507216.57%                       2.669Mi        +221836.62%                    1.783Ki        +44.81%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt            │    bench_eatonphil_direct.txt    │           bench_glebarez_driver.txt            │             bench_mattn_driver.txt             │            bench_modernc_driver.txt            │           bench_tailscale_driver.txt           │      bench_zombiezen_direct.txt       │
                         │        allocs/op         │    allocs/op      vs base                      │  allocs/op   vs base             │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │  allocs/op    vs base                 │
Query/OrderByParallel                  20.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349771.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   34.000 ± ∞ ¹        ~ (p=0.400 n=3) ²
Query/OrderByParallel-2                 8.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349771.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   16.000 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-4                 8.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449782.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349772.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449781.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   16.000 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-8                 9.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449783.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349776.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449784.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   18.000 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-16               10.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449791.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349791.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449789.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149769.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   19.000 ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                 10.29             349.8k        +3400086.50%                                ?               ³ ⁴       449.8k        +4372228.00%                       349.8k        +3400063.17%                       449.8k        +4372224.11%                       149.8k        +1455776.29%                    19.71        +91.63%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean
⁴ ratios must be >0 to compute geomean
```
<!--END_BENCHMARK-->

## Query - RecursiveCTE

<!--SQL:gsb_common_sql_query_recursivecte.sql-->
```sql
WITH RECURSIVE dates(day) AS (
  SELECT date('now', '-30 days')
  UNION ALL
  SELECT date(day, '+1 day')
  FROM dates
  WHERE day < date('now')
)
SELECT day,
  (SELECT COUNT(*) FROM posts WHERE date(created) = day) as day_total
FROM dates
ORDER BY day
```
<!--END_SQL-->

### - Sequential & Parallel Performance

<!--BENCHMARK:slow/benchstat_query_recursivecte.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                              │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │      bench_eatonphil_direct.txt       │      bench_glebarez_driver.txt       │        bench_mattn_driver.txt         │       bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │      bench_zombiezen_direct.txt      │
                              │          sec/op          │    sec/op     vs base                │    sec/op     vs base                 │    sec/op     vs base                │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                │
Query/RecursiveCTE-12                       6.276m ± ∞ ¹   6.256m ± ∞ ¹       ~ (p=0.400 n=3) ²   2.503m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.287m ± ∞ ¹       ~ (p=0.100 n=3) ²   2.474m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.313m ± ∞ ¹        ~ (p=0.100 n=3) ²   2.431m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.268m ± ∞ ¹       ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-12               670.0µ ± ∞ ¹   687.9µ ± ∞ ¹       ~ (p=0.100 n=3) ²   276.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   657.3µ ± ∞ ¹       ~ (p=0.400 n=3) ²   275.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   633.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²   272.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   655.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²
geomean                                     2.051m         2.074m        +1.16%                   832.5µ        -59.40%                   1.864m        -9.09%                   824.8µ        -59.78%                   1.835m        -10.53%                   813.2µ        -60.34%                   1.859m        -9.36%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt           │    bench_eatonphil_direct.txt    │          bench_glebarez_driver.txt          │           bench_mattn_driver.txt            │          bench_modernc_driver.txt           │         bench_tailscale_driver.txt          │        bench_zombiezen_direct.txt         │
                              │           B/op           │      B/op       vs base                     │    B/op      vs base             │      B/op       vs base                     │      B/op       vs base                     │      B/op       vs base                     │      B/op       vs base                     │     B/op       vs base                    │
Query/RecursiveCTE-12                        1.000 ± ∞ ¹   2488.000 ± ∞ ¹            ~ (p=0.600 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.300 n=3) ²     2354.000 ± ∞ ¹            ~ (p=0.400 n=3) ²   6859.000 ± ∞ ¹            ~ (p=0.600 n=3) ²   2293.000 ± ∞ ¹            ~ (p=0.600 n=3) ²   1505.000 ± ∞ ¹            ~ (p=0.400 n=3) ²   389.000 ± ∞ ¹           ~ (p=0.600 n=3) ²
Query/RecursiveCTEParallel-12                2.000 ± ∞ ¹   2483.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.800 n=3) ²     2358.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6850.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   2286.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1496.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   368.000 ± ∞ ¹           ~ (p=0.100 n=3) ²
geomean                                      1.414          2.427Ki        +175651.30%                                ?               ³ ⁴    2.301Ki        +166494.30%                    6.694Ki        +484586.24%                    2.236Ki        +161791.91%                    1.465Ki        +106000.90%                     378.4        +26653.69%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean
⁴ ratios must be >0 to compute geomean

                              │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                              │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/RecursiveCTE-12                        0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ²   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   143.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-12                0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   48.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                ⁴     110.0        ?                                +0.00%               ⁴     113.0        ?                     142.5        ?                     112.0        ?                    48.50        ?                   6.000        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

### - Parallel Performance

<!--BENCHMARK:slow/benchstat_query_recursivecte_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                              │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                              │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │
Query/RecursiveCTEParallel                  6.152m ± ∞ ¹    6.186m ± ∞ ¹       ~ (p=0.100 n=3) ²   2.512m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.289m ± ∞ ¹        ~ (p=0.100 n=3) ²   2.464m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.281m ± ∞ ¹        ~ (p=0.100 n=3) ²   2.422m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.239m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                3.122m ± ∞ ¹    3.117m ± ∞ ¹       ~ (p=0.700 n=3) ²   1.269m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.672m ± ∞ ¹        ~ (p=0.100 n=3) ²   1.249m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.667m ± ∞ ¹        ~ (p=0.100 n=3) ²   1.230m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.650m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4               1593.6µ ± ∞ ¹   1597.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²   645.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1374.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   632.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1368.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   623.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1362.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                799.0µ ± ∞ ¹    805.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²   324.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    749.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²   324.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    742.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   317.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    746.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16               674.0µ ± ∞ ¹    678.3µ ± ∞ ¹       ~ (p=0.400 n=3) ²   280.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    647.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   279.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    642.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   274.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    663.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                     1.752m          1.759m        +0.42%                   715.1µ        -59.17%                    1.566m        -10.58%                   706.6µ        -59.66%                    1.558m        -11.04%                   694.6µ        -60.34%                    1.564m        -10.68%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │   bench_eatonphil_direct.txt   │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                              │           B/op           │      B/op       vs base           │    B/op      vs base           │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │     B/op       vs base           │
Query/RecursiveCTEParallel                   1.000 ± ∞ ¹   2482.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2260.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6857.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2258.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1504.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   362.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                 1.000 ± ∞ ¹   2481.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.400 n=3) ²   2263.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6825.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2258.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1502.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   366.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4                 1.000 ± ∞ ¹   2481.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.400 n=3) ²   2280.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6842.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2259.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1502.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   373.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                   0.0 ± ∞ ¹     2482.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ²     2356.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     6857.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     2280.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1504.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     363.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16                0.000 ± ∞ ¹   2489.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.700 n=3) ²   2361.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6862.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2284.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1506.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   371.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                ³    2.425Ki        ?                                ?               ³    2.250Ki        ?                    6.688Ki        ?                    2.215Ki        ?                    1.468Ki        ?                     367.0        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean

                              │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                              │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/RecursiveCTEParallel                   0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   143.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   48.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   48.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   143.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16                0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   143.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                ⁴     110.0        ?                                +0.00%               ⁴     112.4        ?                     142.6        ?                     112.0        ?                    48.60        ?                   6.000        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

## Query - Window

<!--SQL:gsb_common_sql_query_window.sql-->
```sql
WITH day_totals AS (
  SELECT date(created) as day, COUNT(*) as day_total
  FROM posts
  GROUP BY day
)
SELECT day, day_total,
  SUM(day_total) OVER (ORDER BY day) as running_total
FROM day_totals
ORDER BY day
```
<!--END_SQL-->

### - Sequential & Parallel Performance

<!--BENCHMARK:slow/benchstat_query_window.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                        │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │         bench_mattn_driver.txt          │        bench_modernc_driver.txt        │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                        │          sec/op          │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op      vs base                  │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │
Query/Window-12                      2166.2µ ± ∞ ¹   2448.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   917.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1896.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   2078.7µ ± ∞ ¹         ~ (p=0.100 n=3) ²   1900.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1055.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1634.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-12               238.7µ ± ∞ ¹    321.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   122.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    899.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1069.2µ ± ∞ ¹         ~ (p=0.100 n=3) ²    908.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    133.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    864.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                               719.1µ          887.2µ        +23.37%                   334.5µ        -53.48%                    1.307m        +81.68%                    1.491m        +107.31%                    1.314m        +82.74%                    375.1µ        -47.84%                    1.189m        +65.33%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                        │ bench_ncruces_direct.txt │      bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt      │       bench_mattn_driver.txt        │      bench_modernc_driver.txt      │     bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                        │           B/op           │      B/op        vs base           │    B/op      vs base                │      B/op        vs base           │       B/op        vs base           │      B/op        vs base           │      B/op        vs base           │     B/op       vs base           │
Query/Window-12                          0.0 ± ∞ ¹     62767.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹       ~ (p=1.000 n=3) ²     54897.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198937.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54911.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30798.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     376.0 ± ∞ ¹  ~ (p=0.600 n=3) ²
Query/WindowParallel-12                1.000 ± ∞ ¹   62763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹       ~ (p=1.000 n=3) ²   54886.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   198953.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54907.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   30765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   372.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                          ³     61.29Ki        ?                                +0.00%               ³     53.60Ki        ?                      194.3Ki        ?                     53.62Ki        ?                     30.06Ki        ?                     374.0        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean

                        │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │   bench_zombiezen_direct.txt   │
                        │        allocs/op         │   allocs/op     vs base           │  allocs/op   vs base                │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │  allocs/op   vs base           │
Query/Window-12                        0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-12                0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                          ⁴     3.763k        ?                                +0.00%               ⁴     3.766k        ?                     4.765k        ?                     3.766k        ?                     1.763k        ?                   6.000        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

### - Parallel Performance

<!--BENCHMARK:slow/benchstat_query_window_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                        │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │         bench_mattn_driver.txt         │        bench_modernc_driver.txt        │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                        │          sec/op          │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │
Query/WindowParallel                 2119.8µ ± ∞ ¹   2406.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   930.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1884.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1568.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1876.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1060.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1614.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-2               1070.4µ ± ∞ ¹   1234.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   476.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    993.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    849.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    977.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    544.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    846.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-4                544.7µ ± ∞ ¹    643.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   254.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    541.9µ ± ∞ ¹        ~ (p=0.700 n=3) ²    503.4µ ± ∞ ¹        ~ (p=0.700 n=3) ²    539.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    283.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    465.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-8                277.1µ ± ∞ ¹    352.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   142.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    864.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1071.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    879.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    149.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    854.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-16               236.8µ ± ∞ ¹    314.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   122.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    915.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1112.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    913.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    144.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    874.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                               605.1µ          733.4µ        +21.20%                   287.5µ        -52.48%                    957.0µ        +58.17%                    956.3µ        +58.05%                    955.2µ        +57.86%                    323.1µ        -46.60%                    861.7µ        +42.42%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                        │ bench_ncruces_direct.txt │      bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt      │       bench_mattn_driver.txt        │      bench_modernc_driver.txt      │     bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                        │           B/op           │      B/op        vs base           │    B/op      vs base                │      B/op        vs base           │       B/op        vs base           │      B/op        vs base           │      B/op        vs base           │     B/op       vs base           │
Query/WindowParallel                     0.0 ± ∞ ¹     62762.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³     54775.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198936.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54789.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30792.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     360.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-2                   0.0 ± ∞ ¹     62761.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³     54831.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198863.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54815.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30793.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     365.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-4                   0.0 ± ∞ ¹     62761.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³     54882.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198853.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54883.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30794.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     363.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-8                   0.0 ± ∞ ¹     62763.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³     54884.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198906.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54904.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30793.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     365.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-16                1.000 ± ∞ ¹   62766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹       ~ (p=1.000 n=3) ²   54894.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   198976.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54917.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   30773.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   375.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                          ⁴     61.29Ki        ?                                +0.00%               ⁴     53.57Ki        ?                      194.2Ki        ?                     53.58Ki        ?                     30.07Ki        ?                     365.6        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                        │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │   bench_zombiezen_direct.txt   │
                        │        allocs/op         │   allocs/op     vs base           │  allocs/op   vs base                │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │  allocs/op   vs base           │
Query/WindowParallel                   0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-2                 0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-4                 0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-8                 0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-16                0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                          ⁴     3.763k        ?                                +0.00%               ⁴     3.766k        ?                     4.764k        ?                     3.765k        ?                     1.764k        ?                   6.000        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

# Reviewing the Implementations

Please note that the benchmark framework sets the `busy_timeout`, `foreign_keys`, `journal_mode` and `synchronous` pragmas; what you're seeing below isn't their default values.

## eatonphil

### Compile-time Options

<!--GREP:tests/test_eatonphil_direct.txt:OPTION-->
```
ATOMIC_INTRINSICS=1
COMPILER=clang-17.0.0
DEFAULT_AUTOVACUUM
DEFAULT_CACHE_SIZE=-2000
DEFAULT_FILE_FORMAT=4
DEFAULT_JOURNAL_SIZE_LIMIT=-1
DEFAULT_MEMSTATUS=0
DEFAULT_MMAP_SIZE=0
DEFAULT_PAGE_SIZE=4096
DEFAULT_PCACHE_INITSZ=20
DEFAULT_RECURSIVE_TRIGGERS
DEFAULT_SECTOR_SIZE=4096
DEFAULT_SYNCHRONOUS=2
DEFAULT_WAL_AUTOCHECKPOINT=1000
DEFAULT_WAL_SYNCHRONOUS=2
DEFAULT_WORKER_THREADS=0
DIRECT_OVERFLOW_READ
ENABLE_API_ARMOR
ENABLE_FTS3
ENABLE_FTS3_PARENTHESIS
ENABLE_FTS4
ENABLE_FTS5
ENABLE_GEOPOLY
ENABLE_PREUPDATE_HOOK
ENABLE_RTREE
ENABLE_SESSION
ENABLE_STAT4
ENABLE_UNLOCK_NOTIFY
ENABLE_UPDATE_DELETE_LIMIT
MALLOC_SOFT_LIMIT=1024
MAX_ATTACHED=10
MAX_COLUMN=2000
MAX_COMPOUND_SELECT=500
MAX_DEFAULT_PAGE_SIZE=8192
MAX_EXPR_DEPTH=1000
MAX_FUNCTION_ARG=127
MAX_LENGTH=1000000000
MAX_LIKE_PATTERN_LENGTH=50000
MAX_MMAP_SIZE=0x7fff0000
MAX_PAGE_COUNT=0xfffffffe
MAX_PAGE_SIZE=65536
MAX_SQL_LENGTH=1000000000
MAX_TRIGGER_DEPTH=1000
MAX_VARIABLE_NUMBER=32766
MAX_VDBE_OP=250000000
MAX_WORKER_THREADS=8
MUTEX_PTHREADS
OMIT_AUTOINIT
OMIT_DEPRECATED
OMIT_LOAD_EXTENSION
OMIT_PROGRESS_CALLBACK
OMIT_UTF16
SOUNDEX
SYSTEM_MALLOC
TEMP_STORE=2
THREADSAFE=2
USE_ALLOCA
USE_URI
```
<!--END_GREP-->

### Pragmas

<!--GREP:tests/test_eatonphil_direct.txt:PRAGMA-->
```
auto_vacuum=0
automatic_index=1
busy_timeout=5000
cache_size=-2000
cache_spill=483
cell_size_check=0
checkpoint_fullfsync=0
defer_foreign_keys=0
foreign_keys=1
fullfsync=0
hard_heap_limit=0
journal_mode=wal
journal_size_limit=-1
locking_mode=normal
mmap_size=0
page_size=4096
query_only=0
read_uncommitted=0
recursive_triggers=0
reverse_unordered_selects=0
secure_delete=0
soft_heap_limit=0
synchronous=1
temp_store=0
threads=0
wal_autocheckpoint=1000
```
<!--END_GREP-->

### SQLite Version

<!--GREP:tests/test_eatonphil_direct.txt:VERSION-->
```
3.46.0
```
<!--END_GREP-->

## glebarez

* **Configuration**

  glebarez uses some additional configuration:

    * `_time_format=sqlite` (default is none)
    * `_txlock=immediate` (default is `deferred`)

* **Notes**

  If a transaction is `ReadOnly` it's always `deferred`. If it isn't, it uses the value of `_txlock`.

### Compile-time Options

<!--GREP:tests/test_glebarez_driver.txt:OPTION-->
```
ATOMIC_INTRINSICS=0
COMPILER=clang-16.0.0
DEFAULT_AUTOVACUUM
DEFAULT_CACHE_SIZE=-2000
DEFAULT_FILE_FORMAT=4
DEFAULT_JOURNAL_SIZE_LIMIT=-1
DEFAULT_MEMSTATUS=0
DEFAULT_MMAP_SIZE=0
DEFAULT_PAGE_SIZE=4096
DEFAULT_PCACHE_INITSZ=20
DEFAULT_RECURSIVE_TRIGGERS
DEFAULT_SECTOR_SIZE=4096
DEFAULT_SYNCHRONOUS=2
DEFAULT_WAL_AUTOCHECKPOINT=1000
DEFAULT_WAL_SYNCHRONOUS=2
DEFAULT_WORKER_THREADS=0
DIRECT_OVERFLOW_READ
ENABLE_COLUMN_METADATA
ENABLE_DBSTAT_VTAB
ENABLE_FTS5
ENABLE_GEOPOLY
ENABLE_MATH_FUNCTIONS
ENABLE_MEMORY_MANAGEMENT
ENABLE_OFFSET_SQL_FUNC
ENABLE_PREUPDATE_HOOK
ENABLE_RBU
ENABLE_RTREE
ENABLE_SESSION
ENABLE_SNAPSHOT
ENABLE_STAT4
ENABLE_UNLOCK_NOTIFY
LIKE_DOESNT_MATCH_BLOBS
MALLOC_SOFT_LIMIT=1024
MAX_ATTACHED=10
MAX_COLUMN=2000
MAX_COMPOUND_SELECT=500
MAX_DEFAULT_PAGE_SIZE=8192
MAX_EXPR_DEPTH=1000
MAX_FUNCTION_ARG=1000
MAX_LENGTH=1000000000
MAX_LIKE_PATTERN_LENGTH=50000
MAX_MMAP_SIZE=0x7fff0000
MAX_PAGE_COUNT=0xfffffffe
MAX_PAGE_SIZE=65536
MAX_SQL_LENGTH=1000000000
MAX_TRIGGER_DEPTH=1000
MAX_VARIABLE_NUMBER=32766
MAX_VDBE_OP=250000000
MAX_WORKER_THREADS=8
MUTEX_NOOP
SOUNDEX
SYSTEM_MALLOC
TEMP_STORE=1
THREADSAFE=1
```
<!--END_GREP-->

### Pragmas

<!--GREP:tests/test_glebarez_driver.txt:PRAGMA-->
```
auto_vacuum=0
automatic_index=1
busy_timeout=5000
cache_size=-2000
cache_spill=483
cell_size_check=0
checkpoint_fullfsync=0
defer_foreign_keys=0
encoding=UTF-8
foreign_keys=1
fullfsync=0
hard_heap_limit=0
journal_mode=wal
journal_size_limit=-1
locking_mode=normal
mmap_size=0
page_size=4096
query_only=0
read_uncommitted=0
recursive_triggers=0
reverse_unordered_selects=0
secure_delete=0
soft_heap_limit=0
synchronous=1
temp_store=0
threads=0
wal_autocheckpoint=1000
```
<!--END_GREP-->

### SQLite Version

<!--GREP:tests/test_glebarez_driver.txt:VERSION-->
```
3.49.1
```
<!--END_GREP-->

## mattn

* **Configuration**

  mattn uses some additional configuration:

    * `_mutex=no` (default is `full`)

* **Notes**

  mattn doesn't currently support read and write transactions on the same connection because it [ignores `TxOptions` completely](https://github.com/mattn/go-sqlite3/blob/7658c06970ecf5588d8cd930ed1f2de7223f1010/sqlite3_go18.go#L42). Instead, mattn [uses the value of `_txlock`](https://github.com/mattn/go-sqlite3/blob/7658c06970ecf5588d8cd930ed1f2de7223f1010/sqlite3.go#L969) whenever it starts a transaction. If you change the value of `_txlock` to `immediate` you make all transactions write transactions. If you leave the value of `_txlock` at its default, `deferred`, you risk immediate `SQLITE_BUSY` errors when a read transaction is upgraded to a write transaction. See more [here](https://berthub.eu/articles/posts/a-brief-post-on-sqlite3-database-locked-despite-timeout/).

### Compile-time Options

<!--GREP:tests/test_mattn_driver.txt:OPTION-->
```
ATOMIC_INTRINSICS=1
COMPILER=clang-17.0.0
DEFAULT_AUTOVACUUM
DEFAULT_CACHE_SIZE=-2000
DEFAULT_FILE_FORMAT=4
DEFAULT_JOURNAL_SIZE_LIMIT=-1
DEFAULT_MMAP_SIZE=0
DEFAULT_PAGE_SIZE=4096
DEFAULT_PCACHE_INITSZ=20
DEFAULT_RECURSIVE_TRIGGERS
DEFAULT_SECTOR_SIZE=4096
DEFAULT_SYNCHRONOUS=2
DEFAULT_WAL_AUTOCHECKPOINT=1000
DEFAULT_WAL_SYNCHRONOUS=1
DEFAULT_WORKER_THREADS=0
DIRECT_OVERFLOW_READ
ENABLE_FTS3
ENABLE_FTS3_PARENTHESIS
ENABLE_RTREE
ENABLE_UPDATE_DELETE_LIMIT
MALLOC_SOFT_LIMIT=1024
MAX_ATTACHED=10
MAX_COLUMN=2000
MAX_COMPOUND_SELECT=500
MAX_DEFAULT_PAGE_SIZE=8192
MAX_EXPR_DEPTH=1000
MAX_FUNCTION_ARG=1000
MAX_LENGTH=1000000000
MAX_LIKE_PATTERN_LENGTH=50000
MAX_MMAP_SIZE=0x7fff0000
MAX_PAGE_COUNT=0xfffffffe
MAX_PAGE_SIZE=65536
MAX_SQL_LENGTH=1000000000
MAX_TRIGGER_DEPTH=1000
MAX_VARIABLE_NUMBER=32766
MAX_VDBE_OP=250000000
MAX_WORKER_THREADS=8
MUTEX_PTHREADS
OMIT_DEPRECATED
SYSTEM_MALLOC
TEMP_STORE=1
THREADSAFE=1
```
<!--END_GREP-->

### Pragmas

<!--GREP:tests/test_mattn_driver.txt:PRAGMA-->
```
auto_vacuum=0
automatic_index=1
busy_timeout=5000
cache_size=-2000
cache_spill=483
cell_size_check=0
checkpoint_fullfsync=0
defer_foreign_keys=0
encoding=UTF-8
foreign_keys=1
fullfsync=0
hard_heap_limit=0
journal_mode=wal
journal_size_limit=-1
locking_mode=normal
mmap_size=0
page_size=4096
query_only=0
read_uncommitted=0
recursive_triggers=0
reverse_unordered_selects=0
secure_delete=0
soft_heap_limit=0
synchronous=1
temp_store=0
threads=0
wal_autocheckpoint=1000
```
<!--END_GREP-->

### SQLite Version

<!--GREP:tests/test_mattn_driver.txt:VERSION-->
```
3.49.1
```
<!--END_GREP-->

## modernc

* **Configuration**

  modernc uses some additional configuration:

    * `_time_format=sqlite` (default is none)
    * `_txlock=immediate` (default is `deferred`)

* **Notes**

  If a transaction is `ReadOnly` it's always `deferred`. If it isn't, it uses the value of `_txlock`.

### Compile-time Options

<!--GREP:tests/test_modernc_driver.txt:OPTION-->
```
ATOMIC_INTRINSICS=0
COMPILER=clang-16.0.0
DEFAULT_AUTOVACUUM
DEFAULT_CACHE_SIZE=-2000
DEFAULT_FILE_FORMAT=4
DEFAULT_JOURNAL_SIZE_LIMIT=-1
DEFAULT_MEMSTATUS=0
DEFAULT_MMAP_SIZE=0
DEFAULT_PAGE_SIZE=4096
DEFAULT_PCACHE_INITSZ=20
DEFAULT_RECURSIVE_TRIGGERS
DEFAULT_SECTOR_SIZE=4096
DEFAULT_SYNCHRONOUS=2
DEFAULT_WAL_AUTOCHECKPOINT=1000
DEFAULT_WAL_SYNCHRONOUS=2
DEFAULT_WORKER_THREADS=0
DIRECT_OVERFLOW_READ
ENABLE_COLUMN_METADATA
ENABLE_DBSTAT_VTAB
ENABLE_FTS5
ENABLE_GEOPOLY
ENABLE_MATH_FUNCTIONS
ENABLE_MEMORY_MANAGEMENT
ENABLE_OFFSET_SQL_FUNC
ENABLE_PREUPDATE_HOOK
ENABLE_RBU
ENABLE_RTREE
ENABLE_SESSION
ENABLE_SNAPSHOT
ENABLE_STAT4
ENABLE_UNLOCK_NOTIFY
LIKE_DOESNT_MATCH_BLOBS
MALLOC_SOFT_LIMIT=1024
MAX_ATTACHED=10
MAX_COLUMN=2000
MAX_COMPOUND_SELECT=500
MAX_DEFAULT_PAGE_SIZE=8192
MAX_EXPR_DEPTH=1000
MAX_FUNCTION_ARG=1000
MAX_LENGTH=1000000000
MAX_LIKE_PATTERN_LENGTH=50000
MAX_MMAP_SIZE=0x7fff0000
MAX_PAGE_COUNT=0xfffffffe
MAX_PAGE_SIZE=65536
MAX_SQL_LENGTH=1000000000
MAX_TRIGGER_DEPTH=1000
MAX_VARIABLE_NUMBER=32766
MAX_VDBE_OP=250000000
MAX_WORKER_THREADS=8
MUTEX_NOOP
SOUNDEX
SYSTEM_MALLOC
TEMP_STORE=1
THREADSAFE=1
```
<!--END_GREP-->

### Pragmas

<!--GREP:tests/test_modernc_driver.txt:PRAGMA-->
```
auto_vacuum=0
automatic_index=1
busy_timeout=5000
cache_size=-2000
cache_spill=483
cell_size_check=0
checkpoint_fullfsync=0
defer_foreign_keys=0
encoding=UTF-8
foreign_keys=1
fullfsync=0
hard_heap_limit=0
journal_mode=wal
journal_size_limit=-1
locking_mode=normal
mmap_size=0
page_size=4096
query_only=0
read_uncommitted=0
recursive_triggers=0
reverse_unordered_selects=0
secure_delete=0
soft_heap_limit=0
synchronous=1
temp_store=0
threads=0
wal_autocheckpoint=1000
```
<!--END_GREP-->

### SQLite Version

<!--GREP:tests/test_modernc_driver.txt:VERSION-->
```
3.49.1
```
<!--END_GREP-->

## ncruces

### Compile-time Options

<!--GREP:tests/test_ncruces_direct.txt:OPTION-->
```
ALLOW_URI_AUTHORITY
ATOMIC_INTRINSICS=1
COMPILER=clang-19.1.5
DEFAULT_AUTOVACUUM
DEFAULT_CACHE_SIZE=-2000
DEFAULT_FILE_FORMAT=4
DEFAULT_FOREIGN_KEYS
DEFAULT_JOURNAL_SIZE_LIMIT=-1
DEFAULT_MMAP_SIZE=0
DEFAULT_PAGE_SIZE=4096
DEFAULT_PCACHE_INITSZ=20
DEFAULT_RECURSIVE_TRIGGERS
DEFAULT_SECTOR_SIZE=4096
DEFAULT_SYNCHRONOUS=2
DEFAULT_WAL_AUTOCHECKPOINT=1000
DEFAULT_WAL_SYNCHRONOUS=1
DEFAULT_WORKER_THREADS=0
DIRECT_OVERFLOW_READ
DQS=0
ENABLE_API_ARMOR
ENABLE_ATOMIC_WRITE
ENABLE_BATCH_ATOMIC_WRITE
ENABLE_COLUMN_METADATA
ENABLE_FTS5
ENABLE_GEOPOLY
ENABLE_MATH_FUNCTIONS
ENABLE_RTREE
ENABLE_STAT4
HAVE_ISNAN
LIKE_DOESNT_MATCH_BLOBS
MALLOC_SOFT_LIMIT=1024
MAX_ATTACHED=10
MAX_COLUMN=2000
MAX_COMPOUND_SELECT=500
MAX_DEFAULT_PAGE_SIZE=8192
MAX_EXPR_DEPTH=1000
MAX_FUNCTION_ARG=1000
MAX_LENGTH=1000000000
MAX_LIKE_PATTERN_LENGTH=50000
MAX_MMAP_SIZE=0
MAX_PAGE_COUNT=0xfffffffe
MAX_PAGE_SIZE=65536
MAX_SQL_LENGTH=1000000000
MAX_TRIGGER_DEPTH=1000
MAX_VARIABLE_NUMBER=32766
MAX_VDBE_OP=250000000
MAX_WORKER_THREADS=0
MUTEX_OMIT
OMIT_AUTOINIT
OMIT_DEPRECATED
OMIT_DESERIALIZE
OMIT_LOAD_EXTENSION
OMIT_SHARED_CACHE
SOUNDEX
SYSTEM_MALLOC
TEMP_STORE=1
THREADSAFE=0
UNTESTABLE
```
<!--END_GREP-->

### Pragmas

<!--GREP:tests/test_ncruces_direct.txt:PRAGMA-->
```
auto_vacuum=0
automatic_index=1
busy_timeout=5000
cache_size=-2000
cache_spill=489
cell_size_check=0
checkpoint_fullfsync=0
defer_foreign_keys=0
encoding=UTF-8
foreign_keys=1
fullfsync=0
hard_heap_limit=0
journal_mode=wal
journal_size_limit=-1
locking_mode=normal
mmap_size=0
page_size=4096
query_only=0
read_uncommitted=0
recursive_triggers=0
reverse_unordered_selects=0
secure_delete=0
soft_heap_limit=0
synchronous=1
temp_store=0
threads=0
wal_autocheckpoint=1000
```
<!--END_GREP-->

### SQLite Version

<!--GREP:tests/test_ncruces_direct.txt:VERSION-->
```
3.49.1
```
<!--END_GREP-->

## tailscale

### Compile-time Options

<!--GREP:tests/test_tailscale_driver.txt:OPTION-->
```
ATOMIC_INTRINSICS=1
COMPILER=clang-17.0.0
DEFAULT_AUTOVACUUM
DEFAULT_CACHE_SIZE=-2000
DEFAULT_FILE_FORMAT=4
DEFAULT_JOURNAL_SIZE_LIMIT=-1
DEFAULT_MEMSTATUS=0
DEFAULT_MMAP_SIZE=0
DEFAULT_PAGE_SIZE=4096
DEFAULT_PCACHE_INITSZ=20
DEFAULT_RECURSIVE_TRIGGERS
DEFAULT_SECTOR_SIZE=4096
DEFAULT_SYNCHRONOUS=2
DEFAULT_WAL_AUTOCHECKPOINT=1000
DEFAULT_WAL_SYNCHRONOUS=1
DEFAULT_WORKER_THREADS=0
DIRECT_OVERFLOW_READ
DQS=0
ENABLE_COLUMN_METADATA
ENABLE_DBSTAT_VTAB
ENABLE_FTS5
ENABLE_PREUPDATE_HOOK
ENABLE_RTREE
ENABLE_SESSION
ENABLE_SNAPSHOT
ENABLE_STAT4
LIKE_DOESNT_MATCH_BLOBS
MALLOC_SOFT_LIMIT=1024
MAX_ATTACHED=10
MAX_COLUMN=2000
MAX_COMPOUND_SELECT=500
MAX_DEFAULT_PAGE_SIZE=8192
MAX_EXPR_DEPTH=0
MAX_FUNCTION_ARG=127
MAX_LENGTH=1000000000
MAX_LIKE_PATTERN_LENGTH=50000
MAX_MMAP_SIZE=0x7fff0000
MAX_PAGE_COUNT=0xfffffffe
MAX_PAGE_SIZE=65536
MAX_SQL_LENGTH=1000000000
MAX_TRIGGER_DEPTH=1000
MAX_VARIABLE_NUMBER=32766
MAX_VDBE_OP=250000000
MAX_WORKER_THREADS=8
MUTEX_PTHREADS
OMIT_AUTOINIT
OMIT_DEPRECATED
OMIT_LOAD_EXTENSION
OMIT_PROGRESS_CALLBACK
OMIT_SHARED_CACHE
SYSTEM_MALLOC
TEMP_STORE=1
THREADSAFE=2
USE_ALLOCA
```
<!--END_GREP-->

### Pragmas

<!--GREP:tests/test_tailscale_driver.txt:PRAGMA-->
```
auto_vacuum=0
automatic_index=1
busy_timeout=5000
cache_size=-2000
cache_spill=483
cell_size_check=0
checkpoint_fullfsync=0
defer_foreign_keys=0
encoding=UTF-8
foreign_keys=1
fullfsync=0
hard_heap_limit=0
journal_mode=wal
journal_size_limit=-1
locking_mode=normal
mmap_size=0
page_size=4096
query_only=0
read_uncommitted=0
recursive_triggers=0
reverse_unordered_selects=0
secure_delete=0
soft_heap_limit=0
synchronous=1
temp_store=0
threads=0
wal_autocheckpoint=1000
```
<!--END_GREP-->

### SQLite Version

<!--GREP:tests/test_tailscale_driver.txt:VERSION-->
```
3.46.1
```
<!--END_GREP-->

## zombiezen

### Compile-time Options

<!--GREP:tests/test_zombiezen_direct.txt:OPTION-->
```
ATOMIC_INTRINSICS=0
COMPILER=clang-16.0.0
DEFAULT_AUTOVACUUM
DEFAULT_CACHE_SIZE=-2000
DEFAULT_FILE_FORMAT=4
DEFAULT_JOURNAL_SIZE_LIMIT=-1
DEFAULT_MEMSTATUS=0
DEFAULT_MMAP_SIZE=0
DEFAULT_PAGE_SIZE=4096
DEFAULT_PCACHE_INITSZ=20
DEFAULT_RECURSIVE_TRIGGERS
DEFAULT_SECTOR_SIZE=4096
DEFAULT_SYNCHRONOUS=2
DEFAULT_WAL_AUTOCHECKPOINT=1000
DEFAULT_WAL_SYNCHRONOUS=2
DEFAULT_WORKER_THREADS=0
DIRECT_OVERFLOW_READ
ENABLE_COLUMN_METADATA
ENABLE_DBSTAT_VTAB
ENABLE_FTS5
ENABLE_GEOPOLY
ENABLE_MATH_FUNCTIONS
ENABLE_MEMORY_MANAGEMENT
ENABLE_OFFSET_SQL_FUNC
ENABLE_PREUPDATE_HOOK
ENABLE_RBU
ENABLE_RTREE
ENABLE_SESSION
ENABLE_SNAPSHOT
ENABLE_STAT4
ENABLE_UNLOCK_NOTIFY
LIKE_DOESNT_MATCH_BLOBS
MALLOC_SOFT_LIMIT=1024
MAX_ATTACHED=10
MAX_COLUMN=2000
MAX_COMPOUND_SELECT=500
MAX_DEFAULT_PAGE_SIZE=8192
MAX_EXPR_DEPTH=1000
MAX_FUNCTION_ARG=1000
MAX_LENGTH=1000000000
MAX_LIKE_PATTERN_LENGTH=50000
MAX_MMAP_SIZE=0x7fff0000
MAX_PAGE_COUNT=0xfffffffe
MAX_PAGE_SIZE=65536
MAX_SQL_LENGTH=1000000000
MAX_TRIGGER_DEPTH=1000
MAX_VARIABLE_NUMBER=32766
MAX_VDBE_OP=250000000
MAX_WORKER_THREADS=8
MUTEX_NOOP
SOUNDEX
SYSTEM_MALLOC
TEMP_STORE=1
THREADSAFE=1
```
<!--END_GREP-->

### Pragmas

<!--GREP:tests/test_zombiezen_direct.txt:PRAGMA-->
```
auto_vacuum=0
automatic_index=1
busy_timeout=5000
cache_size=-2000
cache_spill=483
cell_size_check=0
checkpoint_fullfsync=0
defer_foreign_keys=0
encoding=UTF-8
foreign_keys=1
fullfsync=0
hard_heap_limit=0
journal_mode=wal
journal_size_limit=-1
locking_mode=normal
mmap_size=0
page_size=4096
query_only=0
read_uncommitted=0
recursive_triggers=0
reverse_unordered_selects=0
secure_delete=0
soft_heap_limit=0
synchronous=1
temp_store=0
threads=0
wal_autocheckpoint=1000
```
<!--END_GREP-->

### SQLite Version

<!--GREP:tests/test_zombiezen_direct.txt:VERSION-->
```
3.49.1
```
<!--END_GREP-->
