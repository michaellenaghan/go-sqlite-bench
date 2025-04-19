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
Baseline/Conn-12                                     190.7n ± ∞ ¹    221.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   186.9n ± ∞ ¹        ~ (p=0.600 n=3) ²    224.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    220.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    227.6n ± ∞ ¹        ~ (p=0.100 n=3) ²    224.7n ± ∞ ¹       ~ (p=0.100 n=3) ²   1104.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-12                             225.0n ± ∞ ¹    536.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   223.7n ± ∞ ¹        ~ (p=0.700 n=3) ²    542.3n ± ∞ ¹        ~ (p=0.100 n=3) ²    541.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    543.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    544.1n ± ∞ ¹       ~ (p=0.100 n=3) ²   1342.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1-12                                 1936.0n ± ∞ ¹   1968.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   923.8n ± ∞ ¹        ~ (p=0.100 n=3) ²   2044.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2617.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2038.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1118.0n ± ∞ ¹       ~ (p=0.100 n=3) ²   3160.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-12                         1392.0n ± ∞ ¹   1165.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   729.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   2262.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1479.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2314.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    911.3n ± ∞ ¹       ~ (p=0.100 n=3) ²   2639.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePrepared-12                       536.6n ± ∞ ¹    529.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   488.5n ± ∞ ¹        ~ (p=0.100 n=3) ²   1997.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1308.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2010.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    411.2n ± ∞ ¹       ~ (p=0.100 n=3) ²   1508.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-12               678.5n ± ∞ ¹    934.5n ± ∞ ¹        ~ (p=0.100 n=3) ²   475.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   1769.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1077.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1798.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    766.4n ± ∞ ¹       ~ (p=0.100 n=3) ²   1549.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                              589.8n          716.1n        +21.41%                   432.6n        -26.66%                    1.122µ        +90.15%                    931.2n        +57.89%                    1.132µ        +91.96%                    583.0n        -1.16%                    1.751µ        +196.94%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                       │ bench_ncruces_direct.txt │    bench_ncruces_driver.txt    │      bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt    │     bench_mattn_driver.txt      │    bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt    │
                                       │           B/op           │    B/op      vs base           │    B/op      vs base                 │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │
Baseline/Conn-12                                       0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   353.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-12                               0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1-12                                   48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   705.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-12                           48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePrepared-12                         0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-12                 0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                         ⁴   46.15        ?                                +25.99%               ⁴    159.6        ?                    164.2        ?                    159.6        ?                    90.34        ?                    443.8        ?
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
Baseline/ConnParallel                                187.0n ± ∞ ¹    234.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   186.1n ± ∞ ¹        ~ (p=0.600 n=3) ²    241.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    225.9n ± ∞ ¹        ~ (p=0.100 n=3) ²    241.9n ± ∞ ¹        ~ (p=0.100 n=3) ²    230.7n ± ∞ ¹       ~ (p=0.100 n=3) ²   1372.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-2                              160.2n ± ∞ ¹    229.4n ± ∞ ¹        ~ (p=0.100 n=3) ²   155.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    219.3n ± ∞ ¹        ~ (p=0.100 n=3) ²    212.6n ± ∞ ¹        ~ (p=0.100 n=3) ²    220.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    215.8n ± ∞ ¹       ~ (p=0.100 n=3) ²    697.2n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-4                              208.1n ± ∞ ¹    305.4n ± ∞ ¹        ~ (p=0.100 n=3) ²   205.2n ± ∞ ¹        ~ (p=0.400 n=3) ²    262.9n ± ∞ ¹        ~ (p=0.100 n=3) ²    262.1n ± ∞ ¹        ~ (p=0.100 n=3) ²    261.5n ± ∞ ¹        ~ (p=0.100 n=3) ²    263.0n ± ∞ ¹       ~ (p=0.100 n=3) ²    668.3n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-8                              222.5n ± ∞ ¹    510.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   218.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    516.9n ± ∞ ¹        ~ (p=0.100 n=3) ²    512.4n ± ∞ ¹        ~ (p=0.100 n=3) ²    517.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    514.5n ± ∞ ¹       ~ (p=0.100 n=3) ²   1233.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-16                             226.8n ± ∞ ¹    539.6n ± ∞ ¹        ~ (p=0.100 n=3) ²   221.5n ± ∞ ¹        ~ (p=0.100 n=3) ²    550.1n ± ∞ ¹        ~ (p=0.100 n=3) ²    547.6n ± ∞ ¹        ~ (p=0.100 n=3) ²    553.4n ± ∞ ¹        ~ (p=0.100 n=3) ²    546.4n ± ∞ ¹       ~ (p=0.100 n=3) ²   1331.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel                            1995.0n ± ∞ ¹   1900.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   936.8n ± ∞ ¹        ~ (p=0.100 n=3) ²   1870.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1649.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1872.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1148.0n ± ∞ ¹       ~ (p=0.100 n=3) ²   2859.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-2                          1056.0n ± ∞ ¹   1084.0n ± ∞ ¹        ~ (p=0.700 n=3) ²   547.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   1088.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    981.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   1073.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    646.7n ± ∞ ¹       ~ (p=0.100 n=3) ²   1519.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-4                           602.8n ± ∞ ¹    584.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   337.9n ± ∞ ¹        ~ (p=0.100 n=3) ²    970.3n ± ∞ ¹        ~ (p=0.100 n=3) ²    605.7n ± ∞ ¹        ~ (p=1.000 n=3) ²    970.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    388.9n ± ∞ ¹       ~ (p=0.100 n=3) ²   1088.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-8                          1385.0n ± ∞ ¹   1118.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   750.2n ± ∞ ¹        ~ (p=0.100 n=3) ²   2193.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1451.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2243.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    880.7n ± ∞ ¹       ~ (p=0.100 n=3) ²   2512.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-16                         1397.0n ± ∞ ¹   1171.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   740.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   2189.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1483.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2304.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    936.1n ± ∞ ¹       ~ (p=0.100 n=3) ²   2702.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel                  532.8n ± ∞ ¹    522.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   485.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   1859.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    848.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   1852.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    418.2n ± ∞ ¹       ~ (p=0.100 n=3) ²   1608.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-2                347.6n ± ∞ ¹    455.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   337.4n ± ∞ ¹        ~ (p=0.100 n=3) ²   1039.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    602.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   1030.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    342.7n ± ∞ ¹       ~ (p=0.100 n=3) ²    824.5n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-4                427.5n ± ∞ ¹    359.4n ± ∞ ¹        ~ (p=0.100 n=3) ²   301.5n ± ∞ ¹        ~ (p=0.100 n=3) ²    740.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    542.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    744.5n ± ∞ ¹        ~ (p=0.100 n=3) ²    276.6n ± ∞ ¹       ~ (p=0.100 n=3) ²    693.9n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-8                676.4n ± ∞ ¹    879.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   471.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   1716.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1072.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1741.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    736.8n ± ∞ ¹       ~ (p=0.100 n=3) ²   1448.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-16               682.0n ± ∞ ¹    942.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   470.4n ± ∞ ¹        ~ (p=0.100 n=3) ²   1759.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1105.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1787.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    774.8n ± ∞ ¹       ~ (p=0.100 n=3) ²   1512.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                              497.4n          602.9n        +21.21%                   367.7n        -26.08%                    885.6n        +78.05%                    669.7n        +34.64%                    891.3n        +79.18%                    485.8n        -2.33%                    1.325µ        +166.29%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                       │ bench_ncruces_direct.txt │    bench_ncruces_driver.txt    │      bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt    │     bench_mattn_driver.txt      │    bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt    │
                                       │           B/op           │    B/op      vs base           │    B/op      vs base                 │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │
Baseline/ConnParallel                                  0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   387.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-2                                0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-4                                0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-8                                0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-16                               0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel                              48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   729.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-2                            48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-4                            48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-8                            48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-16                           48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel                    0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   378.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-2                  0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-4                  0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-8                  0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-16                 0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                         ⁴   46.15        ?                                +25.99%               ⁴    159.6        ?                    164.2        ?                    159.6        ?                    90.34        ?                    449.5        ?
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
                              │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                │
Populate/PopulateDB-12                       2.532 ± ∞ ¹     2.503 ± ∞ ¹       ~ (p=0.400 n=3) ²    2.034 ± ∞ ¹        ~ (p=0.100 n=3) ²     2.742 ± ∞ ¹        ~ (p=0.100 n=3) ²     2.216 ± ∞ ¹       ~ (p=0.100 n=3) ²     2.763 ± ∞ ¹        ~ (p=0.100 n=3) ²    2.069 ± ∞ ¹        ~ (p=0.100 n=3) ²     2.414 ± ∞ ¹       ~ (p=0.400 n=3) ²
Populate/PopulateDBWithTx-12               1069.4m ± ∞ ¹   1162.5m ± ∞ ¹       ~ (p=0.700 n=3) ²   916.0m ± ∞ ¹        ~ (p=0.100 n=3) ²   1399.2m ± ∞ ¹        ~ (p=0.100 n=3) ²   1099.8m ± ∞ ¹       ~ (p=0.700 n=3) ²   1422.1m ± ∞ ¹        ~ (p=0.100 n=3) ²   979.0m ± ∞ ¹        ~ (p=0.200 n=3) ²   1179.7m ± ∞ ¹       ~ (p=0.700 n=3) ²
Populate/PopulateDBWithTxs-12                1.061 ± ∞ ¹     1.102 ± ∞ ¹       ~ (p=0.400 n=3) ²    1.030 ± ∞ ¹        ~ (p=0.700 n=3) ²     1.324 ± ∞ ¹        ~ (p=0.100 n=3) ²     1.150 ± ∞ ¹       ~ (p=0.200 n=3) ²     1.330 ± ∞ ¹        ~ (p=0.100 n=3) ²    1.027 ± ∞ ¹        ~ (p=0.700 n=3) ²     1.106 ± ∞ ¹       ~ (p=0.100 n=3) ²
geomean                                      1.422           1.475        +3.70%                    1.243        -12.60%                     1.719        +20.88%                     1.410        -0.84%                     1.735        +22.03%                    1.277        -10.20%                     1.466        +3.08%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │         bench_ncruces_driver.txt          │       bench_eatonphil_direct.txt        │         bench_glebarez_driver.txt         │           bench_mattn_driver.txt            │         bench_modernc_driver.txt          │        bench_tailscale_driver.txt         │       bench_zombiezen_direct.txt        │
                              │           B/op           │      B/op       vs base                   │     B/op       vs base                  │      B/op       vs base                   │      B/op        vs base                    │      B/op       vs base                   │      B/op       vs base                   │     B/op       vs base                  │
Populate/PopulateDB-12                     2.362Mi ± ∞ ¹   18.968Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   5.771Mi ± ∞ ¹         ~ (p=0.100 n=3) ²   31.489Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   273.297Mi ± ∞ ¹           ~ (p=0.100 n=3) ²   31.477Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   18.963Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   6.167Mi ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTx-12               2.541Mi ± ∞ ¹   32.135Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   5.771Mi ± ∞ ¹         ~ (p=0.100 n=3) ²   44.542Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   286.581Mi ± ∞ ¹           ~ (p=0.100 n=3) ²   44.814Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   32.945Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   6.167Mi ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTxs-12              2.373Mi ± ∞ ¹   31.225Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   5.772Mi ± ∞ ¹         ~ (p=0.100 n=3) ²   43.665Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   285.969Mi ± ∞ ¹           ~ (p=0.100 n=3) ²   43.853Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   31.271Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   6.213Mi ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                    2.424Mi          26.70Mi        +1001.46%                   5.771Mi        +138.09%                    39.42Mi        +1526.13%                     281.9Mi        +11528.84%                    39.55Mi        +1531.57%                    26.93Mi        +1011.08%                   6.182Mi        +155.04%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │         bench_mattn_driver.txt          │        bench_modernc_driver.txt         │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                              │        allocs/op         │  allocs/op    vs base                  │  allocs/op    vs base                  │   allocs/op    vs base                  │   allocs/op    vs base                  │   allocs/op    vs base                  │  allocs/op    vs base                  │  allocs/op    vs base                  │
Populate/PopulateDB-12                      140.0k ± ∞ ¹   598.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   394.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1209.4k ± ∞ ¹         ~ (p=0.100 n=3) ²   1057.5k ± ∞ ¹         ~ (p=0.100 n=3) ²   1209.2k ± ∞ ¹         ~ (p=0.100 n=3) ²   598.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   445.1k ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTx-12                140.1k ± ∞ ¹   751.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   394.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1362.4k ± ∞ ¹         ~ (p=0.100 n=3) ²   1210.5k ± ∞ ¹         ~ (p=0.100 n=3) ²   1362.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   751.2k ± ∞ ¹         ~ (p=0.100 n=3) ²   445.1k ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTxs-12               140.0k ± ∞ ¹   765.2k ± ∞ ¹         ~ (p=0.100 n=3) ²   394.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1376.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   1238.9k ± ∞ ¹         ~ (p=0.100 n=3) ²   1380.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   767.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   448.1k ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                     140.1k         700.5k        +400.15%                   394.0k        +181.34%                    1.314M        +838.04%                    1.166M        +732.63%                    1.315M        +838.90%                   701.2k        +400.62%                   446.1k        +218.50%
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
                                                                    │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                                                                    │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                 │
ReadWrite/ReadPost-12                                                            13.914µ ± ∞ ¹   21.994µ ± ∞ ¹       ~ (p=0.100 n=3) ²   9.442µ ± ∞ ¹        ~ (p=0.100 n=3) ²   36.311µ ± ∞ ¹        ~ (p=0.100 n=3) ²   17.645µ ± ∞ ¹       ~ (p=0.100 n=3) ²   38.671µ ± ∞ ¹        ~ (p=0.100 n=3) ²   16.846µ ± ∞ ¹       ~ (p=0.100 n=3) ²   27.154µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadPostWithTx-12                                                      15.757µ ± ∞ ¹   37.727µ ± ∞ ¹       ~ (p=0.100 n=3) ²   9.957µ ± ∞ ¹        ~ (p=0.100 n=3) ²   39.493µ ± ∞ ¹        ~ (p=0.100 n=3) ²   21.886µ ± ∞ ¹       ~ (p=0.100 n=3) ²   42.550µ ± ∞ ¹        ~ (p=0.100 n=3) ²   24.356µ ± ∞ ¹       ~ (p=0.100 n=3) ²   27.257µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndComments-12                                                  93.62µ ± ∞ ¹   115.15µ ± ∞ ¹       ~ (p=0.100 n=3) ²   71.32µ ± ∞ ¹        ~ (p=0.100 n=3) ²   180.73µ ± ∞ ¹        ~ (p=0.100 n=3) ²   172.22µ ± ∞ ¹       ~ (p=0.100 n=3) ²   200.37µ ± ∞ ¹        ~ (p=0.100 n=3) ²    97.28µ ± ∞ ¹       ~ (p=0.100 n=3) ²   106.70µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndCommentsWithTx-12                                            93.08µ ± ∞ ¹   142.57µ ± ∞ ¹       ~ (p=0.100 n=3) ²   70.87µ ± ∞ ¹        ~ (p=0.100 n=3) ²   180.31µ ± ∞ ¹        ~ (p=0.100 n=3) ²   176.30µ ± ∞ ¹       ~ (p=0.100 n=3) ²   202.00µ ± ∞ ¹        ~ (p=0.100 n=3) ²   111.70µ ± ∞ ¹       ~ (p=0.100 n=3) ²   105.64µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePost-12                                                            141.1µ ± ∞ ¹    128.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²   122.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    149.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    122.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²    149.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    113.1µ ± ∞ ¹       ~ (p=0.100 n=3) ²    172.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePostWithTx-12                                                      148.6µ ± ∞ ¹    139.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²   124.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    152.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    130.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²    155.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    121.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²    175.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndComments-12                                                 2.633m ± ∞ ¹    2.577m ± ∞ ¹       ~ (p=0.100 n=3) ²   2.128m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.782m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.304m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.803m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.111m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.556m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndCommentsWithTx-12                                          1001.9µ ± ∞ ¹    982.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²   813.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1225.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    998.2µ ± ∞ ¹       ~ (p=0.700 n=3) ²   1225.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    821.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²   1092.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=10-12                             430.0µ ± ∞ ¹    354.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²   319.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    461.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    376.3µ ± ∞ ¹       ~ (p=0.200 n=3) ²    482.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    308.9µ ± ∞ ¹       ~ (p=0.100 n=3) ²    411.9µ ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=90-12                             2.376m ± ∞ ¹    2.343m ± ∞ ¹       ~ (p=0.200 n=3) ²   1.923m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.558m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.101m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.596m ± ∞ ¹        ~ (p=0.100 n=3) ²    1.948m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.333m ± ∞ ¹        ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-12                     402.9µ ± ∞ ¹    399.7µ ± ∞ ¹       ~ (p=1.000 n=3) ²   310.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    400.2µ ± ∞ ¹        ~ (p=1.000 n=3) ²    348.3µ ± ∞ ¹       ~ (p=0.200 n=3) ²    419.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    322.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²    335.2µ ± ∞ ¹        ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-12                     3.175m ± ∞ ¹    3.063m ± ∞ ¹       ~ (p=0.700 n=3) ²   2.429m ± ∞ ¹        ~ (p=0.200 n=3) ²    2.903m ± ∞ ¹        ~ (p=0.400 n=3) ²    2.488m ± ∞ ¹       ~ (p=0.200 n=3) ²    3.058m ± ∞ ¹        ~ (p=1.000 n=3) ²    2.393m ± ∞ ¹       ~ (p=0.100 n=3) ²    3.678m ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=10-12                       256.7µ ± ∞ ¹    219.9µ ± ∞ ¹       ~ (p=0.100 n=3) ²   191.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    296.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    261.9µ ± ∞ ¹       ~ (p=0.100 n=3) ²    307.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    177.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²    268.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=90-12                       916.7µ ± ∞ ¹    918.4µ ± ∞ ¹       ~ (p=0.700 n=3) ²   748.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1133.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    912.8µ ± ∞ ¹       ~ (p=0.700 n=3) ²   1145.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    753.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²   1014.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-12               114.2µ ± ∞ ¹    122.7µ ± ∞ ¹       ~ (p=0.700 n=3) ²   123.4µ ± ∞ ¹        ~ (p=0.700 n=3) ²    221.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    173.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²    238.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    127.6µ ± ∞ ¹       ~ (p=0.700 n=3) ²    192.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-12               890.5µ ± ∞ ¹    904.0µ ± ∞ ¹       ~ (p=0.700 n=3) ²   887.2µ ± ∞ ¹        ~ (p=0.700 n=3) ²   3002.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1164.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²   4098.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    881.8µ ± ∞ ¹       ~ (p=1.000 n=3) ²   3147.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                                                           295.1µ          323.0µ        +9.46%                   235.7µ        -20.12%                    424.4µ        +43.81%                    323.4µ        +9.58%                    451.0µ        +52.84%                    268.7µ        -8.95%                    372.7µ        +26.28%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                                                    │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │    bench_eatonphil_direct.txt    │     bench_glebarez_driver.txt     │       bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                                                                    │           B/op           │      B/op       vs base           │     B/op       vs base           │      B/op       vs base           │      B/op        vs base           │      B/op       vs base           │      B/op       vs base           │     B/op       vs base           │
ReadWrite/ReadPost-12                                                            40.17Ki ± ∞ ¹    41.23Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.18Ki ± ∞ ¹  ~ (p=0.400 n=3) ²    81.35Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     41.39Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    81.31Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    41.19Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.56Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostWithTx-12                                                      40.17Ki ± ∞ ¹    42.29Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.18Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    82.22Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     42.60Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    82.37Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    42.13Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.61Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndComments-12                                                 250.2Ki ± ∞ ¹    260.3Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.2Ki ± ∞ ¹  ~ (p=0.600 n=3) ²    501.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     263.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    501.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    257.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndCommentsWithTx-12                                           250.1Ki ± ∞ ¹    261.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    503.1Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     265.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    503.3Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    258.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePost-12                                                              0.00 ± ∞ ¹     240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     48.00 ± ∞ ¹  ~ (p=0.400 n=3) ²     496.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    42288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     496.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    409.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostWithTx-12                                                        0.00 ± ∞ ¹     859.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    1104.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    43253.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    1290.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     881.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    457.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndComments-12                                                0.000Ki ± ∞ ¹   13.906Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   2.781Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.661Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   308.876Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.673Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   13.906Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   3.532Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndCommentsWithTx-12                                          0.000Ki ± ∞ ¹   26.451Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   2.781Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   39.199Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   321.751Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   39.376Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.473Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   3.573Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=10-12                            226.2Ki ± ∞ ¹    236.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.6Ki ± ∞ ¹  ~ (p=1.000 n=3) ²    453.5Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     268.0Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    452.0Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    232.0Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   226.2Ki ± ∞ ¹  ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=90-12                            28.86Ki ± ∞ ¹    38.35Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   28.72Ki ± ∞ ¹  ~ (p=1.000 n=3) ²    73.58Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    303.76Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    77.87Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    36.87Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   28.55Ki ± ∞ ¹  ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-12                    225.2Ki ± ∞ ¹    236.3Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.1Ki ± ∞ ¹  ~ (p=0.400 n=3) ²    453.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     268.5Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    454.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    232.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   227.9Ki ± ∞ ¹  ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-12                    25.09Ki ± ∞ ¹    35.90Ki ± ∞ ¹  ~ (p=0.400 n=3) ²   27.52Ki ± ∞ ¹  ~ (p=1.000 n=3) ²    75.80Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    304.88Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    58.27Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    38.67Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   28.83Ki ± ∞ ¹  ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=10-12                      224.2Ki ± ∞ ¹    238.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   224.8Ki ± ∞ ¹  ~ (p=0.400 n=3) ²    455.0Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     271.0Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    457.4Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    235.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   224.9Ki ± ∞ ¹  ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=90-12                      25.57Ki ± ∞ ¹    51.03Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   28.99Ki ± ∞ ¹  ~ (p=0.200 n=3) ²    88.39Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    315.60Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    90.70Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    51.06Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   28.68Ki ± ∞ ¹  ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-12              225.5Ki ± ∞ ¹    238.3Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.6Ki ± ∞ ¹  ~ (p=0.700 n=3) ²    458.3Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     271.4Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    457.3Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    235.5Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.0Ki ± ∞ ¹  ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-12              27.68Ki ± ∞ ¹    50.41Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   29.31Ki ± ∞ ¹  ~ (p=0.700 n=3) ²    81.64Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    316.30Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    77.46Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    49.33Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   23.81Ki ± ∞ ¹  ~ (p=0.700 n=3) ²
geomean                                                                                      ³    42.08Ki        ?                   21.98Ki        ?                    77.20Ki        ?                     178.4Ki        ?                    76.86Ki        ?                    41.95Ki        ?                   29.57Ki        ?
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
ReadWrite/ReadOrWritePostAndComments/write_rate=10-12                              236.0 ± ∞ ¹    737.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    696.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=90-12                              30.00 ± ∞ ¹   443.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   800.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   421.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-12                      235.0 ± ∞ ¹    737.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    697.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-12                      26.00 ± ∞ ¹   439.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   803.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   422.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=10-12                        234.0 ± ∞ ¹    778.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1011.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    745.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    274.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=90-12                        26.00 ± ∞ ¹   597.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1119.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   965.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1123.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-12                236.0 ± ∞ ¹    778.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1010.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    747.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-12                28.00 ± ∞ ¹   597.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1122.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   969.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1127.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                                                      ³    273.4        ?                    85.22        ?                     431.1        ?                    368.2        ?                     436.0        ?                    237.4        ?                    122.4        ?
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
                                                                    │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                                                                    │          sec/op          │    sec/op     vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10                        435.8µ ± ∞ ¹   350.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²   314.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    442.6µ ± ∞ ¹        ~ (p=1.000 n=3) ²    346.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²    456.7µ ± ∞ ¹        ~ (p=0.700 n=3) ²   288.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    423.5µ ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-2                      359.5µ ± ∞ ¹   325.5µ ± ∞ ¹       ~ (p=0.200 n=3) ²   299.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    390.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    312.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²    403.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   275.0µ ± ∞ ¹        ~ (p=0.700 n=3) ²    361.8µ ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-4                      349.4µ ± ∞ ¹   336.2µ ± ∞ ¹       ~ (p=1.000 n=3) ²   268.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    368.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    318.5µ ± ∞ ¹       ~ (p=0.200 n=3) ²    386.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   263.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    344.0µ ± ∞ ¹        ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-8                      351.9µ ± ∞ ¹   368.5µ ± ∞ ¹       ~ (p=0.200 n=3) ²   275.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    389.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    321.9µ ± ∞ ¹       ~ (p=0.200 n=3) ²    413.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   278.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    352.1µ ± ∞ ¹        ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-16                     438.5µ ± ∞ ¹   424.8µ ± ∞ ¹       ~ (p=1.000 n=3) ²   369.4µ ± ∞ ¹        ~ (p=0.700 n=3) ²    448.5µ ± ∞ ¹        ~ (p=0.400 n=3) ²    352.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²    432.3µ ± ∞ ¹        ~ (p=1.000 n=3) ²   309.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    395.8µ ± ∞ ¹        ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90                        2.428m ± ∞ ¹   2.399m ± ∞ ¹       ~ (p=1.000 n=3) ²   1.955m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.430m ± ∞ ¹        ~ (p=0.700 n=3) ²    2.037m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.536m ± ∞ ¹        ~ (p=0.100 n=3) ²   1.954m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.295m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-2                      2.538m ± ∞ ¹   2.542m ± ∞ ¹       ~ (p=0.700 n=3) ²   1.999m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.630m ± ∞ ¹        ~ (p=0.700 n=3) ²    2.056m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.738m ± ∞ ¹        ~ (p=0.200 n=3) ²   1.966m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.345m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-4                      2.634m ± ∞ ¹   2.596m ± ∞ ¹       ~ (p=1.000 n=3) ²   2.023m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.593m ± ∞ ¹        ~ (p=0.400 n=3) ²    2.102m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.759m ± ∞ ¹        ~ (p=0.100 n=3) ²   2.123m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.360m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-8                      2.848m ± ∞ ¹   2.905m ± ∞ ¹       ~ (p=0.700 n=3) ²   2.265m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.727m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.292m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.674m ± ∞ ¹        ~ (p=0.200 n=3) ²   2.226m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.412m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-16                     3.247m ± ∞ ¹   3.388m ± ∞ ¹       ~ (p=0.100 n=3) ²   2.547m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.892m ± ∞ ¹        ~ (p=0.700 n=3) ²    2.645m ± ∞ ¹       ~ (p=0.100 n=3) ²    5.238m ± ∞ ¹        ~ (p=0.100 n=3) ²   2.790m ± ∞ ¹        ~ (p=0.100 n=3) ²   11.591m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10                  269.3µ ± ∞ ¹   222.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²   192.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    293.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    232.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²    295.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   184.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    278.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-2                159.8µ ± ∞ ¹   138.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²   140.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    215.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    167.6µ ± ∞ ¹       ~ (p=0.100 n=3) ²    219.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   125.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    199.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-4                119.9µ ± ∞ ¹   119.2µ ± ∞ ¹       ~ (p=1.000 n=3) ²   116.0µ ± ∞ ¹        ~ (p=1.000 n=3) ²    189.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    143.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²    191.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   116.6µ ± ∞ ¹        ~ (p=1.000 n=3) ²    177.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-8                110.9µ ± ∞ ¹   110.2µ ± ∞ ¹       ~ (p=0.700 n=3) ²   118.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    212.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    175.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²    208.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   127.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    179.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-16               123.2µ ± ∞ ¹   128.1µ ± ∞ ¹       ~ (p=0.400 n=3) ²   129.5µ ± ∞ ¹        ~ (p=0.700 n=3) ²    241.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    215.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²    239.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   144.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    194.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90                  945.0µ ± ∞ ¹   933.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²   781.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1134.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    872.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²   1129.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   754.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1023.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-2                818.9µ ± ∞ ¹   787.1µ ± ∞ ¹       ~ (p=0.200 n=3) ²   756.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1165.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    886.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²   1177.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   748.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    970.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-4                781.5µ ± ∞ ¹   805.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²   775.6µ ± ∞ ¹        ~ (p=1.000 n=3) ²   1199.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1014.6µ ± ∞ ¹       ~ (p=0.100 n=3) ²   1188.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   774.9µ ± ∞ ¹        ~ (p=0.400 n=3) ²    998.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-8                829.5µ ± ∞ ¹   868.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²   833.0µ ± ∞ ¹        ~ (p=1.000 n=3) ²   1200.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1129.1µ ± ∞ ¹       ~ (p=0.100 n=3) ²   1935.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   857.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1067.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-16               1.019m ± ∞ ¹   1.080m ± ∞ ¹       ~ (p=0.200 n=3) ²   1.080m ± ∞ ¹        ~ (p=0.700 n=3) ²   10.169m ± ∞ ¹        ~ (p=0.100 n=3) ²    1.216m ± ∞ ¹       ~ (p=0.100 n=3) ²   10.519m ± ∞ ¹        ~ (p=0.100 n=3) ²   1.002m ± ∞ ¹        ~ (p=0.700 n=3) ²   10.602m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                                                           606.7µ         590.8µ        -2.62%                   522.6µ        -13.87%                    816.2µ        +34.53%                    607.9µ        +0.19%                    874.3µ        +44.11%                   514.5µ        -15.19%                    794.6µ        +30.97%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                                                    │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │          bench_mattn_driver.txt          │        bench_modernc_driver.txt         │       bench_tailscale_driver.txt       │      bench_zombiezen_direct.txt       │
                                                                    │           B/op           │     B/op       vs base                 │     B/op       vs base                │     B/op       vs base                  │      B/op       vs base                  │     B/op       vs base                  │     B/op       vs base                 │     B/op       vs base                │
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10                       224.5Ki ± ∞ ¹   234.8Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.1Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   453.7Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   451.7Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   234.0Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.8Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-2                     224.3Ki ± ∞ ¹   236.0Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.1Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   455.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   452.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   231.6Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.9Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-4                     223.9Ki ± ∞ ¹   235.1Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.9Ki ± ∞ ¹       ~ (p=0.200 n=3) ²   454.6Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.7Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   453.0Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   232.0Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.6Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-8                     225.6Ki ± ∞ ¹   237.5Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.0Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   454.9Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   452.0Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   232.7Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.4Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-16                    224.1Ki ± ∞ ¹   236.8Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.7Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   453.9Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   455.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   232.2Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.3Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90                       23.63Ki ± ∞ ¹   36.06Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.37Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   73.19Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.92Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   74.97Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   41.38Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   30.75Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-2                     21.73Ki ± ∞ ¹   40.08Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   25.16Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   77.77Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.07Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   72.56Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   35.53Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   26.36Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-4                     25.45Ki ± ∞ ¹   39.08Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.17Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   73.95Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.09Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   68.39Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   36.97Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.58Ki ± ∞ ¹       ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-8                     26.11Ki ± ∞ ¹   39.28Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   25.86Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   78.24Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.85Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   80.09Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   35.71Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.12Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-16                    27.96Ki ± ∞ ¹   37.56Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   24.07Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   74.62Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.09Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   77.47Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   36.79Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   33.23Ki ± ∞ ¹       ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10                 225.5Ki ± ∞ ¹   237.3Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.6Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   458.0Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   457.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   234.4Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.1Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-2               225.6Ki ± ∞ ¹   238.5Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.1Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   456.5Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   455.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   235.3Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.9Ki ± ∞ ¹       ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-4               225.6Ki ± ∞ ¹   237.7Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.2Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   456.6Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   456.6Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   234.4Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.6Ki ± ∞ ¹       ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-8               225.2Ki ± ∞ ¹   238.7Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.7Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   456.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   457.0Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   235.4Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.2Ki ± ∞ ¹       ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-16              224.9Ki ± ∞ ¹   239.1Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.5Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   456.9Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.5Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   456.5Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   235.3Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.9Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90                 23.92Ki ± ∞ ¹   49.71Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   26.33Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   83.66Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.67Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   85.56Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   50.51Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.16Ki ± ∞ ¹       ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-2               23.94Ki ± ∞ ¹   49.55Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   26.28Ki ± ∞ ¹       ~ (p=0.200 n=3) ²   90.27Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.39Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   92.76Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   50.09Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.60Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-4               23.81Ki ± ∞ ¹   50.07Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   25.99Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   87.00Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.74Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   84.37Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   50.46Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.09Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-8               24.70Ki ± ∞ ¹   51.83Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.48Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   84.36Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.13Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   83.98Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   49.82Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.39Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-16              25.74Ki ± ∞ ¹   48.50Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   26.72Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   85.67Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.09Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   82.79Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   49.62Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   25.88Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
geomean                                                                          74.45Ki         101.9Ki        +36.84%                   77.04Ki        +3.47%                   191.7Ki        +157.51%                    289.4Ki        +288.68%                   190.7Ki        +156.17%                   100.5Ki        +34.94%                   80.16Ki        +7.67%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                                                    │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt         │         bench_mattn_driver.txt         │         bench_modernc_driver.txt         │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                                                                    │        allocs/op         │  allocs/op    vs base                  │  allocs/op    vs base                  │   allocs/op    vs base                   │  allocs/op    vs base                  │   allocs/op    vs base                   │  allocs/op    vs base                  │  allocs/op    vs base                  │
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10                         235.0 ± ∞ ¹    733.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     973.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    695.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     973.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    545.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    270.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-2                       234.0 ± ∞ ¹    735.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    696.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    545.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    270.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-4                       234.0 ± ∞ ¹    735.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    698.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-8                       236.0 ± ∞ ¹    739.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    697.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-16                      234.0 ± ∞ ¹    737.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    697.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90                         24.00 ± ∞ ¹   439.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   803.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   424.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-2                       22.00 ± ∞ ¹   445.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   801.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   420.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-4                       26.00 ± ∞ ¹   444.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   801.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   421.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-8                       27.00 ± ∞ ¹   444.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   803.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   420.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-16                      29.00 ± ∞ ¹   442.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   801.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   421.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   261.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10                   236.0 ± ∞ ¹    776.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1009.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1013.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    585.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-2                 236.0 ± ∞ ¹    777.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1010.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    585.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-4                 236.0 ± ∞ ¹    777.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1010.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    745.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    585.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-8                 235.0 ± ∞ ¹    778.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1011.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    747.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-16                235.0 ± ∞ ¹    779.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1011.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    748.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    274.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90                   25.00 ± ∞ ¹   596.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1120.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   971.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1124.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-2                 25.00 ± ∞ ¹   596.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1118.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   969.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1122.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-4                 24.00 ± ∞ ¹   596.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1120.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   971.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1124.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-8                 25.00 ± ∞ ¹   598.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1120.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   968.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1125.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-16                27.00 ± ∞ ¹   595.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1121.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   968.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1125.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                                                            77.17          623.4        +707.77%                    231.6        +200.10%                    1.016k        +1216.90%                    797.3        +933.16%                    1.018k        +1219.17%                    527.8        +583.97%                    266.7        +245.58%
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
Query/Correlated-12                       643.5m ± ∞ ¹    643.2m ± ∞ ¹       ~ (p=1.000 n=3) ²   287.3m ± ∞ ¹        ~ (p=0.100 n=3) ²    303.1m ± ∞ ¹        ~ (p=0.100 n=3) ²   176.0m ± ∞ ¹        ~ (p=0.100 n=3) ²    302.7m ± ∞ ¹        ~ (p=0.100 n=3) ²   170.7m ± ∞ ¹        ~ (p=0.100 n=3) ²    324.9m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-12              119.49m ± ∞ ¹   118.81m ± ∞ ¹       ~ (p=0.100 n=3) ²   52.45m ± ∞ ¹        ~ (p=0.100 n=3) ²   116.03m ± ∞ ¹        ~ (p=0.100 n=3) ²   44.96m ± ∞ ¹        ~ (p=0.100 n=3) ²   113.52m ± ∞ ¹        ~ (p=0.100 n=3) ²   45.14m ± ∞ ¹        ~ (p=0.100 n=3) ²   184.84m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                   277.3m          276.4m        -0.31%                   122.8m        -55.73%                    187.5m        -32.37%                   88.97m        -67.92%                    185.4m        -33.15%                   87.79m        -68.34%                    245.1m        -11.62%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │         bench_ncruces_driver.txt         │      bench_eatonphil_direct.txt      │        bench_glebarez_driver.txt         │          bench_mattn_driver.txt           │         bench_modernc_driver.txt         │       bench_tailscale_driver.txt        │       bench_zombiezen_direct.txt       │
                            │           B/op           │     B/op       vs base                   │    B/op      vs base                 │     B/op       vs base                   │      B/op       vs base                   │     B/op       vs base                   │     B/op       vs base                  │     B/op       vs base                 │
Query/Correlated-12                      57168.0 ± ∞ ¹   71168.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   736.0 ± ∞ ¹        ~ (p=0.100 n=3) ²   63120.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   207200.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   63136.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   47053.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    9328.0 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-12                342.0 ± ∞ ¹   71465.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   318.0 ± ∞ ¹        ~ (p=1.000 n=3) ²   63237.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   207676.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   63214.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   48502.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1622.0 ± ∞ ¹        ~ (p=0.700 n=3) ²
geomean                                  4.318Ki         69.64Ki        +1512.87%                   483.8        -89.06%                   61.70Ki        +1328.83%                    202.6Ki        +4591.36%                   61.69Ki        +1328.75%                   46.65Ki        +980.40%                   3.799Ki        -12.03%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt           │      bench_eatonphil_direct.txt      │          bench_glebarez_driver.txt          │           bench_mattn_driver.txt            │          bench_modernc_driver.txt           │         bench_tailscale_driver.txt         │       bench_zombiezen_direct.txt       │
                            │        allocs/op         │   allocs/op     vs base                     │  allocs/op   vs base                 │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                    │  allocs/op    vs base                  │
Query/Correlated-12                       12.000 ± ∞ ¹   5769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹           ~ (p=0.100 n=3) ²   30.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-12                2.000 ± ∞ ¹   5770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹        ~ (p=0.700 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6773.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3769.000 ± ∞ ¹           ~ (p=0.100 n=3) ²   15.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                    4.899           5.769k        +117669.43%                   3.464        -29.29%                     6.770k        +138102.25%                     6.772k        +138132.87%                     6.770k        +138102.25%                     3.768k        +76813.98%                    21.21        +333.01%
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
Query/CorrelatedParallel                  645.3m ± ∞ ¹    644.0m ± ∞ ¹       ~ (p=0.700 n=3) ²    287.8m ± ∞ ¹        ~ (p=0.100 n=3) ²    302.3m ± ∞ ¹        ~ (p=0.100 n=3) ²   175.1m ± ∞ ¹        ~ (p=0.100 n=3) ²    301.6m ± ∞ ¹        ~ (p=0.100 n=3) ²   172.5m ± ∞ ¹        ~ (p=0.100 n=3) ²    324.4m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-2               435.13m ± ∞ ¹   434.05m ± ∞ ¹       ~ (p=0.400 n=3) ²   171.47m ± ∞ ¹        ~ (p=0.100 n=3) ²   183.27m ± ∞ ¹        ~ (p=0.100 n=3) ²   96.54m ± ∞ ¹        ~ (p=0.100 n=3) ²   183.34m ± ∞ ¹        ~ (p=0.100 n=3) ²   95.95m ± ∞ ¹        ~ (p=0.100 n=3) ²   188.03m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-4               225.63m ± ∞ ¹   225.12m ± ∞ ¹       ~ (p=0.400 n=3) ²    85.29m ± ∞ ¹        ~ (p=0.100 n=3) ²   116.55m ± ∞ ¹        ~ (p=0.100 n=3) ²   59.55m ± ∞ ¹        ~ (p=0.100 n=3) ²   116.59m ± ∞ ¹        ~ (p=0.100 n=3) ²   55.89m ± ∞ ¹        ~ (p=0.100 n=3) ²   124.69m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-8               154.79m ± ∞ ¹   155.58m ± ∞ ¹       ~ (p=0.100 n=3) ²    55.06m ± ∞ ¹        ~ (p=0.100 n=3) ²   137.81m ± ∞ ¹        ~ (p=0.100 n=3) ²   43.02m ± ∞ ¹        ~ (p=0.100 n=3) ²   134.69m ± ∞ ¹        ~ (p=0.100 n=3) ²   42.46m ± ∞ ¹        ~ (p=0.100 n=3) ²   186.16m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-16               78.16m ± ∞ ¹    77.95m ± ∞ ¹       ~ (p=1.000 n=3) ²    52.51m ± ∞ ¹        ~ (p=0.100 n=3) ²   115.75m ± ∞ ¹        ~ (p=0.100 n=3) ²   44.66m ± ∞ ¹        ~ (p=0.100 n=3) ²   112.08m ± ∞ ¹        ~ (p=0.100 n=3) ²   48.39m ± ∞ ¹        ~ (p=0.100 n=3) ²   184.82m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                   238.2m          238.0m        -0.09%                    104.0m        -56.33%                    159.4m        -33.06%                   71.99m        -69.77%                    157.6m        -33.82%                   71.74m        -69.88%                    192.1m        -19.34%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │         bench_ncruces_driver.txt          │      bench_eatonphil_direct.txt       │         bench_glebarez_driver.txt         │           bench_mattn_driver.txt           │         bench_modernc_driver.txt          │        bench_tailscale_driver.txt        │       bench_zombiezen_direct.txt       │
                            │           B/op           │      B/op       vs base                   │     B/op      vs base                 │      B/op       vs base                   │      B/op        vs base                   │      B/op       vs base                   │      B/op       vs base                  │     B/op       vs base                 │
Query/CorrelatedParallel                 57140.0 ± ∞ ¹    71132.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    718.0 ± ∞ ¹        ~ (p=0.100 n=3) ²    63050.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    207185.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    63066.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    47009.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    2216.0 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-2              57125.00 ± ∞ ¹   71146.00 ± ∞ ¹          ~ (p=0.100 n=3) ²    60.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   63013.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   207152.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   63068.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   46978.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    460.00 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-4              11729.00 ± ∞ ¹   71044.00 ± ∞ ¹          ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   63092.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   207231.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   63049.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   46987.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    500.00 ± ∞ ¹        ~ (p=0.700 n=3) ²
Query/CorrelatedParallel-8                150.00 ± ∞ ¹   71048.00 ± ∞ ¹          ~ (p=0.100 n=3) ²    81.00 ± ∞ ¹        ~ (p=0.400 n=3) ²   63181.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   207274.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   63234.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   47202.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1276.00 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-16                601.0 ± ∞ ¹    71504.0 ± ∞ ¹          ~ (p=0.100 n=3) ²   1070.0 ± ∞ ¹        ~ (p=0.700 n=3) ²    63349.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    207756.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    63264.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    48882.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1065.0 ± ∞ ¹        ~ (p=0.200 n=3) ²
geomean                                  4.981Ki          69.51Ki        +1295.49%                    178.1        -96.51%                    61.66Ki        +1137.90%                     202.5Ki        +3964.82%                    61.66Ki        +1137.88%                    46.29Ki        +829.47%                     929.2        -81.78%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt           │      bench_eatonphil_direct.txt      │          bench_glebarez_driver.txt          │           bench_mattn_driver.txt            │          bench_modernc_driver.txt           │         bench_tailscale_driver.txt         │       bench_zombiezen_direct.txt       │
                            │        allocs/op         │   allocs/op     vs base                     │  allocs/op   vs base                 │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                    │  allocs/op    vs base                  │
Query/CorrelatedParallel                  14.000 ± ∞ ¹   5771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   7.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹           ~ (p=0.100 n=3) ²   30.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-2                14.000 ± ∞ ¹   5771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹           ~ (p=0.100 n=3) ²    7.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-4                 4.000 ± ∞ ¹   5769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹           ~ (p=0.100 n=3) ²    7.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-8                 2.000 ± ∞ ¹   5770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹           ~ (p=0.100 n=3) ²   11.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-16                4.000 ± ∞ ¹   5772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹        ~ (p=0.200 n=3) ²   6774.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6775.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6773.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3771.000 ± ∞ ¹           ~ (p=0.100 n=3) ²   13.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                    5.748           5.771k        +100301.40%                   1.695        -70.51%                     6.771k        +117703.65%                     6.772k        +117731.49%                     6.771k        +117703.66%                     3.767k        +65448.15%                    11.60        +101.86%
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
Query/GroupBy-12                       569.1µ ± ∞ ¹   577.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²   276.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    488.2µ ± ∞ ¹         ~ (p=0.100 n=3) ²    334.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    503.2µ ± ∞ ¹         ~ (p=0.100 n=3) ²   261.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    479.6µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-12               61.93µ ± ∞ ¹   64.55µ ± ∞ ¹       ~ (p=0.100 n=3) ²   37.60µ ± ∞ ¹        ~ (p=0.100 n=3) ²   685.36µ ± ∞ ¹         ~ (p=0.100 n=3) ²   403.39µ ± ∞ ¹        ~ (p=0.100 n=3) ²   693.15µ ± ∞ ¹         ~ (p=0.100 n=3) ²   27.99µ ± ∞ ¹        ~ (p=0.100 n=3) ²   695.49µ ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                187.7µ         193.0µ        +2.81%                   101.9µ        -45.69%                    578.4µ        +208.13%                    367.3µ        +95.65%                    590.6µ        +214.59%                   85.54µ        -54.44%                    577.5µ        +207.66%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │    bench_eatonphil_direct.txt     │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │    bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                         │           B/op           │     B/op       vs base           │   B/op     vs base                │     B/op       vs base           │     B/op       vs base           │     B/op       vs base           │     B/op       vs base           │    B/op      vs base           │
Query/GroupBy-12                          0.0 ± ∞ ¹    2381.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ²    1964.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    7164.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1976.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1582.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   364.0 ± ∞ ¹  ~ (p=0.600 n=3) ²
Query/GroupByParallel-12                  0.0 ± ∞ ¹    2368.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1966.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    7176.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1986.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1528.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   370.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ⁴   2.319Ki        ?                              +0.00%               ⁴   1.919Ki        ?                   7.002Ki        ?                   1.935Ki        ?                   1.518Ki        ?                   367.0        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

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
                         │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │         bench_mattn_driver.txt          │        bench_modernc_driver.txt         │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt        │
                         │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op     vs base                 │    sec/op      vs base                  │
Query/GroupByParallel                  561.9µ ± ∞ ¹    558.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²   287.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    482.6µ ± ∞ ¹         ~ (p=0.100 n=3) ²    301.2µ ± ∞ ¹         ~ (p=0.100 n=3) ²    479.9µ ± ∞ ¹         ~ (p=0.100 n=3) ²   263.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    466.2µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-2                283.9µ ± ∞ ¹    285.1µ ± ∞ ¹       ~ (p=0.100 n=3) ²   145.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    400.5µ ± ∞ ¹         ~ (p=0.100 n=3) ²    371.6µ ± ∞ ¹         ~ (p=0.100 n=3) ²    389.8µ ± ∞ ¹         ~ (p=0.100 n=3) ²   130.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    383.5µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-4               145.09µ ± ∞ ¹   147.25µ ± ∞ ¹       ~ (p=0.100 n=3) ²   79.72µ ± ∞ ¹        ~ (p=0.100 n=3) ²   392.92µ ± ∞ ¹         ~ (p=0.100 n=3) ²   448.94µ ± ∞ ¹         ~ (p=0.100 n=3) ²   393.93µ ± ∞ ¹         ~ (p=0.100 n=3) ²   65.75µ ± ∞ ¹        ~ (p=0.100 n=3) ²   392.99µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-8                73.08µ ± ∞ ¹    75.64µ ± ∞ ¹       ~ (p=0.100 n=3) ²   43.00µ ± ∞ ¹        ~ (p=0.100 n=3) ²   655.48µ ± ∞ ¹         ~ (p=0.100 n=3) ²   463.76µ ± ∞ ¹         ~ (p=0.100 n=3) ²   658.65µ ± ∞ ¹         ~ (p=0.100 n=3) ²   33.92µ ± ∞ ¹        ~ (p=0.100 n=3) ²   658.33µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-16               61.01µ ± ∞ ¹    64.37µ ± ∞ ¹       ~ (p=0.100 n=3) ²   37.02µ ± ∞ ¹        ~ (p=0.100 n=3) ²   691.65µ ± ∞ ¹         ~ (p=0.100 n=3) ²   409.47µ ± ∞ ¹         ~ (p=0.100 n=3) ²   698.02µ ± ∞ ¹         ~ (p=0.100 n=3) ²   28.34µ ± ∞ ¹        ~ (p=0.100 n=3) ²   694.41µ ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                159.5µ          162.7µ        +2.03%                   88.07µ        -44.78%                    509.8µ        +219.64%                    394.4µ        +147.28%                    508.1µ        +218.61%                   73.75µ        -53.76%                    502.8µ        +215.23%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │    bench_eatonphil_direct.txt     │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │    bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                         │           B/op           │     B/op       vs base           │   B/op     vs base                │     B/op       vs base           │     B/op       vs base           │     B/op       vs base           │     B/op       vs base           │    B/op      vs base           │
Query/GroupByParallel                     0.0 ± ∞ ¹    2363.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1849.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    7160.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1864.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1581.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   379.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-2                   0.0 ± ∞ ¹    2363.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1960.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    7160.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1900.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1580.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   360.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-4                   0.0 ± ∞ ¹    2364.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1961.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    7134.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1955.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1541.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   361.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-8                   0.0 ± ∞ ¹    2366.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1965.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    7158.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1983.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1565.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   363.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-16                  0.0 ± ∞ ¹    2368.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.0 ± ∞ ¹       ~ (p=1.000 n=3) ³    1969.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    7176.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1989.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1514.0 ± ∞ ¹  ~ (p=0.100 n=3) ²   366.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ⁴   2.309Ki        ?                              +0.00%               ⁴   1.895Ki        ?                   6.990Ki        ?                   1.892Ki        ?                   1.520Ki        ?                   365.7        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                         │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/GroupByParallel                   0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   120.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   155.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   120.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   52.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-2                 0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   155.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   121.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-4                 0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   154.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   121.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-8                 0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   154.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
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
                      │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                      │          sec/op          │    sec/op     vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │
Query/JSON-12                       9.791m ± ∞ ¹   9.935m ± ∞ ¹       ~ (p=0.700 n=3) ²   7.208m ± ∞ ¹        ~ (p=0.100 n=3) ²    9.218m ± ∞ ¹        ~ (p=0.100 n=3) ²   8.203m ± ∞ ¹        ~ (p=0.100 n=3) ²    9.179m ± ∞ ¹        ~ (p=0.100 n=3) ²   7.132m ± ∞ ¹        ~ (p=0.100 n=3) ²   10.051m ± ∞ ¹        ~ (p=0.400 n=3) ²
Query/JSONParallel-12               9.626m ± ∞ ¹   9.409m ± ∞ ¹       ~ (p=0.100 n=3) ²   9.522m ± ∞ ¹        ~ (p=0.100 n=3) ²   15.111m ± ∞ ¹        ~ (p=0.100 n=3) ²   8.782m ± ∞ ¹        ~ (p=0.100 n=3) ²   14.855m ± ∞ ¹        ~ (p=0.100 n=3) ²   9.570m ± ∞ ¹        ~ (p=0.400 n=3) ²   12.971m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                             9.708m         9.668m        -0.41%                   8.285m        -14.66%                    11.80m        +21.58%                   8.488m        -12.57%                    11.68m        +20.28%                   8.261m        -14.90%                    11.42m        +17.62%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                      │ bench_ncruces_direct.txt │           bench_ncruces_driver.txt           │      bench_eatonphil_direct.txt      │          bench_glebarez_driver.txt           │             bench_mattn_driver.txt             │           bench_modernc_driver.txt           │          bench_tailscale_driver.txt          │        bench_zombiezen_direct.txt        │
                      │           B/op           │      B/op        vs base                     │    B/op      vs base                 │      B/op        vs base                     │       B/op        vs base                      │      B/op        vs base                     │      B/op        vs base                     │     B/op       vs base                   │
Query/JSON-12                        2.000 ± ∞ ¹   64902.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.300 n=3) ²   56949.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   201045.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56966.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   32892.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   409.000 ± ∞ ¹          ~ (p=0.600 n=3) ²
Query/JSONParallel-12                49.00 ± ∞ ¹    64880.00 ± ∞ ¹            ~ (p=0.100 n=3) ²   34.00 ± ∞ ¹        ~ (p=0.700 n=3) ²    57009.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    201145.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    57036.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    33000.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    472.00 ± ∞ ¹          ~ (p=0.100 n=3) ²
geomean                              9.899           63.37Ki        +655398.08%                   5.831        -41.10%                     55.64Ki        +575474.74%                      196.4Ki        +2031266.20%                     55.67Ki        +575696.94%                     32.17Ki        +332704.41%                     439.4        +4338.33%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                      │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │   bench_zombiezen_direct.txt   │
                      │        allocs/op         │   allocs/op     vs base           │  allocs/op   vs base                │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │  allocs/op   vs base           │
Query/JSON-12                        0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-12                0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4023.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
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
                      │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                      │          sec/op          │    sec/op     vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │
Query/JSONParallel                  9.824m ± ∞ ¹   9.744m ± ∞ ¹       ~ (p=0.200 n=3) ²   7.220m ± ∞ ¹        ~ (p=0.100 n=3) ²    9.366m ± ∞ ¹        ~ (p=0.100 n=3) ²   7.545m ± ∞ ¹        ~ (p=0.100 n=3) ²    9.206m ± ∞ ¹        ~ (p=0.100 n=3) ²   6.987m ± ∞ ¹        ~ (p=0.100 n=3) ²    9.916m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-2                7.043m ± ∞ ¹   7.072m ± ∞ ¹       ~ (p=1.000 n=3) ²   5.594m ± ∞ ¹        ~ (p=0.100 n=3) ²    7.359m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.919m ± ∞ ¹        ~ (p=0.100 n=3) ²    7.295m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.742m ± ∞ ¹        ~ (p=0.100 n=3) ²    8.562m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-4                4.803m ± ∞ ¹   4.789m ± ∞ ¹       ~ (p=1.000 n=3) ²   4.087m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.418m ± ∞ ¹        ~ (p=0.100 n=3) ²   4.206m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.467m ± ∞ ¹        ~ (p=0.100 n=3) ²   4.337m ± ∞ ¹        ~ (p=0.100 n=3) ²    7.537m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-8                7.424m ± ∞ ¹   7.579m ± ∞ ¹       ~ (p=0.100 n=3) ²   6.938m ± ∞ ¹        ~ (p=0.100 n=3) ²   13.727m ± ∞ ¹        ~ (p=0.100 n=3) ²   6.688m ± ∞ ¹        ~ (p=0.100 n=3) ²   13.448m ± ∞ ¹        ~ (p=0.100 n=3) ²   7.413m ± ∞ ¹        ~ (p=0.400 n=3) ²   14.733m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-16               9.524m ± ∞ ¹   9.651m ± ∞ ¹       ~ (p=0.100 n=3) ²   9.077m ± ∞ ¹        ~ (p=0.100 n=3) ²   16.072m ± ∞ ¹        ~ (p=0.100 n=3) ²   9.112m ± ∞ ¹        ~ (p=0.100 n=3) ²   15.214m ± ∞ ¹        ~ (p=0.100 n=3) ²   9.702m ± ∞ ¹        ~ (p=0.100 n=3) ²   12.915m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                             7.485m         7.526m        +0.54%                   6.359m        -15.05%                    9.620m        +28.52%                   6.483m        -13.40%                    9.444m        +26.17%                   6.599m        -11.84%                    10.40m        +38.96%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                      │ bench_ncruces_direct.txt │           bench_ncruces_driver.txt           │      bench_eatonphil_direct.txt       │          bench_glebarez_driver.txt           │             bench_mattn_driver.txt             │           bench_modernc_driver.txt           │          bench_tailscale_driver.txt          │        bench_zombiezen_direct.txt        │
                      │           B/op           │      B/op        vs base                     │     B/op      vs base                 │      B/op        vs base                     │       B/op        vs base                      │      B/op        vs base                     │      B/op        vs base                     │     B/op       vs base                   │
Query/JSONParallel                   2.000 ± ∞ ¹   64845.000 ± ∞ ¹            ~ (p=0.100 n=3) ²    1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   56877.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   201019.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56860.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   32891.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   364.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-2                 2.000 ± ∞ ¹   64844.000 ± ∞ ¹            ~ (p=0.100 n=3) ²    1.000 ± ∞ ¹        ~ (p=0.400 n=3) ²   56950.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   201008.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56908.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   32885.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   371.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-4                 7.000 ± ∞ ¹   64844.000 ± ∞ ¹            ~ (p=0.100 n=3) ²    5.000 ± ∞ ¹        ~ (p=0.700 n=3) ²   56958.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   201006.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56940.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   32902.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   380.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-8                 5.000 ± ∞ ¹   64850.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   25.000 ± ∞ ¹        ~ (p=1.000 n=3) ²   56997.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   201051.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   57006.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   32970.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   391.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/JSONParallel-16                96.00 ± ∞ ¹    64954.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    41.00 ± ∞ ¹        ~ (p=0.100 n=3) ²    57048.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    201206.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    57012.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    32973.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    513.00 ± ∞ ¹          ~ (p=0.100 n=3) ²
geomean                              6.694           63.35Ki        +968950.79%                    5.520        -17.54%                     55.63Ki        +850911.95%                      196.3Ki        +3003495.70%                     55.61Ki        +850601.20%                     32.15Ki        +491752.71%                     400.4        +5881.78%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                      │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │   bench_zombiezen_direct.txt   │
                      │        allocs/op         │   allocs/op     vs base           │  allocs/op   vs base                │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │  allocs/op   vs base           │
Query/JSONParallel                   0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-2                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-4                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-8                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-16                0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4023.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                        ⁴     4.019k        ?                                +0.00%               ⁴     4.022k        ?                     5.021k        ?                     4.021k        ?                     2.020k        ?                   6.000        ?
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
Query/OrderBy-12                       76.25m ± ∞ ¹   89.37m ± ∞ ¹       ~ (p=0.100 n=3) ²   46.46m ± ∞ ¹        ~ (p=0.100 n=3) ²   80.84m ± ∞ ¹        ~ (p=0.100 n=3) ²   108.58m ± ∞ ¹        ~ (p=0.100 n=3) ²   80.77m ± ∞ ¹        ~ (p=0.100 n=3) ²   61.00m ± ∞ ¹       ~ (p=0.100 n=3) ²   83.14m ± ∞ ¹        ~ (p=0.700 n=3) ²
Query/OrderByParallel-12               48.40m ± ∞ ¹   49.05m ± ∞ ¹       ~ (p=0.700 n=3) ²   49.46m ± ∞ ¹        ~ (p=0.700 n=3) ²   67.35m ± ∞ ¹        ~ (p=0.100 n=3) ²    93.17m ± ∞ ¹        ~ (p=0.100 n=3) ²   69.05m ± ∞ ¹        ~ (p=0.100 n=3) ²   50.97m ± ∞ ¹       ~ (p=0.100 n=3) ²   78.22m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                60.75m         66.21m        +8.98%                   47.94m        -21.09%                   73.79m        +21.47%                    100.6m        +65.57%                   74.68m        +22.93%                   55.76m        -8.22%                   80.64m        +32.74%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │           bench_ncruces_driver.txt           │      bench_eatonphil_direct.txt      │          bench_glebarez_driver.txt           │            bench_mattn_driver.txt             │           bench_modernc_driver.txt           │         bench_tailscale_driver.txt          │       bench_zombiezen_direct.txt       │
                         │           B/op           │      B/op        vs base                     │    B/op      vs base                 │      B/op        vs base                     │       B/op        vs base                     │      B/op        vs base                     │      B/op        vs base                    │     B/op       vs base                 │
Query/OrderBy-12                      57465.0 ± ∞ ¹   6399323.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   118.0 ± ∞ ¹        ~ (p=0.100 n=3) ²   6397790.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   11999142.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   6397816.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   2798850.0 ± ∞ ¹           ~ (p=0.100 n=3) ²    7242.0 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-12                603.0 ± ∞ ¹   6399440.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   488.0 ± ∞ ¹        ~ (p=0.400 n=3) ²   6398236.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   12000645.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   6398321.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   2799087.0 ± ∞ ¹           ~ (p=0.100 n=3) ²    2030.0 ± ∞ ¹        ~ (p=0.700 n=3) ²
geomean                               5.749Ki           6.103Mi        +108612.04%                   240.0        -95.92%                     6.102Mi        +108588.79%                      11.44Mi        +203752.96%                     6.102Mi        +108589.73%                     2.669Mi        +47448.59%                   3.744Ki        -34.86%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt            │      bench_eatonphil_direct.txt      │           bench_glebarez_driver.txt            │             bench_mattn_driver.txt             │            bench_modernc_driver.txt            │           bench_tailscale_driver.txt           │      bench_zombiezen_direct.txt       │
                         │        allocs/op         │    allocs/op      vs base                      │  allocs/op   vs base                 │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │  allocs/op    vs base                 │
Query/OrderBy-12                       20.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   449782.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349772.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449782.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149765.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   37.000 ± ∞ ¹        ~ (p=0.300 n=3) ²
Query/OrderByParallel-12                9.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   449785.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449786.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149768.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   19.000 ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                 13.42             349.8k        +2606995.01%                   1.000        -92.55%                       449.8k        +3352388.27%                       349.8k        +2606968.92%                       449.8k        +3352392.00%                       149.8k        +1116193.58%                    26.51        +97.62%
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
Query/OrderByParallel                  75.90m ± ∞ ¹   89.80m ± ∞ ¹        ~ (p=0.100 n=3) ²   46.16m ± ∞ ¹        ~ (p=0.100 n=3) ²   82.25m ± ∞ ¹        ~ (p=0.100 n=3) ²   87.35m ± ∞ ¹        ~ (p=0.100 n=3) ²   81.94m ± ∞ ¹        ~ (p=0.100 n=3) ²   61.79m ± ∞ ¹       ~ (p=0.100 n=3) ²   84.05m ± ∞ ¹        ~ (p=0.700 n=3) ²
Query/OrderByParallel-2                47.37m ± ∞ ¹   55.96m ± ∞ ¹        ~ (p=0.100 n=3) ²   33.89m ± ∞ ¹        ~ (p=0.100 n=3) ²   56.12m ± ∞ ¹        ~ (p=0.100 n=3) ²   59.77m ± ∞ ¹        ~ (p=0.100 n=3) ²   56.63m ± ∞ ¹        ~ (p=0.100 n=3) ²   43.18m ± ∞ ¹       ~ (p=0.100 n=3) ²   49.58m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-4                30.90m ± ∞ ¹   36.22m ± ∞ ¹        ~ (p=0.100 n=3) ²   23.19m ± ∞ ¹        ~ (p=0.100 n=3) ²   37.06m ± ∞ ¹        ~ (p=0.100 n=3) ²   48.73m ± ∞ ¹        ~ (p=0.100 n=3) ²   37.73m ± ∞ ¹        ~ (p=0.100 n=3) ²   28.59m ± ∞ ¹       ~ (p=0.100 n=3) ²   38.97m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-8                36.73m ± ∞ ¹   42.41m ± ∞ ¹        ~ (p=0.100 n=3) ²   41.06m ± ∞ ¹        ~ (p=0.100 n=3) ²   62.23m ± ∞ ¹        ~ (p=0.100 n=3) ²   74.23m ± ∞ ¹        ~ (p=0.100 n=3) ²   60.93m ± ∞ ¹        ~ (p=0.100 n=3) ²   39.06m ± ∞ ¹       ~ (p=0.100 n=3) ²   75.03m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-16               52.16m ± ∞ ¹   54.09m ± ∞ ¹        ~ (p=0.100 n=3) ²   49.96m ± ∞ ¹        ~ (p=0.400 n=3) ²   73.77m ± ∞ ¹        ~ (p=0.100 n=3) ²   97.65m ± ∞ ¹        ~ (p=0.100 n=3) ²   65.19m ± ∞ ¹        ~ (p=0.100 n=3) ²   50.61m ± ∞ ¹       ~ (p=0.700 n=3) ²   80.34m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                46.30m         52.98m        +14.43%                   37.53m        -18.96%                   60.12m        +29.84%                   71.31m        +54.01%                   58.68m        +26.72%                   43.22m        -6.66%                   62.83m        +35.69%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt            │      bench_eatonphil_direct.txt       │           bench_glebarez_driver.txt            │              bench_mattn_driver.txt              │            bench_modernc_driver.txt            │           bench_tailscale_driver.txt           │       bench_zombiezen_direct.txt        │
                         │           B/op           │       B/op         vs base                     │     B/op      vs base                 │       B/op         vs base                     │        B/op         vs base                      │       B/op         vs base                     │       B/op         vs base                     │      B/op       vs base                 │
Query/OrderByParallel                50759.00 ± ∞ ¹    6399278.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    83.00 ± ∞ ¹        ~ (p=0.100 n=3) ²    6397705.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    11999009.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    6397710.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    2798846.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    6477.00 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-2               373.000 ± ∞ ¹   6399272.000 ± ∞ ¹            ~ (p=0.100 n=3) ²    9.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6397742.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   11999020.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   6397742.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   2798848.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1121.000 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-4                376.00 ± ∞ ¹    6399278.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    11.00 ± ∞ ¹        ~ (p=0.100 n=3) ²    6397857.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    11999174.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    6397882.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    2798881.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    1153.00 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-8                 392.0 ± ∞ ¹     6399448.0 ± ∞ ¹            ~ (p=0.100 n=3) ²    356.0 ± ∞ ¹        ~ (p=0.400 n=3) ²     6398024.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     11999459.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     6398034.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     2798904.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     1283.0 ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-16                636.0 ± ∞ ¹     6400390.0 ± ∞ ¹            ~ (p=0.100 n=3) ²   1168.0 ± ∞ ¹        ~ (p=0.700 n=3) ²     6398247.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     12001167.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     6398680.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     2799186.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     2334.0 ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                               1.095Ki             6.103Mi        +570481.53%                    80.67        -92.81%                       6.102Mi        +570337.26%                        11.44Mi        +1069779.70%                       6.102Mi        +570345.69%                       2.669Mi        +249452.50%                    1.860Ki        +69.82%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt            │    bench_eatonphil_direct.txt    │           bench_glebarez_driver.txt            │             bench_mattn_driver.txt             │            bench_modernc_driver.txt            │           bench_tailscale_driver.txt           │       bench_zombiezen_direct.txt       │
                         │        allocs/op         │    allocs/op      vs base                      │  allocs/op   vs base             │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │  allocs/op    vs base                  │
Query/OrderByParallel                  19.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349771.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   34.000 ± ∞ ¹         ~ (p=0.400 n=3) ²
Query/OrderByParallel-2                 8.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349771.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   16.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-4                 8.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449781.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349773.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449781.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   16.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-8                 8.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449783.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349775.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449784.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   18.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-16               10.000 ± ∞ ¹   349780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   3.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449785.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349788.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449790.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149768.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   20.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                 9.945             349.8k        +3517030.82%                                ?               ³ ⁴       449.8k        +4522593.59%                       349.8k        +3517000.66%                       449.8k        +4522605.66%                       149.8k        +1505846.97%                    19.92        +100.26%
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
Query/RecursiveCTE-12                       6.211m ± ∞ ¹   6.213m ± ∞ ¹       ~ (p=1.000 n=3) ²   2.504m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.281m ± ∞ ¹       ~ (p=0.100 n=3) ²   2.465m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.304m ± ∞ ¹       ~ (p=0.100 n=3) ²   2.431m ± ∞ ¹        ~ (p=0.100 n=3) ²   5.255m ± ∞ ¹       ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-12               678.6µ ± ∞ ¹   678.7µ ± ∞ ¹       ~ (p=1.000 n=3) ²   274.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   656.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²   271.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   648.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²   267.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   651.5µ ± ∞ ¹       ~ (p=0.400 n=3) ²
geomean                                     2.053m         2.053m        +0.02%                   828.9µ        -59.62%                   1.861m        -9.33%                   818.4µ        -60.14%                   1.855m        -9.66%                   805.7µ        -60.75%                   1.850m        -9.88%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt           │    bench_eatonphil_direct.txt    │          bench_glebarez_driver.txt          │           bench_mattn_driver.txt            │          bench_modernc_driver.txt           │         bench_tailscale_driver.txt          │        bench_zombiezen_direct.txt         │
                              │           B/op           │      B/op       vs base                     │    B/op      vs base             │      B/op       vs base                     │      B/op       vs base                     │      B/op       vs base                     │      B/op       vs base                     │     B/op       vs base                    │
Query/RecursiveCTE-12                        1.000 ± ∞ ¹   2482.000 ± ∞ ¹            ~ (p=0.600 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.300 n=3) ²     2354.000 ± ∞ ¹            ~ (p=0.600 n=3) ²   6857.000 ± ∞ ¹            ~ (p=0.600 n=3) ²   2294.000 ± ∞ ¹            ~ (p=0.600 n=3) ²   1517.000 ± ∞ ¹            ~ (p=0.600 n=3) ²   391.000 ± ∞ ¹           ~ (p=0.600 n=3) ²
Query/RecursiveCTEParallel-12                2.000 ± ∞ ¹   2485.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.700 n=3) ²     2357.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6851.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   2296.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1505.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   370.000 ± ∞ ¹           ~ (p=0.100 n=3) ²
geomean                                      1.414          2.425Ki        +175509.94%                                ?               ³ ⁴    2.300Ki        +166458.97%                    6.693Ki        +484550.94%                    2.241Ki        +162180.99%                    1.476Ki        +106742.99%                     380.4        +26795.17%
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
                              │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt       │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                              │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │
Query/RecursiveCTEParallel                  6.116m ± ∞ ¹    6.131m ± ∞ ¹       ~ (p=0.100 n=3) ²   2.496m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.268m ± ∞ ¹       ~ (p=0.100 n=3) ²   2.456m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.246m ± ∞ ¹        ~ (p=0.100 n=3) ²   2.431m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.233m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                3.086m ± ∞ ¹    3.099m ± ∞ ¹       ~ (p=0.100 n=3) ²   1.256m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.668m ± ∞ ¹       ~ (p=0.100 n=3) ²   1.247m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.651m ± ∞ ¹        ~ (p=0.100 n=3) ²   1.237m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.647m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4               1578.4µ ± ∞ ¹   1587.6µ ± ∞ ¹       ~ (p=0.100 n=3) ²   642.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1371.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²   633.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1362.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   625.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1358.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                792.2µ ± ∞ ¹    798.6µ ± ∞ ¹       ~ (p=0.100 n=3) ²   322.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    742.9µ ± ∞ ¹       ~ (p=0.100 n=3) ²   320.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    746.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   316.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    744.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16               667.3µ ± ∞ ¹    676.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²   274.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    651.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²   277.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    644.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   274.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    660.2µ ± ∞ ¹        ~ (p=0.400 n=3) ²
geomean                                     1.736m          1.748m        +0.69%                   708.3µ        -59.19%                    1.563m        -9.96%                   703.6µ        -59.46%                    1.556m        -10.36%                   696.2µ        -59.89%                    1.560m        -10.11%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │    bench_eatonphil_direct.txt    │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                              │           B/op           │      B/op       vs base           │    B/op      vs base             │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │     B/op       vs base           │
Query/RecursiveCTEParallel                   1.000 ± ∞ ¹   2482.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     2260.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6857.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2258.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1504.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   362.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                   0.0 ± ∞ ¹     2481.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ²       2261.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     6817.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     2257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1502.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     369.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4                 1.000 ± ∞ ¹   2481.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.700 n=3) ²     2275.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6840.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2260.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1503.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   373.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                 2.000 ± ∞ ¹   2481.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.800 n=3) ²     2356.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6853.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2295.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1504.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   364.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16                2.000 ± ∞ ¹   2490.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=1.000 n=3) ²     2366.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6868.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2285.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1514.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   369.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                ³    2.425Ki        ?                                ?               ³ ⁴    2.249Ki        ?                    6.687Ki        ?                    2.218Ki        ?                    1.470Ki        ?                     367.4        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean
⁴ ratios must be >0 to compute geomean

                              │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                              │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/RecursiveCTEParallel                   0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   143.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   48.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   48.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   48.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16                0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   143.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                ⁴     110.0        ?                                +0.00%               ⁴     112.4        ?                     142.4        ?                     112.0        ?                    48.40        ?                   6.000        ?
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
Query/Window-12                      2148.4µ ± ∞ ¹   2431.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   914.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1880.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   2089.4µ ± ∞ ¹         ~ (p=0.400 n=3) ²   1890.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1050.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1634.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-12               236.8µ ± ∞ ¹    316.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   121.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    912.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1078.0µ ± ∞ ¹         ~ (p=0.100 n=3) ²    900.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    131.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    841.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                               713.2µ          877.8µ        +23.08%                   332.9µ        -53.32%                    1.310m        +83.70%                    1.501m        +110.42%                    1.304m        +82.90%                    371.2µ        -47.96%                    1.173m        +64.46%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                        │ bench_ncruces_direct.txt │      bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt      │       bench_mattn_driver.txt        │      bench_modernc_driver.txt      │     bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                        │           B/op           │      B/op        vs base           │    B/op      vs base                │      B/op        vs base           │       B/op        vs base           │      B/op        vs base           │      B/op        vs base           │     B/op       vs base           │
Query/Window-12                          0.0 ± ∞ ¹     62839.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹       ~ (p=1.000 n=3) ²     54893.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198934.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54923.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30798.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     371.0 ± ∞ ¹  ~ (p=0.600 n=3) ²
Query/WindowParallel-12                1.000 ± ∞ ¹   62763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹       ~ (p=1.000 n=3) ²   54887.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   198955.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54908.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   30766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   376.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                          ³     61.33Ki        ?                                +0.00%               ³     53.60Ki        ?                      194.3Ki        ?                     53.63Ki        ?                     30.06Ki        ?                     373.5        ?
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
Query/WindowParallel                 2105.4µ ± ∞ ¹   2399.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   937.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1874.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1569.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1869.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1055.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1614.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-2               1061.9µ ± ∞ ¹   1229.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   477.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    977.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    845.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    981.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    541.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    845.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-4                544.0µ ± ∞ ¹    643.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   257.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    532.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    549.2µ ± ∞ ¹        ~ (p=0.700 n=3) ²    531.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    282.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    456.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-8                273.9µ ± ∞ ¹    353.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   142.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    873.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1071.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    867.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    149.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    835.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/WindowParallel-16               234.3µ ± ∞ ¹    315.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   118.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    920.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1083.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    908.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    141.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    847.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                               600.5µ          732.7µ        +22.01%                   287.0µ        -52.21%                    952.8µ        +58.68%                    966.9µ        +61.02%                    948.8µ        +58.00%                    320.8µ        -46.57%                    849.2µ        +41.43%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                        │ bench_ncruces_direct.txt │      bench_ncruces_driver.txt      │      bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt      │       bench_mattn_driver.txt        │      bench_modernc_driver.txt      │     bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                        │           B/op           │      B/op        vs base           │    B/op      vs base                 │      B/op        vs base           │       B/op        vs base           │      B/op        vs base           │      B/op        vs base           │     B/op       vs base           │
Query/WindowParallel                     0.0 ± ∞ ¹     62762.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹        ~ (p=1.000 n=3) ³     54775.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198936.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54789.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30792.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     362.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-2                   0.0 ± ∞ ¹     62761.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹        ~ (p=1.000 n=3) ³     54826.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198852.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54809.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30792.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     360.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-4                   0.0 ± ∞ ¹     62762.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹        ~ (p=1.000 n=3) ³     54882.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198911.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54878.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30793.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     363.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-8                   0.0 ± ∞ ¹     62762.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹        ~ (p=1.000 n=3) ³     54886.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198908.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54900.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30796.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     364.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/WindowParallel-16                1.000 ± ∞ ¹   62771.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹        ~ (p=0.800 n=3) ²   54899.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   198983.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54915.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   30767.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   373.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                          ⁴     61.29Ki        ?                                +14.87%               ⁴     53.57Ki        ?                      194.3Ki        ?                     53.57Ki        ?                     30.07Ki        ?                     364.4        ?
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

  `ReadOnly` transactions are always `deferred`. Not-`ReadOnly` transactions use the value of `_txlock`.

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

  `ReadOnly` transactions are always `deferred`. Not-`ReadOnly` transactions use the value of `_txlock`.

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
