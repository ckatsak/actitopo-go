GO ?= go
SHADOW ?= $(shell $(GO) env GOPATH)/bin/shadow

all: lint

build:
	$(GO) build ./...

lint:
	[ -f $(SHADOW) ]
	$(GO) vet -vettool=$(SHADOW) ./...

test:
	$(GO) test ./...

doc:
	@$(GO) doc -all . | $(PAGER)

.PHONY: all lint test doc

