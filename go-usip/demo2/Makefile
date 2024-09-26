VERSION := 0.0.1

.PHONY: build image

build:
	@mkdir -p bin
	@go build -o bin/server .

image: build
	@docker build -t usip-server:${VERSION} .