ARG BASE=golang:1.17-alpine
FROM ${BASE} AS builder


# set the working directory
WORKDIR /simple-driver

COPY . .

RUN go build -ldflags "-s -w" -o  ./cmd/build/driver-linux-amd64  ./cmd

FROM alpine:3.12

COPY --from=builder /simple-driver/cmd/build/driver-linux-amd64 /app/

RUN mkdir /logs

STOPSIGNAL SIGKILL

# Start
ENTRYPOINT [ "/app/driver-linux-amd64" ]

#docker buildx build --platform linux/amd64 -t registry.cn-shanghai.aliyuncs.com/winc-link/simple-driver:tcp-demo -f docker/Dockerfile . --push

