TAGS := ncruces_direct ncruces_driver modernc_driver zombiezen_direct mattn_driver tailscale_driver

BENCH_COUNT := 1
BENCH_CPU := ''
BENCH_CPU_PARALLEL := 1,2,4,8,16
BENCH_OPTS := -benchmem -short
BENCH_PATTERN := .
BENCH_SKIP := ''
BENCH_TIME := 1s
BENCH_TIMEOUT := 15m

ifdef BENCH_SLOW
	BENCH_COUNT := 10
	BENCH_OPTS := -benchmem
	BENCH_TIMEOUT := 30m
endif

BENCH_MAX_READ_CONNECTIONS := 0
BENCH_MAX_WRITE_CONNECTIONS := 32
BENCH_POSTS := 200
BENCH_POST_PARAGRAPHS := 10
BENCH_COMMENTS := 10
BENCH_COMMENT_PARAGRAPHS := 1

ifdef BENCH_BIG
	BENCH_POSTS := 1000
	BENCH_POST_PARAGRAPHS := 50
	BENCH_COMMENTS := 50
	BENCH_COMMENT_PARAGRAPHS := 5
endif

TEST_COUNT := 1
TEST_CPU := ''
TEST_OPTS := -short -v
TEST_PATTERN := .
TEST_SKIP := ''

ifdef TEST_SLOW
	TEST_OPTS := -v
endif

.NOTPARALLEL:

.PHONY: help
help:
	@echo "Targets:"
	@echo ""
	@echo "  bench-all                                          - Run all benchmarks"
	@echo "  bench-by-category                                  - Run all benchmark categories"
	@echo "  bench-category-baseline                            - Run baseline benchmarks"
	@echo "  bench-category-baseline-parallel                   - Run baseline benchmarks in parallel"
	@echo "  bench-category-populate                            - Run populate benchmarks"
	@echo "  bench-category-readwrite                           - Run readwrite benchmarks"
	@echo "  bench-category-readwrite-parallel                  - Run readwrite benchmarks in parallel"
	@echo "  bench-category-query-correlated                    - Run correlated query benchmarks"
	@echo "  bench-category-query-correlated-parallel           - Run correlated query benchmarks in parallel"
	@echo "  bench-category-query-groupby                       - Run groupby query benchmarks"
	@echo "  bench-category-query-groupby-parallel              - Run groupby query benchmarks in parallel"
	@echo "  bench-category-query-json                          - Run json query benchmarks"
	@echo "  bench-category-query-json-parallel                 - Run json query benchmarks in parallel"
	@echo "  bench-category-query-nonrecursivecte               - Run nonrecursivecte query benchmarks"
	@echo "  bench-category-query-nonrecursivecte-parallel      - Run nonrecursivecte query benchmarks in parallel"
	@echo "  bench-category-query-orderby                       - Run orderby query benchmarks"
	@echo "  bench-category-query-orderby-parallel              - Run orderby query benchmarks in parallel"
	@echo "  bench-category-query-recursivecte                  - Run recursivecte query benchmarks"
	@echo "  bench-category-query-recursivecte-parallel         - Run recursivecte query benchmarks in parallel"
	@echo "  benchstat-all                                      - Compare all benchmarks"
	@echo "  benchstat-by-category                              - Run and compare all benchmark categories"
	@echo "  benchstat-category-baseline                        - Run and compare baseline benchmarks"
	@echo "  benchstat-category-baseline-parallel               - Run and compare baseline benchmarks in parallel"
	@echo "  benchstat-category-populate                        - Run and compare populate benchmarks"
	@echo "  benchstat-category-readwrite                       - Run and compare readwrite benchmarks"
	@echo "  benchstat-category-readwrite-parallel              - Run and compare readwrite benchmarks in parallel"
	@echo "  benchstat-category-query-correlated                - Run and compare correlated query benchmarks"
	@echo "  benchstat-category-query-correlated-parallel       - Run and compare correlated query benchmarks in parallel"
	@echo "  benchstat-category-query-groupby                   - Run and compare groupby query benchmarks"
	@echo "  benchstat-category-query-groupby-parallel          - Run and compare groupby query benchmarks in parallel"
	@echo "  benchstat-category-query-json                      - Run and compare json query benchmarks"
	@echo "  benchstat-category-query-json-parallel             - Run and compare json query benchmarks in parallel"
	@echo "  benchstat-category-query-nonrecursivecte           - Run and compare nonrecursivecte query benchmarks"
	@echo "  benchstat-category-query-nonrecursivecte-parallel  - Run and compare nonrecursivecte query benchmarks in parallel"
	@echo "  benchstat-category-query-orderby                   - Run and compare orderby query benchmarks"
	@echo "  benchstat-category-query-orderby-parallel          - Run and compare orderby query benchmarks in parallel"
	@echo "  benchstat-category-query-recursivecte              - Run and compare recursivecte query benchmarks"
	@echo "  benchstat-category-query-recursivecte-parallel     - Run and compare recursivecte query benchmarks in parallel"
	@echo "  clean                                              - Remove all benchmark, benchstat and test files"
	@echo "  test-all                                           - Run all tests"
	@echo ""
	@echo "Variables:"
	@echo ""
	@echo "  TAGS=\"$(TAGS)\""
	@echo ""
	@echo "  BENCH_COUNT=$(BENCH_COUNT)"
	@echo "  BENCH_CPU=$(BENCH_CPU)"
	@echo "  BENCH_CPU_PARALLEL=$(BENCH_CPU_PARALLEL)"
	@echo "  BENCH_OPTS=\"$(BENCH_OPTS)\""
	@echo "  BENCH_PATTERN=$(BENCH_PATTERN)"
	@echo "  BENCH_SKIP=$(BENCH_SKIP)"
	@echo "  BENCH_TIME=$(BENCH_TIME)"
	@echo "  BENCH_TIMEOUT=$(BENCH_TIMEOUT)"
	@echo ""
	@echo "  BENCH_SLOW=$(BENCH_SLOW)"
	@echo ""
	@echo "    BENCH_SLOW changes some of the values shown above. To see the BENCH_SLOW values,"
	@echo "    try \"make help BENCH_SLOW=1\""
	@echo ""
	@echo "  BENCH_MAX_READ_CONNECTIONS=$(BENCH_MAX_READ_CONNECTIONS)"
	@echo "  BENCH_MAX_WRITE_CONNECTIONS=$(BENCH_MAX_WRITE_CONNECTIONS)"
	@echo "  BENCH_POSTS=$(BENCH_POSTS)"
	@echo "  BENCH_POST_PARAGRAPHS=$(BENCH_POST_PARAGRAPHS)"
	@echo "  BENCH_COMMENTS=$(BENCH_COMMENTS)"
	@echo "  BENCH_COMMENT_PARAGRAPHS=$(BENCH_COMMENT_PARAGRAPHS)"
	@echo ""
	@echo "  BENCH_BIG=$(BENCH_BIG)"
	@echo ""
	@echo "    BENCH_BIG changes some of the values shown above. To see the BENCH_BIG values,"
	@echo "    try \"make help BENCH_BIG=1\""
	@echo ""
	@echo "  TEST_COUNT=$(TEST_COUNT)"
	@echo "  TEST_CPU=$(TEST_CPU)"
	@echo "  TEST_OPTS=\"$(TEST_OPTS)\""
	@echo "  TEST_PATTERN=$(TEST_PATTERN)"
	@echo "  TEST_SKIP=$(TEST_SKIP)"
	@echo ""
	@echo "Examples:"
	@echo ""
	@echo "  make bench-all"
	@echo "  make benchstat-all"
	@echo ""
	@echo "  make bench-all BENCH_COUNT=3"
	@echo "  make benchstat-all"
	@echo ""
	@echo "  make bench-all BENCH_SLOW=1"
	@echo "  make benchstat-all"
	@echo ""
	@echo "  make bench-all BENCH_BIG=1"
	@echo "  make benchstat-all"
	@echo ""
	@echo "  make bench-all BENCH_SLOW=1 BENCH_BIG=1"
	@echo "  make benchstat-all"
	@echo ""
	@echo "  make bench-all TAGS=\"ncruces_direct ncruces_driver\""
	@echo "  make benchstat-all TAGS=\"ncruces_direct ncruces_driver\""
	@echo ""
	@echo "  make benchstat-by-category"
	@echo ""
	@echo "  make benchstat-by-category BENCH_COUNT=3"
	@echo ""
	@echo "  make benchstat-by-category BENCH_SLOW=1"
	@echo ""
	@echo "  make benchstat-by-category BENCH_BIG=1"
	@echo ""
	@echo "  make benchstat-by-category BENCH_SLOW=1 BENCH_BIG=1"
	@echo ""
	@echo "  make benchstat-by-category TAGS=\"ncruces_direct ncruces_driver\""
	@echo ""
	@echo "  make test-all"
	@echo ""
	@echo "  make test-all TAGS=\"ncruces_direct ncruces_driver\""
	@echo ""
	@echo "  make clean"
	@echo ""

.PHONY: all
all: bench-all test-all

.PHONY: bench-all
bench-all: $(addprefix bench-,$(TAGS))

.PHONY: bench-by-category
bench-by-category: \
	bench-category-baseline \
	bench-category-baseline-parallel \
  bench-category-populate \
  bench-category-readwrite \
  bench-category-readwrite-parallel \
  bench-category-query-correlated \
  bench-category-query-correlated-parallel \
  bench-category-query-groupby \
  bench-category-query-groupby-parallel \
  bench-category-query-json \
  bench-category-query-json-parallel \
  bench-category-query-nonrecursivecte \
  bench-category-query-nonrecursivecte-parallel \
  bench-category-query-orderby \
  bench-category-query-orderby-parallel \
  bench-category-query-recursivecte \
  bench-category-query-recursivecte-parallel

.PHONY: bench-category-baseline
bench-category-baseline:
	$(MAKE) bench-all BENCH_PATTERN="Baseline"

.PHONY: bench-category-baseline-parallel
bench-category-baseline-parallel:
	$(MAKE) bench-all BENCH_PATTERN="Baseline/.*Parallel\$$" BENCH_CPU=$(BENCH_CPU_PARALLEL)

.PHONY: bench-category-populate
bench-category-populate:
	$(MAKE) bench-all BENCH_PATTERN="Populate"

.PHONY: bench-category-readwrite
bench-category-readwrite:
	$(MAKE) bench-all BENCH_PATTERN="ReadWrite"

.PHONY: bench-category-readwrite-parallel
bench-category-readwrite-parallel:
	$(MAKE) bench-all BENCH_PATTERN="ReadWrite/.*Parallel\$$" BENCH_CPU=$(BENCH_CPU_PARALLEL)

.PHONY: bench-category-query-correlated
bench-category-query-correlated:
	$(MAKE) bench-all BENCH_PATTERN="Query/Correlated"

.PHONY: bench-category-query-correlated-parallel
bench-category-query-correlated-parallel:
	$(MAKE) bench-all BENCH_PATTERN="Query/CorrelatedParallel\$$" BENCH_CPU=$(BENCH_CPU_PARALLEL)

.PHONY: bench-category-query-groupby
bench-category-query-groupby:
	$(MAKE) bench-all BENCH_PATTERN="Query/GroupBy"

.PHONY: bench-category-query-groupby-parallel
bench-category-query-groupby-parallel:
	$(MAKE) bench-all BENCH_PATTERN="Query/GroupByParallel\$$" BENCH_CPU=$(BENCH_CPU_PARALLEL)

.PHONY: bench-category-query-json
bench-category-query-json:
	$(MAKE) bench-all BENCH_PATTERN="Query/JSON"

.PHONY: bench-category-query-json-parallel
bench-category-query-json-parallel:
	$(MAKE) bench-all BENCH_PATTERN="Query/JSONParallel\$$" BENCH_CPU=$(BENCH_CPU_PARALLEL)

.PHONY: bench-category-query-nonrecursivecte
bench-category-query-nonrecursivecte:
	$(MAKE) bench-all BENCH_PATTERN="Query/NonRecursiveCTE"

.PHONY: bench-category-query-nonrecursivecte-parallel
bench-category-query-nonrecursivecte-parallel:
	$(MAKE) bench-all BENCH_PATTERN="Query/NonRecursiveCTEParallel\$$" BENCH_CPU=$(BENCH_CPU_PARALLEL)

.PHONY: bench-category-query-orderby
bench-category-query-orderby:
	$(MAKE) bench-all BENCH_PATTERN="Query/OrderBy"

.PHONY: bench-category-query-orderby-parallel
bench-category-query-orderby-parallel:
	$(MAKE) bench-all BENCH_PATTERN="Query/OrderByParallel\$$" BENCH_CPU=$(BENCH_CPU_PARALLEL)

.PHONY: bench-category-query-recursivecte
bench-category-query-recursivecte:
	$(MAKE) bench-all BENCH_PATTERN="Query/^RecursiveCTE"

.PHONY: bench-category-query-recursivecte-parallel
bench-category-query-recursivecte-parallel:
	$(MAKE) bench-all BENCH_PATTERN="Query/^RecursiveCTEParallel\$$" BENCH_CPU=$(BENCH_CPU_PARALLEL)

.PHONY: $(addprefix bench-,$(TAGS))
$(addprefix bench-,$(TAGS)):
	go test \
		-tags "$(subst bench-,,$@)" \
		-bench "$(BENCH_PATTERN)" \
		-benchtime $(BENCH_TIME) \
		-count $(BENCH_COUNT) \
		-cpu $(BENCH_CPU) \
		-run "''" \
		-skip "$(BENCH_SKIP)" \
		-timeout $(BENCH_TIMEOUT) \
		$(BENCH_OPTS) \
		-gsb-max-read-connections $(BENCH_MAX_READ_CONNECTIONS) \
		-gsb-max-write-connections $(BENCH_MAX_WRITE_CONNECTIONS) \
		-gsb-posts $(BENCH_POSTS) \
		-gsb-post-paragraphs $(BENCH_POST_PARAGRAPHS) \
		-gsb-comments $(BENCH_COMMENTS) \
		-gsb-comment-paragraphs $(BENCH_COMMENT_PARAGRAPHS) \
		| tee $(addprefix bench_,$(addsuffix .txt,$(subst bench-,,$@)))

.PHONY: benchstat
benchstat:
	@if ! command -v benchstat >/dev/null 2>&1; then \
		echo "benchstat not found. Install with: go install golang.org/x/perf/cmd/benchstat@latest"; \
		exit 1; \
	fi

.PHONY: benchstat-all
benchstat-all: benchstat
	benchstat $(addprefix bench_,$(addsuffix .txt,$(TAGS)))

.PHONY: benchstat-by-category
benchstat-by-category: \
  benchstat-category-baseline \
  benchstat-category-baseline-parallel \
  benchstat-category-populate \
  benchstat-category-readwrite \
  benchstat-category-readwrite-parallel \
  benchstat-category-query-correlated \
  benchstat-category-query-correlated-parallel \
  benchstat-category-query-groupby \
  benchstat-category-query-groupby-parallel \
  benchstat-category-query-json \
  benchstat-category-query-json-parallel \
  benchstat-category-query-nonrecursivecte \
  benchstat-category-query-nonrecursivecte-parallel \
  benchstat-category-query-orderby \
  benchstat-category-query-orderby-parallel \
  benchstat-category-query-recursivecte \
  benchstat-category-query-recursivecte-parallel

.PHONY: benchstat-category-baseline
benchstat-category-baseline: benchstat bench-category-baseline
	$(MAKE) benchstat-all | tee benchstat_baseline.txt

.PHONY: benchstat-category-baseline-parallel
benchstat-category-baseline-parallel: benchstat bench-category-baseline-parallel
	$(MAKE) benchstat-all | tee benchstat_baseline_parallel.txt

.PHONY: benchstat-category-populate
benchstat-category-populate: benchstat bench-category-populate
	$(MAKE) benchstat-all | tee benchstat_populate.txt

.PHONY: benchstat-category-readwrite
benchstat-category-readwrite: benchstat bench-category-readwrite
	$(MAKE) benchstat-all | tee benchstat_readwrite.txt

.PHONY: benchstat-category-readwrite-parallel
benchstat-category-readwrite-parallel: benchstat bench-category-readwrite-parallel
	$(MAKE) benchstat-all | tee benchstat_readwrite_parallel.txt

.PHONY: benchstat-category-query-correlated
benchstat-category-query-correlated: benchstat bench-category-query-correlated
	$(MAKE) benchstat-all | tee benchstat_query_correlated.txt

.PHONY: benchstat-category-query-correlated-parallel
benchstat-category-query-correlated-parallel: benchstat bench-category-query-correlated-parallel
	$(MAKE) benchstat-all | tee benchstat_query_correlated_parallel.txt

.PHONY: benchstat-category-query-groupby
benchstat-category-query-groupby: benchstat bench-category-query-groupby
	$(MAKE) benchstat-all | tee benchstat_query_groupby.txt

.PHONY: benchstat-category-query-groupby-parallel
benchstat-category-query-groupby-parallel: benchstat bench-category-query-groupby-parallel
	$(MAKE) benchstat-all | tee benchstat_query_groupby_parallel.txt

.PHONY: benchstat-category-query-json
benchstat-category-query-json: benchstat bench-category-query-json
	$(MAKE) benchstat-all | tee benchstat_query_json.txt

.PHONY: benchstat-category-query-json-parallel
benchstat-category-query-json-parallel: benchstat bench-category-query-json-parallel
	$(MAKE) benchstat-all | tee benchstat_query_json_parallel.txt

.PHONY: benchstat-category-query-nonrecursivecte
benchstat-category-query-nonrecursivecte: benchstat bench-category-query-nonrecursivecte
	$(MAKE) benchstat-all | tee benchstat_query_nonrecursivecte.txt

.PHONY: benchstat-category-query-nonrecursivecte-parallel
benchstat-category-query-nonrecursivecte-parallel: benchstat bench-category-query-nonrecursivecte-parallel
	$(MAKE) benchstat-all | tee benchstat_query_nonrecursivecte_parallel.txt

.PHONY: benchstat-category-query-orderby
benchstat-category-query-orderby: benchstat bench-category-query-orderby
	$(MAKE) benchstat-all | tee benchstat_query_orderby.txt

.PHONY: benchstat-category-query-orderby-parallel
benchstat-category-query-orderby-parallel: benchstat bench-category-query-orderby-parallel
	$(MAKE) benchstat-all | tee benchstat_query_orderby_parallel.txt

.PHONY: benchstat-category-query-recursivecte
benchstat-category-query-recursivecte: benchstat bench-category-query-recursivecte
	$(MAKE) benchstat-all | tee benchstat_query_recursivecte.txt

.PHONY: benchstat-category-query-recursivecte-parallel
benchstat-category-query-recursivecte-parallel: benchstat bench-category-query-recursivecte-parallel
	$(MAKE) benchstat-all | tee benchstat_query_recursivecte_parallel.txt

.PHONY: test-all
test-all: $(addprefix test-,$(TAGS))

.PHONY: $(addprefix test-,$(TAGS))
$(addprefix test-,$(TAGS)):
	go test \
		-tags "$(subst test-,,$@)" \
		-count $(TEST_COUNT) \
		-cpu $(TEST_CPU) \
		-run "$(TEST_PATTERN)" \
		-skip "$(TEST_SKIP)" \
		$(TEST_OPTS) \
		| tee $(addprefix test_,$(addsuffix .txt,$(subst test-,,$@)))

.PHONY: clean
clean:
	rm -f $(addprefix bench_,$(addsuffix .txt,$(TAGS))) $(addprefix test_,$(addsuffix .txt,$(TAGS))) benchstat_*.txt
