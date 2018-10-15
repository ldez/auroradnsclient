# These env vars have to be set in the CI
# GITHUB_TOKEN

.PHONY: test release check ci-release version help

VERSION := $(shell cat VERSION)
SHA := $(shell git rev-parse --short HEAD)

default: check test

help:
	@echo "make check - run golangci-lint"
	@echo "make test - run tests"
	@echo "make release - tag with version and trigger CI release build"
	@echo "make version - show app version"

test:
	GO111MODULE=on go test -v ./...

check:
	golangci-lint run

release:
	git tag --force -s `cat VERSION` -m `cat VERSION`
	git push --force-with-lease origin master --tags

version:
	@echo $(VERSION) $(SHA)
