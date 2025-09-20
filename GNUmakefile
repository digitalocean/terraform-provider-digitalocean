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

# Release automation targets

BUMP?=patch
LATEST_TAG=$(shell git describe --tags --abbrev=0)

bump_version:
	@echo "Bumping version: $(BUMP) from $(LATEST_TAG)"
	@BASE_VERSION=$$(echo $(LATEST_TAG) | sed 's/^v//'); \
	IFS=. read MAJOR MINOR PATCH <<<"$$BASE_VERSION"; \
	if [ "$(BUMP)" = "major" ] || [ "$(BUMP)" = "breaking" ]; then \
		NEW_VERSION="v$$((MAJOR+1)).0.0"; \
	elif [ "$(BUMP)" = "minor" ] || [ "$(BUMP)" = "feature" ]; then \
		NEW_VERSION="v$$MAJOR.$$((MINOR+1)).0"; \
	else \
		NEW_VERSION="v$$MAJOR.$$MINOR.$$((PATCH+1))"; \
	fi; \
	git tag -a "$$NEW_VERSION" -m "release $$NEW_VERSION"; \
	git push origin tag "$$NEW_VERSION"

COMMIT?=HEAD
ORIGIN?=origin
TAG_VERSION=$(shell git describe --tags --abbrev=0)
tag:
			 @echo "Tagging commit $(COMMIT) as $(TAG_VERSION)"
			 git tag -a "$(TAG_VERSION)" -m "release $(TAG_VERSION)" $(COMMIT)
			 git push $(ORIGIN) tag "$(TAG_VERSION)"
			 @echo "Generating changelog for tag $(TAG_VERSION)"
			 @if ! command -v github_changelog_generator &> /dev/null; then \
				 echo "github_changelog_generator not found. Installing..."; \
				 gem install github_changelog_generator; \
			 fi
			 @github_changelog_generator --user digitalocean --project terraform-provider-digitalocean --future-release $(TAG_VERSION) --output CHANGELOG.md
			 @echo "Creating or updating draft GitHub release for tag $(TAG_VERSION) with changelog"
			 @if gh release view $(TAG_VERSION) > /dev/null 2>&1; then \
				 gh release edit $(TAG_VERSION) --title "$(TAG_VERSION)" --notes-file CHANGELOG.md --draft; \
			 else \
				 gh release create $(TAG_VERSION) --title "$(TAG_VERSION)" --draft --notes-file CHANGELOG.md; \
			 fi

# Changelog generator (all commit logs since last tag)
changes:
	@git log $(shell git describe --tags --abbrev=0)..HEAD --pretty=format:'- %s (%an)'



