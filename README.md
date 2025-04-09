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

Among other things, the tests capture the [compile-time options, pragmas, query plans and SQLite version](#reviewing-the-implementations) used by each implementation. (You can diff one test file against another to get a sense of how and where implementations... well, differ.)

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

  Get one connection from the pool and then return it to the pool. (Every benchmark does that.)

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
Baseline/Conn-12                                     153.4n ± ∞ ¹    222.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   167.6n ± ∞ ¹        ~ (p=0.100 n=3) ²    224.4n ± ∞ ¹        ~ (p=0.100 n=3) ²    220.6n ± ∞ ¹        ~ (p=0.100 n=3) ²    225.4n ± ∞ ¹        ~ (p=0.100 n=3) ²    223.7n ± ∞ ¹       ~ (p=0.100 n=3) ²   1058.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-12                             451.3n ± ∞ ¹    540.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   458.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    546.6n ± ∞ ¹        ~ (p=0.100 n=3) ²    542.1n ± ∞ ¹        ~ (p=0.100 n=3) ²    545.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    541.2n ± ∞ ¹       ~ (p=0.100 n=3) ²   1355.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1-12                                 1773.0n ± ∞ ¹   1925.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   872.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   2036.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2584.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2038.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1106.0n ± ∞ ¹       ~ (p=0.100 n=3) ²   2992.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-12                         1003.0n ± ∞ ¹   1165.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   799.6n ± ∞ ¹        ~ (p=0.100 n=3) ²   2269.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1494.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2317.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1100.0n ± ∞ ¹       ~ (p=0.600 n=3) ²   2638.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePrepared-12                       425.7n ± ∞ ¹    519.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   437.8n ± ∞ ¹        ~ (p=0.100 n=3) ²   2007.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1302.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1986.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    431.7n ± ∞ ¹       ~ (p=0.200 n=3) ²   1486.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-12               741.1n ± ∞ ¹    933.4n ± ∞ ¹        ~ (p=0.100 n=3) ²   692.6n ± ∞ ¹        ~ (p=0.100 n=3) ²   1773.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1034.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1797.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    827.7n ± ∞ ¹       ~ (p=0.100 n=3) ²   1551.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                              581.9n          712.6n        +22.44%                   503.3n        -13.51%                    1.124µ        +93.14%                    923.8n        +58.74%                    1.129µ        +94.04%                    612.2n        +5.19%                    1.722µ        +195.92%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                       │ bench_ncruces_direct.txt │    bench_ncruces_driver.txt    │      bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt    │     bench_mattn_driver.txt      │    bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt    │
                                       │           B/op           │    B/op      vs base           │    B/op      vs base                 │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │
Baseline/Conn-12                                       0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   353.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-12                               0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   353.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1-12                                   48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   705.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-12                           48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePrepared-12                         0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-12                 0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
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
Baseline/ConnParallel                                151.7n ± ∞ ¹    234.5n ± ∞ ¹        ~ (p=0.100 n=3) ²   167.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    243.5n ± ∞ ¹        ~ (p=0.100 n=3) ²    228.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    243.1n ± ∞ ¹        ~ (p=0.100 n=3) ²    235.0n ± ∞ ¹       ~ (p=0.100 n=3) ²   1361.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-2                              154.9n ± ∞ ¹    216.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   158.0n ± ∞ ¹        ~ (p=0.200 n=3) ²    220.5n ± ∞ ¹        ~ (p=0.100 n=3) ²    212.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    218.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    214.8n ± ∞ ¹       ~ (p=0.100 n=3) ²    670.7n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-4                              249.7n ± ∞ ¹    269.4n ± ∞ ¹        ~ (p=0.100 n=3) ²   239.6n ± ∞ ¹        ~ (p=0.200 n=3) ²    263.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    259.0n ± ∞ ¹        ~ (p=0.700 n=3) ²    264.5n ± ∞ ¹        ~ (p=0.100 n=3) ²    263.6n ± ∞ ¹       ~ (p=0.200 n=3) ²    552.8n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-8                              427.1n ± ∞ ¹    512.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   422.4n ± ∞ ¹        ~ (p=0.100 n=3) ²    519.6n ± ∞ ¹        ~ (p=0.100 n=3) ²    513.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    513.1n ± ∞ ¹        ~ (p=0.100 n=3) ²    511.9n ± ∞ ¹       ~ (p=0.100 n=3) ²   1107.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-16                             459.0n ± ∞ ¹    542.5n ± ∞ ¹        ~ (p=0.100 n=3) ²   466.2n ± ∞ ¹        ~ (p=0.200 n=3) ²    548.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    545.1n ± ∞ ¹        ~ (p=0.100 n=3) ²    539.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    547.3n ± ∞ ¹       ~ (p=0.100 n=3) ²   1254.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel                            1769.0n ± ∞ ¹   1900.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   871.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   1882.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1657.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1870.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1141.0n ± ∞ ¹       ~ (p=0.100 n=3) ²   2919.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-2                           923.2n ± ∞ ¹   1010.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   531.4n ± ∞ ¹        ~ (p=0.100 n=3) ²   1086.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    962.4n ± ∞ ¹        ~ (p=0.100 n=3) ²   1069.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    645.6n ± ∞ ¹       ~ (p=0.100 n=3) ²   1537.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-4                           572.4n ± ∞ ¹    561.2n ± ∞ ¹        ~ (p=1.000 n=3) ²   320.4n ± ∞ ¹        ~ (p=0.100 n=3) ²    958.9n ± ∞ ¹        ~ (p=0.100 n=3) ²    603.0n ± ∞ ¹        ~ (p=0.400 n=3) ²    965.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    385.2n ± ∞ ¹       ~ (p=0.100 n=3) ²   1069.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-8                           980.4n ± ∞ ¹    957.2n ± ∞ ¹        ~ (p=0.100 n=3) ²   865.8n ± ∞ ¹        ~ (p=0.100 n=3) ²   2200.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1453.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2227.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    880.7n ± ∞ ¹       ~ (p=0.100 n=3) ²   2472.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-16                         1003.0n ± ∞ ¹   1016.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   898.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   2256.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1587.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2299.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    930.8n ± ∞ ¹       ~ (p=0.100 n=3) ²   2671.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel                  425.6n ± ∞ ¹    521.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   433.4n ± ∞ ¹        ~ (p=0.100 n=3) ²   1853.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    858.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1837.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    424.2n ± ∞ ¹       ~ (p=1.000 n=3) ²   1670.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-2                356.7n ± ∞ ¹    468.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   345.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   1042.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    568.8n ± ∞ ¹        ~ (p=0.100 n=3) ²   1029.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    342.0n ± ∞ ¹       ~ (p=0.100 n=3) ²    814.9n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-4                265.6n ± ∞ ¹    418.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   234.4n ± ∞ ¹        ~ (p=0.100 n=3) ²    735.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    489.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    743.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    273.5n ± ∞ ¹       ~ (p=0.100 n=3) ²    635.8n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-8                710.4n ± ∞ ¹    841.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   725.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   1718.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1021.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1745.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    734.0n ± ∞ ¹       ~ (p=0.100 n=3) ²   1373.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-16               742.4n ± ∞ ¹    883.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   755.6n ± ∞ ¹        ~ (p=0.100 n=3) ²   1763.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1065.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1794.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    774.3n ± ∞ ¹       ~ (p=0.100 n=3) ²   1509.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                              489.9n          582.5n        +18.89%                   426.1n        -13.03%                    888.1n        +81.27%                    661.5n        +35.02%                    888.4n        +81.33%                    485.4n        -0.93%                    1.278µ        +160.92%
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
Baseline/Select1PrePreparedParallel                    0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   372.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-2                  0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-4                  0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-8                  0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-16                 0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                         ⁴   46.15        ?                                +25.99%               ⁴    159.6        ?                    164.2        ?                    159.6        ?                    90.34        ?                    448.4        ?
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
                              │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │       bench_tailscale_driver.txt       │      bench_zombiezen_direct.txt       │
                              │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                │
Populate/PopulateDB-12                       2.506 ± ∞ ¹     2.475 ± ∞ ¹       ~ (p=0.400 n=3) ²    2.005 ± ∞ ¹        ~ (p=0.100 n=3) ²     2.709 ± ∞ ¹        ~ (p=0.100 n=3) ²     2.206 ± ∞ ¹       ~ (p=0.100 n=3) ²     2.697 ± ∞ ¹        ~ (p=0.100 n=3) ²     2.050 ± ∞ ¹        ~ (p=0.100 n=3) ²     2.358 ± ∞ ¹       ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTx-12               1133.2m ± ∞ ¹   1193.1m ± ∞ ¹       ~ (p=0.100 n=3) ²   947.4m ± ∞ ¹        ~ (p=0.100 n=3) ²   1418.1m ± ∞ ¹        ~ (p=0.100 n=3) ²   1179.6m ± ∞ ¹       ~ (p=0.400 n=3) ²   1432.9m ± ∞ ¹        ~ (p=0.100 n=3) ²   1031.1m ± ∞ ¹        ~ (p=0.100 n=3) ²   1102.6m ± ∞ ¹       ~ (p=0.700 n=3) ²
Populate/PopulateDBWithTxs-12                1.094 ± ∞ ¹     1.145 ± ∞ ¹       ~ (p=0.100 n=3) ²    1.076 ± ∞ ¹        ~ (p=0.400 n=3) ²     1.380 ± ∞ ¹        ~ (p=0.100 n=3) ²     1.207 ± ∞ ¹       ~ (p=0.100 n=3) ²     1.363 ± ∞ ¹        ~ (p=0.100 n=3) ²     1.065 ± ∞ ¹        ~ (p=1.000 n=3) ²     1.167 ± ∞ ¹       ~ (p=0.100 n=3) ²
geomean                                      1.459           1.501        +2.87%                    1.269        -13.01%                     1.744        +19.52%                     1.464        +0.38%                     1.740        +19.29%                     1.311        -10.16%                     1.448        -0.76%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │         bench_ncruces_driver.txt          │       bench_eatonphil_direct.txt        │         bench_glebarez_driver.txt         │           bench_mattn_driver.txt            │         bench_modernc_driver.txt          │        bench_tailscale_driver.txt         │       bench_zombiezen_direct.txt        │
                              │           B/op           │      B/op       vs base                   │     B/op       vs base                  │      B/op       vs base                   │      B/op        vs base                    │      B/op       vs base                   │      B/op       vs base                   │     B/op       vs base                  │
Populate/PopulateDB-12                     2.287Mi ± ∞ ¹   18.967Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   5.770Mi ± ∞ ¹         ~ (p=0.100 n=3) ²   31.508Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   273.300Mi ± ∞ ¹           ~ (p=0.100 n=3) ²   31.483Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   18.963Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   6.167Mi ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTx-12               2.276Mi ± ∞ ¹   32.185Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   5.771Mi ± ∞ ¹         ~ (p=0.100 n=3) ²   44.755Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   286.421Mi ± ∞ ¹           ~ (p=0.100 n=3) ²   44.704Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   32.503Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   6.166Mi ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTxs-12              2.275Mi ± ∞ ¹   31.225Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   5.770Mi ± ∞ ¹         ~ (p=0.100 n=3) ²   43.668Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   285.989Mi ± ∞ ¹           ~ (p=0.100 n=3) ²   43.848Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   31.267Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   6.214Mi ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                    2.279Mi          26.71Mi        +1071.93%                   5.770Mi        +153.16%                    39.49Mi        +1632.42%                     281.8Mi        +12264.53%                    39.52Mi        +1633.68%                    26.81Mi        +1076.20%                   6.182Mi        +171.22%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │         bench_mattn_driver.txt          │        bench_modernc_driver.txt         │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                              │        allocs/op         │  allocs/op    vs base                  │  allocs/op    vs base                  │   allocs/op    vs base                  │   allocs/op    vs base                  │   allocs/op    vs base                  │  allocs/op    vs base                  │  allocs/op    vs base                  │
Populate/PopulateDB-12                      140.0k ± ∞ ¹   598.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   394.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1209.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   1057.5k ± ∞ ¹         ~ (p=0.100 n=3) ²   1209.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   598.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   445.1k ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTx-12                140.0k ± ∞ ¹   751.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   394.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1362.5k ± ∞ ¹         ~ (p=0.100 n=3) ²   1210.4k ± ∞ ¹         ~ (p=0.100 n=3) ²   1362.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   751.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   445.1k ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTxs-12               140.0k ± ∞ ¹   765.2k ± ∞ ¹         ~ (p=0.100 n=3) ²   394.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1376.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   1239.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1380.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   767.2k ± ∞ ¹         ~ (p=0.100 n=3) ²   448.1k ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                     140.0k         700.5k        +400.31%                   394.0k        +181.41%                    1.314M        +838.33%                    1.166M        +732.90%                    1.315M        +839.17%                   701.1k        +400.76%                   446.1k        +218.60%
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
                                                                    │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │         bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                                                                    │          sec/op          │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                 │
ReadWrite/ReadPost-12                                                            14.083µ ± ∞ ¹   21.918µ ± ∞ ¹        ~ (p=0.100 n=3) ²   9.154µ ± ∞ ¹        ~ (p=0.100 n=3) ²   28.314µ ± ∞ ¹        ~ (p=0.100 n=3) ²   17.663µ ± ∞ ¹        ~ (p=0.100 n=3) ²   28.630µ ± ∞ ¹        ~ (p=0.100 n=3) ²   16.850µ ± ∞ ¹       ~ (p=0.100 n=3) ²   17.596µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadPostWithTx-12                                                      16.181µ ± ∞ ¹   37.483µ ± ∞ ¹        ~ (p=0.100 n=3) ²   9.834µ ± ∞ ¹        ~ (p=0.100 n=3) ²   31.564µ ± ∞ ¹        ~ (p=0.100 n=3) ²   21.857µ ± ∞ ¹        ~ (p=0.100 n=3) ²   32.214µ ± ∞ ¹        ~ (p=0.100 n=3) ²   24.518µ ± ∞ ¹       ~ (p=0.100 n=3) ²   18.067µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndComments-12                                                  99.29µ ± ∞ ¹   115.58µ ± ∞ ¹        ~ (p=0.100 n=3) ²   70.73µ ± ∞ ¹        ~ (p=0.100 n=3) ²   173.15µ ± ∞ ¹        ~ (p=0.100 n=3) ²   178.31µ ± ∞ ¹        ~ (p=0.100 n=3) ²   174.90µ ± ∞ ¹        ~ (p=0.100 n=3) ²    98.21µ ± ∞ ¹       ~ (p=0.100 n=3) ²    91.48µ ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadPostAndCommentsWithTx-12                                            99.29µ ± ∞ ¹   142.52µ ± ∞ ¹        ~ (p=0.100 n=3) ²   70.36µ ± ∞ ¹        ~ (p=0.100 n=3) ²   173.27µ ± ∞ ¹        ~ (p=0.100 n=3) ²   180.98µ ± ∞ ¹        ~ (p=0.100 n=3) ²   176.50µ ± ∞ ¹        ~ (p=0.100 n=3) ²   110.02µ ± ∞ ¹       ~ (p=0.100 n=3) ²    89.63µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePost-12                                                            129.0µ ± ∞ ¹    130.5µ ± ∞ ¹        ~ (p=1.000 n=3) ²   115.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    143.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    128.1µ ± ∞ ¹        ~ (p=0.400 n=3) ²    143.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    114.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²    166.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePostWithTx-12                                                      136.7µ ± ∞ ¹    145.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   115.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    150.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    131.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    144.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    125.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²    164.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndComments-12                                                 2.545m ± ∞ ¹    2.605m ± ∞ ¹        ~ (p=0.700 n=3) ²   2.765m ± ∞ ¹        ~ (p=0.200 n=3) ²    2.868m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.316m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.848m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.118m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.536m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndCommentsWithTx-12                                           953.5µ ± ∞ ¹   1012.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   949.3µ ± ∞ ¹        ~ (p=1.000 n=3) ²   1252.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    998.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1227.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    824.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²   1095.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=10-12                             351.0µ ± ∞ ¹    368.4µ ± ∞ ¹        ~ (p=0.200 n=3) ²   336.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    422.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    392.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    457.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    300.6µ ± ∞ ¹       ~ (p=0.100 n=3) ²    414.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=90-12                             2.377m ± ∞ ¹    2.401m ± ∞ ¹        ~ (p=0.700 n=3) ²   2.162m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.570m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.122m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.596m ± ∞ ¹        ~ (p=0.100 n=3) ²    1.946m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.294m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-12                     385.4µ ± ∞ ¹    405.8µ ± ∞ ¹        ~ (p=1.000 n=3) ²   297.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    388.1µ ± ∞ ¹        ~ (p=1.000 n=3) ²    319.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    401.8µ ± ∞ ¹        ~ (p=1.000 n=3) ²    325.1µ ± ∞ ¹       ~ (p=0.100 n=3) ²    365.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-12                     3.174m ± ∞ ¹    3.050m ± ∞ ¹        ~ (p=0.200 n=3) ²   2.437m ± ∞ ¹        ~ (p=0.100 n=3) ²    3.018m ± ∞ ¹        ~ (p=0.400 n=3) ²    2.761m ± ∞ ¹        ~ (p=0.400 n=3) ²    2.715m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.587m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.975m ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=10-12                       188.2µ ± ∞ ¹    228.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   150.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    299.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    268.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    285.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    187.1µ ± ∞ ¹       ~ (p=0.700 n=3) ²    262.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=90-12                       892.3µ ± ∞ ¹    936.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   741.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1183.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    939.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1165.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    765.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²   1037.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-12               131.8µ ± ∞ ¹    132.5µ ± ∞ ¹        ~ (p=0.700 n=3) ²   116.8µ ± ∞ ¹        ~ (p=0.700 n=3) ²    206.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    175.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    207.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    138.8µ ± ∞ ¹       ~ (p=0.400 n=3) ²    178.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-12               850.3µ ± ∞ ¹    907.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   820.5µ ± ∞ ¹        ~ (p=0.700 n=3) ²   2904.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1179.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   3512.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    894.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²   2103.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                                                           284.5µ          329.3µ        +15.72%                   235.3µ        -17.30%                    406.7µ        +42.92%                    328.8µ        +15.55%                    410.3µ        +44.22%                    273.5µ        -3.87%                    331.2µ        +16.40%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                                                    │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │    bench_eatonphil_direct.txt    │     bench_glebarez_driver.txt     │       bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                                                                    │           B/op           │      B/op       vs base           │     B/op       vs base           │      B/op       vs base           │      B/op        vs base           │      B/op       vs base           │      B/op       vs base           │     B/op       vs base           │
ReadWrite/ReadPost-12                                                            40.17Ki ± ∞ ¹    41.23Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.18Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    81.34Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     41.39Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    81.28Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    41.19Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.55Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostWithTx-12                                                      40.17Ki ± ∞ ¹    42.29Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.18Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    82.22Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     42.60Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    82.33Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    42.13Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.60Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndComments-12                                                 250.2Ki ± ∞ ¹    260.3Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    501.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     263.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    501.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    257.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndCommentsWithTx-12                                           250.2Ki ± ∞ ¹    261.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    503.1Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     265.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    503.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    258.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePost-12                                                              0.00 ± ∞ ¹     240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     496.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    42288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     496.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    408.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostWithTx-12                                                        0.00 ± ∞ ¹     858.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    1107.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    43252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    1289.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     881.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    457.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndComments-12                                                0.000Ki ± ∞ ¹   13.906Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   2.781Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.663Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   308.876Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.672Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   13.906Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   3.540Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndCommentsWithTx-12                                          0.000Ki ± ∞ ¹   26.449Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   2.781Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   39.189Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   321.744Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   39.373Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.469Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   3.577Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=10-12                            224.7Ki ± ∞ ¹    235.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.0Ki ± ∞ ¹  ~ (p=1.000 n=3) ²    458.0Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     268.3Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    452.0Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    232.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.3Ki ± ∞ ¹  ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=90-12                            20.68Ki ± ∞ ¹    40.59Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   23.90Ki ± ∞ ¹  ~ (p=0.700 n=3) ²    80.82Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    304.24Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    73.75Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    36.35Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   31.94Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-12                    225.6Ki ± ∞ ¹    235.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.6Ki ± ∞ ¹  ~ (p=1.000 n=3) ²    457.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     268.1Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    454.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    233.5Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.7Ki ± ∞ ¹  ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-12                    21.85Ki ± ∞ ¹    39.66Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   30.08Ki ± ∞ ¹  ~ (p=0.200 n=3) ²    71.55Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    304.82Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    75.63Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    40.32Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.77Ki ± ∞ ¹  ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=10-12                      225.7Ki ± ∞ ¹    237.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.0Ki ± ∞ ¹  ~ (p=0.200 n=3) ²    455.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     271.4Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    457.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    234.1Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   227.0Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=90-12                      25.19Ki ± ∞ ¹    49.41Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   27.09Ki ± ∞ ¹  ~ (p=0.400 n=3) ²    87.38Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    316.41Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    84.91Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    49.45Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   29.44Ki ± ∞ ¹  ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-12              225.2Ki ± ∞ ¹    238.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.3Ki ± ∞ ¹  ~ (p=1.000 n=3) ²    457.1Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     271.1Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    457.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    235.5Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.7Ki ± ∞ ¹  ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-12              28.40Ki ± ∞ ¹    47.54Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.00Ki ± ∞ ¹  ~ (p=0.400 n=3) ²    83.35Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    316.24Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    89.08Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    49.98Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   32.07Ki ± ∞ ¹  ~ (p=0.400 n=3) ²
geomean                                                                                      ³    42.22Ki        ?                   21.60Ki        ?                    77.51Ki        ?                     178.4Ki        ?                    78.21Ki        ?                    41.98Ki        ?                   30.24Ki        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean

                                                                    │ bench_ncruces_direct.txt │    bench_ncruces_driver.txt     │   bench_eatonphil_direct.txt    │    bench_glebarez_driver.txt     │     bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt    │
                                                                    │        allocs/op         │  allocs/op    vs base           │  allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op    vs base           │
ReadWrite/ReadPost-12                                                              5.000 ± ∞ ¹   39.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    42.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   42.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    42.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   31.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   12.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostWithTx-12                                                        5.000 ± ∞ ¹   60.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    57.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   70.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    60.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   50.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   15.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndComments-12                                                   262.0 ± ∞ ¹    773.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    264.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     976.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    684.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    562.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    272.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndCommentsWithTx-12                                             262.0 ± ∞ ¹    801.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    264.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     997.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    719.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1000.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    587.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    275.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePost-12                                                             0.000 ± ∞ ¹    7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    3.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    17.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   15.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    17.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    9.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostWithTx-12                                                       0.000 ± ∞ ¹   18.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    3.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    28.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   40.000 ± ∞ ¹  ~ (p=0.100 n=3) ²    32.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   20.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   12.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndComments-12                                                    0.0 ± ∞ ¹    407.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    203.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     967.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    815.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     967.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    407.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    259.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndCommentsWithTx-12                                              0.0 ± ∞ ¹    574.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    203.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1134.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    996.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1138.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    576.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    262.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=10-12                              235.0 ± ∞ ¹    736.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    696.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=90-12                              21.00 ± ∞ ¹   446.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   801.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   421.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-12                      236.0 ± ∞ ¹    736.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    696.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    547.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-12                      23.00 ± ∞ ¹   445.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   803.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   423.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=10-12                        236.0 ± ∞ ¹    777.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1011.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    747.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    274.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=90-12                        26.00 ± ∞ ¹   596.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1119.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   969.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1124.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-12                235.0 ± ∞ ¹    778.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1010.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-12                29.00 ± ∞ ¹   594.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1121.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1123.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                                                      ³    273.5        ?                    85.15        ?                     430.4        ?                    368.3        ?                     435.9        ?                    237.4        ?                    122.4        ?
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
                                                                    │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │         bench_mattn_driver.txt         │        bench_modernc_driver.txt         │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt        │
                                                                    │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │     sec/op      vs base                 │    sec/op      vs base                 │     sec/op      vs base                 │    sec/op      vs base                │     sec/op      vs base                 │
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10                        361.7µ ± ∞ ¹    339.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²   285.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²     465.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    360.2µ ± ∞ ¹        ~ (p=0.700 n=3) ²     449.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    291.9µ ± ∞ ¹       ~ (p=0.100 n=3) ²     431.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-2                      316.5µ ± ∞ ¹    334.3µ ± ∞ ¹       ~ (p=0.400 n=3) ²   266.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²     435.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    322.0µ ± ∞ ¹        ~ (p=0.700 n=3) ²     414.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    262.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²     340.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-4                      344.4µ ± ∞ ¹    327.1µ ± ∞ ¹       ~ (p=0.400 n=3) ²   265.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²     400.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    312.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²     363.1µ ± ∞ ¹        ~ (p=0.400 n=3) ²    292.1µ ± ∞ ¹       ~ (p=0.100 n=3) ²     355.1µ ± ∞ ¹        ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-8                      386.1µ ± ∞ ¹    387.4µ ± ∞ ¹       ~ (p=0.700 n=3) ²   271.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²     406.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    343.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²     398.4µ ± ∞ ¹        ~ (p=0.200 n=3) ²    303.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²     355.9µ ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-16                     418.7µ ± ∞ ¹    427.7µ ± ∞ ¹       ~ (p=0.700 n=3) ²   297.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²     437.7µ ± ∞ ¹        ~ (p=0.400 n=3) ²    387.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²     437.3µ ± ∞ ¹        ~ (p=0.200 n=3) ²    324.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²     401.5µ ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90                        2.328m ± ∞ ¹    2.324m ± ∞ ¹       ~ (p=1.000 n=3) ²   1.970m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.529m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.091m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.523m ± ∞ ¹        ~ (p=0.100 n=3) ²    1.978m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.354m ± ∞ ¹        ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-2                      2.553m ± ∞ ¹    2.782m ± ∞ ¹       ~ (p=0.100 n=3) ²   1.936m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.540m ± ∞ ¹        ~ (p=1.000 n=3) ²    2.146m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.629m ± ∞ ¹        ~ (p=0.200 n=3) ²    2.009m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.376m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-4                      2.637m ± ∞ ¹    2.767m ± ∞ ¹       ~ (p=0.200 n=3) ²   2.112m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.645m ± ∞ ¹        ~ (p=1.000 n=3) ²    2.258m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.863m ± ∞ ¹        ~ (p=0.200 n=3) ²    2.028m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.485m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-8                      2.994m ± ∞ ¹    3.162m ± ∞ ¹       ~ (p=0.100 n=3) ²   2.105m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.623m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.351m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.717m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.201m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.952m ± ∞ ¹        ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-16                     3.501m ± ∞ ¹    3.485m ± ∞ ¹       ~ (p=1.000 n=3) ²   2.753m ± ∞ ¹        ~ (p=0.200 n=3) ²     4.175m ± ∞ ¹        ~ (p=0.700 n=3) ²    2.887m ± ∞ ¹        ~ (p=0.100 n=3) ²     6.322m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.566m ± ∞ ¹       ~ (p=0.200 n=3) ²    11.550m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10                  214.9µ ± ∞ ¹    225.6µ ± ∞ ¹       ~ (p=0.700 n=3) ²   203.0µ ± ∞ ¹        ~ (p=0.700 n=3) ²     296.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    234.3µ ± ∞ ¹        ~ (p=0.400 n=3) ²     306.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    187.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²     279.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-2                128.5µ ± ∞ ¹    138.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²   113.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²     213.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    174.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²     213.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    131.3µ ± ∞ ¹       ~ (p=0.400 n=3) ²     188.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-4                112.6µ ± ∞ ¹    124.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²   112.7µ ± ∞ ¹        ~ (p=1.000 n=3) ²     189.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    144.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²     193.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    120.7µ ± ∞ ¹       ~ (p=0.700 n=3) ²     169.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-8                102.2µ ± ∞ ¹    109.9µ ± ∞ ¹       ~ (p=0.100 n=3) ²   107.6µ ± ∞ ¹        ~ (p=0.700 n=3) ²     199.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    225.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²     204.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    128.6µ ± ∞ ¹       ~ (p=0.100 n=3) ²     176.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-16               125.0µ ± ∞ ¹    148.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²   113.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²     223.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    291.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²     249.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    160.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²     194.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90                  895.0µ ± ∞ ¹    951.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²   765.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1311.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1123.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1158.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    748.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²    1050.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-2                750.0µ ± ∞ ¹    816.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²   739.3µ ± ∞ ¹        ~ (p=1.000 n=3) ²    1213.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1057.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1176.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    752.8µ ± ∞ ¹       ~ (p=0.700 n=3) ²     978.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-4                752.6µ ± ∞ ¹    806.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²   753.6µ ± ∞ ¹        ~ (p=1.000 n=3) ²    1193.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1016.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1239.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    786.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²     990.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-8                843.4µ ± ∞ ¹    899.6µ ± ∞ ¹       ~ (p=0.100 n=3) ²   786.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1817.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1070.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1638.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    819.0µ ± ∞ ¹       ~ (p=0.700 n=3) ²    1053.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-16              1069.1µ ± ∞ ¹   1055.6µ ± ∞ ¹       ~ (p=0.700 n=3) ²   795.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   10463.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1461.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   12563.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1041.4µ ± ∞ ¹       ~ (p=1.000 n=3) ²   10578.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                                                           580.6µ          606.6µ        +4.49%                   490.4µ        -15.54%                     862.5µ        +48.56%                    659.6µ        +13.61%                     884.5µ        +52.35%                    522.6µ        -9.98%                     802.4µ        +38.21%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                                                    │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │          bench_mattn_driver.txt          │        bench_modernc_driver.txt         │       bench_tailscale_driver.txt       │      bench_zombiezen_direct.txt       │
                                                                    │           B/op           │     B/op       vs base                 │     B/op       vs base                │     B/op       vs base                  │      B/op       vs base                  │     B/op       vs base                  │     B/op       vs base                 │     B/op       vs base                │
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10                       224.3Ki ± ∞ ¹   236.1Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.1Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   451.7Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   455.6Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   233.4Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.2Ki ± ∞ ¹       ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-2                     226.0Ki ± ∞ ¹   234.9Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.8Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   453.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   452.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   233.5Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.6Ki ± ∞ ¹       ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-4                     224.1Ki ± ∞ ¹   237.0Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.2Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   453.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.5Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   454.6Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   233.8Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.1Ki ± ∞ ¹       ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-8                     223.7Ki ± ∞ ¹   234.3Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.8Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   454.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.6Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   453.7Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   231.9Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.7Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-16                    224.8Ki ± ∞ ¹   235.9Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.8Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   453.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.7Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   454.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   231.8Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.6Ki ± ∞ ¹       ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90                       27.25Ki ± ∞ ¹   43.02Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   25.82Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   73.95Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.64Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   74.04Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   40.01Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.23Ki ± ∞ ¹       ~ (p=0.500 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-2                     27.07Ki ± ∞ ¹   35.40Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   29.24Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   75.91Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   303.95Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   70.35Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   40.89Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.74Ki ± ∞ ¹       ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-4                     26.82Ki ± ∞ ¹   41.99Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   30.83Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   66.08Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.03Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   70.17Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   36.21Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   21.42Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-8                     21.70Ki ± ∞ ¹   35.69Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   25.46Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   72.69Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.14Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   75.48Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   40.57Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   29.81Ki ± ∞ ¹       ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-16                    24.36Ki ± ∞ ¹   38.42Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.77Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   77.52Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.62Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   74.24Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   39.41Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   25.85Ki ± ∞ ¹       ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10                 224.9Ki ± ∞ ¹   237.8Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.5Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   457.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.0Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   455.7Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   234.2Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.9Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-2               224.9Ki ± ∞ ¹   238.6Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.7Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   458.9Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   455.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   235.2Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.7Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-4               223.9Ki ± ∞ ¹   238.0Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.7Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   457.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   457.6Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   234.5Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.7Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-8               225.0Ki ± ∞ ¹   239.4Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.1Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   456.0Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   457.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   236.5Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.3Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-16              224.3Ki ± ∞ ¹   237.5Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.9Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   456.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   457.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   236.1Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.8Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90                 26.06Ki ± ∞ ¹   50.62Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   26.07Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   83.75Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.28Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   80.32Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   49.78Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.15Ki ± ∞ ¹       ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-2               26.28Ki ± ∞ ¹   48.84Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.58Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   91.96Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.30Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   86.11Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   49.74Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   29.17Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-4               25.42Ki ± ∞ ¹   49.29Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.81Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   85.22Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.12Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   87.72Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   48.93Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   29.16Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-8               23.37Ki ± ∞ ¹   49.54Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   26.51Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   82.72Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.46Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   76.41Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   51.54Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   30.31Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-16              23.34Ki ± ∞ ¹   50.58Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.25Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   76.42Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   317.02Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   92.57Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   47.64Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.38Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
geomean                                                                          75.09Ki         102.0Ki        +35.88%                   78.95Ki        +5.15%                   188.8Ki        +151.41%                    289.3Ki        +285.31%                   189.0Ki        +151.66%                   101.7Ki        +35.40%                   79.10Ki        +5.35%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                                                    │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt         │         bench_mattn_driver.txt         │         bench_modernc_driver.txt         │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                                                                    │        allocs/op         │  allocs/op    vs base                  │  allocs/op    vs base                  │   allocs/op    vs base                   │  allocs/op    vs base                  │   allocs/op    vs base                   │  allocs/op    vs base                  │  allocs/op    vs base                  │
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10                         235.0 ± ∞ ¹    735.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     973.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    695.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     973.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    545.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    270.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-2                       236.0 ± ∞ ¹    734.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    696.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    270.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-4                       234.0 ± ∞ ¹    738.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    697.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    547.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-8                       234.0 ± ∞ ¹    734.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    697.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-16                      235.0 ± ∞ ¹    736.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    698.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90                         28.00 ± ∞ ¹   450.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   802.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   423.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-2                       28.00 ± ∞ ¹   438.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   800.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   424.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-4                       28.00 ± ∞ ¹   448.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   800.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   421.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-8                       22.00 ± ∞ ¹   439.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   801.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   424.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-16                      25.00 ± ∞ ¹   443.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   802.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   423.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   261.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10                   235.0 ± ∞ ¹    777.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1009.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    745.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    585.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-2                 235.0 ± ∞ ¹    777.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1009.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    745.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    585.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-4                 234.0 ± ∞ ¹    777.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1010.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1013.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-8                 235.0 ± ∞ ¹    779.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1011.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-16                234.0 ± ∞ ¹    777.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1011.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90                   27.00 ± ∞ ¹   597.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1121.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   969.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1125.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-2                 27.00 ± ∞ ¹   595.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1118.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   969.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1124.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-4                 26.00 ± ∞ ¹   596.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1120.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   968.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1124.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-8                 24.00 ± ∞ ¹   596.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1121.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   970.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1127.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-16                24.00 ± ∞ ¹   597.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1123.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   972.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1122.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   264.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                                                            77.85          623.6        +700.99%                    231.6        +197.48%                    1.016k        +1205.56%                    796.9        +923.70%                    1.018k        +1207.69%                    528.5        +578.93%                    266.7        +242.58%
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
                            │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │      bench_zombiezen_direct.txt       │
                            │          sec/op          │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                │
Query/Correlated-12                       649.3m ± ∞ ¹    663.2m ± ∞ ¹        ~ (p=0.100 n=3) ²   287.3m ± ∞ ¹        ~ (p=0.100 n=3) ²    314.5m ± ∞ ¹        ~ (p=0.100 n=3) ²   178.2m ± ∞ ¹        ~ (p=0.100 n=3) ²    317.5m ± ∞ ¹        ~ (p=0.100 n=3) ²   174.5m ± ∞ ¹        ~ (p=0.100 n=3) ²    329.1m ± ∞ ¹       ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-12               84.52m ± ∞ ¹   113.60m ± ∞ ¹        ~ (p=0.700 n=3) ²   50.68m ± ∞ ¹        ~ (p=0.100 n=3) ²   118.29m ± ∞ ¹        ~ (p=0.400 n=3) ²   50.03m ± ∞ ¹        ~ (p=0.100 n=3) ²   121.14m ± ∞ ¹        ~ (p=0.200 n=3) ²   45.95m ± ∞ ¹        ~ (p=0.100 n=3) ²   187.24m ± ∞ ¹       ~ (p=0.100 n=3) ²
geomean                                   234.3m          274.5m        +17.17%                   120.7m        -48.49%                    192.9m        -17.66%                   94.43m        -59.69%                    196.1m        -16.28%                   89.55m        -61.77%                    248.2m        +5.97%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt          │      bench_eatonphil_direct.txt      │         bench_glebarez_driver.txt          │           bench_mattn_driver.txt            │          bench_modernc_driver.txt          │         bench_tailscale_driver.txt         │       bench_zombiezen_direct.txt        │
                            │           B/op           │      B/op       vs base                    │    B/op      vs base                 │      B/op       vs base                    │      B/op        vs base                    │      B/op       vs base                    │      B/op       vs base                    │     B/op       vs base                  │
Query/Correlated-12                        96.00 ± ∞ ¹   71216.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   48.00 ± ∞ ¹        ~ (p=0.300 n=3) ²   63072.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   207261.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   63088.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   47024.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   2776.00 ± ∞ ¹         ~ (p=0.600 n=3) ²
Query/CorrelatedParallel-12               1903.0 ± ∞ ¹    71107.0 ± ∞ ¹           ~ (p=0.700 n=3) ²   545.0 ± ∞ ¹        ~ (p=0.400 n=3) ²    63577.0 ± ∞ ¹           ~ (p=0.700 n=3) ²    207829.0 ± ∞ ¹           ~ (p=0.100 n=3) ²    63829.0 ± ∞ ¹           ~ (p=0.700 n=3) ²    47348.0 ± ∞ ¹           ~ (p=0.700 n=3) ²    1585.0 ± ∞ ¹         ~ (p=0.700 n=3) ²
geomean                                    427.4          69.49Ki        +16549.07%                   161.7        -62.16%                    61.84Ki        +14715.40%                     202.7Ki        +48457.56%                    61.97Ki        +14746.61%                    46.08Ki        +10939.66%                   2.048Ki        +390.76%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt           │    bench_eatonphil_direct.txt    │          bench_glebarez_driver.txt          │           bench_mattn_driver.txt            │          bench_modernc_driver.txt           │         bench_tailscale_driver.txt          │       bench_zombiezen_direct.txt        │
                            │        allocs/op         │   allocs/op     vs base                     │  allocs/op   vs base             │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                     │  allocs/op    vs base                   │
Query/Correlated-12                        1.000 ± ∞ ¹   5770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.300 n=3) ²     6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   31.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-12                3.000 ± ∞ ¹   5769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=0.300 n=3) ²     6778.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6773.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6780.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   14.000 ± ∞ ¹          ~ (p=0.600 n=3) ²
geomean                                    1.732           5.769k        +333002.24%                                ?               ³ ⁴     6.773k        +390968.12%                     6.772k        +390910.47%                     6.774k        +391025.81%                     3.766k        +217358.98%                    20.83        +1102.77%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean
⁴ ratios must be >0 to compute geomean
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
Query/CorrelatedParallel                  683.6m ± ∞ ¹    657.1m ± ∞ ¹       ~ (p=0.700 n=3) ²    286.5m ± ∞ ¹        ~ (p=0.100 n=3) ²    300.7m ± ∞ ¹        ~ (p=0.100 n=3) ²   176.6m ± ∞ ¹        ~ (p=0.100 n=3) ²    305.3m ± ∞ ¹        ~ (p=0.100 n=3) ²   185.4m ± ∞ ¹        ~ (p=0.100 n=3) ²    328.7m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-2               439.53m ± ∞ ¹   445.28m ± ∞ ¹       ~ (p=0.400 n=3) ²   172.67m ± ∞ ¹        ~ (p=0.100 n=3) ²   185.33m ± ∞ ¹        ~ (p=0.100 n=3) ²   98.10m ± ∞ ¹        ~ (p=0.100 n=3) ²   193.21m ± ∞ ¹        ~ (p=0.100 n=3) ²   99.00m ± ∞ ¹        ~ (p=0.100 n=3) ²   191.93m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-4               230.90m ± ∞ ¹   231.58m ± ∞ ¹       ~ (p=0.700 n=3) ²    85.30m ± ∞ ¹        ~ (p=0.100 n=3) ²   118.18m ± ∞ ¹        ~ (p=0.100 n=3) ²   60.57m ± ∞ ¹        ~ (p=0.100 n=3) ²   127.83m ± ∞ ¹        ~ (p=0.100 n=3) ²   58.45m ± ∞ ¹        ~ (p=0.100 n=3) ²   135.77m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-8               157.26m ± ∞ ¹   155.57m ± ∞ ¹       ~ (p=0.400 n=3) ²    56.71m ± ∞ ¹        ~ (p=0.100 n=3) ²   138.34m ± ∞ ¹        ~ (p=0.200 n=3) ²   43.44m ± ∞ ¹        ~ (p=0.100 n=3) ²   138.49m ± ∞ ¹        ~ (p=0.100 n=3) ²   44.47m ± ∞ ¹        ~ (p=0.100 n=3) ²   187.32m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-16               82.07m ± ∞ ¹    81.23m ± ∞ ¹       ~ (p=1.000 n=3) ²    55.11m ± ∞ ¹        ~ (p=0.100 n=3) ²   117.09m ± ∞ ¹        ~ (p=0.100 n=3) ²   45.84m ± ∞ ¹        ~ (p=0.100 n=3) ²   114.25m ± ∞ ¹        ~ (p=0.100 n=3) ²   46.35m ± ∞ ¹        ~ (p=0.100 n=3) ²   184.76m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                   245.7m          243.5m        -0.89%                    105.7m        -56.98%                    160.6m        -34.65%                   73.12m        -70.24%                    164.2m        -33.18%                   73.95m        -69.90%                    197.0m        -19.83%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt          │      bench_eatonphil_direct.txt       │         bench_glebarez_driver.txt          │            bench_mattn_driver.txt            │          bench_modernc_driver.txt          │         bench_tailscale_driver.txt         │       bench_zombiezen_direct.txt        │
                            │           B/op           │      B/op       vs base                    │     B/op      vs base                 │      B/op       vs base                    │      B/op        vs base                     │      B/op       vs base                    │      B/op       vs base                    │     B/op       vs base                  │
Query/CorrelatedParallel                   68.00 ± ∞ ¹   71132.00 ± ∞ ¹           ~ (p=0.100 n=3) ²    34.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   63050.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   207185.00 ± ∞ ¹            ~ (p=0.100 n=3) ²   63066.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   47009.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   2216.00 ± ∞ ¹         ~ (p=0.600 n=3) ²
Query/CorrelatedParallel-2                 104.0 ± ∞ ¹    71082.0 ± ∞ ¹           ~ (p=0.100 n=3) ²    108.0 ± ∞ ¹        ~ (p=1.000 n=3) ²    63011.0 ± ∞ ¹           ~ (p=0.100 n=3) ²    207186.0 ± ∞ ¹            ~ (p=0.100 n=3) ²    63029.0 ± ∞ ¹           ~ (p=0.100 n=3) ²    46978.0 ± ∞ ¹           ~ (p=0.100 n=3) ²     555.0 ± ∞ ¹         ~ (p=0.700 n=3) ²
Query/CorrelatedParallel-4                 94.00 ± ∞ ¹   71052.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   188.00 ± ∞ ¹        ~ (p=1.000 n=3) ²   63048.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   207225.00 ± ∞ ¹            ~ (p=0.100 n=3) ²   63117.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   46964.00 ± ∞ ¹           ~ (p=0.100 n=3) ²    596.00 ± ∞ ¹         ~ (p=0.700 n=3) ²
Query/CorrelatedParallel-8                 88.00 ± ∞ ¹   71050.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   165.00 ± ∞ ¹        ~ (p=1.000 n=3) ²   63132.00 ± ∞ ¹           ~ (p=0.400 n=3) ²   207202.00 ± ∞ ¹            ~ (p=0.100 n=3) ²   63122.00 ± ∞ ¹           ~ (p=0.700 n=3) ²   47243.00 ± ∞ ¹           ~ (p=0.700 n=3) ²    662.00 ± ∞ ¹         ~ (p=0.700 n=3) ²
Query/CorrelatedParallel-16                927.0 ± ∞ ¹    71594.0 ± ∞ ¹           ~ (p=0.100 n=3) ²   2642.0 ± ∞ ¹        ~ (p=1.000 n=3) ²    63670.0 ± ∞ ¹           ~ (p=0.100 n=3) ²    209118.0 ± ∞ ¹            ~ (p=0.100 n=3) ²    64344.0 ± ∞ ¹           ~ (p=0.100 n=3) ²    47654.0 ± ∞ ¹           ~ (p=0.100 n=3) ²     774.0 ± ∞ ¹         ~ (p=1.000 n=3) ²
geomean                                    140.2          69.51Ki        +50660.03%                    197.6        +40.88%                    61.70Ki        +44955.20%                     202.7Ki        +147927.62%                    61.85Ki        +45063.51%                    46.06Ki        +33536.36%                     822.1        +486.27%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt           │      bench_eatonphil_direct.txt      │          bench_glebarez_driver.txt          │           bench_mattn_driver.txt            │          bench_modernc_driver.txt           │         bench_tailscale_driver.txt          │       bench_zombiezen_direct.txt       │
                            │        allocs/op         │   allocs/op     vs base                     │  allocs/op   vs base                 │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                     │  allocs/op    vs base                  │
Query/CorrelatedParallel                   3.000 ± ∞ ¹   5771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   29.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-2                 3.000 ± ∞ ¹   5770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹            ~ (p=0.100 n=3) ²    9.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-4                 2.000 ± ∞ ¹   5769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹            ~ (p=0.100 n=3) ²    8.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-8                 2.000 ± ∞ ¹   5770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   10.000 ± ∞ ¹         ~ (p=0.600 n=3) ²
Query/CorrelatedParallel-16                3.000 ± ∞ ¹   5772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   4.000 ± ∞ ¹        ~ (p=1.000 n=3) ²   6779.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6775.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6788.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3768.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   13.000 ± ∞ ¹         ~ (p=0.600 n=3) ²
geomean                                    2.551           5.770k        +226114.88%                   1.320        -48.27%                     6.772k        +265364.52%                     6.773k        +265403.76%                     6.773k        +265434.97%                     3.767k        +147568.48%                    12.21        +378.69%
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
                         │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │      bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │         bench_mattn_driver.txt          │        bench_modernc_driver.txt         │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt        │
                         │          sec/op          │    sec/op     vs base                │    sec/op     vs base                 │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op     vs base                 │    sec/op      vs base                  │
Query/GroupBy-12                       556.4µ ± ∞ ¹   573.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²   278.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    506.3µ ± ∞ ¹         ~ (p=0.100 n=3) ²    341.1µ ± ∞ ¹         ~ (p=0.100 n=3) ²    512.8µ ± ∞ ¹         ~ (p=0.100 n=3) ²   283.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    500.5µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-12               62.00µ ± ∞ ¹   66.32µ ± ∞ ¹       ~ (p=0.100 n=3) ²   40.54µ ± ∞ ¹        ~ (p=0.100 n=3) ²   664.36µ ± ∞ ¹         ~ (p=0.100 n=3) ²   482.27µ ± ∞ ¹         ~ (p=0.100 n=3) ²   670.37µ ± ∞ ¹         ~ (p=0.100 n=3) ²   47.13µ ± ∞ ¹        ~ (p=0.100 n=3) ²   663.89µ ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                185.7µ         195.0µ        +4.98%                   106.2µ        -42.80%                    579.9µ        +212.25%                    405.6µ        +118.36%                    586.3µ        +215.67%                   115.6µ        -37.77%                    576.4µ        +210.35%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │    bench_eatonphil_direct.txt     │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │    bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                         │           B/op           │     B/op       vs base           │   B/op     vs base                │     B/op       vs base           │     B/op       vs base           │     B/op       vs base           │     B/op       vs base           │    B/op      vs base           │
Query/GroupBy-12                          0.0 ± ∞ ¹    2382.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ²    1963.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    7161.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1976.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1580.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   363.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-12                  0.0 ± ∞ ¹    2369.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1966.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    7168.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1980.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1525.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   370.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ⁴   2.320Ki        ?                              +0.00%               ⁴   1.918Ki        ?                   6.997Ki        ?                   1.932Ki        ?                   1.516Ki        ?                   366.5        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                         │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/GroupBy-12                        0.000 ± ∞ ¹   119.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   155.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-12                0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   155.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ⁴     118.5        ?                                +0.00%               ⁴     122.0        ?                     155.0        ?                     122.0        ?                    51.00        ?                   6.000        ?
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
Query/GroupByParallel                  554.2µ ± ∞ ¹    558.6µ ± ∞ ¹       ~ (p=0.700 n=3) ²   291.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    523.7µ ± ∞ ¹         ~ (p=0.100 n=3) ²    315.9µ ± ∞ ¹         ~ (p=0.100 n=3) ²    499.2µ ± ∞ ¹         ~ (p=0.100 n=3) ²   266.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    482.9µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-2                281.9µ ± ∞ ¹    285.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²   148.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    386.5µ ± ∞ ¹         ~ (p=0.100 n=3) ²    391.4µ ± ∞ ¹         ~ (p=0.100 n=3) ²    391.9µ ± ∞ ¹         ~ (p=0.100 n=3) ²   131.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    384.4µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-4               143.90µ ± ∞ ¹   147.08µ ± ∞ ¹       ~ (p=0.100 n=3) ²   84.14µ ± ∞ ¹        ~ (p=0.100 n=3) ²   398.72µ ± ∞ ¹         ~ (p=0.100 n=3) ²   524.50µ ± ∞ ¹         ~ (p=0.100 n=3) ²   387.97µ ± ∞ ¹         ~ (p=0.100 n=3) ²   65.76µ ± ∞ ¹        ~ (p=0.100 n=3) ²   378.39µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-8                75.57µ ± ∞ ¹    78.19µ ± ∞ ¹       ~ (p=0.100 n=3) ²   47.31µ ± ∞ ¹        ~ (p=0.100 n=3) ²   648.56µ ± ∞ ¹         ~ (p=0.100 n=3) ²   549.08µ ± ∞ ¹         ~ (p=0.100 n=3) ²   652.89µ ± ∞ ¹         ~ (p=0.100 n=3) ²   33.86µ ± ∞ ¹        ~ (p=0.100 n=3) ²   646.90µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-16               66.83µ ± ∞ ¹    68.95µ ± ∞ ¹       ~ (p=0.700 n=3) ²   39.91µ ± ∞ ¹        ~ (p=0.100 n=3) ²   664.61µ ± ∞ ¹         ~ (p=0.100 n=3) ²   492.97µ ± ∞ ¹         ~ (p=0.100 n=3) ²   673.88µ ± ∞ ¹         ~ (p=0.100 n=3) ²   29.00µ ± ∞ ¹        ~ (p=0.100 n=3) ²   667.03µ ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                162.6µ          166.1µ        +2.19%                   92.76µ        -42.94%                    510.8µ        +214.23%                    445.5µ        +174.07%                    506.7µ        +211.68%                   74.31µ        -54.29%                    497.0µ        +205.70%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │    bench_eatonphil_direct.txt     │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │    bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                         │           B/op           │     B/op       vs base           │   B/op     vs base                │     B/op       vs base           │     B/op       vs base           │     B/op       vs base           │     B/op       vs base           │    B/op      vs base           │
Query/GroupByParallel                     0.0 ± ∞ ¹    2364.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1849.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    7160.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1863.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1581.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   360.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-2                   0.0 ± ∞ ¹    2364.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1960.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    7157.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1914.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1579.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   360.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-4                   0.0 ± ∞ ¹    2364.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1963.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    7138.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1943.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1541.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   361.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-8                   0.0 ± ∞ ¹    2365.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1965.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    7163.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1980.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1571.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   362.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-16                  0.0 ± ∞ ¹    2368.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1973.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    7174.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1994.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1527.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   369.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ⁴   2.310Ki        ?                              +0.00%               ⁴   1.896Ki        ?                   6.991Ki        ?                   1.893Ki        ?                   1.523Ki        ?                   362.4        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                         │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/GroupByParallel                   0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   120.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   155.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   120.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   52.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-2                 0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   154.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   121.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-4                 0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   154.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   121.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-8                 0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   155.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-16                0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   155.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ⁴     118.0        ?                                +0.00%               ⁴     121.6        ?                     154.6        ?                     121.2        ?                    51.20        ?                   6.000        ?
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
                      │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt        │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt      │       bench_zombiezen_direct.txt       │
                      │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                │    sec/op      vs base                 │    sec/op     vs base                │    sec/op      vs base                 │
Query/JSON-12                      10.379m ± ∞ ¹   10.340m ± ∞ ¹       ~ (p=1.000 n=3) ²   7.279m ± ∞ ¹        ~ (p=0.100 n=3) ²    9.778m ± ∞ ¹        ~ (p=0.100 n=3) ²   8.977m ± ∞ ¹       ~ (p=0.100 n=3) ²   10.099m ± ∞ ¹        ~ (p=0.400 n=3) ²   7.474m ± ∞ ¹       ~ (p=0.100 n=3) ²   11.358m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-12               8.599m ± ∞ ¹    9.055m ± ∞ ¹       ~ (p=0.700 n=3) ²   8.830m ± ∞ ¹        ~ (p=0.700 n=3) ²   14.756m ± ∞ ¹        ~ (p=0.100 n=3) ²   9.650m ± ∞ ¹       ~ (p=0.700 n=3) ²   15.259m ± ∞ ¹        ~ (p=0.100 n=3) ²   9.744m ± ∞ ¹       ~ (p=0.700 n=3) ²   14.074m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                             9.447m          9.676m        +2.42%                   8.017m        -15.14%                    12.01m        +27.15%                   9.308m        -1.48%                    12.41m        +31.40%                   8.534m        -9.67%                    12.64m        +33.83%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                      │ bench_ncruces_direct.txt │           bench_ncruces_driver.txt            │      bench_eatonphil_direct.txt      │          bench_glebarez_driver.txt           │             bench_mattn_driver.txt             │           bench_modernc_driver.txt           │          bench_tailscale_driver.txt          │        bench_zombiezen_direct.txt        │
                      │           B/op           │      B/op        vs base                      │    B/op      vs base                 │      B/op        vs base                     │       B/op        vs base                      │      B/op        vs base                     │      B/op        vs base                     │     B/op       vs base                   │
Query/JSON-12                        1.000 ± ∞ ¹   64882.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=1.000 n=3) ²   56956.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   201032.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56968.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   32891.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   467.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-12                33.00 ± ∞ ¹    64899.00 ± ∞ ¹             ~ (p=0.100 n=3) ²   48.00 ± ∞ ¹        ~ (p=1.000 n=3) ²    57000.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    201262.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    56994.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    33013.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    457.00 ± ∞ ¹          ~ (p=0.100 n=3) ²
geomean                              5.745           63.37Ki        +1129498.60%                   6.928        +20.60%                     55.64Ki        +991759.59%                      196.4Ki        +3501419.25%                     55.65Ki        +991811.87%                     32.18Ki        +573519.71%                     462.0        +7941.92%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                      │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │   bench_zombiezen_direct.txt   │
                      │        allocs/op         │   allocs/op     vs base           │  allocs/op   vs base                │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │  allocs/op   vs base           │
Query/JSON-12                        0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-12                0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                        ⁴     4.019k        ?                                +0.00%               ⁴     4.022k        ?                     5.021k        ?                     4.022k        ?                     2.020k        ?                   6.481        ?
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
                      │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt        │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt      │       bench_zombiezen_direct.txt       │
                      │          sec/op          │    sec/op      vs base                │    sec/op      vs base                │    sec/op      vs base                 │    sec/op     vs base                │    sec/op      vs base                 │    sec/op     vs base                │    sec/op      vs base                 │
Query/JSONParallel                 10.391m ± ∞ ¹   10.744m ± ∞ ¹       ~ (p=0.100 n=3) ²    7.348m ± ∞ ¹       ~ (p=0.100 n=3) ²   10.340m ± ∞ ¹        ~ (p=0.400 n=3) ²   8.393m ± ∞ ¹       ~ (p=0.100 n=3) ²   10.254m ± ∞ ¹        ~ (p=0.200 n=3) ²   8.604m ± ∞ ¹       ~ (p=0.100 n=3) ²   10.945m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-2                7.402m ± ∞ ¹    7.733m ± ∞ ¹       ~ (p=0.200 n=3) ²    6.274m ± ∞ ¹       ~ (p=0.100 n=3) ²    8.130m ± ∞ ¹        ~ (p=0.100 n=3) ²   6.678m ± ∞ ¹       ~ (p=0.100 n=3) ²    7.648m ± ∞ ¹        ~ (p=0.700 n=3) ²   6.552m ± ∞ ¹       ~ (p=0.100 n=3) ²    9.910m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-4                6.210m ± ∞ ¹    6.862m ± ∞ ¹       ~ (p=0.100 n=3) ²    6.344m ± ∞ ¹       ~ (p=0.700 n=3) ²    7.498m ± ∞ ¹        ~ (p=0.100 n=3) ²   6.256m ± ∞ ¹       ~ (p=1.000 n=3) ²    6.651m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.843m ± ∞ ¹       ~ (p=0.700 n=3) ²    8.953m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-8                6.562m ± ∞ ¹    7.572m ± ∞ ¹       ~ (p=0.100 n=3) ²    7.501m ± ∞ ¹       ~ (p=0.100 n=3) ²   13.713m ± ∞ ¹        ~ (p=0.100 n=3) ²   6.706m ± ∞ ¹       ~ (p=0.100 n=3) ²   13.563m ± ∞ ¹        ~ (p=0.100 n=3) ²   6.664m ± ∞ ¹       ~ (p=0.200 n=3) ²   14.897m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-16               8.742m ± ∞ ¹    9.681m ± ∞ ¹       ~ (p=0.400 n=3) ²   10.164m ± ∞ ¹       ~ (p=0.200 n=3) ²   16.527m ± ∞ ¹        ~ (p=0.100 n=3) ²   8.779m ± ∞ ¹       ~ (p=0.700 n=3) ²   14.939m ± ∞ ¹        ~ (p=0.100 n=3) ²   8.875m ± ∞ ¹       ~ (p=0.700 n=3) ²   13.890m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                             7.719m          8.399m        +8.81%                    7.407m        -4.04%                    10.74m        +39.13%                   7.294m        -5.51%                    10.11m        +31.00%                   7.210m        -6.59%                    11.50m        +48.96%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                      │ bench_ncruces_direct.txt │           bench_ncruces_driver.txt            │    bench_eatonphil_direct.txt    │          bench_glebarez_driver.txt           │             bench_mattn_driver.txt             │           bench_modernc_driver.txt           │          bench_tailscale_driver.txt          │        bench_zombiezen_direct.txt        │
                      │           B/op           │      B/op        vs base                      │    B/op      vs base             │      B/op        vs base                     │       B/op        vs base                      │      B/op        vs base                     │      B/op        vs base                     │     B/op       vs base                   │
Query/JSONParallel                   1.000 ± ∞ ¹   64845.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     56886.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   201019.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56862.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   32891.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   365.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-2                 1.000 ± ∞ ¹   64849.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=1.000 n=3) ²     56949.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   201013.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56897.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   32900.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   371.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-4                 4.000 ± ∞ ¹   64846.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=0.400 n=3) ²     56957.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   201026.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56941.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   32895.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   381.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-8                 19.00 ± ∞ ¹    64854.00 ± ∞ ¹             ~ (p=0.100 n=3) ²   13.00 ± ∞ ¹  ~ (p=1.000 n=3) ²      56972.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    201049.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    56992.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    32955.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    427.00 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-16                99.00 ± ∞ ¹    64959.00 ± ∞ ¹             ~ (p=0.100 n=3) ²   87.00 ± ∞ ¹  ~ (p=1.000 n=3) ²      57154.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    201236.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    57061.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    32936.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    442.00 ± ∞ ¹          ~ (p=0.700 n=3) ²
geomean                              5.961           63.35Ki        +1088223.62%                                ?               ³ ⁴     55.65Ki        +955903.72%                      196.4Ki        +3373196.05%                     55.62Ki        +955350.55%                     32.14Ki        +552116.34%                     396.0        +6543.53%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean
⁴ ratios must be >0 to compute geomean

                      │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │   bench_zombiezen_direct.txt   │
                      │        allocs/op         │   allocs/op     vs base           │  allocs/op   vs base                │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │  allocs/op   vs base           │
Query/JSONParallel                   0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-2                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-4                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-8                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-16                0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ²   4023.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4023.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.400 n=3) ²
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
                         │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt       │         bench_mattn_driver.txt         │       bench_modernc_driver.txt        │      bench_tailscale_driver.txt      │      bench_zombiezen_direct.txt       │
                         │          sec/op          │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op     vs base                │    sec/op     vs base                 │
Query/OrderBy-12                       80.01m ± ∞ ¹   93.63m ± ∞ ¹        ~ (p=0.100 n=3) ²   49.54m ± ∞ ¹        ~ (p=0.100 n=3) ²   85.49m ± ∞ ¹        ~ (p=0.100 n=3) ²   114.61m ± ∞ ¹        ~ (p=0.100 n=3) ²   86.46m ± ∞ ¹        ~ (p=0.100 n=3) ²   64.33m ± ∞ ¹       ~ (p=0.100 n=3) ²   88.03m ± ∞ ¹        ~ (p=0.200 n=3) ²
Query/OrderByParallel-12               48.44m ± ∞ ¹   50.20m ± ∞ ¹        ~ (p=0.400 n=3) ²   49.44m ± ∞ ¹        ~ (p=0.700 n=3) ²   70.97m ± ∞ ¹        ~ (p=0.100 n=3) ²    90.99m ± ∞ ¹        ~ (p=0.100 n=3) ²   78.64m ± ∞ ¹        ~ (p=0.100 n=3) ²   49.73m ± ∞ ¹       ~ (p=0.700 n=3) ²   81.31m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                62.26m         68.56m        +10.12%                   49.49m        -20.51%                   77.89m        +25.12%                    102.1m        +64.03%                   82.46m        +32.45%                   56.56m        -9.15%                   84.60m        +35.89%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt             │      bench_eatonphil_direct.txt      │            bench_glebarez_driver.txt            │              bench_mattn_driver.txt              │            bench_modernc_driver.txt             │           bench_tailscale_driver.txt           │        bench_zombiezen_direct.txt        │
                         │           B/op           │       B/op         vs base                      │    B/op      vs base                 │       B/op         vs base                      │        B/op         vs base                      │       B/op         vs base                      │       B/op         vs base                     │      B/op       vs base                  │
Query/OrderBy-12                      374.000 ± ∞ ¹   6399614.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   8.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6397803.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   11999263.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   6397796.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   2798893.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3110.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-12                682.0 ± ∞ ¹     6399522.0 ± ∞ ¹             ~ (p=0.100 n=3) ²   519.0 ± ∞ ¹        ~ (p=0.400 n=3) ²     6398049.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     11999820.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     6398451.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     2799240.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     1551.0 ± ∞ ¹         ~ (p=0.700 n=3) ²
geomean                                 505.0             6.103Mi        +1267034.37%                   64.44        -87.24%                       6.102Mi        +1266709.24%                        11.44Mi        +2375846.53%                       6.102Mi        +1266748.35%                       2.669Mi        +554123.87%                    2.145Ki        +334.87%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt            │    bench_eatonphil_direct.txt    │           bench_glebarez_driver.txt            │             bench_mattn_driver.txt             │            bench_modernc_driver.txt            │           bench_tailscale_driver.txt           │       bench_zombiezen_direct.txt       │
                         │        allocs/op         │    allocs/op      vs base                      │  allocs/op   vs base             │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │  allocs/op    vs base                  │
Query/OrderBy-12                        8.000 ± ∞ ¹   349780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449782.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349773.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449781.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   36.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-12                9.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.300 n=3) ²     449784.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449784.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149767.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   18.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                 8.485             349.8k        +4122090.94%                                ?               ³ ⁴       449.8k        +5300643.49%                       349.8k        +4122043.80%                       449.8k        +5300637.60%                       149.8k        +1764915.13%                    25.46        +200.00%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean
⁴ ratios must be >0 to compute geomean
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
                         │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt       │         bench_mattn_driver.txt         │       bench_modernc_driver.txt        │      bench_tailscale_driver.txt      │      bench_zombiezen_direct.txt       │
                         │          sec/op          │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op     vs base                │    sec/op     vs base                 │
Query/OrderByParallel                  79.85m ± ∞ ¹   94.82m ± ∞ ¹        ~ (p=0.100 n=3) ²   49.05m ± ∞ ¹        ~ (p=0.100 n=3) ²   87.02m ± ∞ ¹        ~ (p=0.100 n=3) ²    91.73m ± ∞ ¹        ~ (p=0.100 n=3) ²   87.06m ± ∞ ¹        ~ (p=0.100 n=3) ²   64.14m ± ∞ ¹       ~ (p=0.100 n=3) ²   90.92m ± ∞ ¹        ~ (p=0.700 n=3) ²
Query/OrderByParallel-2                56.19m ± ∞ ¹   65.56m ± ∞ ¹        ~ (p=0.100 n=3) ²   39.14m ± ∞ ¹        ~ (p=0.100 n=3) ²   58.72m ± ∞ ¹        ~ (p=0.400 n=3) ²    68.32m ± ∞ ¹        ~ (p=0.100 n=3) ²   59.87m ± ∞ ¹        ~ (p=0.400 n=3) ²   48.94m ± ∞ ¹       ~ (p=0.200 n=3) ²   50.61m ± ∞ ¹        ~ (p=0.200 n=3) ²
Query/OrderByParallel-4                42.29m ± ∞ ¹   51.90m ± ∞ ¹        ~ (p=0.100 n=3) ²   36.25m ± ∞ ¹        ~ (p=0.100 n=3) ²   42.15m ± ∞ ¹        ~ (p=0.700 n=3) ²    59.23m ± ∞ ¹        ~ (p=0.100 n=3) ²   42.93m ± ∞ ¹        ~ (p=1.000 n=3) ²   38.80m ± ∞ ¹       ~ (p=0.100 n=3) ²   43.06m ± ∞ ¹        ~ (p=0.700 n=3) ²
Query/OrderByParallel-8                40.31m ± ∞ ¹   45.42m ± ∞ ¹        ~ (p=0.100 n=3) ²   40.43m ± ∞ ¹        ~ (p=0.700 n=3) ²   65.33m ± ∞ ¹        ~ (p=0.100 n=3) ²    80.36m ± ∞ ¹        ~ (p=0.100 n=3) ²   62.51m ± ∞ ¹        ~ (p=0.100 n=3) ²   39.56m ± ∞ ¹       ~ (p=0.200 n=3) ²   75.12m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-16               48.17m ± ∞ ¹   57.97m ± ∞ ¹        ~ (p=0.100 n=3) ²   51.03m ± ∞ ¹        ~ (p=0.700 n=3) ²   66.90m ± ∞ ¹        ~ (p=0.100 n=3) ²   102.96m ± ∞ ¹        ~ (p=0.100 n=3) ²   70.09m ± ∞ ¹        ~ (p=0.100 n=3) ²   50.97m ± ∞ ¹       ~ (p=0.400 n=3) ²   78.11m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                51.67m         61.07m        +18.18%                   42.80m        -17.18%                   62.34m        +20.64%                    78.97m        +52.82%                   62.85m        +21.62%                   47.65m        -7.80%                   65.03m        +25.84%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt             │      bench_eatonphil_direct.txt      │            bench_glebarez_driver.txt            │              bench_mattn_driver.txt              │            bench_modernc_driver.txt             │           bench_tailscale_driver.txt           │        bench_zombiezen_direct.txt        │
                         │           B/op           │       B/op         vs base                      │    B/op      vs base                 │       B/op         vs base                      │        B/op         vs base                      │       B/op         vs base                      │       B/op         vs base                     │      B/op       vs base                  │
Query/OrderByParallel                 373.000 ± ∞ ¹   6399282.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6397702.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   11999015.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   6397705.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   2798849.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   2964.000 ± ∞ ¹         ~ (p=0.200 n=3) ²
Query/OrderByParallel-2               394.000 ± ∞ ¹   6399271.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   9.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6397725.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   11999022.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   6397731.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   2798853.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1137.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-4                387.00 ± ∞ ¹    6399303.00 ± ∞ ¹             ~ (p=0.100 n=3) ²   11.00 ± ∞ ¹        ~ (p=0.100 n=3) ²    6397896.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    11999172.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    6397853.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    2798884.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    1170.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-8                385.00 ± ∞ ¹    6399336.00 ± ∞ ¹             ~ (p=0.100 n=3) ²   40.00 ± ∞ ¹        ~ (p=0.100 n=3) ²    6398095.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    11999725.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    6398010.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    2798925.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    1313.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-16               1353.0 ± ∞ ¹     6399886.0 ± ∞ ¹             ~ (p=0.100 n=3) ²   310.0 ± ∞ ¹        ~ (p=0.400 n=3) ²     6398418.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     12001650.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     6398731.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     2799639.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     1847.0 ± ∞ ¹         ~ (p=0.700 n=3) ²
geomean                                 494.7             6.103Mi        +1293516.08%                   22.78        -95.39%                       6.102Mi        +1293223.29%                        11.44Mi        +2425594.40%                       6.102Mi        +1293231.14%                       2.669Mi        +565712.64%                    1.534Ki        +217.52%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt            │    bench_eatonphil_direct.txt    │           bench_glebarez_driver.txt            │             bench_mattn_driver.txt             │            bench_modernc_driver.txt            │           bench_tailscale_driver.txt           │       bench_zombiezen_direct.txt       │
                         │        allocs/op         │    allocs/op      vs base                      │  allocs/op   vs base             │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │  allocs/op    vs base                  │
Query/OrderByParallel                   8.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349771.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   35.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-2                 8.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349771.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   16.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-4                 8.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449782.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349772.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449781.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   16.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-8                 8.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449784.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449783.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   17.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-16               10.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449784.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349798.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449787.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149769.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   19.000 ± ∞ ¹         ~ (p=0.600 n=3) ²
geomean                                 8.365             349.8k        +4181295.48%                                ?               ³ ⁴       449.8k        +5376776.75%                       349.8k        +4181288.31%                       449.8k        +5376779.14%                       149.8k        +1790270.78%                    19.60        +134.33%
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
                              │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt       │        bench_mattn_driver.txt         │       bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │      bench_zombiezen_direct.txt      │
                              │          sec/op          │    sec/op     vs base                │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                │
Query/RecursiveCTE-12                       6.354m ± ∞ ¹   6.285m ± ∞ ¹       ~ (p=0.200 n=3) ²   2.562m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.392m ± ∞ ¹        ~ (p=0.100 n=3) ²   2.484m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.400m ± ∞ ¹        ~ (p=0.100 n=3) ²   2.736m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.461m ± ∞ ¹       ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-12               689.8µ ± ∞ ¹   686.7µ ± ∞ ¹       ~ (p=0.700 n=3) ²   278.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   650.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²   275.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   649.2µ ± ∞ ¹        ~ (p=0.700 n=3) ²   459.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   706.8µ ± ∞ ¹       ~ (p=0.700 n=3) ²
geomean                                     2.094m         2.077m        -0.77%                   844.2µ        -59.68%                   1.873m        -10.53%                   826.9µ        -60.50%                   1.872m        -10.57%                   1.122m        -46.42%                   1.965m        -6.16%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt           │    bench_eatonphil_direct.txt    │          bench_glebarez_driver.txt          │           bench_mattn_driver.txt            │          bench_modernc_driver.txt           │         bench_tailscale_driver.txt         │        bench_zombiezen_direct.txt         │
                              │           B/op           │      B/op       vs base                     │    B/op      vs base             │      B/op       vs base                     │      B/op       vs base                     │      B/op       vs base                     │      B/op       vs base                    │     B/op       vs base                    │
Query/RecursiveCTE-12                        1.000 ± ∞ ¹   2485.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.300 n=3) ²     2355.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6863.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   2291.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1514.000 ± ∞ ¹           ~ (p=0.100 n=3) ²   391.000 ± ∞ ¹           ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-12                3.000 ± ∞ ¹   2482.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=1.000 n=3) ²     2357.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6855.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   2277.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1496.000 ± ∞ ¹           ~ (p=0.100 n=3) ²   367.000 ± ∞ ¹           ~ (p=0.100 n=3) ²
geomean                                      1.732          2.425Ki        +143284.91%                                ?               ³ ⁴    2.301Ki        +135923.71%                    6.698Ki        +395904.48%                    2.230Ki        +131766.18%                    1.470Ki        +86789.66%                     378.8        +21770.60%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean
⁴ ratios must be >0 to compute geomean

                              │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                              │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/RecursiveCTE-12                        0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   143.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
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
Query/RecursiveCTEParallel                  6.368m ± ∞ ¹    6.394m ± ∞ ¹       ~ (p=1.000 n=3) ²   2.567m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.402m ± ∞ ¹        ~ (p=0.100 n=3) ²   2.480m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.392m ± ∞ ¹        ~ (p=0.100 n=3) ²   2.492m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.323m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                3.167m ± ∞ ¹    3.145m ± ∞ ¹       ~ (p=0.700 n=3) ²   1.287m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.729m ± ∞ ¹        ~ (p=0.100 n=3) ²   1.265m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.721m ± ∞ ¹        ~ (p=0.100 n=3) ²   1.248m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.779m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4               1592.0µ ± ∞ ¹   1597.5µ ± ∞ ¹       ~ (p=0.700 n=3) ²   645.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1391.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   633.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1385.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   623.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1382.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                862.8µ ± ∞ ¹    846.5µ ± ∞ ¹       ~ (p=0.700 n=3) ²   344.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    815.7µ ± ∞ ¹        ~ (p=0.700 n=3) ²   363.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    777.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   335.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    776.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16               756.8µ ± ∞ ¹    917.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²   307.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    705.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   305.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    696.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   289.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    694.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                     1.838m          1.903m        +3.55%                   742.3µ        -59.61%                    1.638m        -10.84%                   739.2µ        -59.78%                    1.616m        -12.08%                   715.8µ        -61.05%                    1.616m        -12.07%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │      bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                              │           B/op           │      B/op       vs base           │    B/op      vs base                 │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │     B/op       vs base           │
Query/RecursiveCTEParallel                     0.0 ± ∞ ¹     2482.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹        ~ (p=1.000 n=3) ³     2263.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     6857.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     2258.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1505.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     362.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                   0.0 ± ∞ ¹     2481.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹        ~ (p=1.000 n=3) ³     2274.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     6823.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     2258.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1503.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     369.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4                   0.0 ± ∞ ¹     2481.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹        ~ (p=1.000 n=3) ²     2349.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     6838.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     2264.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1502.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     367.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                 2.000 ± ∞ ¹   2481.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.300 n=3) ²   2355.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6854.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2276.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1504.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   362.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16                4.000 ± ∞ ¹   2487.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   2368.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6861.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2282.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1507.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   367.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                ⁴    2.424Ki        ?                                -24.21%               ⁴    2.267Ki        ?                    6.686Ki        ?                    2.214Ki        ?                    1.469Ki        ?                     365.4        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                              │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                              │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/RecursiveCTEParallel                   0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   143.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   48.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   48.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16                0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   143.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                ⁴     110.0        ?                                +0.00%               ⁴     112.4        ?                     142.4        ?                     112.0        ?                    48.60        ?                   6.000        ?
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
Query/Window-12                      2180.0µ ± ∞ ¹   2471.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   937.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1924.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   2165.5µ ± ∞ ¹         ~ (p=0.700 n=3) ²   1922.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1053.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1667.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-12               250.9µ ± ∞ ¹    340.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   122.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    912.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1030.2µ ± ∞ ¹         ~ (p=0.100 n=3) ²    903.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    133.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    849.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                               739.5µ          916.7µ        +23.95%                   339.3µ        -54.11%                    1.325m        +79.23%                    1.494m        +101.97%                    1.318m        +78.22%                    374.5µ        -49.36%                    1.190m        +60.94%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                        │ bench_ncruces_direct.txt │      bench_ncruces_driver.txt      │   bench_eatonphil_direct.txt   │     bench_glebarez_driver.txt      │       bench_mattn_driver.txt        │      bench_modernc_driver.txt      │     bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                        │           B/op           │      B/op        vs base           │    B/op      vs base           │      B/op        vs base           │       B/op        vs base           │      B/op        vs base           │      B/op        vs base           │     B/op       vs base           │
Query/Window-12                          0.0 ± ∞ ¹     62792.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ²     54924.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198935.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54899.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30805.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     364.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-12                0.000 ± ∞ ¹   62761.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.700 n=3) ²   54886.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   198942.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54909.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   30777.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   369.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                          ³     61.31Ki        ?                                ?               ³     53.62Ki        ?                      194.3Ki        ?                     53.62Ki        ?                     30.07Ki        ?                     366.5        ?
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
Query/WindowParallel                 2133.7µ ± ∞ ¹   2416.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   952.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1922.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1673.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1911.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1066.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1647.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-2               1089.5µ ± ∞ ¹   1251.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   486.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    982.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    880.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    990.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    542.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    856.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-4                541.1µ ± ∞ ¹    671.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   263.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    585.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    525.4µ ± ∞ ¹        ~ (p=0.700 n=3) ²    568.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    282.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    507.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-8                280.0µ ± ∞ ¹    358.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   142.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    877.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1023.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    867.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    151.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    838.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-16               239.8µ ± ∞ ¹    321.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   124.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    924.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1070.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    909.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    146.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    854.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                               610.0µ          748.1µ        +22.64%                   293.3µ        -51.91%                    978.4µ        +60.40%                    967.5µ        +58.62%                    967.7µ        +58.65%                    325.1µ        -46.70%                    875.1µ        +43.46%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                        │ bench_ncruces_direct.txt │      bench_ncruces_driver.txt      │      bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt      │       bench_mattn_driver.txt        │      bench_modernc_driver.txt      │     bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                        │           B/op           │      B/op        vs base           │    B/op      vs base                 │      B/op        vs base           │       B/op        vs base           │      B/op        vs base           │      B/op        vs base           │     B/op       vs base           │
Query/WindowParallel                     0.0 ± ∞ ¹     62762.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹        ~ (p=1.000 n=3) ³     54776.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198936.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54790.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30792.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     362.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-2                   0.0 ± ∞ ¹     62761.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹        ~ (p=1.000 n=3) ³     54816.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198906.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54816.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30793.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     363.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-4                   0.0 ± ∞ ¹     62761.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹        ~ (p=1.000 n=3) ³     54881.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198870.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54889.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30793.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     361.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-8                   0.0 ± ∞ ¹     62762.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹        ~ (p=1.000 n=3) ³     54884.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198917.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54902.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30791.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     363.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-16                1.000 ± ∞ ¹   62767.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹        ~ (p=0.500 n=3) ²   54906.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   198972.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54909.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   30779.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   371.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                          ⁴     61.29Ki        ?                                +14.87%               ⁴     53.57Ki        ?                      194.3Ki        ?                     53.58Ki        ?                     30.07Ki        ?                     364.0        ?
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
3.46.1
```
<!--END_GREP-->

## modernc

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
MAX_EXPR_DEPTH=0
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
USE_ALLOCA
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
