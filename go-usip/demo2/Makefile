VERSION ?= latest

.PHONY: build image sync-sheet-host install-sheet-host build-sheet-host dev-sheet-host

build:
	@mkdir -p bin
	@GOARCH=amd64 go build -o bin/server-amd64 .
	@GOARCH=arm64 go build -o bin/server-arm64 .

sync-sheet-host:
	@test -n "$(SOURCE)" || (echo "Usage: make sync-sheet-host SOURCE=/abs/path/to/univer-pro-sheet-start-kit/dist" && exit 1)
	@mkdir -p web/public/sheet-host
	@rm -rf web/public/sheet-host/*
	@cp -R $(SOURCE)/* web/public/sheet-host/
	@sed -i 's#src="./main.js"#src="/sheet/main.js"#g' web/public/sheet-host/index.html
	@sed -i 's#href="./assets/#href="/sheet/assets/#g' web/public/sheet-host/index.html
	@echo "synced sheet-host assets from $(SOURCE)"

install-sheet-host:
	@cd sheet-host-src && pnpm install

build-sheet-host:
	@cd sheet-host-src && pnpm build
	@echo "built sheet-host assets into web/public/sheet-host"

dev-sheet-host:
	@cd sheet-host-src && pnpm dev

image: build
	@docker buildx --builder=images-builder build --platform linux/arm64,linux/amd64 -t \
		univer-acr-registry.cn-shenzhen.cr.aliyuncs.com/release/usip-server:${VERSION} \
		--push .
