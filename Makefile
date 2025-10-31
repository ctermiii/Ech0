# Variables
GOHOSTOS:=$(shell go env GOHOSTOS)
GOHOSTARCH:=$(shell go env GOHOSTARCH)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
BUILD_TIME=$(shell date +%Y-%m-%dT%H:%M:%S)
GIT_COMMIT=$(shell git rev-parse HEAD)

# Docker variables
DOCKER_REGISTRY?=sn0wl1n
IMAGE_NAME?=ech0
IMAGE_TAG?=latest
OS?=$(if $(GOHOSTOS),$(GOHOSTOS),linux)
ARCH?=$(if $(GOHOSTARCH),$(GOHOSTARCH),amd64)

.PHONY: wire
wire:
	cd internal/di && wire

.PHONY: build-image
build-image:
	@echo "Building image for platform: $(OS)/$(ARCH)"
	docker build --platform $(OS)/$(ARCH) \
		--build-arg TARGETOS=$(OS) \
		--build-arg TARGETARCH=$(ARCH) \
		-t $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG) -f build.Dockerfile .

.PHONY: push-image
push-image:
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

