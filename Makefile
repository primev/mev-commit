TAG ?= $(shell git describe --tags || git rev-parse --short HEAD)

.PHONY: docker
docker:
	cd infrastructure/docker && TAG=$(TAG) docker buildx bake
