From ubuntu:22.04

WORKDIR /app

RUN mkdir /app/web

COPY configs /app/configs
COPY bin/server /app/server
COPY web/public /app/web/public
COPY web/views /app/web/views

RUN chmod +x /app/server

CMD ["/app/server"]