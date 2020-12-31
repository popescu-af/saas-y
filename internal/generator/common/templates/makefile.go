package templates

// Makefile is the template for the service's Makefile.
const Makefile = `APP_VERSION=v0.0.0
TAG_REVISION=$(shell git rev-parse --short HEAD)
BUILD_TAG=${APP_VERSION}-${TAG_REVISION}
PLATFORM?=linux/amd64

DOCKER_REGISTRY_PORT=5000
DOCKER_REGISTRY=localhost:${DOCKER_REGISTRY_PORT}

ifeq ($(shell uname -s),Darwin)
	DOCKER_REGISTRY=docker.for.mac.localhost:${DOCKER_REGISTRY_PORT}
endif

.PHONY: build
build:
	docker buildx build --network host --platform ${PLATFORM} \
		--build-arg GITHUB_URL="${GITHUB_URL}" \
		-t tutorial-svc:${BUILD_TAG} .

.PHONY: tag
tag: build
	docker tag {{.Name}}:${BUILD_TAG} {{.Name}}:latest

.PHONY: run
run: tag
	docker run -p 8000:{{.Port}} -t {{.Name}}:latest

.PHONY: publish
publish: tag
	docker tag {{.Name}}:${BUILD_TAG} ${DOCKER_REGISTRY}/{{.Name}}:${BUILD_TAG}
	docker tag {{.Name}}:latest ${DOCKER_REGISTRY}/{{.Name}}:latest
	docker push ${DOCKER_REGISTRY}/{{.Name}}:${BUILD_TAG}
	docker push ${DOCKER_REGISTRY}/{{.Name}}:latest

.PHONY: deploy
deploy: publish
	kubectl apply -f deploy/{{.Name}}.yaml`
