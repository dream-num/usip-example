# python usip example

use [`fastapi`](https://fastapi.tiangolo.com/) to apply web server.

## run

1. install python dependencies
```shell
pip install fastapi
pip install "uvicorn[standard]"
```

2. run web server
```shell
uvicorn main:app --reload --port 8080
```
