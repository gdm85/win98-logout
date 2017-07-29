all: build

build:
	@mkdir -p bin/
	@rm -f *.gen.go
	go generate
	GOBIN="$(CURDIR)/bin" go install -ldflags="-s -w" .

.PHONY: all build
