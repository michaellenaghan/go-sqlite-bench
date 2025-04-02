This project, originally inspired by [Benchmarking SQLite Performance in Go](https://www.golang.dk/articles/benchmarking-sqlite-performance-in-go), bechmarks various SQLite implementations:

* [github.com/glebarez/go-sqlite](https://github.com/glebarez/go-sqlite) (aka "glebarez")
* [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) (aka "mattn")
* [github.com/ncruces/go-sqlite3](https://github.com/ncruces/go-sqlite3) (aka "ncruces")
* [github.com/tailscale/sqlite](https://github.com/tailscale/sqlite) (aka "tailscale")
*	[github.com/zombiezen/go-sqlite](https://github.com/zombiezen/go-sqlite) (aka "zombiezen")
*	[gitlab.com/cznic/sqlite](https://gitlab.com/cznic/sqlite) (aka "modernc")

Here are some quick descriptions:

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

It takes ~8m to run the "quick" benchmarks on my laptop.

The benchmark results in the [slow](./slow/) directory were generated using:

```sh
make benchstat-by-category BENCH_BIG=1 BENCH_SLOW=1
```

It takes ~2h to run the "slow" benchmarks on my laptop.

The tests in the [tests](./tests/) directory were generated using:

```sh
make test-all
```

Among other things, the tests capture the compile-time options, pragmas, and SQLite version used by each implementation.

There are lots of ways to play the benchmark game; take a look at the examples below.

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

Variables:

  TAGS="ncruces_direct ncruces_driver modernc_driver zombiezen_direct mattn_driver tailscale_driver"

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

One thing to note: the first implementation in `TAGS` becomes the baseline for `benchstat` comparisons.

For example:

```sh
make benchstat-by-category TAGS="tailscale_driver mattn_driver modernc_driver ncruces_driver"
```

compares `mattn_driver`, `modernc_driver` and `ncruces_driver` against `tailscale_driver`.
