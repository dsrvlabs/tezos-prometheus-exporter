FROM golang:1.16.7-alpine3.14 as builder
MAINTAINER Joonkyo Kim <joonkyo.kim@dsrvlabs.com>

ENV GO111MODULE=on
ADD . $GOPATH/src/dsrvlabs/tezos-prometheus-exporter
WORKDIR $GOPATH/src/dsrvlabs/tezos-prometheus-exporter

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build && go install

FROM alpine:3.14.1
MAINTAINER Joonkyo Kim <joonkyo.kim@dsrvlabs.com>

COPY --from=builder /go/bin/tezos-prometheus-exporter /app/tezos-prometheus-exporter

CMD ["/app/tezos-prometheus-exporter"]
