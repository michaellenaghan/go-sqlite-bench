This project, originally inspired by [Benchmarking SQLite Performance in Go](https://www.golang.dk/articles/benchmarking-sqlite-performance-in-go), bechmarks various SQLite implementations:

* [github.com/eatonphil/gosqlite](https://github.com/eatonphil/gosqlite) (aka "eatonphil")
* [github.com/glebarez/go-sqlite](https://github.com/glebarez/go-sqlite) (aka "glebarez")
* [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) (aka "mattn")
* [github.com/ncruces/go-sqlite3](https://github.com/ncruces/go-sqlite3) (aka "ncruces")
* [github.com/tailscale/sqlite](https://github.com/tailscale/sqlite) (aka "tailscale")
*	[github.com/zombiezen/go-sqlite](https://github.com/zombiezen/go-sqlite) (aka "zombiezen")
*	[gitlab.com/cznic/sqlite](https://gitlab.com/cznic/sqlite) (aka "modernc")

Here are some quick descriptions:

* **eatonphil** is a CGO-based implementation. eatonphil offers a direct interface.

* **glebarez** is a non-CGO transpilation-based implementation, based on modernc. glebarez offers a `database/sql` interface.

* **mattn** is a CGO-based implementation. mattn offers a `database/sql` interface.

* **modernc** is a non-CGO transpilation-based implementation. modernc offers a `database/sql` interface.

* **ncruces** is a non-CGO WASM-based implementation. ncruces offers both direct and `database/sql` interfaces.

* **tailscale** is a CGO-based implementation. tailscale offers both direct and `database/sql` interfaces.

* **zombiezen** is a non-CGO transpilation-based implementation, based on modernc. zombiezen offers a direct interface.

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

The "slow" results are included in this README file, below.

The tests in the [tests](./tests/) directory were generated using:

```sh
make test-all
```

Among other things, the tests capture the compile-time options, pragmas, and SQLite version used by each implementation.

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
  bench-category-query-nonrecursivecte               - Run nonrecursivecte query benchmarks
  bench-category-query-nonrecursivecte-parallel      - Run nonrecursivecte query benchmarks in parallel
  bench-category-query-orderby                       - Run orderby query benchmarks
  bench-category-query-orderby-parallel              - Run orderby query benchmarks in parallel
  bench-category-query-recursivecte                  - Run recursivecte query benchmarks
  bench-category-query-recursivecte-parallel         - Run recursivecte query benchmarks in parallel
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
  benchstat-category-query-nonrecursivecte           - Run and compare nonrecursivecte query benchmarks
  benchstat-category-query-nonrecursivecte-parallel  - Run and compare nonrecursivecte query benchmarks in parallel
  benchstat-category-query-orderby                   - Run and compare orderby query benchmarks
  benchstat-category-query-orderby-parallel          - Run and compare orderby query benchmarks in parallel
  benchstat-category-query-recursivecte              - Run and compare recursivecte query benchmarks
  benchstat-category-query-recursivecte-parallel     - Run and compare recursivecte query benchmarks in parallel
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

## Baseline (Sequential)

<!--BENCHMARK:slow/benchstat_baseline.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                                       │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │         bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt        │
                                       │          sec/op          │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                  │
Baseline/Conn-12                                     152.6n ± ∞ ¹    222.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   167.3n ± ∞ ¹        ~ (p=0.100 n=3) ²    225.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    220.1n ± ∞ ¹        ~ (p=0.100 n=3) ²    225.8n ± ∞ ¹        ~ (p=0.100 n=3) ²    223.8n ± ∞ ¹       ~ (p=0.100 n=3) ²   1124.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-12                             453.4n ± ∞ ¹    540.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   459.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    546.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    540.6n ± ∞ ¹        ~ (p=0.100 n=3) ²    546.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    541.9n ± ∞ ¹       ~ (p=0.100 n=3) ²   1346.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1-12                                 1773.0n ± ∞ ¹   1935.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   878.4n ± ∞ ¹        ~ (p=0.100 n=3) ²   2041.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2602.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2055.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1122.0n ± ∞ ¹       ~ (p=0.100 n=3) ²   3182.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-12                         1006.0n ± ∞ ¹   1162.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   887.6n ± ∞ ¹        ~ (p=0.100 n=3) ²   2286.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1484.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2329.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    914.1n ± ∞ ¹       ~ (p=0.100 n=3) ²   2648.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePrepared-12                       430.1n ± ∞ ¹    522.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   438.4n ± ∞ ¹        ~ (p=0.100 n=3) ²   2019.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1303.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1984.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    410.4n ± ∞ ¹       ~ (p=0.100 n=3) ²   1527.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-12               744.1n ± ∞ ¹    939.8n ± ∞ ¹        ~ (p=0.100 n=3) ²   754.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   1782.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1041.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1799.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    769.0n ± ∞ ¹       ~ (p=0.100 n=3) ²   1563.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                              583.6n          714.5n        +22.44%                   520.2n        -10.86%                    1.128µ        +93.34%                    924.2n        +58.37%                    1.132µ        +94.02%                    583.0n        -0.10%                    1.767µ        +202.77%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                       │ bench_ncruces_direct.txt │    bench_ncruces_driver.txt    │      bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt    │     bench_mattn_driver.txt      │    bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt    │
                                       │           B/op           │    B/op      vs base           │    B/op      vs base                 │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │
Baseline/Conn-12                                       0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   353.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-12                               0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1-12                                   48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   705.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-12                           48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePrepared-12                         0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   353.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
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

 ## Baseline (Parallel)

<!--BENCHMARK:slow/benchstat_baseline_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                                       │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │         bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt      │       bench_zombiezen_direct.txt        │
                                       │          sec/op          │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op     vs base                │    sec/op      vs base                  │
Baseline/ConnParallel                                150.4n ± ∞ ¹    232.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   167.4n ± ∞ ¹        ~ (p=0.100 n=3) ²    243.6n ± ∞ ¹        ~ (p=0.100 n=3) ²    227.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    243.2n ± ∞ ¹        ~ (p=0.100 n=3) ²   235.4n ± ∞ ¹       ~ (p=0.100 n=3) ²   1361.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-2                              156.3n ± ∞ ¹    229.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   154.1n ± ∞ ¹        ~ (p=0.500 n=3) ²    220.5n ± ∞ ¹        ~ (p=0.100 n=3) ²    213.3n ± ∞ ¹        ~ (p=0.100 n=3) ²    220.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   216.6n ± ∞ ¹       ~ (p=0.100 n=3) ²    672.1n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-4                              241.7n ± ∞ ¹    303.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   241.1n ± ∞ ¹        ~ (p=1.000 n=3) ²    263.3n ± ∞ ¹        ~ (p=0.100 n=3) ²    269.1n ± ∞ ¹        ~ (p=0.100 n=3) ²    261.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   265.6n ± ∞ ¹       ~ (p=0.100 n=3) ²    670.6n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-8                              455.7n ± ∞ ¹    512.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   426.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    518.7n ± ∞ ¹        ~ (p=0.100 n=3) ²    515.1n ± ∞ ¹        ~ (p=0.100 n=3) ²    518.9n ± ∞ ¹        ~ (p=0.100 n=3) ²   511.7n ± ∞ ¹       ~ (p=0.100 n=3) ²   1127.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/ConnParallel-16                             486.3n ± ∞ ¹    540.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   465.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    556.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    548.2n ± ∞ ¹        ~ (p=0.100 n=3) ²    554.6n ± ∞ ¹        ~ (p=0.100 n=3) ²   547.7n ± ∞ ¹       ~ (p=0.100 n=3) ²   1372.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel                             1.776µ ± ∞ ¹    1.932µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1.034µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1.870µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1.628µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1.879µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1.139µ ± ∞ ¹       ~ (p=0.100 n=3) ²    2.825µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-2                           942.0n ± ∞ ¹   1001.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   537.2n ± ∞ ¹        ~ (p=0.100 n=3) ²   1100.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    956.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   1094.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   634.7n ± ∞ ¹       ~ (p=0.100 n=3) ²   1517.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-4                           540.4n ± ∞ ¹    577.5n ± ∞ ¹        ~ (p=0.100 n=3) ²   367.5n ± ∞ ¹        ~ (p=0.100 n=3) ²    958.5n ± ∞ ¹        ~ (p=0.100 n=3) ²    601.5n ± ∞ ¹        ~ (p=0.100 n=3) ²    998.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   381.5n ± ∞ ¹       ~ (p=0.100 n=3) ²   1077.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-8                           887.9n ± ∞ ¹   1111.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   762.5n ± ∞ ¹        ~ (p=0.100 n=3) ²   2206.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1449.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2236.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   877.4n ± ∞ ¹       ~ (p=0.100 n=3) ²   2461.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-16                          928.2n ± ∞ ¹   1165.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   823.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   2263.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1510.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   2304.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   924.2n ± ∞ ¹       ~ (p=0.100 n=3) ²   2680.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel                  423.7n ± ∞ ¹    521.1n ± ∞ ¹        ~ (p=0.100 n=3) ²   437.7n ± ∞ ¹        ~ (p=0.100 n=3) ²   1846.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    853.6n ± ∞ ¹        ~ (p=0.100 n=3) ²   1824.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   417.8n ± ∞ ¹       ~ (p=0.100 n=3) ²   1597.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-2                363.5n ± ∞ ¹    443.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   317.4n ± ∞ ¹        ~ (p=0.100 n=3) ²   1051.0n ± ∞ ¹        ~ (p=0.100 n=3) ²    576.6n ± ∞ ¹        ~ (p=0.100 n=3) ²   1028.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   342.1n ± ∞ ¹       ~ (p=0.100 n=3) ²    838.1n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-4                292.3n ± ∞ ¹    345.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   239.9n ± ∞ ¹        ~ (p=0.100 n=3) ²    838.4n ± ∞ ¹        ~ (p=0.100 n=3) ²    470.9n ± ∞ ¹        ~ (p=0.100 n=3) ²    749.2n ± ∞ ¹        ~ (p=0.100 n=3) ²   273.9n ± ∞ ¹       ~ (p=0.400 n=3) ²    631.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-8                697.3n ± ∞ ¹    875.2n ± ∞ ¹        ~ (p=0.100 n=3) ²   673.3n ± ∞ ¹        ~ (p=0.100 n=3) ²   1723.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1047.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1741.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   731.4n ± ∞ ¹       ~ (p=0.100 n=3) ²   1370.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-16               724.4n ± ∞ ¹    941.5n ± ∞ ¹        ~ (p=0.100 n=3) ²   716.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1772.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1059.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   1792.0n ± ∞ ¹        ~ (p=0.100 n=3) ²   773.2n ± ∞ ¹       ~ (p=0.100 n=3) ²   1515.0n ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                              488.2n          595.9n        +22.06%                   423.6n        -13.24%                    898.1n        +83.96%                    659.5n        +35.09%                    894.3n        +83.17%                   484.1n        -0.84%                    1.299µ        +166.08%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                       │ bench_ncruces_direct.txt │    bench_ncruces_driver.txt    │      bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt    │     bench_mattn_driver.txt      │    bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt    │
                                       │           B/op           │    B/op      vs base           │    B/op      vs base                 │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │     B/op      vs base           │
Baseline/ConnParallel                                  0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   387.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-2                                0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-4                                0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-8                                0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/ConnParallel-16                               0.00 ± ∞ ¹   64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    64.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel                              48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   727.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-2                            48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-4                            48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-8                            48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1Parallel-16                           48.00 ± ∞ ¹   32.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   96.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   704.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel                    0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   369.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-2                  0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-4                  0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-8                  0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
Baseline/Select1PrePreparedParallel-16                 0.00 ± ∞ ¹   48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    0.00 ± ∞ ¹        ~ (p=1.000 n=3) ³   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   252.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   352.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                         ⁴   46.15        ?                                +25.99%               ⁴    159.6        ?                    164.2        ?                    159.6        ?                    90.34        ?                    448.7        ?
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

<!--BENCHMARK:slow/benchstat_populate.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                              │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │      bench_zombiezen_direct.txt       │
                              │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                │
Populate/PopulateDB-12                       2.615 ± ∞ ¹     2.476 ± ∞ ¹       ~ (p=0.100 n=3) ²    2.000 ± ∞ ¹        ~ (p=0.100 n=3) ²     2.758 ± ∞ ¹        ~ (p=0.200 n=3) ²     2.193 ± ∞ ¹       ~ (p=0.100 n=3) ²     2.697 ± ∞ ¹        ~ (p=0.200 n=3) ²    2.025 ± ∞ ¹        ~ (p=0.100 n=3) ²     2.369 ± ∞ ¹       ~ (p=0.700 n=3) ²
Populate/PopulateDBWithTx-12               1091.9m ± ∞ ¹   1167.3m ± ∞ ¹       ~ (p=0.700 n=3) ²   937.9m ± ∞ ¹        ~ (p=0.100 n=3) ²   1473.0m ± ∞ ¹        ~ (p=0.100 n=3) ²   1142.0m ± ∞ ¹       ~ (p=0.700 n=3) ²   1486.9m ± ∞ ¹        ~ (p=0.100 n=3) ²   991.4m ± ∞ ¹        ~ (p=0.700 n=3) ²   1128.8m ± ∞ ¹       ~ (p=0.700 n=3) ²
Populate/PopulateDBWithTxs-12                1.077 ± ∞ ¹     1.156 ± ∞ ¹       ~ (p=0.200 n=3) ²    1.075 ± ∞ ¹        ~ (p=0.700 n=3) ²     1.317 ± ∞ ¹        ~ (p=0.100 n=3) ²     1.205 ± ∞ ¹       ~ (p=0.100 n=3) ²     1.354 ± ∞ ¹        ~ (p=0.100 n=3) ²    1.075 ± ∞ ¹        ~ (p=1.000 n=3) ²     1.206 ± ∞ ¹       ~ (p=0.100 n=3) ²
geomean                                      1.454           1.495        +2.80%                    1.263        -13.11%                     1.749        +20.27%                     1.445        -0.64%                     1.758        +20.87%                    1.292        -11.13%                     1.477        +1.60%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │         bench_ncruces_driver.txt          │       bench_eatonphil_direct.txt        │         bench_glebarez_driver.txt         │           bench_mattn_driver.txt            │         bench_modernc_driver.txt          │        bench_tailscale_driver.txt         │       bench_zombiezen_direct.txt        │
                              │           B/op           │      B/op       vs base                   │     B/op       vs base                  │      B/op       vs base                   │      B/op        vs base                    │      B/op       vs base                   │      B/op       vs base                   │     B/op       vs base                  │
Populate/PopulateDB-12                     2.287Mi ± ∞ ¹   18.968Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   5.770Mi ± ∞ ¹         ~ (p=0.100 n=3) ²   31.472Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   273.295Mi ± ∞ ¹           ~ (p=0.100 n=3) ²   31.478Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   18.963Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   6.167Mi ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTx-12               2.276Mi ± ∞ ¹   32.185Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   5.771Mi ± ∞ ¹         ~ (p=0.100 n=3) ²   44.672Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   286.682Mi ± ∞ ¹           ~ (p=0.100 n=3) ²   44.664Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   33.021Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   6.167Mi ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTxs-12              2.276Mi ± ∞ ¹   31.224Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   5.770Mi ± ∞ ¹         ~ (p=0.100 n=3) ²   43.673Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   285.984Mi ± ∞ ¹           ~ (p=0.100 n=3) ²   43.855Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   31.266Mi ± ∞ ¹          ~ (p=0.100 n=3) ²   6.212Mi ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                    2.280Mi          26.71Mi        +1071.77%                   5.770Mi        +153.12%                    39.45Mi        +1630.54%                     281.9Mi        +12266.52%                    39.51Mi        +1632.94%                    26.95Mi        +1082.25%                   6.182Mi        +171.18%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │         bench_mattn_driver.txt          │        bench_modernc_driver.txt         │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                              │        allocs/op         │  allocs/op    vs base                  │  allocs/op    vs base                  │   allocs/op    vs base                  │   allocs/op    vs base                  │   allocs/op    vs base                  │  allocs/op    vs base                  │  allocs/op    vs base                  │
Populate/PopulateDB-12                      140.0k ± ∞ ¹   598.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   394.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1209.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   1057.5k ± ∞ ¹         ~ (p=0.100 n=3) ²   1209.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   598.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   445.1k ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTx-12                140.0k ± ∞ ¹   751.1k ± ∞ ¹         ~ (p=0.100 n=3) ²   394.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1362.4k ± ∞ ¹         ~ (p=0.100 n=3) ²   1210.5k ± ∞ ¹         ~ (p=0.100 n=3) ²   1362.4k ± ∞ ¹         ~ (p=0.100 n=3) ²   751.2k ± ∞ ¹         ~ (p=0.100 n=3) ²   445.1k ± ∞ ¹         ~ (p=0.100 n=3) ²
Populate/PopulateDBWithTxs-12               140.0k ± ∞ ¹   765.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   394.0k ± ∞ ¹         ~ (p=0.100 n=3) ²   1376.4k ± ∞ ¹         ~ (p=0.100 n=3) ²   1238.8k ± ∞ ¹         ~ (p=0.100 n=3) ²   1380.4k ± ∞ ¹         ~ (p=0.100 n=3) ²   767.3k ± ∞ ¹         ~ (p=0.100 n=3) ²   448.1k ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                     140.0k         700.5k        +400.31%                   394.0k        +181.40%                    1.314M        +838.29%                    1.166M        +732.84%                    1.315M        +839.20%                   701.2k        +400.77%                   446.1k        +218.59%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
```
<!--END_BENCHMARK-->

## ReadWrite (Sequential)

<!--BENCHMARK:slow/benchstat_readwrite.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                                                                    │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │         bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                                                                    │          sec/op          │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                │    sec/op      vs base                 │
ReadWrite/ReadPost-12                                                            13.853µ ± ∞ ¹   21.773µ ± ∞ ¹        ~ (p=0.100 n=3) ²   9.151µ ± ∞ ¹        ~ (p=0.100 n=3) ²   28.435µ ± ∞ ¹        ~ (p=0.100 n=3) ²   17.812µ ± ∞ ¹        ~ (p=0.100 n=3) ²   28.971µ ± ∞ ¹        ~ (p=0.100 n=3) ²   18.539µ ± ∞ ¹       ~ (p=0.100 n=3) ²   19.126µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadPostWithTx-12                                                      16.013µ ± ∞ ¹   37.814µ ± ∞ ¹        ~ (p=0.100 n=3) ²   9.675µ ± ∞ ¹        ~ (p=0.100 n=3) ²   31.798µ ± ∞ ¹        ~ (p=0.100 n=3) ²   22.077µ ± ∞ ¹        ~ (p=0.100 n=3) ²   32.867µ ± ∞ ¹        ~ (p=0.100 n=3) ²   25.754µ ± ∞ ¹       ~ (p=0.100 n=3) ²   19.681µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndComments-12                                                  98.31µ ± ∞ ¹   113.92µ ± ∞ ¹        ~ (p=0.100 n=3) ²   71.03µ ± ∞ ¹        ~ (p=0.100 n=3) ²   169.54µ ± ∞ ¹        ~ (p=0.100 n=3) ²   176.30µ ± ∞ ¹        ~ (p=0.100 n=3) ²   179.98µ ± ∞ ¹        ~ (p=0.100 n=3) ²   104.26µ ± ∞ ¹       ~ (p=0.100 n=3) ²    98.10µ ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadPostAndCommentsWithTx-12                                            99.17µ ± ∞ ¹   142.44µ ± ∞ ¹        ~ (p=0.100 n=3) ²   70.58µ ± ∞ ¹        ~ (p=0.100 n=3) ²   171.79µ ± ∞ ¹        ~ (p=0.100 n=3) ²   179.70µ ± ∞ ¹        ~ (p=0.100 n=3) ²   177.26µ ± ∞ ¹        ~ (p=0.100 n=3) ²   118.70µ ± ∞ ¹       ~ (p=0.100 n=3) ²    95.55µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePost-12                                                            131.6µ ± ∞ ¹    127.0µ ± ∞ ¹        ~ (p=0.400 n=3) ²   115.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    140.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    123.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    140.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    116.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²    161.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePostWithTx-12                                                      140.1µ ± ∞ ¹    145.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   116.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    147.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    130.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    149.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    125.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²    164.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndComments-12                                                 2.559m ± ∞ ¹    2.577m ± ∞ ¹        ~ (p=0.400 n=3) ²   2.095m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.809m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.283m ± ∞ ¹        ~ (p=0.100 n=3) ²    3.408m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.139m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.541m ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/WritePostAndCommentsWithTx-12                                           933.1µ ± ∞ ¹   1022.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²   772.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1256.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    997.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1556.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    829.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²   1085.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=10-12                             354.5µ ± ∞ ¹    349.9µ ± ∞ ¹        ~ (p=1.000 n=3) ²   268.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    448.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    412.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    495.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    304.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²    402.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=90-12                             2.310m ± ∞ ¹    2.390m ± ∞ ¹        ~ (p=0.100 n=3) ²   1.867m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.585m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.104m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.629m ± ∞ ¹        ~ (p=0.100 n=3) ²    1.937m ± ∞ ¹       ~ (p=0.100 n=3) ²    2.318m ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-12                     400.3µ ± ∞ ¹    406.1µ ± ∞ ¹        ~ (p=0.700 n=3) ²   267.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    406.4µ ± ∞ ¹        ~ (p=1.000 n=3) ²    344.8µ ± ∞ ¹        ~ (p=0.200 n=3) ²    404.8µ ± ∞ ¹        ~ (p=1.000 n=3) ²    274.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²    382.9µ ± ∞ ¹        ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-12                     3.210m ± ∞ ¹    3.090m ± ∞ ¹        ~ (p=0.700 n=3) ²   2.419m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.850m ± ∞ ¹        ~ (p=0.200 n=3) ²    2.655m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.837m ± ∞ ¹        ~ (p=0.700 n=3) ²    2.459m ± ∞ ¹       ~ (p=0.100 n=3) ²    3.034m ± ∞ ¹        ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=10-12                       182.8µ ± ∞ ¹    225.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   146.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    283.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    282.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    306.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    182.1µ ± ∞ ¹       ~ (p=0.700 n=3) ²    254.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=90-12                       860.3µ ± ∞ ¹    948.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   712.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1160.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    931.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1186.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    768.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²   1007.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-12               113.7µ ± ∞ ¹    128.0µ ± ∞ ¹        ~ (p=0.200 n=3) ²   115.9µ ± ∞ ¹        ~ (p=0.700 n=3) ²    219.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    184.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    209.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    141.2µ ± ∞ ¹       ~ (p=0.100 n=3) ²    177.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-12               872.5µ ± ∞ ¹    917.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   939.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   4739.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1383.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   5886.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    912.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²   3606.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                                                           281.6µ          327.0µ        +16.13%                   222.3µ        -21.07%                    418.3µ        +48.55%                    334.3µ        +18.72%                    443.9µ        +57.63%                    275.5µ        -2.16%                    347.8µ        +23.50%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                                                    │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │    bench_eatonphil_direct.txt    │     bench_glebarez_driver.txt     │       bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                                                                    │           B/op           │      B/op       vs base           │     B/op       vs base           │      B/op       vs base           │      B/op        vs base           │      B/op       vs base           │      B/op       vs base           │     B/op       vs base           │
ReadWrite/ReadPost-12                                                            40.17Ki ± ∞ ¹    41.23Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.18Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    81.34Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     41.39Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    81.28Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    41.19Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.55Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostWithTx-12                                                      40.17Ki ± ∞ ¹    42.29Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.18Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    82.22Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     42.60Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    82.33Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    42.13Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   40.60Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndComments-12                                                 250.2Ki ± ∞ ¹    260.3Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    501.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     263.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    501.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    257.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadPostAndCommentsWithTx-12                                           250.2Ki ± ∞ ¹    261.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    503.1Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     265.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    503.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    258.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   250.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePost-12                                                              0.00 ± ∞ ¹     240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     496.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    42288.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     496.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     240.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    409.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostWithTx-12                                                        0.00 ± ∞ ¹     857.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     48.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    1106.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    43253.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    1289.00 ± ∞ ¹  ~ (p=0.100 n=3) ²     880.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    456.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndComments-12                                                0.000Ki ± ∞ ¹   13.906Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   2.781Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.667Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   308.876Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.686Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   13.906Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   3.535Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/WritePostAndCommentsWithTx-12                                          0.000Ki ± ∞ ¹   26.454Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   2.781Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   39.192Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   321.743Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   39.365Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   26.469Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   3.580Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=10-12                            225.0Ki ± ∞ ¹    237.1Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   226.0Ki ± ∞ ¹  ~ (p=0.400 n=3) ²    453.8Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     268.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    452.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    232.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.8Ki ± ∞ ¹  ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=90-12                            25.84Ki ± ∞ ¹    37.72Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   29.65Ki ± ∞ ¹  ~ (p=0.700 n=3) ²    71.22Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    304.51Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    73.97Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    36.89Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   28.31Ki ± ∞ ¹  ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-12                    226.3Ki ± ∞ ¹    234.5Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   226.8Ki ± ∞ ¹  ~ (p=0.700 n=3) ²    453.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     268.4Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    457.3Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    233.4Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.9Ki ± ∞ ¹  ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-12                    23.86Ki ± ∞ ¹    39.02Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   22.35Ki ± ∞ ¹  ~ (p=1.000 n=3) ²    71.75Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    303.63Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    65.64Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    37.83Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   24.37Ki ± ∞ ¹  ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=10-12                      225.6Ki ± ∞ ¹    237.7Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   226.1Ki ± ∞ ¹  ~ (p=0.400 n=3) ²    457.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     271.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    456.0Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    234.6Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.5Ki ± ∞ ¹  ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=90-12                      24.85Ki ± ∞ ¹    47.78Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   25.59Ki ± ∞ ¹  ~ (p=0.700 n=3) ²    89.08Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    316.30Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    86.17Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    50.35Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   28.36Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-12              224.8Ki ± ∞ ¹    237.9Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   225.5Ki ± ∞ ¹  ~ (p=1.000 n=3) ²    457.1Ki ± ∞ ¹  ~ (p=0.100 n=3) ²     270.8Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    457.2Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    235.5Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   226.0Ki ± ∞ ¹  ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-12              24.34Ki ± ∞ ¹    49.81Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   27.46Ki ± ∞ ¹  ~ (p=0.200 n=3) ²    79.46Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    315.78Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    82.02Ki ± ∞ ¹  ~ (p=0.100 n=3) ²    50.63Ki ± ∞ ¹  ~ (p=0.100 n=3) ²   31.46Ki ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                                                      ³    42.02Ki        ?                   21.50Ki        ?                    76.72Ki        ?                     178.3Ki        ?                    77.24Ki        ?                    41.93Ki        ?                   29.73Ki        ?
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
ReadWrite/ReadOrWritePostAndComments/write_rate=10-12                              235.0 ± ∞ ¹    738.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    697.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndComments/write_rate=90-12                              27.00 ± ∞ ¹   442.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   802.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   421.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-12                      237.0 ± ∞ ¹    734.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    697.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    547.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-12                      25.00 ± ∞ ¹   444.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   207.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   799.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   422.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=10-12                        236.0 ± ∞ ¹    777.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1010.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    274.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTx/write_rate=90-12                        26.00 ± ∞ ¹   594.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1119.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   969.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1124.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-12                235.0 ± ∞ ¹    777.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1010.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    744.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹  ~ (p=0.100 n=3) ²    274.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-12                25.00 ± ∞ ¹   596.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1122.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   966.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   1125.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                                                      ³    273.3        ?                    85.21        ?                     430.4        ?                    368.2        ?                     435.9        ?                    237.4        ?                    122.4        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

## ReadWrite (Parallel)

<!--BENCHMARK:slow/benchstat_readwrite_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                                                                    │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt         │      bench_tailscale_driver.txt      │       bench_zombiezen_direct.txt        │
                                                                    │          sec/op          │    sec/op      vs base                │    sec/op      vs base                 │     sec/op      vs base                 │    sec/op      vs base                │     sec/op      vs base                 │    sec/op     vs base                │     sec/op      vs base                 │
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10                        378.9µ ± ∞ ¹    346.3µ ± ∞ ¹       ~ (p=0.200 n=3) ²    268.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²     428.5µ ± ∞ ¹        ~ (p=0.700 n=3) ²    357.9µ ± ∞ ¹       ~ (p=0.700 n=3) ²     423.5µ ± ∞ ¹        ~ (p=0.700 n=3) ²   288.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²     426.2µ ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-2                      334.5µ ± ∞ ¹    338.4µ ± ∞ ¹       ~ (p=1.000 n=3) ²    247.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²     375.8µ ± ∞ ¹        ~ (p=0.700 n=3) ²    310.9µ ± ∞ ¹       ~ (p=0.700 n=3) ²     370.9µ ± ∞ ¹        ~ (p=0.700 n=3) ²   265.1µ ± ∞ ¹       ~ (p=0.100 n=3) ²     353.7µ ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-4                      330.0µ ± ∞ ¹    331.6µ ± ∞ ¹       ~ (p=1.000 n=3) ²    265.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²     367.2µ ± ∞ ¹        ~ (p=0.200 n=3) ²    299.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²     352.8µ ± ∞ ¹        ~ (p=0.700 n=3) ²   270.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²     315.1µ ± ∞ ¹        ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-8                      371.7µ ± ∞ ¹    358.4µ ± ∞ ¹       ~ (p=0.700 n=3) ²    268.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²     374.6µ ± ∞ ¹        ~ (p=0.700 n=3) ²    325.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²     389.3µ ± ∞ ¹        ~ (p=0.200 n=3) ²   288.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²     345.8µ ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-16                     398.4µ ± ∞ ¹    415.1µ ± ∞ ¹       ~ (p=0.400 n=3) ²    295.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²     425.8µ ± ∞ ¹        ~ (p=0.700 n=3) ²    370.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²     403.9µ ± ∞ ¹        ~ (p=1.000 n=3) ²   307.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²     546.6µ ± ∞ ¹        ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90                        2.387m ± ∞ ¹    2.352m ± ∞ ¹       ~ (p=0.700 n=3) ²    1.945m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.455m ± ∞ ¹        ~ (p=0.400 n=3) ²    2.049m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.380m ± ∞ ¹        ~ (p=1.000 n=3) ²   1.914m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.773m ± ∞ ¹        ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-2                      2.539m ± ∞ ¹    2.612m ± ∞ ¹       ~ (p=0.200 n=3) ²    1.895m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.522m ± ∞ ¹        ~ (p=0.700 n=3) ²    2.108m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.509m ± ∞ ¹        ~ (p=1.000 n=3) ²   1.898m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.458m ± ∞ ¹        ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-4                      2.632m ± ∞ ¹    2.675m ± ∞ ¹       ~ (p=0.400 n=3) ²    2.068m ± ∞ ¹        ~ (p=0.400 n=3) ²     2.621m ± ∞ ¹        ~ (p=1.000 n=3) ²    2.157m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.571m ± ∞ ¹        ~ (p=0.700 n=3) ²   1.989m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.484m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-8                      2.923m ± ∞ ¹    3.013m ± ∞ ¹       ~ (p=0.200 n=3) ²    2.416m ± ∞ ¹        ~ (p=0.100 n=3) ²     2.909m ± ∞ ¹        ~ (p=0.200 n=3) ²    2.557m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.606m ± ∞ ¹        ~ (p=0.400 n=3) ²   2.173m ± ∞ ¹       ~ (p=0.100 n=3) ²     2.598m ± ∞ ¹        ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-16                     3.494m ± ∞ ¹    3.515m ± ∞ ¹       ~ (p=1.000 n=3) ²    2.655m ± ∞ ¹        ~ (p=0.100 n=3) ²     8.327m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.683m ± ∞ ¹       ~ (p=0.100 n=3) ²     4.946m ± ∞ ¹        ~ (p=0.700 n=3) ²   2.499m ± ∞ ¹       ~ (p=0.100 n=3) ²    10.548m ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10                  198.9µ ± ∞ ¹    224.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²    145.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²     282.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    230.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²     282.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   179.6µ ± ∞ ¹       ~ (p=0.100 n=3) ²     281.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-2                127.3µ ± ∞ ¹    140.0µ ± ∞ ¹       ~ (p=0.100 n=3) ²    110.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²     202.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    166.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²     205.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   126.7µ ± ∞ ¹       ~ (p=0.700 n=3) ²     184.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-4                106.7µ ± ∞ ¹    122.9µ ± ∞ ¹       ~ (p=0.700 n=3) ²    109.2µ ± ∞ ¹        ~ (p=1.000 n=3) ²     183.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    146.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²     170.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   120.1µ ± ∞ ¹       ~ (p=0.700 n=3) ²     170.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-8                101.8µ ± ∞ ¹    113.7µ ± ∞ ¹       ~ (p=0.200 n=3) ²    108.9µ ± ∞ ¹        ~ (p=0.400 n=3) ²     197.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    173.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²     199.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   230.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²     168.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-16               124.0µ ± ∞ ¹    130.3µ ± ∞ ¹       ~ (p=0.100 n=3) ²    134.7µ ± ∞ ¹        ~ (p=0.700 n=3) ²     227.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    220.9µ ± ∞ ¹       ~ (p=0.100 n=3) ²     226.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   146.1µ ± ∞ ¹       ~ (p=0.100 n=3) ²     185.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90                  878.3µ ± ∞ ¹    913.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²    719.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1103.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    863.3µ ± ∞ ¹       ~ (p=0.200 n=3) ²    1107.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   793.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²    1075.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-2                742.4µ ± ∞ ¹    799.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²    696.9µ ± ∞ ¹        ~ (p=0.700 n=3) ²    1125.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    909.1µ ± ∞ ¹       ~ (p=0.100 n=3) ²    1149.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   741.1µ ± ∞ ¹       ~ (p=1.000 n=3) ²    1019.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-4                763.0µ ± ∞ ¹    816.5µ ± ∞ ¹       ~ (p=0.700 n=3) ²    709.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    1128.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    984.6µ ± ∞ ¹       ~ (p=0.100 n=3) ²    1136.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   816.3µ ± ∞ ¹       ~ (p=0.400 n=3) ²    1093.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-8                845.1µ ± ∞ ¹    870.9µ ± ∞ ¹       ~ (p=0.700 n=3) ²    788.6µ ± ∞ ¹        ~ (p=0.400 n=3) ²    1282.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1071.9µ ± ∞ ¹       ~ (p=0.100 n=3) ²    1395.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   905.1µ ± ∞ ¹       ~ (p=0.100 n=3) ²    1209.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-16              1060.4µ ± ∞ ¹   1075.3µ ± ∞ ¹       ~ (p=0.700 n=3) ²   1127.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²   11574.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1398.8µ ± ∞ ¹       ~ (p=0.100 n=3) ²   10699.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²   958.4µ ± ∞ ¹       ~ (p=0.100 n=3) ²   10608.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                                                           574.5µ          596.1µ        +3.75%                    486.7µ        -15.28%                     847.1µ        +47.43%                    617.1µ        +7.40%                     814.5µ        +41.77%                   527.3µ        -8.22%                     817.3µ        +42.26%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                                                    │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt        │          bench_mattn_driver.txt          │        bench_modernc_driver.txt         │       bench_tailscale_driver.txt       │      bench_zombiezen_direct.txt       │
                                                                    │           B/op           │     B/op       vs base                 │     B/op       vs base                │     B/op       vs base                  │      B/op       vs base                  │     B/op       vs base                  │     B/op       vs base                 │     B/op       vs base                │
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10                       224.1Ki ± ∞ ¹   235.3Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.0Ki ± ∞ ¹       ~ (p=0.200 n=3) ²   455.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   455.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   232.4Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.9Ki ± ∞ ¹       ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-2                     226.1Ki ± ∞ ¹   234.9Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.1Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   453.7Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   455.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   232.3Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.4Ki ± ∞ ¹       ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-4                     224.3Ki ± ∞ ¹   235.0Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.2Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   452.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   454.8Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   232.0Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.7Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-8                     224.9Ki ± ∞ ¹   236.1Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.0Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   455.0Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.5Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   453.2Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   232.4Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.1Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-16                    225.3Ki ± ∞ ¹   235.6Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.0Ki ± ∞ ¹       ~ (p=0.700 n=3) ²   452.9Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    268.5Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   453.9Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   233.6Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.0Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90                       21.05Ki ± ∞ ¹   37.93Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   26.15Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   66.55Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.61Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   81.83Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   36.80Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   29.03Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-2                     24.69Ki ± ∞ ¹   40.85Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   29.85Ki ± ∞ ¹       ~ (p=0.200 n=3) ²   75.70Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.46Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   68.57Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   41.87Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   26.79Ki ± ∞ ¹       ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-4                     23.21Ki ± ∞ ¹   43.67Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.89Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   71.76Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.05Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   76.24Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   40.21Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   26.49Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-8                     28.15Ki ± ∞ ¹   40.66Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   30.99Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   65.55Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   304.46Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   73.25Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   42.00Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   30.46Ki ± ∞ ¹       ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-16                    27.45Ki ± ∞ ¹   37.00Ki ± ∞ ¹        ~ (p=0.200 n=3) ²   28.29Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   74.55Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   302.81Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   77.15Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   41.68Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   30.74Ki ± ∞ ¹       ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10                 225.7Ki ± ∞ ¹   237.8Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.4Ki ± ∞ ¹       ~ (p=1.000 n=3) ²   457.9Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   457.9Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   235.0Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.6Ki ± ∞ ¹       ~ (p=1.000 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-2               225.2Ki ± ∞ ¹   238.6Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.0Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   456.4Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   455.5Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   235.0Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.1Ki ± ∞ ¹       ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-4               224.1Ki ± ∞ ¹   237.4Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.2Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   455.0Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   457.5Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   235.1Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   224.7Ki ± ∞ ¹       ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-8               225.2Ki ± ∞ ¹   238.1Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.9Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   457.5Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.3Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   455.7Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   236.0Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.7Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-16              224.9Ki ± ∞ ¹   238.6Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   225.4Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   455.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²    271.1Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   454.8Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   235.1Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   226.1Ki ± ∞ ¹       ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90                 26.64Ki ± ∞ ¹   50.10Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.96Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   83.21Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   315.84Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   84.61Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   48.46Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   29.41Ki ± ∞ ¹       ~ (p=0.200 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-2               24.78Ki ± ∞ ¹   50.20Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.07Ki ± ∞ ¹       ~ (p=0.400 n=3) ²   85.48Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.28Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   83.75Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   48.52Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   28.04Ki ± ∞ ¹       ~ (p=0.400 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-4               24.34Ki ± ∞ ¹   49.95Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   25.87Ki ± ∞ ¹       ~ (p=0.200 n=3) ²   87.95Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.16Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   89.78Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   51.08Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   27.91Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-8               24.86Ki ± ∞ ¹   53.17Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   26.51Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   86.35Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.49Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   86.69Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   50.83Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   29.44Ki ± ∞ ¹       ~ (p=0.700 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-16              24.71Ki ± ∞ ¹   48.74Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   30.02Ki ± ∞ ¹       ~ (p=0.100 n=3) ²   76.41Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   316.14Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   95.14Ki ± ∞ ¹         ~ (p=0.100 n=3) ²   50.05Ki ± ∞ ¹        ~ (p=0.100 n=3) ²   35.77Ki ± ∞ ¹       ~ (p=0.100 n=3) ²
geomean                                                                          74.86Ki         103.1Ki        +37.68%                   79.62Ki        +6.36%                   187.2Ki        +150.00%                    289.2Ki        +286.37%                   192.5Ki        +157.11%                   102.5Ki        +36.86%                   81.35Ki        +8.67%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                                                    │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │        bench_glebarez_driver.txt         │         bench_mattn_driver.txt         │         bench_modernc_driver.txt         │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                                                                    │        allocs/op         │  allocs/op    vs base                  │  allocs/op    vs base                  │   allocs/op    vs base                   │  allocs/op    vs base                  │   allocs/op    vs base                   │  allocs/op    vs base                  │  allocs/op    vs base                  │
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10                         234.0 ± ∞ ¹    734.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     973.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    696.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     973.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    544.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    270.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-2                       236.0 ± ∞ ¹    734.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    696.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    545.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    270.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-4                       235.0 ± ∞ ¹    735.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    697.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-8                       235.0 ± ∞ ¹    737.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    697.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    546.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=10-16                      236.0 ± ∞ ¹    736.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     975.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    697.0 ± ∞ ¹         ~ (p=0.100 n=3) ²     974.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    547.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    271.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90                         22.00 ± ∞ ¹   442.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   802.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   421.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-2                       25.00 ± ∞ ¹   446.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   802.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    967.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   424.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-4                       24.00 ± ∞ ¹   451.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   800.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   423.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-8                       29.00 ± ∞ ¹   446.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   802.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   425.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   260.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsParallel/write_rate=90-16                      28.00 ± ∞ ¹   441.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   797.00 ± ∞ ¹         ~ (p=0.100 n=3) ²    968.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   424.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   261.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10                   236.0 ± ∞ ¹    777.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1009.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1013.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    585.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-2                 235.0 ± ∞ ¹    777.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1010.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    585.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-4                 234.0 ± ∞ ¹    776.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1011.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1013.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-8                 235.0 ± ∞ ¹    778.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    258.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1010.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    747.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1014.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    585.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    274.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=10-16                235.0 ± ∞ ¹    778.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    257.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1011.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    746.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    1015.0 ± ∞ ¹          ~ (p=0.100 n=3) ²    586.0 ± ∞ ¹         ~ (p=0.100 n=3) ²    273.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90                   27.00 ± ∞ ¹   596.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1120.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   966.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1124.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   576.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-2                 25.00 ± ∞ ¹   596.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1120.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   969.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1124.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-4                 25.00 ± ∞ ¹   596.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1119.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   968.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1123.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-8                 26.00 ± ∞ ¹   599.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   208.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1120.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   970.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1124.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   263.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
ReadWrite/ReadOrWritePostAndCommentsWithTxParallel/write_rate=90-16                25.00 ± ∞ ¹   595.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   209.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1123.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   968.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   1122.00 ± ∞ ¹          ~ (p=0.100 n=3) ²   577.00 ± ∞ ¹         ~ (p=0.100 n=3) ²   264.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                                                            77.47          624.1        +705.59%                    231.8        +199.16%                    1.016k        +1212.01%                    796.7        +928.41%                    1.018k        +1213.78%                    528.5        +582.16%                    266.7        +244.32%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
```
<!--END_BENCHMARK-->

## Query - Correlated (Sequential)

<!--BENCHMARK:slow/benchstat_query_correlated.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                            │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                            │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │
Query/Correlated-12                       650.0m ± ∞ ¹    644.7m ± ∞ ¹       ~ (p=1.000 n=3) ²   288.3m ± ∞ ¹        ~ (p=0.100 n=3) ²    303.6m ± ∞ ¹        ~ (p=0.100 n=3) ²   175.0m ± ∞ ¹        ~ (p=0.100 n=3) ²    305.1m ± ∞ ¹        ~ (p=0.100 n=3) ²   172.2m ± ∞ ¹        ~ (p=0.100 n=3) ²    325.9m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-12              118.49m ± ∞ ¹   119.66m ± ∞ ¹       ~ (p=0.100 n=3) ²   50.17m ± ∞ ¹        ~ (p=0.100 n=3) ²   118.92m ± ∞ ¹        ~ (p=0.400 n=3) ²   44.91m ± ∞ ¹        ~ (p=0.100 n=3) ²   115.22m ± ∞ ¹        ~ (p=0.700 n=3) ²   53.44m ± ∞ ¹        ~ (p=0.100 n=3) ²   190.64m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                   277.5m          277.7m        +0.08%                   120.3m        -56.66%                    190.0m        -31.53%                   88.66m        -68.05%                    187.5m        -32.44%                   95.93m        -65.43%                    249.3m        -10.19%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt          │      bench_eatonphil_direct.txt      │         bench_glebarez_driver.txt          │            bench_mattn_driver.txt            │          bench_modernc_driver.txt          │         bench_tailscale_driver.txt         │       bench_zombiezen_direct.txt        │
                            │           B/op           │      B/op       vs base                    │    B/op      vs base                 │      B/op       vs base                    │      B/op        vs base                     │      B/op       vs base                    │      B/op       vs base                    │     B/op       vs base                  │
Query/Correlated-12                        96.00 ± ∞ ¹   71168.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   48.00 ± ∞ ¹        ~ (p=0.300 n=3) ²   63068.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   207200.00 ± ∞ ¹            ~ (p=0.100 n=3) ²   63092.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   47072.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   2732.00 ± ∞ ¹         ~ (p=0.600 n=3) ²
Query/CorrelatedParallel-12                430.0 ± ∞ ¹    71086.0 ± ∞ ¹           ~ (p=0.100 n=3) ²   294.0 ± ∞ ¹        ~ (p=0.400 n=3) ²    63728.0 ± ∞ ¹           ~ (p=0.700 n=3) ²    207937.0 ± ∞ ¹            ~ (p=0.100 n=3) ²    63178.0 ± ∞ ¹           ~ (p=0.700 n=3) ²    48384.0 ± ∞ ¹           ~ (p=0.700 n=3) ²     769.0 ± ∞ ¹         ~ (p=0.700 n=3) ²
geomean                                    203.2          69.46Ki        +34907.78%                   118.8        -41.53%                    61.91Ki        +31103.25%                     202.7Ki        +102062.36%                    61.66Ki        +30974.22%                    46.60Ki        +23388.88%                   1.415Ki        +613.40%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt           │    bench_eatonphil_direct.txt    │          bench_glebarez_driver.txt          │           bench_mattn_driver.txt            │          bench_modernc_driver.txt           │         bench_tailscale_driver.txt          │       bench_zombiezen_direct.txt        │
                            │        allocs/op         │   allocs/op     vs base                     │  allocs/op   vs base             │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                     │  allocs/op    vs base                   │
Query/Correlated-12                        1.000 ± ∞ ¹   5769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.300 n=3) ²     6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   29.000 ± ∞ ¹          ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-12                3.000 ± ∞ ¹   5770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.300 n=3) ²     6779.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6773.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   13.000 ± ∞ ¹          ~ (p=0.700 n=3) ²
geomean                                    1.732           5.769k        +333002.24%                                ?               ³ ⁴     6.774k        +390996.97%                     6.772k        +390881.60%                     6.770k        +390766.13%                     3.768k        +217445.57%                    19.42        +1021.01%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean
⁴ ratios must be >0 to compute geomean
```
<!--END_BENCHMARK-->

## Query - Correlated (Parallel)

<!--BENCHMARK:slow/benchstat_query_correlated_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                            │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                            │          sec/op          │    sec/op      vs base                │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │
Query/CorrelatedParallel                  643.9m ± ∞ ¹    649.3m ± ∞ ¹       ~ (p=0.100 n=3) ²    287.5m ± ∞ ¹        ~ (p=0.100 n=3) ²    307.7m ± ∞ ¹        ~ (p=0.100 n=3) ²   174.6m ± ∞ ¹        ~ (p=0.100 n=3) ²    302.4m ± ∞ ¹        ~ (p=0.100 n=3) ²   173.2m ± ∞ ¹        ~ (p=0.100 n=3) ²    323.7m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-2               435.50m ± ∞ ¹   436.30m ± ∞ ¹       ~ (p=1.000 n=3) ²   171.02m ± ∞ ¹        ~ (p=0.100 n=3) ²   183.41m ± ∞ ¹        ~ (p=0.100 n=3) ²   97.42m ± ∞ ¹        ~ (p=0.100 n=3) ²   184.56m ± ∞ ¹        ~ (p=0.100 n=3) ²   94.18m ± ∞ ¹        ~ (p=0.100 n=3) ²   188.54m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-4               224.22m ± ∞ ¹   226.14m ± ∞ ¹       ~ (p=0.200 n=3) ²    93.17m ± ∞ ¹        ~ (p=0.100 n=3) ²   117.13m ± ∞ ¹        ~ (p=0.100 n=3) ²   59.58m ± ∞ ¹        ~ (p=0.100 n=3) ²   116.34m ± ∞ ¹        ~ (p=0.100 n=3) ²   55.90m ± ∞ ¹        ~ (p=0.100 n=3) ²   118.68m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-8               153.72m ± ∞ ¹   154.07m ± ∞ ¹       ~ (p=1.000 n=3) ²    54.96m ± ∞ ¹        ~ (p=0.100 n=3) ²   137.36m ± ∞ ¹        ~ (p=0.100 n=3) ²   42.34m ± ∞ ¹        ~ (p=0.100 n=3) ²   135.71m ± ∞ ¹        ~ (p=0.100 n=3) ²   44.18m ± ∞ ¹        ~ (p=0.100 n=3) ²   192.56m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-16               77.61m ± ∞ ¹    78.47m ± ∞ ¹       ~ (p=0.700 n=3) ²    54.36m ± ∞ ¹        ~ (p=0.100 n=3) ²   115.70m ± ∞ ¹        ~ (p=0.100 n=3) ²   48.45m ± ∞ ¹        ~ (p=0.100 n=3) ²   114.09m ± ∞ ¹        ~ (p=0.100 n=3) ²   43.05m ± ∞ ¹        ~ (p=0.100 n=3) ²   184.13m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                   237.1m          238.7m        +0.64%                    106.5m        -55.10%                    160.1m        -32.51%                   73.04m        -69.20%                    158.7m        -33.10%                   70.44m        -70.30%                    191.4m        -19.29%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt          │      bench_eatonphil_direct.txt      │         bench_glebarez_driver.txt          │           bench_mattn_driver.txt            │          bench_modernc_driver.txt          │         bench_tailscale_driver.txt         │       bench_zombiezen_direct.txt        │
                            │           B/op           │      B/op       vs base                    │    B/op      vs base                 │      B/op       vs base                    │      B/op        vs base                    │      B/op       vs base                    │      B/op       vs base                    │     B/op       vs base                  │
Query/CorrelatedParallel                   68.00 ± ∞ ¹   71132.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   34.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   63050.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   207185.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   63066.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   47009.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   2216.00 ± ∞ ¹         ~ (p=0.600 n=3) ²
Query/CorrelatedParallel-2                 72.00 ± ∞ ¹   71082.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   30.00 ± ∞ ¹        ~ (p=0.100 n=3) ²   63025.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   207232.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   63027.00 ± ∞ ¹           ~ (p=0.100 n=3) ²   47002.00 ± ∞ ¹           ~ (p=0.100 n=3) ²    460.00 ± ∞ ¹         ~ (p=0.600 n=3) ²
Query/CorrelatedParallel-4                 153.0 ± ∞ ¹    71012.0 ± ∞ ¹           ~ (p=0.100 n=3) ²   147.0 ± ∞ ¹        ~ (p=0.700 n=3) ²    63076.0 ± ∞ ¹           ~ (p=0.100 n=3) ²    207210.0 ± ∞ ¹           ~ (p=0.100 n=3) ²    63151.0 ± ∞ ¹           ~ (p=0.100 n=3) ²    46966.0 ± ∞ ¹           ~ (p=0.100 n=3) ²     570.0 ± ∞ ¹         ~ (p=0.700 n=3) ²
Query/CorrelatedParallel-8                1720.0 ± ∞ ¹    71050.0 ± ∞ ¹           ~ (p=0.100 n=3) ²   427.0 ± ∞ ¹        ~ (p=0.700 n=3) ²    63199.0 ± ∞ ¹           ~ (p=0.400 n=3) ²    207289.0 ± ∞ ¹           ~ (p=0.100 n=3) ²    63229.0 ± ∞ ¹           ~ (p=0.700 n=3) ²    47217.0 ± ∞ ¹           ~ (p=0.700 n=3) ²     913.0 ± ∞ ¹         ~ (p=0.700 n=3) ²
Query/CorrelatedParallel-16                998.0 ± ∞ ¹    97498.0 ± ∞ ¹           ~ (p=0.200 n=3) ²   105.0 ± ∞ ¹        ~ (p=0.400 n=3) ²    63810.0 ± ∞ ¹           ~ (p=0.700 n=3) ²    207778.0 ± ∞ ¹           ~ (p=0.100 n=3) ²    63308.0 ± ∞ ¹           ~ (p=0.700 n=3) ²    49608.0 ± ∞ ¹           ~ (p=0.700 n=3) ²     779.0 ± ∞ ¹         ~ (p=1.000 n=3) ²
geomean                                    264.1          73.93Ki        +28561.89%                   92.36        -65.03%                    61.75Ki        +23838.33%                     202.5Ki        +78394.99%                    61.68Ki        +23809.86%                    46.44Ki        +17901.46%                     838.0        +217.25%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                            │ bench_ncruces_direct.txt │          bench_ncruces_driver.txt           │      bench_eatonphil_direct.txt      │          bench_glebarez_driver.txt          │           bench_mattn_driver.txt            │          bench_modernc_driver.txt           │         bench_tailscale_driver.txt          │       bench_zombiezen_direct.txt       │
                            │        allocs/op         │   allocs/op     vs base                     │  allocs/op   vs base                 │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                     │   allocs/op     vs base                     │  allocs/op    vs base                  │
Query/CorrelatedParallel                   3.000 ± ∞ ¹   5771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   29.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-2                 2.000 ± ∞ ¹   5770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹            ~ (p=0.100 n=3) ²    7.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-4                 2.000 ± ∞ ¹   5769.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.400 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹            ~ (p=0.100 n=3) ²    7.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/CorrelatedParallel-8                 3.000 ± ∞ ¹   5770.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6773.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6772.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3767.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   10.000 ± ∞ ¹         ~ (p=0.600 n=3) ²
Query/CorrelatedParallel-16                4.000 ± ∞ ¹   5828.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹        ~ (p=0.300 n=3) ²   6781.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6774.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   6773.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   3771.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   12.000 ± ∞ ¹         ~ (p=0.700 n=3) ²
geomean                                    2.702           5.782k        +213879.45%                   1.000        -62.99%                     6.772k        +250543.93%                     6.773k        +250558.78%                     6.771k        +250492.16%                     3.768k        +139341.56%                    11.13        +311.80%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
```
<!--END_BENCHMARK-->

## Query - GroupBy (Sequential)

<!--BENCHMARK:slow/benchstat_query_groupby.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                         │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │       bench_eatonphil_direct.txt        │        bench_glebarez_driver.txt        │         bench_mattn_driver.txt          │        bench_modernc_driver.txt         │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt        │
                         │          sec/op          │    sec/op     vs base                │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op     vs base                 │    sec/op      vs base                  │
Query/GroupBy-12                       556.8µ ± ∞ ¹   570.5µ ± ∞ ¹       ~ (p=0.100 n=3) ²    327.7µ ± ∞ ¹         ~ (p=0.100 n=3) ²    490.8µ ± ∞ ¹         ~ (p=0.100 n=3) ²    340.1µ ± ∞ ¹         ~ (p=0.100 n=3) ²    500.7µ ± ∞ ¹         ~ (p=0.100 n=3) ²   260.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    478.8µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-12               61.23µ ± ∞ ¹   64.92µ ± ∞ ¹       ~ (p=0.100 n=3) ²   737.24µ ± ∞ ¹         ~ (p=0.100 n=3) ²   660.68µ ± ∞ ¹         ~ (p=0.100 n=3) ²   482.83µ ± ∞ ¹         ~ (p=0.100 n=3) ²   670.14µ ± ∞ ¹         ~ (p=0.100 n=3) ²   27.91µ ± ∞ ¹        ~ (p=0.100 n=3) ²   664.87µ ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                184.6µ         192.4µ        +4.22%                    491.5µ        +166.19%                    569.4µ        +208.39%                    405.2µ        +119.47%                    579.3µ        +213.74%                   85.33µ        -53.78%                    564.2µ        +205.56%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │   bench_eatonphil_direct.txt   │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                         │           B/op           │      B/op       vs base           │    B/op      vs base           │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │     B/op       vs base           │
Query/GroupBy-12                          0.0 ± ∞ ¹     2376.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ²     1969.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     7160.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1976.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1583.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     361.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-12                0.000 ± ∞ ¹   2368.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1965.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7166.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1982.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1529.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   371.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ³    2.316Ki        ?                                ?               ³    1.921Ki        ?                    6.995Ki        ?                    1.933Ki        ?                    1.519Ki        ?                     366.0        ?
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

## Query - GroupBy (Parallel)

<!--BENCHMARK:slow/benchstat_query_groupby_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                         │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt        │        bench_glebarez_driver.txt        │         bench_mattn_driver.txt          │        bench_modernc_driver.txt         │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt        │
                         │          sec/op          │    sec/op      vs base                │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op      vs base                  │    sec/op     vs base                 │    sec/op      vs base                  │
Query/GroupByParallel                  551.1µ ± ∞ ¹    558.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²    338.1µ ± ∞ ¹         ~ (p=0.100 n=3) ²    477.5µ ± ∞ ¹         ~ (p=0.100 n=3) ²    308.8µ ± ∞ ¹         ~ (p=0.100 n=3) ²    477.9µ ± ∞ ¹         ~ (p=0.100 n=3) ²   265.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    463.7µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-2                279.9µ ± ∞ ¹    286.7µ ± ∞ ¹       ~ (p=0.100 n=3) ²    357.1µ ± ∞ ¹         ~ (p=0.100 n=3) ²    384.3µ ± ∞ ¹         ~ (p=0.100 n=3) ²    401.2µ ± ∞ ¹         ~ (p=0.100 n=3) ²    377.6µ ± ∞ ¹         ~ (p=0.100 n=3) ²   131.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    370.2µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-4               143.96µ ± ∞ ¹   147.07µ ± ∞ ¹       ~ (p=0.100 n=3) ²   734.55µ ± ∞ ¹         ~ (p=0.100 n=3) ²   369.62µ ± ∞ ¹         ~ (p=0.100 n=3) ²   504.77µ ± ∞ ¹         ~ (p=0.100 n=3) ²   374.50µ ± ∞ ¹         ~ (p=0.100 n=3) ²   65.98µ ± ∞ ¹        ~ (p=0.100 n=3) ²   360.51µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-8                72.61µ ± ∞ ¹    75.56µ ± ∞ ¹       ~ (p=0.100 n=3) ²   799.80µ ± ∞ ¹         ~ (p=0.100 n=3) ²   640.18µ ± ∞ ¹         ~ (p=0.100 n=3) ²   526.11µ ± ∞ ¹         ~ (p=0.100 n=3) ²   647.80µ ± ∞ ¹         ~ (p=0.100 n=3) ²   34.37µ ± ∞ ¹        ~ (p=0.100 n=3) ²   640.74µ ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/GroupByParallel-16               61.62µ ± ∞ ¹    64.90µ ± ∞ ¹       ~ (p=0.100 n=3) ²   721.82µ ± ∞ ¹         ~ (p=0.100 n=3) ²   661.87µ ± ∞ ¹         ~ (p=0.100 n=3) ²   484.17µ ± ∞ ¹         ~ (p=0.100 n=3) ²   669.78µ ± ∞ ¹         ~ (p=0.100 n=3) ²   28.25µ ± ∞ ¹        ~ (p=0.100 n=3) ²   661.29µ ± ∞ ¹         ~ (p=0.100 n=3) ²
geomean                                158.3µ          163.1µ        +3.06%                    551.9µ        +248.70%                    491.7µ        +210.66%                    437.0µ        +176.08%                    493.7µ        +211.89%                   74.10µ        -53.18%                    482.8µ        +205.01%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │   bench_eatonphil_direct.txt   │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                         │           B/op           │      B/op       vs base           │    B/op      vs base           │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │     B/op       vs base           │
Query/GroupByParallel                     0.0 ± ∞ ¹     2364.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ³     1849.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     7160.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1864.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1581.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     364.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-2                   0.0 ± ∞ ¹     2364.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ³     1960.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     7154.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1916.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1581.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     360.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-4                   0.0 ± ∞ ¹     2364.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ³     1961.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     7135.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1937.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1539.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     362.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-8                 0.000 ± ∞ ¹   2366.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1963.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7160.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1980.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1567.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   362.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-16                0.000 ± ∞ ¹   2368.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1969.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   7183.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1985.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1515.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   370.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ⁴    2.310Ki        ?                                ?               ⁴    1.894Ki        ?                    6.991Ki        ?                    1.891Ki        ?                    1.520Ki        ?                     363.6        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                         │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                         │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/GroupByParallel                   0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   120.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   155.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   120.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   52.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-2                 0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   154.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   121.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   52.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-4                 0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   154.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   121.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-8                 0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   155.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/GroupByParallel-16                0.000 ± ∞ ¹   118.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   155.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   122.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   51.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                           ⁴     118.0        ?                                +0.00%               ⁴     121.6        ?                     154.6        ?                     121.2        ?                    51.40        ?                   6.000        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

## Query - JSON (Sequential)

<!--BENCHMARK:slow/benchstat_query_json_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                      │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt        │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                      │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │
Query/JSONParallel                  9.659m ± ∞ ¹   10.091m ± ∞ ¹       ~ (p=0.400 n=3) ²   6.922m ± ∞ ¹        ~ (p=0.100 n=3) ²    9.483m ± ∞ ¹        ~ (p=0.100 n=3) ²   7.649m ± ∞ ¹       ~ (p=0.100 n=3) ²    9.415m ± ∞ ¹        ~ (p=0.400 n=3) ²   7.081m ± ∞ ¹        ~ (p=0.100 n=3) ²    9.919m ± ∞ ¹        ~ (p=0.400 n=3) ²
Query/JSONParallel-2                6.982m ± ∞ ¹    7.029m ± ∞ ¹       ~ (p=0.700 n=3) ²   5.454m ± ∞ ¹        ~ (p=0.100 n=3) ²    7.209m ± ∞ ¹        ~ (p=0.700 n=3) ²   6.090m ± ∞ ¹       ~ (p=0.100 n=3) ²    7.419m ± ∞ ¹        ~ (p=0.700 n=3) ²   5.606m ± ∞ ¹        ~ (p=0.100 n=3) ²    8.349m ± ∞ ¹        ~ (p=0.700 n=3) ²
Query/JSONParallel-4                4.542m ± ∞ ¹    4.419m ± ∞ ¹       ~ (p=0.200 n=3) ²   3.753m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.454m ± ∞ ¹        ~ (p=0.100 n=3) ²   4.227m ± ∞ ¹       ~ (p=0.100 n=3) ²    5.466m ± ∞ ¹        ~ (p=0.100 n=3) ²   4.157m ± ∞ ¹        ~ (p=0.100 n=3) ²    7.550m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-8                7.195m ± ∞ ¹    7.126m ± ∞ ¹       ~ (p=0.200 n=3) ²   6.686m ± ∞ ¹        ~ (p=0.100 n=3) ²   13.704m ± ∞ ¹        ~ (p=0.100 n=3) ²   7.299m ± ∞ ¹       ~ (p=0.700 n=3) ²   13.514m ± ∞ ¹        ~ (p=0.100 n=3) ²   7.153m ± ∞ ¹        ~ (p=0.400 n=3) ²   14.727m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-16               8.999m ± ∞ ¹    9.005m ± ∞ ¹       ~ (p=1.000 n=3) ²   9.341m ± ∞ ¹        ~ (p=0.400 n=3) ²   15.119m ± ∞ ¹        ~ (p=0.100 n=3) ²   9.070m ± ∞ ¹       ~ (p=0.100 n=3) ²   15.334m ± ∞ ¹        ~ (p=0.100 n=3) ²   9.295m ± ∞ ¹        ~ (p=0.100 n=3) ²   13.293m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                             7.236m          7.256m        +0.28%                   6.157m        -14.91%                    9.497m        +31.25%                   6.653m        -8.05%                    9.542m        +31.88%                   6.428m        -11.17%                    10.41m        +43.91%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                      │ bench_ncruces_direct.txt │           bench_ncruces_driver.txt            │    bench_eatonphil_direct.txt    │           bench_glebarez_driver.txt           │             bench_mattn_driver.txt             │           bench_modernc_driver.txt            │          bench_tailscale_driver.txt           │        bench_zombiezen_direct.txt         │
                      │           B/op           │      B/op        vs base                      │    B/op      vs base             │      B/op        vs base                      │       B/op        vs base                      │      B/op        vs base                      │      B/op        vs base                      │     B/op       vs base                    │
Query/JSONParallel                   1.000 ± ∞ ¹   64845.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     56879.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   201019.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56859.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   32891.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   364.000 ± ∞ ¹           ~ (p=0.100 n=3) ²
Query/JSONParallel-2                 1.000 ± ∞ ¹   64844.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.600 n=3) ²     56952.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   201015.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56897.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   32904.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   364.000 ± ∞ ¹           ~ (p=0.600 n=3) ²
Query/JSONParallel-4                 1.000 ± ∞ ¹   64852.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.700 n=3) ²     56953.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   201014.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56922.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   32890.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   378.000 ± ∞ ¹           ~ (p=0.100 n=3) ²
Query/JSONParallel-8                 6.000 ± ∞ ¹   64849.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   4.000 ± ∞ ¹  ~ (p=0.300 n=3) ²     56981.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   201069.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   57005.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   32932.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   411.000 ± ∞ ¹           ~ (p=0.100 n=3) ²
Query/JSONParallel-16                61.00 ± ∞ ¹    64942.00 ± ∞ ¹             ~ (p=0.100 n=3) ²   64.00 ± ∞ ¹  ~ (p=0.700 n=3) ²      57106.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    201147.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    57047.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    32932.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    487.00 ± ∞ ¹           ~ (p=0.100 n=3) ²
geomean                              3.256           63.35Ki        +1992056.74%                                ?               ³ ⁴     55.64Ki        +1749672.76%                      196.3Ki        +6174572.03%                     55.61Ki        +1748806.86%                     32.14Ki        +1010615.57%                     398.3        +12132.55%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean
⁴ ratios must be >0 to compute geomean

                      │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │   bench_zombiezen_direct.txt   │
                      │        allocs/op         │   allocs/op     vs base           │  allocs/op   vs base                │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │  allocs/op   vs base           │
Query/JSONParallel                   0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-2                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.400 n=3) ²
Query/JSONParallel-4                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-8                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-16                0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ²   4023.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4023.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                        ⁴     4.019k        ?                                +0.00%               ⁴     4.022k        ?                     5.021k        ?                     4.022k        ?                     2.020k        ?                   6.000        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

## Query - JSON (Parallel)

<!--BENCHMARK:slow/benchstat_query_json_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                      │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt        │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                      │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │
Query/JSONParallel                  9.659m ± ∞ ¹   10.091m ± ∞ ¹       ~ (p=0.400 n=3) ²   6.922m ± ∞ ¹        ~ (p=0.100 n=3) ²    9.483m ± ∞ ¹        ~ (p=0.100 n=3) ²   7.649m ± ∞ ¹       ~ (p=0.100 n=3) ²    9.415m ± ∞ ¹        ~ (p=0.400 n=3) ²   7.081m ± ∞ ¹        ~ (p=0.100 n=3) ²    9.919m ± ∞ ¹        ~ (p=0.400 n=3) ²
Query/JSONParallel-2                6.982m ± ∞ ¹    7.029m ± ∞ ¹       ~ (p=0.700 n=3) ²   5.454m ± ∞ ¹        ~ (p=0.100 n=3) ²    7.209m ± ∞ ¹        ~ (p=0.700 n=3) ²   6.090m ± ∞ ¹       ~ (p=0.100 n=3) ²    7.419m ± ∞ ¹        ~ (p=0.700 n=3) ²   5.606m ± ∞ ¹        ~ (p=0.100 n=3) ²    8.349m ± ∞ ¹        ~ (p=0.700 n=3) ²
Query/JSONParallel-4                4.542m ± ∞ ¹    4.419m ± ∞ ¹       ~ (p=0.200 n=3) ²   3.753m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.454m ± ∞ ¹        ~ (p=0.100 n=3) ²   4.227m ± ∞ ¹       ~ (p=0.100 n=3) ²    5.466m ± ∞ ¹        ~ (p=0.100 n=3) ²   4.157m ± ∞ ¹        ~ (p=0.100 n=3) ²    7.550m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-8                7.195m ± ∞ ¹    7.126m ± ∞ ¹       ~ (p=0.200 n=3) ²   6.686m ± ∞ ¹        ~ (p=0.100 n=3) ²   13.704m ± ∞ ¹        ~ (p=0.100 n=3) ²   7.299m ± ∞ ¹       ~ (p=0.700 n=3) ²   13.514m ± ∞ ¹        ~ (p=0.100 n=3) ²   7.153m ± ∞ ¹        ~ (p=0.400 n=3) ²   14.727m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/JSONParallel-16               8.999m ± ∞ ¹    9.005m ± ∞ ¹       ~ (p=1.000 n=3) ²   9.341m ± ∞ ¹        ~ (p=0.400 n=3) ²   15.119m ± ∞ ¹        ~ (p=0.100 n=3) ²   9.070m ± ∞ ¹       ~ (p=0.100 n=3) ²   15.334m ± ∞ ¹        ~ (p=0.100 n=3) ²   9.295m ± ∞ ¹        ~ (p=0.100 n=3) ²   13.293m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                             7.236m          7.256m        +0.28%                   6.157m        -14.91%                    9.497m        +31.25%                   6.653m        -8.05%                    9.542m        +31.88%                   6.428m        -11.17%                    10.41m        +43.91%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                      │ bench_ncruces_direct.txt │           bench_ncruces_driver.txt            │    bench_eatonphil_direct.txt    │           bench_glebarez_driver.txt           │             bench_mattn_driver.txt             │           bench_modernc_driver.txt            │          bench_tailscale_driver.txt           │        bench_zombiezen_direct.txt         │
                      │           B/op           │      B/op        vs base                      │    B/op      vs base             │      B/op        vs base                      │       B/op        vs base                      │      B/op        vs base                      │      B/op        vs base                      │     B/op       vs base                    │
Query/JSONParallel                   1.000 ± ∞ ¹   64845.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     56879.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   201019.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56859.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   32891.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   364.000 ± ∞ ¹           ~ (p=0.100 n=3) ²
Query/JSONParallel-2                 1.000 ± ∞ ¹   64844.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.600 n=3) ²     56952.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   201015.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56897.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   32904.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   364.000 ± ∞ ¹           ~ (p=0.600 n=3) ²
Query/JSONParallel-4                 1.000 ± ∞ ¹   64852.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.700 n=3) ²     56953.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   201014.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   56922.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   32890.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   378.000 ± ∞ ¹           ~ (p=0.100 n=3) ²
Query/JSONParallel-8                 6.000 ± ∞ ¹   64849.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   4.000 ± ∞ ¹  ~ (p=0.300 n=3) ²     56981.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   201069.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   57005.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   32932.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   411.000 ± ∞ ¹           ~ (p=0.100 n=3) ²
Query/JSONParallel-16                61.00 ± ∞ ¹    64942.00 ± ∞ ¹             ~ (p=0.100 n=3) ²   64.00 ± ∞ ¹  ~ (p=0.700 n=3) ²      57106.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    201147.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    57047.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    32932.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    487.00 ± ∞ ¹           ~ (p=0.100 n=3) ²
geomean                              3.256           63.35Ki        +1992056.74%                                ?               ³ ⁴     55.64Ki        +1749672.76%                      196.3Ki        +6174572.03%                     55.61Ki        +1748806.86%                     32.14Ki        +1010615.57%                     398.3        +12132.55%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean
⁴ ratios must be >0 to compute geomean

                      │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │   bench_zombiezen_direct.txt   │
                      │        allocs/op         │   allocs/op     vs base           │  allocs/op   vs base                │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │  allocs/op   vs base           │
Query/JSONParallel                   0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-2                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.400 n=3) ²
Query/JSONParallel-4                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-8                 0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4022.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/JSONParallel-16                0.000 ± ∞ ¹   4019.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ²   4023.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   5021.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4023.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2020.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                        ⁴     4.019k        ?                                +0.00%               ⁴     4.022k        ?                     5.021k        ?                     4.022k        ?                     2.020k        ?                   6.000        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

## Query - NonRecursiveCTE (Sequential)

<!--BENCHMARK:slow/benchstat_query_nonrecursivecte_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                                 │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │         bench_mattn_driver.txt         │        bench_modernc_driver.txt        │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                                 │          sec/op          │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │
Query/NonRecursiveCTEParallel                 2100.4µ ± ∞ ¹   2377.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    993.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1894.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1593.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1882.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1049.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1611.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-2               1060.2µ ± ∞ ¹   1217.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    617.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    972.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    828.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    975.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    537.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    833.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-4                542.4µ ± ∞ ¹    629.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    888.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    525.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    470.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    531.1µ ± ∞ ¹        ~ (p=0.400 n=3) ²    280.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    452.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-8                272.6µ ± ∞ ¹    343.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1636.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    877.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    992.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    871.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    149.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    842.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-16               231.7µ ± ∞ ¹    313.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1448.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    918.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1027.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    906.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    138.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    854.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                        597.7µ          722.3µ        +20.84%                    1.053m        +76.13%                    951.4µ        +59.18%                    912.8µ        +52.71%                    949.3µ        +58.82%                    318.6µ        -46.70%                    847.7µ        +41.82%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                 │ bench_ncruces_direct.txt │      bench_ncruces_driver.txt      │   bench_eatonphil_direct.txt   │     bench_glebarez_driver.txt      │       bench_mattn_driver.txt        │      bench_modernc_driver.txt      │     bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                                 │           B/op           │      B/op        vs base           │    B/op      vs base           │      B/op        vs base           │       B/op        vs base           │      B/op        vs base           │      B/op        vs base           │     B/op       vs base           │
Query/NonRecursiveCTEParallel                     0.0 ± ∞ ¹     62762.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ³     54775.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198936.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54788.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30792.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     360.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-2                   0.0 ± ∞ ¹     62761.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ³     54821.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198838.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54812.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30793.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     361.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-4                 0.000 ± ∞ ¹   62762.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.400 n=3) ²   54878.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   198852.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54883.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   30793.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   363.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-8                 0.000 ± ∞ ¹   62766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   8.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54881.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   198906.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54906.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   30794.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   364.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-16                 0.00 ± ∞ ¹    62768.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   11.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    54893.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    198971.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    54922.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    30772.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    375.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                   ⁴     61.29Ki        ?                                ?               ⁴     53.56Ki        ?                      194.2Ki        ?                     53.58Ki        ?                     30.07Ki        ?                     364.6        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                                 │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │   bench_zombiezen_direct.txt   │
                                 │        allocs/op         │   allocs/op     vs base           │  allocs/op   vs base                │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │  allocs/op   vs base           │
Query/NonRecursiveCTEParallel                   0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-2                 0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-4                 0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-8                 0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-16                0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                   ⁴     3.763k        ?                                +0.00%               ⁴     3.766k        ?                     4.764k        ?                     3.765k        ?                     1.764k        ?                   6.000        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

## Query - NonRecursiveCTE (Parallel)

<!--BENCHMARK:slow/benchstat_query_nonrecursivecte_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                                 │ bench_ncruces_direct.txt │        bench_ncruces_driver.txt        │       bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │         bench_mattn_driver.txt         │        bench_modernc_driver.txt        │       bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                                 │          sec/op          │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │    sec/op      vs base                 │
Query/NonRecursiveCTEParallel                 2100.4µ ± ∞ ¹   2377.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    993.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1894.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1593.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1882.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1049.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1611.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-2               1060.2µ ± ∞ ¹   1217.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    617.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    972.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    828.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    975.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    537.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    833.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-4                542.4µ ± ∞ ¹    629.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    888.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    525.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    470.2µ ± ∞ ¹        ~ (p=0.100 n=3) ²    531.1µ ± ∞ ¹        ~ (p=0.400 n=3) ²    280.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    452.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-8                272.6µ ± ∞ ¹    343.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1636.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    877.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    992.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²    871.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²    149.5µ ± ∞ ¹        ~ (p=0.100 n=3) ²    842.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-16               231.7µ ± ∞ ¹    313.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1448.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    918.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1027.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²    906.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    138.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²    854.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                        597.7µ          722.3µ        +20.84%                    1.053m        +76.13%                    951.4µ        +59.18%                    912.8µ        +52.71%                    949.3µ        +58.82%                    318.6µ        -46.70%                    847.7µ        +41.82%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                                 │ bench_ncruces_direct.txt │      bench_ncruces_driver.txt      │   bench_eatonphil_direct.txt   │     bench_glebarez_driver.txt      │       bench_mattn_driver.txt        │      bench_modernc_driver.txt      │     bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                                 │           B/op           │      B/op        vs base           │    B/op      vs base           │      B/op        vs base           │       B/op        vs base           │      B/op        vs base           │      B/op        vs base           │     B/op       vs base           │
Query/NonRecursiveCTEParallel                     0.0 ± ∞ ¹     62762.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ³     54775.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198936.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54788.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30792.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     360.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-2                   0.0 ± ∞ ¹     62761.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ³     54821.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     198838.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     54812.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     30793.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     361.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-4                 0.000 ± ∞ ¹   62762.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.400 n=3) ²   54878.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   198852.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54883.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   30793.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   363.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-8                 0.000 ± ∞ ¹   62766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   8.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54881.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   198906.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   54906.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   30794.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   364.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-16                 0.00 ± ∞ ¹    62768.00 ± ∞ ¹  ~ (p=0.100 n=3) ²   11.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    54893.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    198971.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    54922.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    30772.00 ± ∞ ¹  ~ (p=0.100 n=3) ²    375.00 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                   ⁴     61.29Ki        ?                                ?               ⁴     53.56Ki        ?                      194.2Ki        ?                     53.58Ki        ?                     30.07Ki        ?                     364.6        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                                 │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │     bench_eatonphil_direct.txt      │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │   bench_zombiezen_direct.txt   │
                                 │        allocs/op         │   allocs/op     vs base           │  allocs/op   vs base                │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │   allocs/op     vs base           │  allocs/op   vs base           │
Query/NonRecursiveCTEParallel                   0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-2                 0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-4                 0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-8                 0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1764.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/NonRecursiveCTEParallel-16                0.000 ± ∞ ¹   3763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   4765.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   3766.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1763.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                   ⁴     3.763k        ?                                +0.00%               ⁴     3.766k        ?                     4.764k        ?                     3.765k        ?                     1.764k        ?                   6.000        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

## Query - OrderBy (Sequential)

<!--BENCHMARK:slow/benchstat_query_orderby_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                         │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │      bench_eatonphil_direct.txt      │       bench_glebarez_driver.txt       │        bench_mattn_driver.txt         │       bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │      bench_zombiezen_direct.txt       │
                         │          sec/op          │    sec/op     vs base                │    sec/op     vs base                │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │
Query/OrderByParallel                  75.52m ± ∞ ¹   89.33m ± ∞ ¹       ~ (p=0.100 n=3) ²   46.87m ± ∞ ¹       ~ (p=0.100 n=3) ²   82.49m ± ∞ ¹        ~ (p=0.100 n=3) ²   87.01m ± ∞ ¹        ~ (p=0.100 n=3) ²   82.53m ± ∞ ¹        ~ (p=0.100 n=3) ²   62.02m ± ∞ ¹        ~ (p=0.100 n=3) ²   85.00m ± ∞ ¹        ~ (p=0.700 n=3) ²
Query/OrderByParallel-2                47.97m ± ∞ ¹   56.51m ± ∞ ¹       ~ (p=0.100 n=3) ²   38.00m ± ∞ ¹       ~ (p=0.100 n=3) ²   56.20m ± ∞ ¹        ~ (p=0.100 n=3) ²   59.59m ± ∞ ¹        ~ (p=0.100 n=3) ²   56.03m ± ∞ ¹        ~ (p=0.100 n=3) ²   43.09m ± ∞ ¹        ~ (p=0.100 n=3) ²   48.96m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-4                34.45m ± ∞ ¹   36.56m ± ∞ ¹       ~ (p=0.700 n=3) ²   32.55m ± ∞ ¹       ~ (p=0.700 n=3) ²   38.67m ± ∞ ¹        ~ (p=0.100 n=3) ²   43.74m ± ∞ ¹        ~ (p=0.100 n=3) ²   37.45m ± ∞ ¹        ~ (p=0.400 n=3) ²   28.49m ± ∞ ¹        ~ (p=0.100 n=3) ²   39.21m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-8                42.79m ± ∞ ¹   42.87m ± ∞ ¹       ~ (p=1.000 n=3) ²   58.23m ± ∞ ¹       ~ (p=0.100 n=3) ²   63.01m ± ∞ ¹        ~ (p=0.100 n=3) ²   74.91m ± ∞ ¹        ~ (p=0.100 n=3) ²   60.74m ± ∞ ¹        ~ (p=0.100 n=3) ²   41.14m ± ∞ ¹        ~ (p=0.100 n=3) ²   63.18m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-16               50.67m ± ∞ ¹   51.22m ± ∞ ¹       ~ (p=0.700 n=3) ²   64.30m ± ∞ ¹       ~ (p=0.100 n=3) ²   67.28m ± ∞ ¹        ~ (p=0.100 n=3) ²   97.72m ± ∞ ¹        ~ (p=0.100 n=3) ²   67.30m ± ∞ ¹        ~ (p=0.100 n=3) ²   50.16m ± ∞ ¹        ~ (p=1.000 n=3) ²   81.29m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                48.58m         52.67m        +8.41%                   46.49m        -4.31%                   59.73m        +22.94%                   69.83m        +43.74%                   58.88m        +21.21%                   43.58m        -10.30%                   60.91m        +25.37%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt             │      bench_eatonphil_direct.txt      │            bench_glebarez_driver.txt            │              bench_mattn_driver.txt              │            bench_modernc_driver.txt             │           bench_tailscale_driver.txt           │        bench_zombiezen_direct.txt        │
                         │           B/op           │       B/op         vs base                      │    B/op      vs base                 │       B/op         vs base                      │        B/op         vs base                      │       B/op         vs base                      │       B/op         vs base                     │      B/op       vs base                  │
Query/OrderByParallel                 369.000 ± ∞ ¹   6399278.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6397701.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   11999011.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   6397705.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   2798848.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   2939.000 ± ∞ ¹         ~ (p=0.200 n=3) ²
Query/OrderByParallel-2               372.000 ± ∞ ¹   6399268.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6397713.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   11999018.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   6397726.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   2798850.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1129.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-4                370.00 ± ∞ ¹    6399289.00 ± ∞ ¹             ~ (p=0.100 n=3) ²   38.00 ± ∞ ¹        ~ (p=0.100 n=3) ²    6397817.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    11999140.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    6397881.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    2798878.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    1149.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-8                 390.0 ± ∞ ¹     6399297.0 ± ∞ ¹             ~ (p=0.100 n=3) ²   499.0 ± ∞ ¹        ~ (p=0.700 n=3) ²     6398041.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     12000527.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     6398044.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     2799121.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     1476.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-16                806.0 ± ∞ ¹     6399552.0 ± ∞ ¹             ~ (p=0.100 n=3) ²   583.0 ± ∞ ¹        ~ (p=1.000 n=3) ²     6398822.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     12002484.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     6398741.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     2799457.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     1686.0 ± ∞ ¹         ~ (p=0.700 n=3) ²
geomean                                 437.2             6.103Mi        +1463764.34%                   50.60        -88.43%                       6.102Mi        +1463462.84%                        11.44Mi        +2744938.31%                       6.102Mi        +1463462.98%                       2.669Mi        +640185.31%                    1.532Ki        +258.75%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt            │    bench_eatonphil_direct.txt    │           bench_glebarez_driver.txt            │             bench_mattn_driver.txt             │            bench_modernc_driver.txt            │           bench_tailscale_driver.txt           │       bench_zombiezen_direct.txt       │
                         │        allocs/op         │    allocs/op      vs base                      │  allocs/op   vs base             │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │  allocs/op    vs base                  │
Query/OrderByParallel                   8.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349771.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   35.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-2                 8.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349771.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   16.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-4                 8.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449781.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349772.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449781.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   16.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-8                 8.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449784.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349776.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449784.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149767.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   18.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-16               10.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449788.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349794.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449788.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149769.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   20.000 ± ∞ ¹         ~ (p=0.700 n=3) ²
geomean                                 8.365             349.8k        +4181293.09%                                ?               ³ ⁴       449.8k        +5376783.92%                       349.8k        +4181273.96%                       449.8k        +5376783.92%                       149.8k        +1790273.17%                    20.03        +139.47%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean
⁴ ratios must be >0 to compute geomean
```
<!--END_BENCHMARK-->

## Query - OrderBy (Parallel)

<!--BENCHMARK:slow/benchstat_query_orderby_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                         │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt       │      bench_eatonphil_direct.txt      │       bench_glebarez_driver.txt       │        bench_mattn_driver.txt         │       bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │      bench_zombiezen_direct.txt       │
                         │          sec/op          │    sec/op     vs base                │    sec/op     vs base                │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │    sec/op     vs base                 │
Query/OrderByParallel                  75.52m ± ∞ ¹   89.33m ± ∞ ¹       ~ (p=0.100 n=3) ²   46.87m ± ∞ ¹       ~ (p=0.100 n=3) ²   82.49m ± ∞ ¹        ~ (p=0.100 n=3) ²   87.01m ± ∞ ¹        ~ (p=0.100 n=3) ²   82.53m ± ∞ ¹        ~ (p=0.100 n=3) ²   62.02m ± ∞ ¹        ~ (p=0.100 n=3) ²   85.00m ± ∞ ¹        ~ (p=0.700 n=3) ²
Query/OrderByParallel-2                47.97m ± ∞ ¹   56.51m ± ∞ ¹       ~ (p=0.100 n=3) ²   38.00m ± ∞ ¹       ~ (p=0.100 n=3) ²   56.20m ± ∞ ¹        ~ (p=0.100 n=3) ²   59.59m ± ∞ ¹        ~ (p=0.100 n=3) ²   56.03m ± ∞ ¹        ~ (p=0.100 n=3) ²   43.09m ± ∞ ¹        ~ (p=0.100 n=3) ²   48.96m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-4                34.45m ± ∞ ¹   36.56m ± ∞ ¹       ~ (p=0.700 n=3) ²   32.55m ± ∞ ¹       ~ (p=0.700 n=3) ²   38.67m ± ∞ ¹        ~ (p=0.100 n=3) ²   43.74m ± ∞ ¹        ~ (p=0.100 n=3) ²   37.45m ± ∞ ¹        ~ (p=0.400 n=3) ²   28.49m ± ∞ ¹        ~ (p=0.100 n=3) ²   39.21m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-8                42.79m ± ∞ ¹   42.87m ± ∞ ¹       ~ (p=1.000 n=3) ²   58.23m ± ∞ ¹       ~ (p=0.100 n=3) ²   63.01m ± ∞ ¹        ~ (p=0.100 n=3) ²   74.91m ± ∞ ¹        ~ (p=0.100 n=3) ²   60.74m ± ∞ ¹        ~ (p=0.100 n=3) ²   41.14m ± ∞ ¹        ~ (p=0.100 n=3) ²   63.18m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/OrderByParallel-16               50.67m ± ∞ ¹   51.22m ± ∞ ¹       ~ (p=0.700 n=3) ²   64.30m ± ∞ ¹       ~ (p=0.100 n=3) ²   67.28m ± ∞ ¹        ~ (p=0.100 n=3) ²   97.72m ± ∞ ¹        ~ (p=0.100 n=3) ²   67.30m ± ∞ ¹        ~ (p=0.100 n=3) ²   50.16m ± ∞ ¹        ~ (p=1.000 n=3) ²   81.29m ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                48.58m         52.67m        +8.41%                   46.49m        -4.31%                   59.73m        +22.94%                   69.83m        +43.74%                   58.88m        +21.21%                   43.58m        -10.30%                   60.91m        +25.37%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt             │      bench_eatonphil_direct.txt      │            bench_glebarez_driver.txt            │              bench_mattn_driver.txt              │            bench_modernc_driver.txt             │           bench_tailscale_driver.txt           │        bench_zombiezen_direct.txt        │
                         │           B/op           │       B/op         vs base                      │    B/op      vs base                 │       B/op         vs base                      │        B/op         vs base                      │       B/op         vs base                      │       B/op         vs base                     │      B/op       vs base                  │
Query/OrderByParallel                 369.000 ± ∞ ¹   6399278.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   5.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6397701.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   11999011.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   6397705.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   2798848.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   2939.000 ± ∞ ¹         ~ (p=0.200 n=3) ²
Query/OrderByParallel-2               372.000 ± ∞ ¹   6399268.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹        ~ (p=0.100 n=3) ²   6397713.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   11999018.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   6397726.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   2798850.000 ± ∞ ¹            ~ (p=0.100 n=3) ²   1129.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-4                370.00 ± ∞ ¹    6399289.00 ± ∞ ¹             ~ (p=0.100 n=3) ²   38.00 ± ∞ ¹        ~ (p=0.100 n=3) ²    6397817.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    11999140.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    6397881.00 ± ∞ ¹             ~ (p=0.100 n=3) ²    2798878.00 ± ∞ ¹            ~ (p=0.100 n=3) ²    1149.00 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-8                 390.0 ± ∞ ¹     6399297.0 ± ∞ ¹             ~ (p=0.100 n=3) ²   499.0 ± ∞ ¹        ~ (p=0.700 n=3) ²     6398041.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     12000527.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     6398044.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     2799121.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     1476.0 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-16                806.0 ± ∞ ¹     6399552.0 ± ∞ ¹             ~ (p=0.100 n=3) ²   583.0 ± ∞ ¹        ~ (p=1.000 n=3) ²     6398822.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     12002484.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     6398741.0 ± ∞ ¹             ~ (p=0.100 n=3) ²     2799457.0 ± ∞ ¹            ~ (p=0.100 n=3) ²     1686.0 ± ∞ ¹         ~ (p=0.700 n=3) ²
geomean                                 437.2             6.103Mi        +1463764.34%                   50.60        -88.43%                       6.102Mi        +1463462.84%                        11.44Mi        +2744938.31%                       6.102Mi        +1463462.98%                       2.669Mi        +640185.31%                    1.532Ki        +258.75%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                         │ bench_ncruces_direct.txt │            bench_ncruces_driver.txt            │    bench_eatonphil_direct.txt    │           bench_glebarez_driver.txt            │             bench_mattn_driver.txt             │            bench_modernc_driver.txt            │           bench_tailscale_driver.txt           │       bench_zombiezen_direct.txt       │
                         │        allocs/op         │    allocs/op      vs base                      │  allocs/op   vs base             │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │    allocs/op      vs base                      │  allocs/op    vs base                  │
Query/OrderByParallel                   8.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349771.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   35.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-2                 8.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349771.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449780.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   16.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-4                 8.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449781.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349772.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449781.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149766.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   16.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-8                 8.000 ± ∞ ¹   349778.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449784.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349776.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449784.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149767.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   18.000 ± ∞ ¹         ~ (p=0.100 n=3) ²
Query/OrderByParallel-16               10.000 ± ∞ ¹   349779.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=0.100 n=3) ²     449788.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   349794.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   449788.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   149769.000 ± ∞ ¹             ~ (p=0.100 n=3) ²   20.000 ± ∞ ¹         ~ (p=0.700 n=3) ²
geomean                                 8.365             349.8k        +4181293.09%                                ?               ³ ⁴       449.8k        +5376783.92%                       349.8k        +4181273.96%                       449.8k        +5376783.92%                       149.8k        +1790273.17%                    20.03        +139.47%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean
⁴ ratios must be >0 to compute geomean
```
<!--END_BENCHMARK-->

## Query - RecursiveCTE (Sequential)

<!--BENCHMARK:slow/benchstat_query_recursivecte_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                              │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                              │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │
Query/RecursiveCTEParallel                  6.134m ± ∞ ¹    6.069m ± ∞ ¹       ~ (p=0.100 n=3) ²   2.505m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.285m ± ∞ ¹        ~ (p=0.100 n=3) ²   2.400m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.304m ± ∞ ¹        ~ (p=0.100 n=3) ²   2.434m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.227m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                3.083m ± ∞ ¹    3.076m ± ∞ ¹       ~ (p=0.700 n=3) ²   1.266m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.671m ± ∞ ¹        ~ (p=0.100 n=3) ²   1.223m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.668m ± ∞ ¹        ~ (p=0.100 n=3) ²   1.231m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.645m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4               1577.8µ ± ∞ ¹   1584.4µ ± ∞ ¹       ~ (p=0.400 n=3) ²   646.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1374.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   620.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1371.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   626.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1358.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                853.6µ ± ∞ ¹    797.2µ ± ∞ ¹       ~ (p=0.200 n=3) ²   328.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    762.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   318.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    754.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   318.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    745.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16               717.7µ ± ∞ ¹    705.8µ ± ∞ ¹       ~ (p=0.700 n=3) ²   286.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    662.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   280.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    651.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   276.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    662.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                     1.788m          1.755m        -1.86%                   719.6µ        -59.75%                    1.578m        -11.74%                   695.0µ        -61.13%                    1.570m        -12.21%                   697.9µ        -60.97%                    1.561m        -12.69%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │    bench_eatonphil_direct.txt    │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                              │           B/op           │      B/op       vs base           │    B/op      vs base             │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │     B/op       vs base           │
Query/RecursiveCTEParallel                     0.0 ± ∞ ¹     2482.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ³       2261.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     6857.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     2258.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1505.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     362.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                   0.0 ± ∞ ¹     2481.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ²       2263.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     6826.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     2257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1503.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     367.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4                 2.000 ± ∞ ¹   2481.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.400 n=3) ²     2272.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6840.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2258.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1502.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   371.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                 3.000 ± ∞ ¹   2484.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.600 n=3) ²     2355.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6856.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2296.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1507.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   365.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16                4.000 ± ∞ ¹   2488.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=0.600 n=3) ²     2361.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6858.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2286.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1509.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   368.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                ⁴    2.425Ki        ?                                ?               ⁴ ⁵    2.248Ki        ?                    6.687Ki        ?                    2.218Ki        ?                    1.470Ki        ?                     366.6        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
⁵ ratios must be >0 to compute geomean

                              │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                              │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/RecursiveCTEParallel                   0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   143.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   48.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   143.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16                0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                ⁴     110.0        ?                                +0.00%               ⁴     112.4        ?                     142.4        ?                     112.0        ?                    48.80        ?                   6.000        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->

## Query - RecursiveCTE (Parallel)

<!--BENCHMARK:slow/benchstat_query_recursivecte_parallel.txt-->
```
benchstat bench_ncruces_direct.txt bench_ncruces_driver.txt bench_eatonphil_direct.txt bench_glebarez_driver.txt bench_mattn_driver.txt bench_modernc_driver.txt bench_tailscale_driver.txt bench_zombiezen_direct.txt
goos: darwin
goarch: arm64
pkg: github.com/michaellenaghan/go-sqlite-bench
cpu: Apple M2 Pro
                              │ bench_ncruces_direct.txt │       bench_ncruces_driver.txt        │      bench_eatonphil_direct.txt       │       bench_glebarez_driver.txt        │        bench_mattn_driver.txt         │        bench_modernc_driver.txt        │      bench_tailscale_driver.txt       │       bench_zombiezen_direct.txt       │
                              │          sec/op          │    sec/op      vs base                │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │    sec/op     vs base                 │    sec/op      vs base                 │
Query/RecursiveCTEParallel                  6.134m ± ∞ ¹    6.069m ± ∞ ¹       ~ (p=0.100 n=3) ²   2.505m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.285m ± ∞ ¹        ~ (p=0.100 n=3) ²   2.400m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.304m ± ∞ ¹        ~ (p=0.100 n=3) ²   2.434m ± ∞ ¹        ~ (p=0.100 n=3) ²    5.227m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                3.083m ± ∞ ¹    3.076m ± ∞ ¹       ~ (p=0.700 n=3) ²   1.266m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.671m ± ∞ ¹        ~ (p=0.100 n=3) ²   1.223m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.668m ± ∞ ¹        ~ (p=0.100 n=3) ²   1.231m ± ∞ ¹        ~ (p=0.100 n=3) ²    2.645m ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4               1577.8µ ± ∞ ¹   1584.4µ ± ∞ ¹       ~ (p=0.400 n=3) ²   646.7µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1374.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   620.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1371.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   626.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   1358.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                853.6µ ± ∞ ¹    797.2µ ± ∞ ¹       ~ (p=0.200 n=3) ²   328.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    762.4µ ± ∞ ¹        ~ (p=0.100 n=3) ²   318.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    754.8µ ± ∞ ¹        ~ (p=0.100 n=3) ²   318.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²    745.6µ ± ∞ ¹        ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16               717.7µ ± ∞ ¹    705.8µ ± ∞ ¹       ~ (p=0.700 n=3) ²   286.1µ ± ∞ ¹        ~ (p=0.100 n=3) ²    662.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   280.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²    651.0µ ± ∞ ¹        ~ (p=0.100 n=3) ²   276.9µ ± ∞ ¹        ~ (p=0.100 n=3) ²    662.3µ ± ∞ ¹        ~ (p=0.100 n=3) ²
geomean                                     1.788m          1.755m        -1.86%                   719.6µ        -59.75%                    1.578m        -11.74%                   695.0µ        -61.13%                    1.570m        -12.21%                   697.9µ        -60.97%                    1.561m        -12.69%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                              │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt      │    bench_eatonphil_direct.txt    │     bench_glebarez_driver.txt     │      bench_mattn_driver.txt       │     bench_modernc_driver.txt      │    bench_tailscale_driver.txt     │    bench_zombiezen_direct.txt    │
                              │           B/op           │      B/op       vs base           │    B/op      vs base             │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │      B/op       vs base           │     B/op       vs base           │
Query/RecursiveCTEParallel                     0.0 ± ∞ ¹     2482.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ³       2261.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     6857.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     2258.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1505.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     362.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                   0.0 ± ∞ ¹     2481.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     0.0 ± ∞ ¹  ~ (p=1.000 n=3) ²       2263.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     6826.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     2257.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     1503.0 ± ∞ ¹  ~ (p=0.100 n=3) ²     367.0 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4                 2.000 ± ∞ ¹   2481.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹  ~ (p=0.400 n=3) ²     2272.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6840.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2258.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1502.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   371.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                 3.000 ± ∞ ¹   2484.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1.000 ± ∞ ¹  ~ (p=0.600 n=3) ²     2355.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6856.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2296.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1507.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   365.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16                4.000 ± ∞ ¹   2488.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2.000 ± ∞ ¹  ~ (p=0.600 n=3) ²     2361.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6858.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   2286.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   1509.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   368.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                ⁴    2.425Ki        ?                                ?               ⁴ ⁵    2.248Ki        ?                    6.687Ki        ?                    2.218Ki        ?                    1.470Ki        ?                     366.6        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
⁵ ratios must be >0 to compute geomean

                              │ bench_ncruces_direct.txt │     bench_ncruces_driver.txt     │     bench_eatonphil_direct.txt      │    bench_glebarez_driver.txt     │      bench_mattn_driver.txt      │     bench_modernc_driver.txt     │   bench_tailscale_driver.txt    │   bench_zombiezen_direct.txt   │
                              │        allocs/op         │   allocs/op    vs base           │  allocs/op   vs base                │   allocs/op    vs base           │   allocs/op    vs base           │   allocs/op    vs base           │  allocs/op    vs base           │  allocs/op   vs base           │
Query/RecursiveCTEParallel                   0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   143.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-2                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-4                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   48.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-8                 0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   143.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
Query/RecursiveCTEParallel-16                0.000 ± ∞ ¹   110.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   0.000 ± ∞ ¹       ~ (p=1.000 n=3) ³   113.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   142.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   112.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   49.000 ± ∞ ¹  ~ (p=0.100 n=3) ²   6.000 ± ∞ ¹  ~ (p=0.100 n=3) ²
geomean                                                ⁴     110.0        ?                                +0.00%               ⁴     112.4        ?                     142.4        ?                     112.0        ?                    48.80        ?                   6.000        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean
```
<!--END_BENCHMARK-->
