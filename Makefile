#! /usr/bin/make
#
# Makefile for shogoa
#
# Targets:
# - "depend" retrieves the Go packages needed to run the linter and tests
# - "lint" runs the linter and checks the code format using goimports
# - "test" runs the tests
#
# Meta targets:
# - "all" is the default target, it runs all the targets in the order above.
#

all: depend lint shogoagen

.PHONY: depend
depend:
	go mod download
	go install github.com/onsi/ginkgo/ginkgo@latest

.PHONY: test
test:
	ginkgo -r --randomizeAllSpecs --failOnPending --randomizeSuites -race
	go test -v github.com/shogo82148/shogoa/_integration_tests

.PHONY: shogoagen
shogoagen:
	@cd shogoagen && \
	go install
