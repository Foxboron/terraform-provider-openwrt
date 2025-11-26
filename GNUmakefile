SHELL := /bin/bash
OUT := $(shell pwd)/_out

EXCLUDED_PACKAGES := \
	github.com/foxboron/terraform-provider-openwrt \
	github.com/foxboron/terraform-provider-openwrt/internal/api \
	github.com/foxboron/terraform-provider-openwrt/internal/types \
	github.com/foxboron/terraform-provider-openwrt/mocks

PACKAGES := $(shell go list ./... | grep -Fvx -f <(printf '%s\n' $(EXCLUDED_PACKAGES)))

.PHONY: clean fmt lint install generate build test snapshot release
default: clean fmt install generate lint

clean:
	rm -Rfv $(OUT) dist mocks
	mkdir -p $(OUT)

fmt:
	gofmt -s -w -e .

lint:
	golangci-lint run

install:
	go mod download

generate:
	go generate -v ./...

tf-generate: generate
	cd tools; go generate -v ./...

build: tf-generate
	go build -v

test:
	go test -tags=test -race -parallel=10 -timeout 120s -cover -coverprofile=_out/.coverage -v $(PACKAGES);
	go tool cover -html=_out/.coverage -o=./_out/coverage.html

snapshot: tf-generate
	goreleaser build --clean --snapshot

release:
	goreleaser release --clean
