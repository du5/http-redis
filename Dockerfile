FROM golang:alpine as build

WORKDIR /root/http-redis/

COPY . /root/http-redis/

RUN go get && go build .

FROM alpine

ENV HOST=127.0.0.1 \
    PORT=6379 \
    PASS= \
    NAME=0

WORKDIR /root/http-redis/

COPY --from=build /root/http-redis/http-redis /root/http-redis/
COPY --from=build /root/http-redis/config.toml /root/http-redis/

CMD sed -i "s|rdbhost|${HOST}|" config.toml && \
    sed -i "s|rdbport|${PORT}|" config.toml && \
    sed -i "s|rdbpass|${PASS}|" config.toml && \
    sed -i "s|rdbname|${NAME}|" config.toml && \
    ./http-redis
