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
	BENCH_COUNT := 6
	BENCH_OPTS := -benchmem
	BENCH_TIME := 2s
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
TEST_CPU := 4
TEST_OPTS := -short -v
TEST_PATTERN := .
TEST_SKIP := ''

ifdef TEST_SLOW
	TEST_OPTS := -v
endif

.NOTPARALLEL:

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  bench-all                                - Run all benchmarks"
	@echo "  bench-by-category                        - Run all benchmark categories"
	@echo "  bench-category-baseline                  - Run baseline benchmarks"
	@echo "  bench-category-baseline-parallel         - Run baseline benchmarks in parallel"
	@echo "  bench-category-populate                  - Run populate benchmarks"
	@echo "  bench-category-readwrite                 - Run readwrite benchmarks"
	@echo "  bench-category-readwrite-parallel        - Run readwrite benchmarks in parallel"
	@echo "  bench-category-query-correlated          - Run correlated query benchmarks"
	@echo "  bench-category-query-groupby             - Run groupby query benchmarks"
	@echo "  bench-category-query-json                - Run json query benchmarks"
	@echo "  bench-category-query-nonrecursivecte     - Run nonrecursivecte query benchmarks"
	@echo "  bench-category-query-orderby             - Run orderby query benchmarks"
	@echo "  bench-category-query-recursivecte        - Run recursivecte query benchmarks"
	@echo "  benchstat-all                            - Compare all benchmarks"
	@echo "  benchstat-by-category                    - Run and compare all benchmark categories"
	@echo "  benchstat-category-baseline              - Run and compare baseline benchmarks"
	@echo "  benchstat-category-baseline-parallel     - Run and compare baseline benchmarks in parallel"
	@echo "  benchstat-category-populate              - Run and compare populate benchmarks"
	@echo "  benchstat-category-readwrite             - Run and compare readwrite benchmarks"
	@echo "  benchstat-category-readwrite-parallel    - Run and compare readwrite benchmarks in parallel"
	@echo "  benchstat-category-query-correlated      - Run and compare correlated query benchmarks"
	@echo "  benchstat-category-query-groupby         - Run and compare groupby query benchmarks"
	@echo "  benchstat-category-query-json            - Run and compare json query benchmarks"
	@echo "  benchstat-category-query-nonrecursivecte - Run and compare nonrecursivecte query benchmarks"
	@echo "  benchstat-category-query-orderby         - Run and compare orderby query benchmarks"
	@echo "  benchstat-category-query-recursivecte    - Run and compare recursivecte query benchmarks"
	@echo "  clean                                    - Remove all benchmark, benchstat and test files"
	@echo "  test-all                                 - Run all tests"

.PHONY: all
all: bench-all test-all

.PHONY: bench-all
bench-all: $(addprefix bench-,$(TAGS))

.PHONY: bench-by-category
bench-by-category: bench-category-baseline bench-category-baseline-parallel bench-category-populate bench-category-readwrite bench-category-readwrite-parallel bench-category-query-correlated bench-category-query-groupby bench-category-query-json bench-category-query-nonrecursivecte bench-category-query-orderby bench-category-query-recursivecte

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

.PHONY: bench-category-query-groupby
bench-category-query-groupby:
	$(MAKE) bench-all BENCH_PATTERN="Query/GroupBy"

.PHONY: bench-category-query-json
bench-category-query-json:
	$(MAKE) bench-all BENCH_PATTERN="Query/JSON"

.PHONY: bench-category-query-nonrecursivecte
bench-category-query-nonrecursivecte:
	$(MAKE) bench-all BENCH_PATTERN="Query/NonRecursiveCTE"

.PHONY: bench-category-query-orderby
bench-category-query-orderby:
	$(MAKE) bench-all BENCH_PATTERN="Query/OrderBy"

.PHONY: bench-category-query-recursivecte
bench-category-query-recursivecte:
	$(MAKE) bench-all BENCH_PATTERN="Query/^RecursiveCTE"

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
		echo "benchstat tool not found. Install with: go install golang.org/x/perf/cmd/benchstat@latest"; \
		exit 1; \
	fi

.PHONY: benchstat-all
benchstat-all: benchstat
	benchstat $(addprefix bench_,$(addsuffix .txt,$(TAGS)))

.PHONY: benchstat-by-category
benchstat-by-category: benchstat-category-baseline benchstat-category-baseline-parallel benchstat-category-populate benchstat-category-readwrite benchstat-category-readwrite-parallel benchstat-category-query-correlated benchstat-category-query-groupby benchstat-category-query-json benchstat-category-query-nonrecursivecte benchstat-category-query-orderby benchstat-category-query-recursivecte

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

.PHONY: benchstat-category-query-groupby
benchstat-category-query-groupby: benchstat bench-category-query-groupby
	$(MAKE) benchstat-all | tee benchstat_query_groupby.txt

.PHONY: benchstat-category-query-json
benchstat-category-query-json: benchstat bench-category-query-json
	$(MAKE) benchstat-all | tee benchstat_query_json.txt

.PHONY: benchstat-category-query-nonrecursivecte
benchstat-category-query-nonrecursivecte: benchstat bench-category-query-nonrecursivecte
	$(MAKE) benchstat-all | tee benchstat_query_nonrecursivecte.txt

.PHONY: benchstat-category-query-orderby
benchstat-category-query-orderby: benchstat bench-category-query-orderby
	$(MAKE) benchstat-all | tee benchstat_query_orderby.txt

.PHONY: benchstat-category-query-recursivecte
benchstat-category-query-recursivecte: benchstat bench-category-query-recursivecte
	$(MAKE) benchstat-all | tee benchstat_query_recursivecte.txt

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
