# golang usip example

Contain:
- Account System: register, login, logout, session
- File System: create, list, open

Dependencies:

- univer server and demo, see: https://www.univer.ai/zh-CN/guides/sheet/server/docker , use `bash -c "$(curl -fsSL https://get.univer.ai)"` to quick intall.
- [iris](https://www.iris-go.com/): apply web server
- postgresql: use to save user, file data
- redis: use to session

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

2. config the postgresql and redis in `configs/config.yaml`

   
3. run web server
```shell
go run .
```

## Display
<image src="./doc/image.png" />