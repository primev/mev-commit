LDFLAGS ?= -s -w \
-X github.com/primev/mev-commit/p2p.version=$(shell git describe --tags || git rev-parse --short HEAD)

.PHONY: build
build: export CGO_ENABLED=0
build: bin
	go build -ldflags '$(LDFLAGS)' -o bin/mev-commit ./cmd

bin:
	mkdir $@

bufgen:
	cd rpc && buf mod clear-cache && buf generate -o ../
	cd messages && buf mod clear-cache && buf generate -o ../
