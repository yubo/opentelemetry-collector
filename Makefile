include ./Makefile.Common

RUN_CONFIG?=local/config.yaml
CMD?=
OTEL_VERSION=main

BUILD_INFO_IMPORT_PATH=github.com/yubo/opentelemetry-collector/internal/version
VERSION=$(shell git describe --always --match "v[0-9]*" HEAD)
BUILD_INFO=-ldflags "-X $(BUILD_INFO_IMPORT_PATH).Version=$(VERSION)"

COMP_REL_PATH=internal/components/components.go
MOD_NAME=github.com/yubo/opentelemetry-collector


.DEFAULT_GOAL := all

.PHONY: all
all: build

# Build the Collector executable.
.PHONY: build
build:
	VERSION=$(VERSION) ${GORELEASER} build --single-target --snapshot --rm-dist

#	GO111MODULE=on CGO_ENABLED=0 $(GOCMD) build -trimpath -o ./bin/otelcol_$(GOOS)_$(GOARCH)$(EXTENSION) \
		$(BUILD_INFO) -tags $(GO_BUILD_TAGS) ./cmd/otelcol

.PHONY: run
run:
	GO111MODULE=on $(GOCMD) run --race ./cmd/otelcol/... --config ${RUN_CONFIG} ${RUN_ARGS}

# https://goreleaser.com/quick-start/
.PHONY: release
release: goreleaser
	VERSION=$(VERSION) $(GORELEASER) release 

.PHONY: pkg
pkg: goreleaser
	VERSION=$(VERSION) ${GORELEASER} release --snapshot --rm-dist

.PHONY: goreleaser-verify
goreleaser-verify: goreleaser
	VERSION=$(VERSION) ${GORELEASER} release --snapshot --rm-dist


.PHONY: go
go:
	@{ \
		if ! command -v '$(GOCMD)' >/dev/null 2>/dev/null; then \
			echo >&2 '$(GOCMD) command not found. Please install golang. https://go.dev/doc/install'; \
			exit 1; \
		fi \
	}

.PHONY: goreleaser
goreleaser:
	@{ \
		if ! command -v '$(GORELEASER)' >/dev/null 2>/dev/null; then \
			echo >&2 '$(GORELEASER) command not found. Please install goreleaser. https://goreleaser.com/install/'; \
			exit 1; \
		fi \
	}

.PHONY: add-tag
add-tag:
	@[ "${TAG}" ] || ( echo ">> env var TAG is not set"; exit 1 )
	@echo "Adding tag ${TAG}"
	@git tag -a ${TAG} -m "Version ${TAG}"

.PHONY: push-tag
push-tag:
	@[ "${TAG}" ] || ( echo ">> env var TAG is not set"; exit 1 )
	@echo "Pushing tag ${TAG}"
	@git push git@github.com:yubo/opentelemetry-collector.git ${TAG}


.PHONY: clean
clean:
	rm -rf dist
