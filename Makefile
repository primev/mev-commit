BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD | sed 's/[^a-zA-Z0-9._-]/-/g')
COMMIT ?= $(shell git rev-parse --short HEAD)
TAG ?= $(BRANCH)-$(COMMIT)
REGISTRY ?= ghcr.io/primev

.PHONY: docker
docker:
	cd infrastructure/docker && \
	TAG=$(TAG) \
	REGISTRY=$(REGISTRY) \
	GIT_BRANCH=$(BRANCH) \
	GIT_COMMIT=$(COMMIT) \
	docker buildx bake
