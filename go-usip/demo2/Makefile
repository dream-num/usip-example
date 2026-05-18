VERSION ?= latest

.PHONY: build image sync-sheet-host install-sheet-host build-sheet-host dev-sheet-host

build:
	@mkdir -p bin
	@GOARCH=amd64 go build -o bin/server-amd64 .
	@GOARCH=arm64 go build -o bin/server-arm64 .

install-sheet-host:
	@cd sheet-host-src && pnpm install

build-sheet-host:
	@cd sheet-host-src && pnpm build
	@echo "built sheet-host assets into web/public/sheet-host"

dev-sheet-host:
	@cd sheet-host-src && pnpm dev

image: build build-sheet-host
	@docker buildx --builder=images-builder build --platform linux/arm64,linux/amd64 -t \
		univer-acr-registry.cn-shenzhen.cr.aliyuncs.com/release/usip-server:${VERSION} \
		--provenance=false \
		--push .
