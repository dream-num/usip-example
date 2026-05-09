# golang usip example

Contain:
- Account System: register, login, logout, session
- File System: create, list, open

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
