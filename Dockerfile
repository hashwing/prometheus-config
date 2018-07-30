FROM golang:alpine AS build-stage

ENV GOROOT=/usr/local/go \
    GOPATH=/gopath \
    GOBIN=/gopath/bin \
    PROJPATH=/gopath/src/github.com/hashwing/prometheus-config

RUN apk add -U -q --no-progress build-base git glide

WORKDIR /gopath/src/github.com/hashwing/prometheus-config

RUN glide up
RUN go build -v



FROM alpine:latest

COPY --from=build-stage /gopath/src/github.com/hashwing/prometheus-config/prometheus-config /usr/local/bin/

ENV PATH=$PATH:/usr/local/bin

WORKDIR /usr/local/bin

CMD ["./prometheus-config"]
