TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME?=digitalocean
ACCTEST_TIMEOUT?=120m
ACCTEST_PARALLELISM?=2

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
	go test ./digitalocean/sweep/... -v -sweep=1

goimports:
	@echo "==> Fixing imports code with goimports..."
	@find . -name '*.go' | grep -v vendor | grep -v generator-resource-id | while read f; do goimports -w "$$f"; done

install-golangci-lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

lint: install-golangci-lint
	@golangci-lint run -v ./...

fmt:
	gofmt -s -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

install-terrafmt:
	@go install github.com/katbyte/terrafmt@latest

terrafmt: install-terrafmt # Formats Terraform configuration blocks in tests.
	@terrafmt fmt --fmtcompat digitalocean/
	@terrafmt fmt --fmtcompat docs/

terrafmt-check: install-terrafmt # Returns non-0 exit code if terrafmt would make a change.
	@terrafmt diff --check --fmtcompat digitalocean/
	@terrafmt diff --check --fmtcompat docs/

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

website:
	@echo "Use this site to preview markdown rendering: https://registry.terraform.io/tools/doc-preview"

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile website sweep

.PHONY: _upgrade_godo
_upgrade_godo:
	go get -u github.com/digitalocean/godo

.PHONY: upgrade_godo
upgrade_godo: _upgrade_godo vendor
	@echo "==> upgrade the godo version"
	@echo ""

.PHONY: vendor
vendor:
	@echo "==> vendor dependencies"
	@echo ""
	go mod vendor
	go mod tidy



changes:
	@if ! command -v github-changelog-generator &> /dev/null; then \
		echo "github-changelog-generator not found. Installing..."; \
		go install github.com/digitalocean/github-changelog-generator@latest; \
	fi
	@github-changelog-generator -org digitalocean -repo terraform-provider-digitalocean

tag:
	@echo "==> BUMP=$(BUMP) tag"
	@echo ""
	bash scripts/release.sh $(BUMP)

