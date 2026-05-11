# golang usip example

Contain:
- Account System: register, login, logout, session
- File System: create, list, open
- Embedded sheet host (`/sheet`) served by this project
- Sheet-only mode: no `docHost` or doc/docx open flow

Dependencies:

- univer server and demo, see: https://www.univer.ai/zh-CN/guides/sheet/server/docker , use `bash -c "$(curl -fsSL https://get.univer.ai)"` to quick intall.
- [iris](https://www.iris-go.com/): apply web server
- postgresql: use to save user, file data
- redis: optional, used for session storage when enabled

## Prepare
Install univer server and demo
```shell
bash -c "$(curl -fsSL https://get.univer.ai)"
```

## Embedded sheet host assets

This project now serves sheet host at `/sheet` from `web/public/sheet-host`.

The editable source workspace is in `sheet-host-src/`.

Install and build from in-repo source:

```shell
make install-sheet-host
make build-sheet-host
```

For frontend dev server:

```shell
make dev-sheet-host
```

If you still want to refresh assets from external `univer-pro-sheet-start-kit` dist, run:

```shell
make sync-sheet-host SOURCE=/absolute/path/to/univer-pro-sheet-start-kit/dist
```

The sync command also rewrites `index.html` asset links to `/sheet/...` paths to avoid `main.js` 404 when opening `/sheet?unit=...`.


## Run

1. init the enviroment
```shell
go mod tidy
```

2. configure database/server/redis in `configs/config.yaml`

   Key options:
   - `server.port`: default listen port (default `8090`)
   - `redis.enabled`: enable redis-backed session store (default `true`)
   - `redis.addr`: required when `redis.enabled=true`
   - `univer.sheetHost`: defaults to `/sheet` (embedded route in this project)
   - `universer.host`: backend target for `/universer-api` proxy (default `http://localhost:8000`)

   Breaking behavior:
   - `docHost` is removed from demo2 configuration.
   - New/Open/Import flow in demo2 supports sheet only (`.xlsx` import).

   Port priority at runtime:
   - `PORT` environment variable
   - `server.port` in config file
   - default `8090`

   Redis enable priority at runtime:
   - `REDIS_ENABLED` environment variable
   - `redis.enabled` in config file

   Examples:
   ```shell
   PORT=9000 go run .
   REDIS_ENABLED=false go run .
   ```

   
3. run web server
```shell
go run .
```

## Display
<image src="./doc/image.png" />
