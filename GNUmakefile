TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME:=digitalocean
ACCTEST_TIMEOUT:=120m
ACCTEST_PARALLELISM:=2

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test -v ./$(PKG_NAME)/... $(TESTARGS) -timeout $(ACCTEST_TIMEOUT) -parallel=$(ACCTEST_PARALLELISM)

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test ./digitalocean/... -v -sweep=1

goimports:
	@echo "==> Fixing imports code with goimports..."
	@find . -name '*.go' | grep -v vendor | grep -v generator-resource-id | while read f; do goimports -w "$$f"; done

install-golangci-lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint: install-golangci-lint
	@golangci-lint run -v ./...

fmt:
	gofmt -s -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

website:
	@echo "Use this site to preview markdown rendering: https://registry.terraform.io/tools/doc-preview"

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile website sweep
