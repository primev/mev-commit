TAG ?= $(shell git describe --tags || git rev-parse --short HEAD)
REGISTRY ?= ghcr.io/primev

.PHONY: docker
docker:
	cd infrastructure/docker && TAG=$(TAG) REGISTRY=$(REGISTRY) docker buildx bake
