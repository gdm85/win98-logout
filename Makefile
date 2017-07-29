all: build

build:
	@rm -f *.gen.go
	go generate
	go install .

.PHONY: all build
