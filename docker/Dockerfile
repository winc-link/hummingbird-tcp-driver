ARG BASE=registry.cn-shanghai.aliyuncs.com/winc-edgex/edgex:go1.23.1-alpine3.16
FROM ${BASE} AS builder

# set the working directory
WORKDIR /simple-driver

COPY . .

RUN go build -ldflags "-s -w" -o  ./cmd/build/driver-linux-amd64  ./cmd

FROM registry.cn-shanghai.aliyuncs.com/winc-edgex/alpine:3.13

COPY --from=builder /simple-driver/cmd/build/driver-linux-amd64 /app/

EXPOSE 7653

RUN mkdir /logs

# Start
CMD [ "/bin/sh", "-c", "/app/driver-linux-amd64" ]


