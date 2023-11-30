# FROM golang:alpine3.15 AS development
FROM docker.io/golang:1.21-alpine3.18 AS development
ARG arch=x86_64

# ENV CGO_ENABLED=0
WORKDIR /go/src/app/
COPY . /go/src/app/

RUN apk update && apk add --no-cache \
    llvm \
    clang \
    llvm-static \
    llvm-dev \
    make \
    libbpf \
    libbpf-dev \
    musl-dev

RUN mkdir -p /build/ && \
    make all && \
    cp ./bin/* /build/

ENV PATH=$PATH:/build
# This entrypoint is just to keep the container running 
# for development and debuging purposes
ENTRYPOINT ["tail", "-f", "/dev/null"]

#----------------------------#

FROM docker.io/alpine:3.18.4 AS test

WORKDIR /app/
COPY --from=development /build .
RUN apk update && apk add iperf3 iproute2 curl

ENTRYPOINT ["./bittwister"]

#----------------------------#

FROM docker.io/alpine:3.18.4 AS production

WORKDIR /app/
COPY --from=development /build .
RUN apk update && apk add iproute2 curl

ENTRYPOINT ["./bittwister", "start", "-d", "eth0", "-p", "50"]
