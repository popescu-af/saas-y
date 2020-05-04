package templates

// Makefile is the template for the service's Makefile.
const Makefile = `TAG_REVISION=$(shell git rev-parse HEAD)

.PHONY: build
build:
	docker build -t {{.Name}}:${TAG_REVISION} .

.PHONY: test
test: 
	go test ./...

.PHONY: clean
clean: 
	go clean ./...

.PHONY: tag
tag: build
	docker tag {{.Name}}:${TAG_REVISION} {{.Name}}:latest

.PHONY: run
run: tag
	docker run -p 8000:{{.Port}} -t {{.Name}}:latest

.PHONY: publish
publish: tag
	docker tag {{.Name}}:${TAG_REVISION} localhost:5000/{{.Name}}:${TAG_REVISION}
	docker tag {{.Name}}:latest localhost:5000/{{.Name}}:latest
	docker push localhost:5000/{{.Name}}:${TAG_REVISION}
	docker push localhost:5000/{{.Name}}:latest`
