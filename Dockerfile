FROM golang:latest as builder
COPY . /srv/watson
RUN cd /srv/watson &&\
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags \'-s -w\' -o bin/watson &&\
    ls -l bin 

#
FROM alpine:latest
LABEL source.url="https://github.com/fimreal/watson"

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache tzdata ca-certificates &&\
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime &&\
    echo "Asia/Shanghai" > /etc/timezone

COPY --from=builder /srv/watson/bin/watson /watson
ENTRYPOINT [ "/watson" ]
