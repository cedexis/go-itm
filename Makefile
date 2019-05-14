TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=go-itm
DOCKER_IMAGE_NAME=$(PKG_NAME)
DOCKER_CONTAINER_NAME=$(PKG_NAME)-container

.PHONY: test

fmt:
	gofmt -w $(GOFMT_FILES)

test:
	go test ./...

docker-build:
	@docker build -t $(DOCKER_IMAGE_NAME) .

docker-run:
	@docker run -it --rm --mount type=bind,readonly=1,src=$(PWD),dst=/go-itm $(DOCKER_IMAGE_NAME) /bin/bash
