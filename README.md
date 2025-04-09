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
Baseline/Conn-12                                     151.2n ± ∞ ¹    222.4n ± ∞ ¹        ~ (p=0.100 n=3) ²   169.5n ± ∞ ¹        ~ (p=0.100 n=3) ²    226.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    221.4n ± ∞ ¹        ~ (p=0.100 n=3) ²    229.5n ± ∞ ¹        ~ (p=0.100 n=3) ²    224.5n ± ∞ ¹       ~ (p=0.100 n=3) ²   1091.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-12                             452.3n ± ∞ ¹    541.6n ± ∞ ¹        ~ (p=0.100 n=3) ²   452.8n ± ∞ ¹        ~ (p=1.000 n=3) ²    543.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    540.1n ± ∞ ¹        ~ (p=0.100 n=3) ²    547.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    542.3n ± ∞ ¹       ~ (p=0.100 n=3) ²   1368.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1-12                                 1801.0n ± ∞ ¹   1943.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   878.2n ± ∞ ¹        ~ (p=0.100 n=3) ²   2024.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2565.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2002.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1114.0n ± ∞ ¹       ~ (p=0.100 n=3) ²   3191.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-12                         1001.0n ± ∞ ¹   1163.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   763.8n ± ∞ ¹        ~ (p=0.100 n=3) ²   2258.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1505.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2295.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    929.2n ± ∞ ¹       ~ (p=0.100 n=3) ²   2649.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePrepared-12                       423.2n ± ∞ ¹    529.8n ± ∞ ¹        ~ (p=0.100 n=3) ²   436.8n ± ∞ ¹        ~ (p=0.100 n=3) ²   1982.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1305.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1985.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    410.2n ± ∞ ¹       ~ (p=0.100 n=3) ²   1524.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-12               734.0n ± ∞ ¹    928.5n ± ∞ ¹        ~ (p=0.100 n=3) ²   675.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   1788.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1099.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1781.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    770.0n ± ∞ ¹       ~ (p=0.100 n=3) ²   1624.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                              580.6n          715.3n        +23.20%                   497.6n        -14.29%                    1.122µ        +93.21%                    933.6n        +60.80%                    1.126µ        +93.97%                    584.3n        +0.65%                    1.775µ        +205.65%
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
Baseline/ConnParallel                                149.7n ± ∞ ¹    233.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   210.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    247.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    232.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    248.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    234.6n ± ∞ ¹       ~ (p=0.100 n=3) ²   1409.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-2                              155.6n ± ∞ ¹    232.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   168.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    224.9n ± ∞ ¹        ~ (p=0.100 n=3) ²    214.5n ± ∞ ¹        ~ (p=0.100 n=3) ²    229.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    221.8n ± ∞ ¹       ~ (p=0.100 n=3) ²    711.2n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-4                              247.5n ± ∞ ¹    309.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   267.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    265.1n ± ∞ ¹        ~ (p=0.100 n=3) ²    260.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    267.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    265.4n ± ∞ ¹       ~ (p=0.100 n=3) ²    576.6n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-8                              452.5n ± ∞ ¹    509.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   439.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    511.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    512.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    520.9n ± ∞ ¹        ~ (p=0.100 n=3) ²    523.1n ± ∞ ¹       ~ (p=0.100 n=3) ²   1265.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-16                             479.1n ± ∞ ¹    535.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   462.1n ± ∞ ¹        ~ (p=0.100 n=3) ²    551.3n ± ∞ ¹        ~ (p=0.100 n=3) ²    546.3n ± ∞ ¹        ~ (p=0.100 n=3) ²    554.3n ± ∞ ¹        ~ (p=0.100 n=3) ²    548.6n ± ∞ ¹       ~ (p=0.100 n=3) ²   1303.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel                            1792.0n ± ∞ ¹   1948.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   890.5n ± ∞ ¹        ~ (p=0.100 n=3) ²   1865.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1625.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2045.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1145.0n ± ∞ ¹       ~ (p=0.100 n=3) ²   2833.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-2                           943.3n ± ∞ ¹   1038.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   520.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   1092.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    999.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   1097.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    650.3n ± ∞ ¹       ~ (p=0.100 n=3) ²   1534.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-4                           514.7n ± ∞ ¹    624.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   317.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1172.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    620.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   1017.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    394.7n ± ∞ ¹       ~ (p=0.100 n=3) ²   1108.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-8                           879.0n ± ∞ ¹   1108.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   733.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   2200.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1452.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2207.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    890.2n ± ∞ ¹       ~ (p=0.100 n=3) ²   2475.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-16                          925.2n ± ∞ ¹   1163.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   770.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   2270.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1531.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2264.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    935.5n ± ∞ ¹       ~ (p=0.100 n=3) ²   2683.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel                  425.2n ± ∞ ¹    523.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   434.4n ± ∞ ¹        ~ (p=0.700 n=3) ²   1843.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    847.6n ± ∞ ¹        ~ (p=0.100 n=3) ²   2015.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    421.4n ± ∞ ¹       ~ (p=0.100 n=3) ²   1628.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-2                473.8n ± ∞ ¹    457.1n ± ∞ ¹        ~ (p=1.000 n=3) ²   310.2n ± ∞ ¹        ~ (p=0.100 n=3) ²   1045.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    609.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   1045.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    344.5n ± ∞ ¹       ~ (p=0.100 n=3) ²    848.1n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-4                294.5n ± ∞ ¹    370.8n ± ∞ ¹        ~ (p=0.100 n=3) ²   226.3n ± ∞ ¹        ~ (p=0.100 n=3) ²    767.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    487.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    759.1n ± ∞ ¹        ~ (p=0.100 n=3) ²    280.7n ± ∞ ¹       ~ (p=0.100 n=3) ²    640.8n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-8                686.7n ± ∞ ¹    876.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   647.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1731.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1082.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1740.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    735.2n ± ∞ ¹       ~ (p=0.100 n=3) ²   1403.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-16               719.0n ± ∞ ¹    926.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   667.5n ± ∞ ¹        ~ (p=0.100 n=3) ²   1793.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1102.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1786.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    773.8n ± ∞ ¹       ~ (p=0.100 n=3) ²   1611.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                              494.6n          605.4n        +22.42%                   418.4n        -15.41%                    906.1n        +83.22%                    669.8n        +35.43%                    911.4n        +84.28%                    489.9n        -0.94%                    1.315µ        +165.92%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                       │ bench_ncruces_direct.txt │    bench_ncruces_driver.txt    │      bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt    │     bench_mattn_driver.txt      │    bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt    │
                                       │           B/op           │    B/op      vs base           │    B/op      vs base                 │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │
Baseline/ConnParallel                                  0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   395.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-2                                0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-4                                0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-8                                0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-16                               0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel                              48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   723.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-2                            48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-4                            48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-8                            48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-16                           48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel                    0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   371.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-2                  0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-4                  0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-8                  0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-16                 0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                         ⁴   46.15        ?                                +25.99%               ⁴    159.6        ?                    164.2        ?                    159.6        ?                    90.34        ?                    449.3        ?
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
                              │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt      │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │      bench_zombiezen_direct.txt       │
                              │          sec/op          │    sec/op      vs base                │    sec/op     vs base                │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                │
Populate/PopulateDB-12                       2.492 ± ∞ ¹     2.472 ± ∞ ¹       ~ (p=0.700 n=3) ²    2.019 ± ∞ ¹       ~ (p=0.100 n=3) ²     2.744 ± ∞ ¹        ~ (p=0.100 n=3) ²     2.225 ± ∞ ¹       ~ (p=0.100 n=3) ²     2.714 ± ∞ ¹        ~ (p=0.100 n=3) ²     2.068 ± ∞ ¹       ~ (p=0.100 n=3) ²     2.390 ± ∞ ¹       ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTx-12               1137.1m ± ∞ ¹   1196.8m ± ∞ ¹       ~ (p=0.400 n=3) ²   951.9m ± ∞ ¹       ~ (p=0.100 n=3) ²   1400.8m ± ∞ ¹        ~ (p=0.100 n=3) ²   1126.0m ± ∞ ¹       ~ (p=1.000 n=3) ²   1414.1m ± ∞ ¹        ~ (p=0.100 n=3) ²   1015.4m ± ∞ ¹       ~ (p=0.200 n=3) ²   1092.2m ± ∞ ¹       ~ (p=1.000 n=3) ²
Populate/PopulateDBWithTxs-12                1.061 ± ∞ ¹     1.139 ± ∞ ¹       ~ (p=0.400 n=3) ²    1.143 ± ∞ ¹       ~ (p=0.700 n=3) ²     1.416 ± ∞ ¹        ~ (p=0.100 n=3) ²     1.224 ± ∞ ¹       ~ (p=0.100 n=3) ²     1.359 ± ∞ ¹        ~ (p=0.100 n=3) ²     1.090 ± ∞ ¹       ~ (p=0.700 n=3) ²     1.184 ± ∞ ¹       ~ (p=0.200 n=3) ²
geomean                                      1.443           1.499        +3.89%                    1.300        -9.92%                     1.759        +21.89%                     1.453        +0.66%                     1.734        +20.16%                     1.318        -8.70%                     1.457        +0.94%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │         bench_ncruces_driver.txt          │       bench_eatonphil_direct.txt        │         bench_glebarez_driver.txt         │           bench_mattn_driver.txt            │         bench_modernc_driver.txt          │        bench_tailscale_driver.txt         │       bench_zombiezen_direct.txt        │
                              │           B/op           │      B/op       vs base                   │     B/op       vs base                  │      B/op       vs base                   │      B/op        vs base                    │      B/op       vs base                   │      B/op       vs base                   │     B/op       vs base                  │
Populate/PopulateDB-12                     2.287Mi ± ∞ ¹   18.967Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   5.770Mi ± ∞ ¹         ~ (p=0.100 n=3) ²   31.479Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   273.293Mi ± ∞ ¹           ~ (p=0.100 n=3) ²   31.500Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   18.962Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   6.165Mi ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTx-12               2.276Mi ± ∞ ¹   32.290Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   5.771Mi ± ∞ ¹         ~ (p=0.100 n=3) ²   44.547Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   286.377Mi ± ∞ ¹           ~ (p=0.100 n=3) ²   44.867Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   32.555Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   6.166Mi ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTxs-12              2.276Mi ± ∞ ¹   31.237Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   5.770Mi ± ∞ ¹         ~ (p=0.100 n=3) ²   43.672Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   285.997Mi ± ∞ ¹           ~ (p=0.100 n=3) ²   43.851Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   31.258Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   6.213Mi ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                    2.279Mi          26.75Mi        +1073.31%                   5.770Mi        +153.14%                    39.42Mi        +1629.22%                     281.8Mi        +12263.59%                    39.57Mi        +1636.10%                    26.82Mi        +1076.69%                   6.181Mi        +171.18%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │         bench_mattn_driver.txt          │        bench_modernc_driver.txt         │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                              │        allocs/op         │  allocs/op    vs base                  │  allocs/op    vs base                  │   allocs/op    vs base                  │   allocs/op    vs base                  │   allocs/op    vs base                  │  allocs/op    vs base                  │  allocs/op    vs base                  │
Populate/PopulateDB-12                      140.0k ± ∞ ¹   598.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   394.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1209.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   1057.5k ± ∞ ¹         ~ (p=0.100 n=3) ²   1209.4k ± ∞ ¹         ~ (p=0.100 n=3) ²   598.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   445.1k ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTx-12                140.0k ± ∞ ¹   751.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   394.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1362.4k ± ∞ ¹         ~ (p=0.100 n=3) ²   1210.5k ± ∞ ¹         ~ (p=0.100 n=3) ²   1362.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   751.2k ± ∞ ¹         ~ (p=0.100 n=3) ²   445.1k ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTxs-12               140.0k ± ∞ ¹   765.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   394.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1376.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   1238.9k ± ∞ ¹         ~ (p=0.100 n=3) ²   1380.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   767.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   448.1k ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                     140.0k         700.5k        +400.32%                   394.0k        +181.41%                    1.314M        +838.32%                    1.166M        +732.87%                    1.315M        +839.22%                   701.1k        +400.76%                   446.1k        +218.59%
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
ReadWrite/ReadPost-12                                                            13.917µ ± ∞ ¹   21.806µ ± ∞ ¹        ~ (p=0.100 n=3) ²   9.239µ ± ∞ ¹        ~ (p=0.100 n=3) ²   29.492µ ± ∞ ¹        ~ (p=0.100 n=3) ²   17.797µ ± ∞ ¹        ~ (p=0.100 n=3) ²   29.368µ ± ∞ ¹        ~ (p=0.100 n=3) ²   18.278µ ± ∞ ¹       ~ (p=0.100 n=3) ²   17.867µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadPostWithTx-12                                                      15.997µ ± ∞ ¹   37.586µ ± ∞ ¹        ~ (p=0.100 n=3) ²   9.772µ ± ∞ ¹        ~ (p=0.100 n=3) ²   33.718µ ± ∞ ¹        ~ (p=0.100 n=3) ²   22.049µ ± ∞ ¹        ~ (p=0.100 n=3) ²   48.190µ ± ∞ ¹        ~ (p=0.100 n=3) ²   25.993µ ± ∞ ¹       ~ (p=0.100 n=3) ²   18.372µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndComments-12                                                  98.71µ ± ∞ ¹   113.58µ ± ∞ ¹        ~ (p=0.100 n=3) ²   71.07µ ± ∞ ¹        ~ (p=0.100 n=3) ²   182.65µ ± ∞ ¹        ~ (p=0.100 n=3) ²   177.86µ ± ∞ ¹        ~ (p=0.100 n=3) ²   279.86µ ± ∞ ¹        ~ (p=0.100 n=3) ²   103.47µ ± ∞ ¹       ~ (p=0.100 n=3) ²    92.02µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndCommentsWithTx-12                                            98.42µ ± ∞ ¹   141.32µ ± ∞ ¹        ~ (p=0.100 n=3) ²   70.78µ ± ∞ ¹        ~ (p=0.100 n=3) ²   185.53µ ± ∞ ¹        ~ (p=0.100 n=3) ²   180.65µ ± ∞ ¹        ~ (p=0.100 n=3) ²   203.43µ ± ∞ ¹        ~ (p=0.100 n=3) ²   117.86µ ± ∞ ¹       ~ (p=0.100 n=3) ²    91.03µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePost-12                                                            130.6µ ± ∞ ¹    128.4µ ± ∞ ¹        ~ (p=1.000 n=3) ²   116.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    140.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    131.4µ ± ∞ ¹        ~ (p=0.700 n=3) ²    156.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    116.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²    161.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePostWithTx-12                                                      134.3µ ± ∞ ¹    141.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   117.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    149.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    131.6µ ± ∞ ¹        ~ (p=0.400 n=3) ²    150.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    123.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²    164.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndComments-12                                                 2.531m ± ∞ ¹    2.598m ± ∞ ¹        ~ (p=0.200 n=3) ²   2.095m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.823m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.338m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.908m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.141m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.564m ± ∞ ¹        ~ (p=0.400 n=3) ²
ReadWrite/WritePostAndCommentsWithTx-12                                           930.9µ ± ∞ ¹    996.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   780.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1231.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1004.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1275.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    830.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²   1104.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=10-12                             354.6µ ± ∞ ¹    375.1µ ± ∞ ¹        ~ (p=0.700 n=3) ²   285.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    461.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    402.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    458.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    302.6µ ± ∞ ¹       ~ (p=0.100 n=3) ²    403.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=90-12                             2.334m ± ∞ ¹    2.342m ± ∞ ¹        ~ (p=0.700 n=3) ²   1.925m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.548m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.082m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.602m ± ∞ ¹        ~ (p=0.100 n=3) ²    1.949m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.337m ± ∞ ¹        ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-12                     387.6µ ± ∞ ¹    393.0µ ± ∞ ¹        ~ (p=1.000 n=3) ²   287.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    412.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    347.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    408.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    314.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²    392.2µ ± ∞ ¹        ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-12                     3.116m ± ∞ ¹    3.086m ± ∞ ¹        ~ (p=1.000 n=3) ²   2.257m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.756m ± ∞ ¹        ~ (p=0.400 n=3) ²    2.335m ± ∞ ¹        ~ (p=0.100 n=3) ²    3.502m ± ∞ ¹        ~ (p=0.400 n=3) ²    2.174m ± ∞ ¹       ~ (p=0.100 n=3) ²    4.039m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=10-12                       184.0µ ± ∞ ¹    223.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   144.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    302.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    267.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    311.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    178.2µ ± ∞ ¹       ~ (p=0.700 n=3) ²    249.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=90-12                       862.6µ ± ∞ ¹    950.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²   718.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1135.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    932.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1166.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    758.9µ ± ∞ ¹       ~ (p=0.100 n=3) ²   1042.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-12               124.1µ ± ∞ ¹    121.9µ ± ∞ ¹        ~ (p=1.000 n=3) ²   114.0µ ± ∞ ¹        ~ (p=0.400 n=3) ²    211.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    170.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    227.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    126.3µ ± ∞ ¹       ~ (p=0.700 n=3) ²    187.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-12               863.8µ ± ∞ ¹    924.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   813.5µ ± ∞ ¹        ~ (p=0.700 n=3) ²   4104.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1125.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   4587.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    855.9µ ± ∞ ¹       ~ (p=0.700 n=3) ²   2433.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                                                           281.1µ          325.3µ        +15.73%                   222.1µ        -21.00%                    421.1µ        +49.79%                    326.5µ        +16.13%                    463.7µ        +64.96%                    271.5µ        -3.42%                    343.0µ        +22.00%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                                                    │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │    bench_eatonphil_direct.txt    │     bench_glebarez_driver.txt     │       bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                                                                    │           B/op           │      B/op       vs base           │     B/op       vs base           │      B/op       vs base           │      B/op        vs base           │      B/op       vs base           │      B/op       vs base           │     B/op       vs base           │
ReadWrite/ReadPost-12                                                            40.17Ki ± ∞ ¹    41.23Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.18Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    81.34Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     41.39Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    81.28Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    41.19Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.55Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostWithTx-12                                                      40.17Ki ± ∞ ¹    42.29Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.18Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    82.21Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     42.60Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    82.33Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    42.13Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.60Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndComments-12                                                 250.2Ki ± ∞ ¹    260.3Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    501.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     263.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    501.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    257.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndCommentsWithTx-12                                           250.2Ki ± ∞ ¹    261.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    503.1Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     265.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    503.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    258.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePost-12                                                              0.00 ± ∞ ¹     240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     496.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    42288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     496.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    409.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostWithTx-12                                                        0.00 ± ∞ ¹     857.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    1106.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    43252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    1289.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     880.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    457.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndComments-12                                                0.000Ki ± ∞ ¹   13.906Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   2.781Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.663Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   308.876Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.665Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   13.906Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   3.528Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndCommentsWithTx-12                                          0.000Ki ± ∞ ¹   26.433Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   2.781Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   39.199Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   321.739Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   39.377Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.462Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   3.573Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=10-12                            224.6Ki ± ∞ ¹    234.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   224.5Ki ± ∞ ¹  ~ (p=1.000 n=3) ²    452.8Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     268.5Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    454.1Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    232.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   226.4Ki ± ∞ ¹  ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=90-12                            23.85Ki ± ∞ ¹    38.75Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   24.56Ki ± ∞ ¹  ~ (p=1.000 n=3) ²    81.00Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    303.71Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    77.63Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    37.26Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   27.53Ki ± ∞ ¹  ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-12                    225.5Ki ± ∞ ¹    236.0Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.8Ki ± ∞ ¹  ~ (p=1.000 n=3) ²    453.3Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     268.4Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    453.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    232.4Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   224.8Ki ± ∞ ¹  ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-12                    21.18Ki ± ∞ ¹    42.11Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   25.13Ki ± ∞ ¹  ~ (p=0.700 n=3) ²    70.87Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    303.66Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    68.68Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    35.75Ki ± ∞ ¹  ~ (p=0.400 n=3) ²   28.37Ki ± ∞ ¹  ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=10-12                      225.0Ki ± ∞ ¹    237.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.6Ki ± ∞ ¹  ~ (p=0.400 n=3) ²    456.5Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     271.5Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    455.5Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    236.0Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   226.7Ki ± ∞ ¹  ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=90-12                      24.88Ki ± ∞ ¹    48.09Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   28.50Ki ± ∞ ¹  ~ (p=0.400 n=3) ²    85.63Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    315.97Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    85.53Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    50.04Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   29.14Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-12              224.4Ki ± ∞ ¹    238.8Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.5Ki ± ∞ ¹  ~ (p=0.400 n=3) ²    456.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     271.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    455.4Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    235.0Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.9Ki ± ∞ ¹  ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-12              25.37Ki ± ∞ ¹    51.54Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   27.04Ki ± ∞ ¹  ~ (p=0.200 n=3) ²    80.02Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    316.50Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    71.65Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    50.95Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   24.91Ki ± ∞ ¹  ~ (p=0.700 n=3) ²
geomean                                                                                      ³    42.40Ki        ?                   21.51Ki        ?                    77.08Ki        ?                     178.4Ki        ?                    76.95Ki        ?                    41.81Ki        ?                   29.58Ki        ?
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
ReadWrite/ReadOrWritePostAndComments/write_rate=10-12                              235.0 ± ∞ ¹    735.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    697.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=90-12                              25.00 ± ∞ ¹   443.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   799.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   421.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-12                      236.0 ± ∞ ¹    736.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    697.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-12                      22.00 ± ∞ ¹   448.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   799.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   420.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=10-12                        235.0 ± ∞ ¹    778.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1011.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    748.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    274.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=90-12                        26.00 ± ∞ ¹   595.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1120.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   967.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1124.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-12                235.0 ± ∞ ¹    778.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1011.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    274.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-12                26.00 ± ∞ ¹   598.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1122.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   970.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1128.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                                                      ³    273.6        ?                    85.19        ?                     430.5        ?                    368.2        ?                     436.0        ?                    237.3        ?                    122.4        ?
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
                                                                    │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │         bench_mattn_driver.txt         │        bench_modernc_driver.txt         │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt        │
                                                                    │          sec/op          │    sec/op      vs base                │    sec/op      vs base                 │     sec/op      vs base                 │    sec/op      vs base                 │     sec/op      vs base                 │    sec/op      vs base                │     sec/op      vs base                 │
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10                        351.0µ ± ∞ ¹    332.4µ ± ∞ ¹       ~ (p=0.200 n=3) ²    277.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²     440.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    353.1µ ± ∞ ¹        ~ (p=0.400 n=3) ²     457.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    296.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²     441.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-2                      310.9µ ± ∞ ¹    330.1µ ± ∞ ¹       ~ (p=0.400 n=3) ²    258.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²     382.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    324.6µ ± ∞ ¹        ~ (p=0.700 n=3) ²     392.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    286.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²     348.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-4                      348.4µ ± ∞ ¹    340.5µ ± ∞ ¹       ~ (p=1.000 n=3) ²    254.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²     358.9µ ± ∞ ¹        ~ (p=0.400 n=3) ²    308.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²     383.3µ ± ∞ ¹        ~ (p=0.400 n=3) ²    257.1µ ± ∞ ¹       ~ (p=0.100 n=3) ²     363.0µ ± ∞ ¹        ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-8                      354.2µ ± ∞ ¹    360.3µ ± ∞ ¹       ~ (p=0.400 n=3) ²    269.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²     406.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    327.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²     397.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    282.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²     357.2µ ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-16                     399.5µ ± ∞ ¹    423.2µ ± ∞ ¹       ~ (p=0.700 n=3) ²    297.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²     424.1µ ± ∞ ¹        ~ (p=0.700 n=3) ²    359.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²     436.7µ ± ∞ ¹        ~ (p=0.200 n=3) ²    323.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²     372.8µ ± ∞ ¹        ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90                        2.356m ± ∞ ¹    2.346m ± ∞ ¹       ~ (p=1.000 n=3) ²    1.908m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.406m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.085m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.440m ± ∞ ¹        ~ (p=0.100 n=3) ²    1.910m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.327m ± ∞ ¹        ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-2                      2.673m ± ∞ ¹    2.531m ± ∞ ¹       ~ (p=0.200 n=3) ²    1.890m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.491m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.187m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.591m ± ∞ ¹        ~ (p=0.200 n=3) ²    1.924m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.334m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-4                      2.698m ± ∞ ¹    2.659m ± ∞ ¹       ~ (p=1.000 n=3) ²    1.963m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.556m ± ∞ ¹        ~ (p=0.700 n=3) ²    2.304m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.562m ± ∞ ¹        ~ (p=0.200 n=3) ²    1.992m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.378m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-8                      2.776m ± ∞ ¹    2.941m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.027m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.564m ± ∞ ¹        ~ (p=0.700 n=3) ²    2.381m ± ∞ ¹        ~ (p=0.400 n=3) ²     2.696m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.221m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.722m ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-16                     3.200m ± ∞ ¹    3.400m ± ∞ ¹       ~ (p=0.400 n=3) ²    2.468m ± ∞ ¹        ~ (p=0.100 n=3) ²     6.183m ± ∞ ¹        ~ (p=0.700 n=3) ²    2.543m ± ∞ ¹        ~ (p=0.100 n=3) ²     4.783m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.881m ± ∞ ¹       ~ (p=0.400 n=3) ²     8.262m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10                  195.7µ ± ∞ ¹    221.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²    147.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²     286.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    235.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²     294.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    180.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²     284.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-2                126.2µ ± ∞ ¹    138.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²    111.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²     204.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    173.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²     207.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    130.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²     182.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-4                100.2µ ± ∞ ¹    108.0µ ± ∞ ¹       ~ (p=0.700 n=3) ²    108.0µ ± ∞ ¹        ~ (p=0.700 n=3) ²     175.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    152.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²     182.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    118.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²     172.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-8                98.76µ ± ∞ ¹   113.15µ ± ∞ ¹       ~ (p=0.100 n=3) ²   113.17µ ± ∞ ¹        ~ (p=0.100 n=3) ²    211.57µ ± ∞ ¹        ~ (p=0.100 n=3) ²   169.63µ ± ∞ ¹        ~ (p=0.100 n=3) ²    203.56µ ± ∞ ¹        ~ (p=0.100 n=3) ²   130.73µ ± ∞ ¹       ~ (p=0.100 n=3) ²    169.00µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-16               123.7µ ± ∞ ¹    136.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²    120.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²     229.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    200.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²     230.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    262.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²     186.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90                  886.5µ ± ∞ ¹    929.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²    728.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1115.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    891.4µ ± ∞ ¹        ~ (p=0.700 n=3) ²    1120.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    818.4µ ± ∞ ¹       ~ (p=0.700 n=3) ²    1019.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-2                748.5µ ± ∞ ¹    787.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²    694.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1128.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    904.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1142.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    776.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²     999.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-4                739.3µ ± ∞ ¹    793.9µ ± ∞ ¹       ~ (p=0.700 n=3) ²    768.0µ ± ∞ ¹        ~ (p=1.000 n=3) ²    1176.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1016.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1153.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    826.6µ ± ∞ ¹       ~ (p=0.400 n=3) ²    1054.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-8                787.8µ ± ∞ ¹    853.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²    777.2µ ± ∞ ¹        ~ (p=1.000 n=3) ²    1742.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1130.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1602.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    819.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²    1176.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-16               997.6µ ± ∞ ¹   1057.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²    822.3µ ± ∞ ¹        ~ (p=0.700 n=3) ²   11581.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1225.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   12613.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    918.5µ ± ∞ ¹       ~ (p=0.700 n=3) ²   11610.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                                                           560.8µ          587.7µ        +4.79%                    472.7µ        -15.70%                     848.8µ        +51.35%                    617.0µ        +10.02%                     850.2µ        +51.60%                    533.9µ        -4.80%                     789.9µ        +40.84%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                                                    │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │          bench_mattn_driver.txt          │        bench_modernc_driver.txt         │       bench_tailscale_driver.txt       │      bench_zombiezen_direct.txt       │
                                                                    │           B/op           │     B/op       vs base                 │     B/op       vs base                │     B/op       vs base                  │      B/op       vs base                  │     B/op       vs base                  │     B/op       vs base                 │     B/op       vs base                │
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10                       225.4Ki ± ∞ ¹   236.6Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.0Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   454.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   451.7Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   233.5Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.6Ki ± ∞ ¹       ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-2                     226.9Ki ± ∞ ¹   235.5Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.9Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   454.5Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   452.9Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   231.9Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.2Ki ± ∞ ¹       ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-4                     224.3Ki ± ∞ ¹   235.0Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.8Ki ± ∞ ¹       ~ (p=0.200 n=3) ²   454.0Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   452.5Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   233.3Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   223.9Ki ± ∞ ¹       ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-8                     225.4Ki ± ∞ ¹   236.7Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.2Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   451.9Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   452.9Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   232.1Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.3Ki ± ∞ ¹       ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-16                    225.0Ki ± ∞ ¹   236.2Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.6Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   455.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   455.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   231.8Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.4Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90                       26.77Ki ± ∞ ¹   41.21Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   29.51Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   73.95Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   303.64Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   78.14Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   39.39Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   30.65Ki ± ∞ ¹       ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-2                     22.44Ki ± ∞ ¹   39.34Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.87Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   73.96Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.47Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   70.52Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   37.68Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   25.07Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-4                     19.33Ki ± ∞ ¹   39.76Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   25.97Ki ± ∞ ¹       ~ (p=0.200 n=3) ²   76.91Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.89Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   77.87Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   36.63Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   31.38Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-8                     28.38Ki ± ∞ ¹   41.43Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   29.53Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   73.97Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.94Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   72.65Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   38.02Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   26.49Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-16                    26.20Ki ± ∞ ¹   36.76Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   26.57Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   76.70Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   305.31Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   71.80Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   37.66Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   24.76Ki ± ∞ ¹       ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10                 224.7Ki ± ∞ ¹   238.5Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.0Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   456.7Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   458.9Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   235.9Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.1Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-2               225.4Ki ± ∞ ¹   238.4Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.3Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   457.6Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   456.8Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   234.9Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.3Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-4               225.6Ki ± ∞ ¹   238.7Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.0Ki ± ∞ ¹       ~ (p=0.200 n=3) ²   455.8Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   456.6Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   234.9Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.3Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-8               226.0Ki ± ∞ ¹   237.7Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.9Ki ± ∞ ¹       ~ (p=0.200 n=3) ²   454.9Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   456.5Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   234.4Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.7Ki ± ∞ ¹       ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-16              224.9Ki ± ∞ ¹   237.6Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.9Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   456.9Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   456.8Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   234.8Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.5Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90                 23.78Ki ± ∞ ¹   49.05Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.20Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   89.52Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.20Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   88.68Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   48.73Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.19Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-2               22.70Ki ± ∞ ¹   49.10Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.22Ki ± ∞ ¹       ~ (p=0.200 n=3) ²   84.55Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.38Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   89.44Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   49.35Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   29.54Ki ± ∞ ¹       ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-4               25.14Ki ± ∞ ¹   48.75Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   24.65Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   85.49Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.58Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   85.19Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   51.19Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   29.19Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-8               27.05Ki ± ∞ ¹   49.39Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   29.58Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   81.70Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.83Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   89.33Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   48.31Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   23.43Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-16              23.94Ki ± ∞ ¹   50.29Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.35Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   85.66Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.25Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   81.15Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   50.51Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   25.90Ki ± ∞ ¹       ~ (p=1.000 n=3) ²
geomean                                                                          74.21Ki         102.4Ki        +37.99%                   78.85Ki        +6.26%                   190.9Ki        +157.22%                    289.4Ki        +289.97%                   191.0Ki        +157.40%                   100.7Ki        +35.64%                   78.36Ki        +5.60%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                                                    │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt         │         bench_mattn_driver.txt         │         bench_modernc_driver.txt         │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                                                                    │        allocs/op         │  allocs/op    vs base                  │  allocs/op    vs base                  │   allocs/op    vs base                   │  allocs/op    vs base                  │   allocs/op    vs base                   │  allocs/op    vs base                  │  allocs/op    vs base                  │
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10                         236.0 ± ∞ ¹    736.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     973.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    695.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     973.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    545.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    270.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-2                       237.0 ± ∞ ¹    735.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    696.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    545.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    270.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-4                       235.0 ± ∞ ¹    735.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    696.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    270.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-8                       236.0 ± ∞ ¹    737.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    697.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-16                      235.0 ± ∞ ¹    737.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    696.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90                         28.00 ± ∞ ¹   447.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   799.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   423.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-2                       23.00 ± ∞ ¹   444.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   802.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   422.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-4                       20.00 ± ∞ ¹   445.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   803.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   421.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-8                       29.00 ± ∞ ¹   447.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   803.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   422.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-16                      27.00 ± ∞ ¹   441.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   804.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   422.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10                   235.0 ± ∞ ¹    777.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1009.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1013.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    585.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-2                 236.0 ± ∞ ¹    777.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1009.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1013.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    585.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-4                 236.0 ± ∞ ¹    778.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1010.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-8                 236.0 ± ∞ ¹    777.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1011.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    274.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-16                235.0 ± ∞ ¹    777.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1011.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    585.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    274.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90                   24.00 ± ∞ ¹   595.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1119.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   968.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1123.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-2                 23.00 ± ∞ ¹   595.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1120.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   969.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1123.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-4                 26.00 ± ∞ ¹   595.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1120.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   970.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1124.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-8                 28.00 ± ∞ ¹   596.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1121.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   971.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1123.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-16                25.00 ± ∞ ¹   597.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1121.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   969.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1126.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   264.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                                                            77.00          623.9        +710.34%                    231.7        +200.91%                    1.016k        +1219.99%                    797.1        +935.31%                    1.018k        +1221.96%                    528.1        +585.87%                    266.7        +246.38%
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
Query/Correlated-12                       654.5m ± ∞ ¹    645.6m ± ∞ ¹       ~ (p=1.000 n=3) ²   293.4m ± ∞ ¹        ~ (p=0.100 n=3) ²    304.1m ± ∞ ¹        ~ (p=0.100 n=3) ²   175.1m ± ∞ ¹        ~ (p=0.100 n=3) ²    304.6m ± ∞ ¹        ~ (p=0.100 n=3) ²   171.9m ± ∞ ¹        ~ (p=0.100 n=3) ²    323.5m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-12              120.37m ± ∞ ¹   113.73m ± ∞ ¹       ~ (p=0.200 n=3) ²   49.41m ± ∞ ¹        ~ (p=0.100 n=3) ²   122.68m ± ∞ ¹        ~ (p=0.100 n=3) ²   44.61m ± ∞ ¹        ~ (p=0.100 n=3) ²   119.95m ± ∞ ¹        ~ (p=1.000 n=3) ²   47.05m ± ∞ ¹        ~ (p=0.100 n=3) ²   184.74m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                   280.7m          271.0m        -3.46%                   120.4m        -57.10%                    193.2m        -31.18%                   88.37m        -68.52%                    191.2m        -31.89%                   89.94m        -67.96%                    244.5m        -12.90%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt          │     bench_eatonphil_direct.txt      │         bench_glebarez_driver.txt          │            bench_mattn_driver.txt            │          bench_modernc_driver.txt          │         bench_tailscale_driver.txt         │        bench_zombiezen_direct.txt        │
                            │           B/op           │      B/op       vs base                    │    B/op      vs base                │      B/op       vs base                    │      B/op        vs base                     │      B/op       vs base                    │      B/op       vs base                    │     B/op       vs base                   │
Query/Correlated-12                        96.00 ± ∞ ¹   71168.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   52.00 ± ∞ ¹       ~ (p=0.400 n=3) ²   63068.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   207216.00 ± ∞ ¹            ~ (p=0.100 n=3) ²   63092.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   47026.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   2656.00 ± ∞ ¹          ~ (p=0.600 n=3) ²
Query/CorrelatedParallel-12                143.0 ± ∞ ¹    71056.0 ± ∞ ¹           ~ (p=0.100 n=3) ²   249.0 ± ∞ ¹       ~ (p=0.700 n=3) ²    64097.0 ± ∞ ¹           ~ (p=0.700 n=3) ²    207468.0 ± ∞ ¹            ~ (p=0.100 n=3) ²    63604.0 ± ∞ ¹           ~ (p=0.700 n=3) ²    47464.0 ± ∞ ¹           ~ (p=0.700 n=3) ²    1610.0 ± ∞ ¹          ~ (p=0.700 n=3) ²
geomean                                    117.2          69.45Ki        +60593.07%                   113.8        -2.88%                    62.09Ki        +54164.99%                     202.5Ki        +176863.45%                    61.86Ki        +53966.19%                    46.14Ki        +40222.51%                   2.019Ki        +1664.91%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt           │    bench_eatonphil_direct.txt    │          bench_glebarez_driver.txt          │           bench_mattn_driver.txt            │          bench_modernc_driver.txt           │         bench_tailscale_driver.txt          │       bench_zombiezen_direct.txt        │
                            │        allocs/op         │   allocs/op     vs base                     │  allocs/op   vs base             │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                     │  allocs/op    vs base                   │
Query/Correlated-12                        1.000 ± ∞ ¹   5769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.300 n=3) ²     6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   30.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-12                2.000 ± ∞ ¹   5769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.300 n=3) ²     6783.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6773.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6778.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3768.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   14.000 ± ∞ ¹          ~ (p=0.600 n=3) ²
geomean                                    1.414           5.769k        +407829.90%                                ?               ³ ⁴     6.776k        +479035.30%                     6.772k        +478788.07%                     6.773k        +478858.67%                     3.767k        +266302.48%                    20.49        +1349.14%
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
Query/CorrelatedParallel                  641.9m ± ∞ ¹    644.5m ± ∞ ¹       ~ (p=0.400 n=3) ²    286.7m ± ∞ ¹        ~ (p=0.100 n=3) ²    302.9m ± ∞ ¹        ~ (p=0.100 n=3) ²   175.7m ± ∞ ¹        ~ (p=0.100 n=3) ²    303.1m ± ∞ ¹        ~ (p=0.100 n=3) ²   174.4m ± ∞ ¹        ~ (p=0.100 n=3) ²    324.2m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-2               434.31m ± ∞ ¹   437.90m ± ∞ ¹       ~ (p=0.100 n=3) ²   170.24m ± ∞ ¹        ~ (p=0.100 n=3) ²   184.12m ± ∞ ¹        ~ (p=0.100 n=3) ²   96.16m ± ∞ ¹        ~ (p=0.100 n=3) ²   184.05m ± ∞ ¹        ~ (p=0.100 n=3) ²   94.01m ± ∞ ¹        ~ (p=0.100 n=3) ²   188.90m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-4               236.70m ± ∞ ¹   230.93m ± ∞ ¹       ~ (p=0.400 n=3) ²   100.41m ± ∞ ¹        ~ (p=0.100 n=3) ²   122.30m ± ∞ ¹        ~ (p=0.100 n=3) ²   61.09m ± ∞ ¹        ~ (p=0.100 n=3) ²   122.60m ± ∞ ¹        ~ (p=0.100 n=3) ²   57.95m ± ∞ ¹        ~ (p=0.100 n=3) ²   132.94m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-8               154.21m ± ∞ ¹   122.80m ± ∞ ¹       ~ (p=0.700 n=3) ²    54.67m ± ∞ ¹        ~ (p=0.100 n=3) ²   139.76m ± ∞ ¹        ~ (p=0.700 n=3) ²   41.78m ± ∞ ¹        ~ (p=0.100 n=3) ²   137.43m ± ∞ ¹        ~ (p=0.700 n=3) ²   44.59m ± ∞ ¹        ~ (p=0.100 n=3) ²   187.45m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-16               79.80m ± ∞ ¹    79.69m ± ∞ ¹       ~ (p=1.000 n=3) ²    53.53m ± ∞ ¹        ~ (p=0.100 n=3) ²   117.75m ± ∞ ¹        ~ (p=0.100 n=3) ²   46.59m ± ∞ ¹        ~ (p=0.100 n=3) ²   115.30m ± ∞ ¹        ~ (p=0.100 n=3) ²   44.00m ± ∞ ¹        ~ (p=0.100 n=3) ²   184.92m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                   240.9m          229.6m        -4.72%                    107.5m        -55.39%                    162.2m        -32.68%                   72.54m        -69.89%                    161.1m        -33.15%                   71.47m        -70.34%                    195.0m        -19.05%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt          │       bench_eatonphil_direct.txt       │         bench_glebarez_driver.txt          │           bench_mattn_driver.txt            │          bench_modernc_driver.txt          │         bench_tailscale_driver.txt         │       bench_zombiezen_direct.txt        │
                            │           B/op           │      B/op       vs base                    │     B/op       vs base                 │      B/op       vs base                    │      B/op        vs base                    │      B/op       vs base                    │      B/op       vs base                    │     B/op       vs base                  │
Query/CorrelatedParallel                   68.00 ± ∞ ¹   71132.00 ± ∞ ¹           ~ (p=0.100 n=3) ²     34.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   63050.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   207185.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   63066.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   47009.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   2216.00 ± ∞ ¹         ~ (p=0.600 n=3) ²
Query/CorrelatedParallel-2                104.00 ± ∞ ¹   71082.00 ± ∞ ¹           ~ (p=0.100 n=3) ²     30.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   63011.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   207152.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   63041.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   46978.00 ± ∞ ¹           ~ (p=0.100 n=3) ²    534.00 ± ∞ ¹         ~ (p=0.700 n=3) ²
Query/CorrelatedParallel-4                 62.00 ± ∞ ¹   71012.00 ± ∞ ¹           ~ (p=0.100 n=3) ²    147.00 ± ∞ ¹        ~ (p=1.000 n=3) ²   63106.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   207155.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   63058.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   46972.00 ± ∞ ¹           ~ (p=0.100 n=3) ²    503.00 ± ∞ ¹         ~ (p=0.600 n=3) ²
Query/CorrelatedParallel-8                1371.0 ± ∞ ¹    70976.0 ± ∞ ¹           ~ (p=0.100 n=3) ²     133.0 ± ∞ ¹        ~ (p=0.600 n=3) ²    63154.0 ± ∞ ¹           ~ (p=0.700 n=3) ²    207347.0 ± ∞ ¹           ~ (p=0.100 n=3) ²    63168.0 ± ∞ ¹           ~ (p=0.400 n=3) ²    47209.0 ± ∞ ¹           ~ (p=0.700 n=3) ²     606.0 ± ∞ ¹         ~ (p=0.700 n=3) ²
Query/CorrelatedParallel-16              2.389Ki ± ∞ ¹   70.377Ki ± ∞ ¹           ~ (p=0.400 n=3) ²   2.122Ki ± ∞ ¹        ~ (p=1.000 n=3) ²   62.410Ki ± ∞ ¹           ~ (p=0.700 n=3) ²   202.821Ki ± ∞ ¹           ~ (p=0.100 n=3) ²   61.914Ki ± ∞ ¹           ~ (p=0.700 n=3) ²   46.542Ki ± ∞ ¹           ~ (p=0.700 n=3) ²   1.309Ki ± ∞ ¹         ~ (p=0.700 n=3) ²
geomean                                    271.3          69.58Ki        +26161.17%                     134.1        -50.58%                    61.76Ki        +23209.88%                     202.4Ki        +76305.61%                    61.67Ki        +23173.59%                    46.06Ki        +17283.26%                     864.7        +218.69%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt           │      bench_eatonphil_direct.txt      │          bench_glebarez_driver.txt          │           bench_mattn_driver.txt            │          bench_modernc_driver.txt           │         bench_tailscale_driver.txt          │       bench_zombiezen_direct.txt       │
                            │        allocs/op         │   allocs/op     vs base                     │  allocs/op   vs base                 │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                     │  allocs/op    vs base                  │
Query/CorrelatedParallel                   3.000 ± ∞ ¹   5771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   29.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-2                 3.000 ± ∞ ¹   5770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹            ~ (p=0.100 n=3) ²    8.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-4                 2.000 ± ∞ ¹   5769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.400 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹            ~ (p=0.100 n=3) ²    8.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-8                 2.000 ± ∞ ¹   5768.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.400 n=3) ²   6771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6773.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹            ~ (p=0.100 n=3) ²    8.000 ± ∞ ¹         ~ (p=0.600 n=3) ²
Query/CorrelatedParallel-16                4.000 ± ∞ ¹   5773.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹        ~ (p=1.000 n=3) ²   6782.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6773.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6774.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3768.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   13.000 ± ∞ ¹         ~ (p=0.700 n=3) ²
geomean                                    2.702           5.770k        +213459.23%                   1.380        -48.94%                     6.772k        +250551.32%                     6.772k        +250536.58%                     6.771k        +250484.75%                     3.767k        +139311.97%                    11.41        +322.13%
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
                         │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │       bench_eatonphil_direct.txt        │        bench_glebarez_driver.txt        │         bench_mattn_driver.txt          │        bench_modernc_driver.txt         │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt        │
                         │          sec/op          │    sec/op     vs base                │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op     vs base                 │    sec/op      vs base                  │
Query/GroupBy-12                       557.1µ ± ∞ ¹   575.9µ ± ∞ ¹       ~ (p=0.100 n=3) ²    327.6µ ± ∞ ¹         ~ (p=0.100 n=3) ²    487.9µ ± ∞ ¹         ~ (p=0.100 n=3) ²    339.4µ ± ∞ ¹         ~ (p=0.100 n=3) ²    502.6µ ± ∞ ¹         ~ (p=0.100 n=3) ²   261.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    473.2µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-12               62.24µ ± ∞ ¹   66.14µ ± ∞ ¹       ~ (p=0.100 n=3) ²   712.58µ ± ∞ ¹         ~ (p=0.100 n=3) ²   659.13µ ± ∞ ¹         ~ (p=0.100 n=3) ²   459.22µ ± ∞ ¹         ~ (p=0.100 n=3) ²   664.17µ ± ∞ ¹         ~ (p=0.100 n=3) ²   28.18µ ± ∞ ¹        ~ (p=0.100 n=3) ²   661.74µ ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                186.2µ         195.2µ        +4.81%                    483.1µ        +159.44%                    567.1µ        +204.54%                    394.8µ        +112.00%                    577.8µ        +210.27%                   85.82µ        -53.91%                    559.6µ        +200.48%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │   bench_eatonphil_direct.txt   │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                         │           B/op           │      B/op       vs base           │    B/op      vs base           │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │     B/op       vs base           │
Query/GroupBy-12                          0.0 ± ∞ ¹     2381.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ²     1963.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     7163.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1976.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1583.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     360.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-12                0.000 ± ∞ ¹   2368.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1966.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7175.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1982.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1527.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   371.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ³    2.319Ki        ?                                ?               ³    1.918Ki        ?                    7.001Ki        ?                    1.933Ki        ?                    1.518Ki        ?                     365.5        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                         │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/GroupBy-12                        0.000 ± ∞ ¹   119.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   155.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   52.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-12                0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   155.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ⁴     118.5        ?                                +0.00%               ⁴     122.0        ?                     155.0        ?                     122.0        ?                    51.50        ?                   6.000        ?
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
                         │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt        │        bench_glebarez_driver.txt        │         bench_mattn_driver.txt          │        bench_modernc_driver.txt         │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt        │
                         │          sec/op          │    sec/op      vs base                │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op     vs base                 │    sec/op      vs base                  │
Query/GroupByParallel                  549.9µ ± ∞ ¹    558.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²    339.0µ ± ∞ ¹         ~ (p=0.100 n=3) ²    479.0µ ± ∞ ¹         ~ (p=0.100 n=3) ²    307.4µ ± ∞ ¹         ~ (p=0.100 n=3) ²    478.3µ ± ∞ ¹         ~ (p=0.100 n=3) ²   264.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    462.0µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-2                280.4µ ± ∞ ¹    285.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²    372.2µ ± ∞ ¹         ~ (p=0.100 n=3) ²    384.6µ ± ∞ ¹         ~ (p=0.100 n=3) ²    209.7µ ± ∞ ¹         ~ (p=0.100 n=3) ²    380.5µ ± ∞ ¹         ~ (p=0.100 n=3) ²   130.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    369.8µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-4               143.45µ ± ∞ ¹   147.68µ ± ∞ ¹       ~ (p=0.100 n=3) ²   753.01µ ± ∞ ¹         ~ (p=0.100 n=3) ²   373.95µ ± ∞ ¹         ~ (p=0.100 n=3) ²   376.39µ ± ∞ ¹         ~ (p=0.100 n=3) ²   374.18µ ± ∞ ¹         ~ (p=0.100 n=3) ²   65.87µ ± ∞ ¹        ~ (p=0.100 n=3) ²   362.86µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-8                72.86µ ± ∞ ¹    75.94µ ± ∞ ¹       ~ (p=0.100 n=3) ²   809.49µ ± ∞ ¹         ~ (p=0.100 n=3) ²   640.34µ ± ∞ ¹         ~ (p=0.100 n=3) ²   485.95µ ± ∞ ¹         ~ (p=0.100 n=3) ²   642.68µ ± ∞ ¹         ~ (p=0.100 n=3) ²   34.03µ ± ∞ ¹        ~ (p=0.100 n=3) ²   641.21µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-16               62.32µ ± ∞ ¹    65.95µ ± ∞ ¹       ~ (p=0.100 n=3) ²   755.44µ ± ∞ ¹         ~ (p=0.100 n=3) ²   662.20µ ± ∞ ¹         ~ (p=0.100 n=3) ²   418.40µ ± ∞ ¹         ~ (p=0.100 n=3) ²   658.90µ ± ∞ ¹         ~ (p=0.100 n=3) ²   28.29µ ± ∞ ¹        ~ (p=0.100 n=3) ²   655.20µ ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                158.6µ          163.9µ        +3.29%                    566.1µ        +256.85%                    493.3µ        +210.99%                    345.7µ        +117.90%                    492.0µ        +210.18%                   73.86µ        -53.44%                    482.1µ        +203.92%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │   bench_eatonphil_direct.txt   │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                         │           B/op           │      B/op       vs base           │    B/op      vs base           │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │     B/op       vs base           │
Query/GroupByParallel                     0.0 ± ∞ ¹     2363.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ³     1849.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     7160.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1864.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1581.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     374.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-2                   0.0 ± ∞ ¹     2364.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ³     1960.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     7162.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1915.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1581.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     360.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-4                   0.0 ± ∞ ¹     2364.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ³     1962.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     7122.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1941.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1540.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     361.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-8                 0.000 ± ∞ ¹   2365.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=0.400 n=3) ²   1964.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7157.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1980.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1568.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   363.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-16                0.000 ± ∞ ¹   2368.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1970.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7172.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1988.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1516.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   369.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ⁴    2.309Ki        ?                                ?               ⁴    1.895Ki        ?                    6.987Ki        ?                    1.892Ki        ?                    1.520Ki        ?                     365.4        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                         │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/GroupByParallel                   0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   120.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   155.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   120.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   52.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-2                 0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   155.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   121.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   52.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-4                 0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   154.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   121.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-8                 0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   154.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-16                0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   155.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ⁴     118.0        ?                                +0.00%               ⁴     121.6        ?                     154.6        ?                     121.2        ?                    51.40        ?                   6.000        ?
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
                      │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                      │          sec/op          │    sec/op     vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │
Query/JSON-12                       9.515m ± ∞ ¹   9.745m ± ∞ ¹       ~ (p=0.100 n=3) ²   6.935m ± ∞ ¹        ~ (p=0.100 n=3) ²    9.174m ± ∞ ¹        ~ (p=0.100 n=3) ²   11.125m ± ∞ ¹       ~ (p=0.100 n=3) ²   10.055m ± ∞ ¹        ~ (p=0.100 n=3) ²   7.305m ± ∞ ¹        ~ (p=0.100 n=3) ²   10.542m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-12               8.881m ± ∞ ¹   8.769m ± ∞ ¹       ~ (p=0.100 n=3) ²   8.986m ± ∞ ¹        ~ (p=0.400 n=3) ²   15.190m ± ∞ ¹        ~ (p=0.100 n=3) ²    8.691m ± ∞ ¹       ~ (p=0.700 n=3) ²   15.180m ± ∞ ¹        ~ (p=0.100 n=3) ²   8.159m ± ∞ ¹        ~ (p=0.100 n=3) ²   14.915m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                             9.192m         9.244m        +0.56%                   7.894m        -14.12%                    11.81m        +28.43%                    9.833m        +6.97%                    12.35m        +34.40%                   7.720m        -16.02%                    12.54m        +36.41%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                      │ bench_ncruces_direct.txt │           bench_ncruces_driver.txt            │      bench_eatonphil_direct.txt      │           bench_glebarez_driver.txt           │             bench_mattn_driver.txt             │           bench_modernc_driver.txt            │          bench_tailscale_driver.txt          │        bench_zombiezen_direct.txt        │
                      │           B/op           │      B/op        vs base                      │    B/op      vs base                 │      B/op        vs base                      │       B/op        vs base                      │      B/op        vs base                      │      B/op        vs base                     │     B/op       vs base                   │
Query/JSON-12                        1.000 ± ∞ ¹   64911.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=1.000 n=3) ²   56951.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   201046.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56967.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   32891.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   424.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-12                27.00 ± ∞ ¹    64867.00 ± ∞ ¹             ~ (p=0.100 n=3) ²   21.00 ± ∞ ¹        ~ (p=0.700 n=3) ²    56990.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    201104.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    57024.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    32963.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    441.00 ± ∞ ¹          ~ (p=0.100 n=3) ²
geomean                              5.196           63.37Ki        +1248689.32%                   4.583        -11.81%                     55.64Ki        +1096297.72%                      196.4Ki        +3869590.14%                     55.66Ki        +1096778.77%                     32.16Ki        +633580.03%                     432.4        +8221.86%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                      │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │   bench_zombiezen_direct.txt   │
                      │        allocs/op         │   allocs/op     vs base           │  allocs/op   vs base                │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │  allocs/op   vs base           │
Query/JSON-12                        0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
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
Query/JSONParallel                  9.740m ± ∞ ¹   9.695m ± ∞ ¹       ~ (p=1.000 n=3) ²   7.201m ± ∞ ¹        ~ (p=0.100 n=3) ²    9.410m ± ∞ ¹        ~ (p=0.200 n=3) ²   7.559m ± ∞ ¹       ~ (p=0.100 n=3) ²    9.508m ± ∞ ¹        ~ (p=0.700 n=3) ²   7.241m ± ∞ ¹        ~ (p=0.100 n=3) ²   10.020m ± ∞ ¹        ~ (p=0.400 n=3) ²
Query/JSONParallel-2                6.887m ± ∞ ¹   7.091m ± ∞ ¹       ~ (p=0.100 n=3) ²   5.663m ± ∞ ¹        ~ (p=0.100 n=3) ²    7.396m ± ∞ ¹        ~ (p=0.100 n=3) ²   6.044m ± ∞ ¹       ~ (p=0.100 n=3) ²    7.387m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.629m ± ∞ ¹        ~ (p=0.100 n=3) ²    8.317m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-4                6.296m ± ∞ ¹   6.329m ± ∞ ¹       ~ (p=1.000 n=3) ²   5.073m ± ∞ ¹        ~ (p=0.100 n=3) ²    6.092m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.947m ± ∞ ¹       ~ (p=0.100 n=3) ²    6.094m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.806m ± ∞ ¹        ~ (p=0.100 n=3) ²    7.847m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-8                7.375m ± ∞ ¹   7.123m ± ∞ ¹       ~ (p=0.700 n=3) ²   7.282m ± ∞ ¹        ~ (p=0.100 n=3) ²   13.884m ± ∞ ¹        ~ (p=0.100 n=3) ²   7.692m ± ∞ ¹       ~ (p=0.100 n=3) ²   13.481m ± ∞ ¹        ~ (p=0.100 n=3) ²   6.868m ± ∞ ¹        ~ (p=0.100 n=3) ²   14.558m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-16               9.543m ± ∞ ¹   8.986m ± ∞ ¹       ~ (p=0.400 n=3) ²   9.099m ± ∞ ¹        ~ (p=0.100 n=3) ²   15.469m ± ∞ ¹        ~ (p=0.100 n=3) ²   9.654m ± ∞ ¹       ~ (p=0.700 n=3) ²   14.834m ± ∞ ¹        ~ (p=0.100 n=3) ²   9.045m ± ∞ ¹        ~ (p=0.400 n=3) ²   13.239m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                             7.846m         7.744m        -1.29%                   6.720m        -14.34%                    9.814m        +25.09%                   7.260m        -7.46%                    9.694m        +23.56%                   6.815m        -13.14%                    10.47m        +33.50%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                      │ bench_ncruces_direct.txt │           bench_ncruces_driver.txt            │    bench_eatonphil_direct.txt    │           bench_glebarez_driver.txt           │             bench_mattn_driver.txt             │           bench_modernc_driver.txt            │          bench_tailscale_driver.txt          │        bench_zombiezen_direct.txt        │
                      │           B/op           │      B/op        vs base                      │    B/op      vs base             │      B/op        vs base                      │       B/op        vs base                      │      B/op        vs base                      │      B/op        vs base                     │     B/op       vs base                   │
Query/JSONParallel                   1.000 ± ∞ ¹   64844.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     56878.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   201020.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56859.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   32891.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   365.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-2                 1.000 ± ∞ ¹   64844.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=1.000 n=3) ³     56948.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   201017.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56901.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   32896.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   366.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-4                 8.000 ± ∞ ¹   64844.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     56953.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   201027.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56931.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   32891.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   373.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-8                 16.00 ± ∞ ¹    64850.00 ± ∞ ¹             ~ (p=0.100 n=3) ²   12.00 ± ∞ ¹  ~ (p=1.000 n=3) ²      56992.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    201043.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    57000.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    32906.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    421.00 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-16                44.00 ± ∞ ¹    64912.00 ± ∞ ¹             ~ (p=0.100 n=3) ²   64.00 ± ∞ ¹  ~ (p=0.700 n=3) ²      57083.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    201228.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    57033.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    32987.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    532.00 ± ∞ ¹          ~ (p=0.700 n=3) ²
geomean                              5.625           63.34Ki        +1152918.68%                                ?               ⁴ ⁵     55.64Ki        +1012690.19%                      196.4Ki        +3574341.83%                     55.61Ki        +1012228.05%                     32.14Ki        +585027.49%                     406.9        +7134.40%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
⁵ ratios must be >0 to compute geomean

                      │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │   bench_zombiezen_direct.txt   │
                      │        allocs/op         │   allocs/op     vs base           │  allocs/op   vs base                │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │  allocs/op   vs base           │
Query/JSONParallel                   0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-2                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-4                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-8                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-16                0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4023.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.600 n=3) ²
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
                         │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt       │         bench_mattn_driver.txt         │       bench_modernc_driver.txt        │      bench_tailscale_driver.txt      │      bench_zombiezen_direct.txt       │
                         │          sec/op          │    sec/op     vs base                │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op     vs base                │    sec/op     vs base                 │
Query/OrderBy-12                       75.62m ± ∞ ¹   89.76m ± ∞ ¹       ~ (p=0.100 n=3) ²   47.18m ± ∞ ¹        ~ (p=0.100 n=3) ²   81.55m ± ∞ ¹        ~ (p=0.100 n=3) ²   109.11m ± ∞ ¹        ~ (p=0.100 n=3) ²   80.84m ± ∞ ¹        ~ (p=0.100 n=3) ²   61.42m ± ∞ ¹       ~ (p=0.100 n=3) ²   84.05m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-12               51.83m ± ∞ ¹   50.96m ± ∞ ¹       ~ (p=0.200 n=3) ²   65.34m ± ∞ ¹        ~ (p=0.100 n=3) ²   76.13m ± ∞ ¹        ~ (p=0.100 n=3) ²    96.50m ± ∞ ¹        ~ (p=0.100 n=3) ²   66.59m ± ∞ ¹        ~ (p=0.100 n=3) ²   51.88m ± ∞ ¹       ~ (p=1.000 n=3) ²   74.01m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                62.61m         67.63m        +8.03%                   55.52m        -11.31%                   78.79m        +25.85%                    102.6m        +63.90%                   73.37m        +17.19%                   56.45m        -9.83%                   78.87m        +25.98%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt             │      bench_eatonphil_direct.txt      │            bench_glebarez_driver.txt            │              bench_mattn_driver.txt              │            bench_modernc_driver.txt             │           bench_tailscale_driver.txt           │        bench_zombiezen_direct.txt        │
                         │           B/op           │       B/op         vs base                      │    B/op      vs base                 │       B/op         vs base                      │        B/op         vs base                      │       B/op         vs base                      │       B/op         vs base                     │      B/op       vs base                  │
Query/OrderBy-12                      372.000 ± ∞ ¹   6399302.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   9.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6397798.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   11999099.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   6397824.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   2798860.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3148.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-12                510.0 ± ∞ ¹     6399422.0 ± ∞ ¹             ~ (p=0.100 n=3) ²   342.0 ± ∞ ¹        ~ (p=0.400 n=3) ²     6398517.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     12002898.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     6398202.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     2799068.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     1662.0 ± ∞ ¹         ~ (p=0.700 n=3) ²
geomean                                 435.6             6.103Mi        +1469097.30%                   55.48        -87.26%                       6.102Mi        +1468820.76%                        11.45Mi        +2755148.78%                       6.102Mi        +1468787.58%                       2.669Mi        +642500.05%                    2.234Ki        +425.14%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt            │    bench_eatonphil_direct.txt    │           bench_glebarez_driver.txt            │             bench_mattn_driver.txt             │            bench_modernc_driver.txt            │           bench_tailscale_driver.txt           │       bench_zombiezen_direct.txt       │
                         │        allocs/op         │    allocs/op      vs base                      │  allocs/op   vs base             │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │  allocs/op    vs base                  │
Query/OrderBy-12                        8.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449782.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349771.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449782.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   36.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-12                9.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.300 n=3) ²     449785.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349787.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449784.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149767.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   19.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                 8.485             349.8k        +4122079.15%                                ?               ³ ⁴       449.8k        +5300649.38%                       349.8k        +4122085.05%                       449.8k        +5300643.49%                       149.8k        +1764915.13%                    26.15        +208.22%
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
                         │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt      │       bench_glebarez_driver.txt       │         bench_mattn_driver.txt         │       bench_modernc_driver.txt        │      bench_tailscale_driver.txt      │      bench_zombiezen_direct.txt       │
                         │          sec/op          │    sec/op     vs base                 │    sec/op     vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op     vs base                │    sec/op     vs base                 │
Query/OrderByParallel                  75.07m ± ∞ ¹   89.27m ± ∞ ¹        ~ (p=0.100 n=3) ²   46.68m ± ∞ ¹       ~ (p=0.100 n=3) ²   82.31m ± ∞ ¹        ~ (p=0.100 n=3) ²    88.14m ± ∞ ¹        ~ (p=0.100 n=3) ²   82.29m ± ∞ ¹        ~ (p=0.100 n=3) ²   61.45m ± ∞ ¹       ~ (p=0.100 n=3) ²   85.05m ± ∞ ¹        ~ (p=0.700 n=3) ²
Query/OrderByParallel-2                48.34m ± ∞ ¹   58.61m ± ∞ ¹        ~ (p=0.100 n=3) ²   37.98m ± ∞ ¹       ~ (p=0.100 n=3) ²   56.12m ± ∞ ¹        ~ (p=0.100 n=3) ²    59.66m ± ∞ ¹        ~ (p=0.100 n=3) ²   56.22m ± ∞ ¹        ~ (p=0.100 n=3) ²   42.95m ± ∞ ¹       ~ (p=0.100 n=3) ²   48.94m ± ∞ ¹        ~ (p=0.700 n=3) ²
Query/OrderByParallel-4                39.99m ± ∞ ¹   44.65m ± ∞ ¹        ~ (p=0.100 n=3) ²   34.34m ± ∞ ¹       ~ (p=0.100 n=3) ²   40.43m ± ∞ ¹        ~ (p=0.700 n=3) ²    53.57m ± ∞ ¹        ~ (p=0.100 n=3) ²   39.65m ± ∞ ¹        ~ (p=1.000 n=3) ²   37.11m ± ∞ ¹       ~ (p=0.100 n=3) ²   40.65m ± ∞ ¹        ~ (p=0.400 n=3) ²
Query/OrderByParallel-8                40.45m ± ∞ ¹   46.37m ± ∞ ¹        ~ (p=0.200 n=3) ²   58.19m ± ∞ ¹       ~ (p=0.100 n=3) ²   63.63m ± ∞ ¹        ~ (p=0.100 n=3) ²    82.79m ± ∞ ¹        ~ (p=0.100 n=3) ²   63.47m ± ∞ ¹        ~ (p=0.100 n=3) ²   41.79m ± ∞ ¹       ~ (p=0.700 n=3) ²   72.04m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-16               47.22m ± ∞ ¹   53.20m ± ∞ ¹        ~ (p=0.100 n=3) ²   65.00m ± ∞ ¹       ~ (p=0.100 n=3) ²   70.81m ± ∞ ¹        ~ (p=0.100 n=3) ²   102.97m ± ∞ ¹        ~ (p=0.100 n=3) ²   67.88m ± ∞ ¹        ~ (p=0.100 n=3) ²   51.64m ± ∞ ¹       ~ (p=0.100 n=3) ²   76.08m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                48.81m         56.51m        +15.77%                   47.04m        -3.64%                   60.95m        +24.87%                    75.18m        +54.01%                   60.19m        +23.31%                   46.24m        -5.27%                   62.15m        +27.33%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt             │      bench_eatonphil_direct.txt      │            bench_glebarez_driver.txt            │              bench_mattn_driver.txt              │            bench_modernc_driver.txt             │           bench_tailscale_driver.txt           │        bench_zombiezen_direct.txt        │
                         │           B/op           │       B/op         vs base                      │    B/op      vs base                 │       B/op         vs base                      │        B/op         vs base                      │       B/op         vs base                      │       B/op         vs base                     │      B/op       vs base                  │
Query/OrderByParallel                 369.000 ± ∞ ¹   6399278.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6397698.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   11999011.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   6397706.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   2798848.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   2939.000 ± ∞ ¹         ~ (p=0.200 n=3) ²
Query/OrderByParallel-2               369.000 ± ∞ ¹   6399283.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6397712.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   11999020.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   6397728.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   2798868.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1140.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-4                417.00 ± ∞ ¹    6399277.00 ± ∞ ¹             ~ (p=0.100 n=3) ²   38.00 ± ∞ ¹        ~ (p=0.100 n=3) ²    6397856.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    11999090.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    6397879.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    2798865.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    1173.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-8                 490.0 ± ∞ ¹     6399329.0 ± ∞ ¹             ~ (p=0.100 n=3) ²   103.0 ± ∞ ¹        ~ (p=0.700 n=3) ²     6398003.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     11999729.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     6398063.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     2798933.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     1235.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-16                968.0 ± ∞ ¹     6399791.0 ± ∞ ¹             ~ (p=0.100 n=3) ²   360.0 ± ∞ ¹        ~ (p=0.400 n=3) ²     6398565.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     12002509.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     6398518.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     2799209.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     1704.0 ± ∞ ¹         ~ (p=0.700 n=3) ²
geomean                                 485.3             6.103Mi        +1318419.76%                   33.51        -93.10%                       6.102Mi        +1318126.19%                        11.44Mi        +2472333.15%                       6.102Mi        +1318128.67%                       2.669Mi        +576589.78%                    1.490Ki        +214.38%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt            │    bench_eatonphil_direct.txt    │           bench_glebarez_driver.txt            │             bench_mattn_driver.txt             │            bench_modernc_driver.txt            │           bench_tailscale_driver.txt           │       bench_zombiezen_direct.txt       │
                         │        allocs/op         │    allocs/op      vs base                      │  allocs/op   vs base             │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │  allocs/op    vs base                  │
Query/OrderByParallel                   8.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349771.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   34.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-2                 8.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349771.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   16.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-4                 8.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449782.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349772.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449781.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   16.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-8                 9.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449783.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349775.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449783.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   17.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-16               10.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449788.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349787.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449788.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149768.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   19.000 ± ∞ ¹         ~ (p=0.600 n=3) ²
geomean                                 8.565             349.8k        +4083947.09%                                ?               ³ ⁴       449.8k        +5251602.99%                       349.8k        +4083907.39%                       449.8k        +5251600.65%                       149.8k        +1748586.26%                    19.49        +127.55%
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
                              │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │      bench_eatonphil_direct.txt       │      bench_glebarez_driver.txt       │        bench_mattn_driver.txt         │       bench_modernc_driver.txt       │      bench_tailscale_driver.txt       │      bench_zombiezen_direct.txt      │
                              │          sec/op          │    sec/op     vs base                │    sec/op     vs base                 │    sec/op     vs base                │    sec/op     vs base                 │    sec/op     vs base                │    sec/op     vs base                 │    sec/op     vs base                │
Query/RecursiveCTE-12                       6.197m ± ∞ ¹   6.212m ± ∞ ¹       ~ (p=0.100 n=3) ²   2.533m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.277m ± ∞ ¹       ~ (p=0.100 n=3) ²   2.427m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.295m ± ∞ ¹       ~ (p=0.100 n=3) ²   2.411m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.266m ± ∞ ¹       ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-12               686.2µ ± ∞ ¹   687.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²   286.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   673.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²   276.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   656.9µ ± ∞ ¹       ~ (p=0.100 n=3) ²   271.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   660.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²
geomean                                     2.062m         2.067m        +0.24%                   851.8µ        -58.69%                   1.885m        -8.58%                   819.1µ        -60.28%                   1.865m        -9.55%                   809.0µ        -60.77%                   1.865m        -9.53%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                              │           B/op           │      B/op       vs base           │    B/op      vs base                │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │     B/op       vs base           │
Query/RecursiveCTE-12                          0.0 ± ∞ ¹     2492.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹       ~ (p=1.000 n=3) ²     2354.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     6857.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     2296.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1505.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     376.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-12                1.000 ± ∞ ¹   2484.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹       ~ (p=1.000 n=3) ²   2358.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6856.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2297.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1506.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   369.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                ³    2.430Ki        ?                                +0.00%               ³    2.301Ki        ?                    6.696Ki        ?                    2.243Ki        ?                    1.470Ki        ?                     372.5        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean

                              │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                              │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/RecursiveCTE-12                        0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   143.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-12                0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                ⁴     110.0        ?                                +0.00%               ⁴     113.0        ?                     142.5        ?                     112.0        ?                    49.00        ?                   6.000        ?
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
                              │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │       bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                              │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │
Query/RecursiveCTEParallel                  6.105m ± ∞ ¹    6.107m ± ∞ ¹       ~ (p=0.700 n=3) ²   2.494m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.258m ± ∞ ¹        ~ (p=0.100 n=3) ²   2.399m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.275m ± ∞ ¹       ~ (p=0.100 n=3) ²   2.418m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.232m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                3.081m ± ∞ ¹    3.094m ± ∞ ¹       ~ (p=0.100 n=3) ²   1.261m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.665m ± ∞ ¹        ~ (p=0.100 n=3) ²   1.229m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.671m ± ∞ ¹       ~ (p=0.100 n=3) ²   1.231m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.646m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4               1571.2µ ± ∞ ¹   1580.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²   643.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1390.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²   621.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1389.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²   620.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1385.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                794.2µ ± ∞ ¹    800.5µ ± ∞ ¹       ~ (p=0.700 n=3) ²   326.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    738.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   315.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    746.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²   314.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    748.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16               683.0µ ± ∞ ¹    690.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²   283.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    650.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   279.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    652.6µ ± ∞ ¹       ~ (p=0.100 n=3) ²   276.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    652.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                     1.742m          1.752m        +0.58%                   715.6µ        -58.92%                    1.564m        -10.21%                   694.4µ        -60.13%                    1.570m        -9.86%                   693.6µ        -60.18%                    1.564m        -10.21%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │      bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                              │           B/op           │      B/op       vs base           │    B/op      vs base                 │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │     B/op       vs base           │
Query/RecursiveCTEParallel                     0.0 ± ∞ ¹     2482.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹        ~ (p=1.000 n=3) ³     2261.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     6857.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     2258.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1504.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     362.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                   0.0 ± ∞ ¹     2481.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹        ~ (p=1.000 n=3) ²     2263.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     6822.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     2257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1501.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     375.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4                   0.0 ± ∞ ¹     2481.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹        ~ (p=1.000 n=3) ²     2353.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     6841.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     2266.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1503.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     365.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                 2.000 ± ∞ ¹   2483.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.500 n=3) ²   2356.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6858.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2298.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1505.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   363.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16                2.000 ± ∞ ¹   2489.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹        ~ (p=1.000 n=3) ²   2366.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6863.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2286.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1506.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   368.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                ⁴    2.425Ki        ?                                -12.94%               ⁴    2.265Ki        ?                    6.688Ki        ?                    2.220Ki        ?                    1.469Ki        ?                     366.6        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                              │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                              │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/RecursiveCTEParallel                   0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   143.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   48.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   48.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16                0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   143.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                ⁴     110.0        ?                                +0.00%               ⁴     112.6        ?                     142.4        ?                     112.0        ?                    48.60        ?                   6.000        ?
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
                        │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │         bench_mattn_driver.txt          │        bench_modernc_driver.txt        │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                        │          sec/op          │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                  │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │
Query/Window-12                      2134.1µ ± ∞ ¹   2430.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    973.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1869.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   2098.1µ ± ∞ ¹         ~ (p=0.100 n=3) ²   1891.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1048.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1715.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-12               235.5µ ± ∞ ¹    328.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1487.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    904.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1072.4µ ± ∞ ¹         ~ (p=0.100 n=3) ²    900.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    248.1µ ± ∞ ¹        ~ (p=0.700 n=3) ²    856.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                               708.9µ          892.9µ        +25.95%                    1.203m        +69.73%                    1.300m        +83.43%                    1.500m        +111.59%                    1.305m        +84.12%                    509.9µ        -28.07%                    1.212m        +71.02%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                        │ bench_ncruces_direct.txt │      bench_ncruces_driver.txt      │      bench_eatonphil_direct.txt       │     bench_glebarez_driver.txt      │       bench_mattn_driver.txt        │      bench_modernc_driver.txt      │     bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                        │           B/op           │      B/op        vs base           │    B/op      vs base                  │      B/op        vs base           │       B/op        vs base           │      B/op        vs base           │      B/op        vs base           │     B/op       vs base           │
Query/Window-12                          0.0 ± ∞ ¹     62777.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹         ~ (p=1.000 n=3) ²     54883.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198935.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54917.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30793.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     375.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-12                1.000 ± ∞ ¹   62763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹         ~ (p=0.100 n=3) ²   54889.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   198952.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54909.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   30745.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   373.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                          ³     61.30Ki        ?                                +164.58%               ³     53.60Ki        ?                      194.3Ki        ?                     53.63Ki        ?                     30.05Ki        ?                     374.0        ?
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
                        │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │         bench_mattn_driver.txt         │        bench_modernc_driver.txt        │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                        │          sec/op          │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │
Query/WindowParallel                 2106.9µ ± ∞ ¹   2384.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    989.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1870.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1615.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1871.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1067.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1613.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-2               1060.3µ ± ∞ ¹   1217.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    625.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    968.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    860.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    971.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    546.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    833.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-4                539.7µ ± ∞ ¹    660.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    966.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    540.1µ ± ∞ ¹        ~ (p=0.700 n=3) ²    520.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    545.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    283.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    468.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-8                273.6µ ± ∞ ¹    347.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1629.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    874.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1055.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    867.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    148.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    845.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-16               247.0µ ± ∞ ¹    306.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1386.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    921.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1088.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    912.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    143.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    857.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                               605.7µ          727.7µ        +20.15%                    1.062m        +75.34%                    953.3µ        +57.40%                    963.6µ        +59.09%                    952.7µ        +57.30%                    323.2µ        -46.63%                    855.0µ        +41.17%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                        │ bench_ncruces_direct.txt │      bench_ncruces_driver.txt      │   bench_eatonphil_direct.txt   │     bench_glebarez_driver.txt      │       bench_mattn_driver.txt        │      bench_modernc_driver.txt      │     bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                        │           B/op           │      B/op        vs base           │    B/op      vs base           │      B/op        vs base           │       B/op        vs base           │      B/op        vs base           │      B/op        vs base           │     B/op       vs base           │
Query/WindowParallel                     0.0 ± ∞ ¹     62762.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ³     54775.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198936.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54789.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30792.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     360.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-2                   0.0 ± ∞ ¹     62761.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ³     54807.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198854.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54806.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30793.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     363.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-4                 0.000 ± ∞ ¹   62762.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.400 n=3) ²   54880.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   198856.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54885.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   30793.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   361.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-8                 0.000 ± ∞ ¹   62769.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4.000 ± ∞ ¹  ~ (p=0.400 n=3) ²   54884.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   198910.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54903.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   30794.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   366.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-16                1.000 ± ∞ ¹   62767.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   9.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54898.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   198971.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54912.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   30769.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   370.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                          ⁴     61.29Ki        ?                                ?               ⁴     53.56Ki        ?                      194.2Ki        ?                     53.57Ki        ?                     30.07Ki        ?                     364.0        ?
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
3.45.1
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
