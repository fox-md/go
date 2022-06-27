FROM golang:1.18-buster as builder

ENV CGO_ENABLED 1
ENV GOOS linux

ARG Version=1.0.0

WORKDIR /go/src/app

ADD . /go/src/app

RUN go build -ldflags "-s -w -linkmode external -extldflags -static \
             -X main.GitHash=`git rev-parse --short HEAD || echo N/A` \
             -X main.BuildDate=`date +%FT%T%z` \
             -X main.Version=$Version" -o /go/bin/app
RUN apt update && apt install -y ca-certificates && cat /etc/ssl/certs/* > /tmp/ca-certificates.crt


FROM scratch

COPY --from=builder /tmp/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY server.crt /
COPY server.key /

USER 65534:65524

COPY --from=builder --chown=65534:65534 /go/bin/app /

ENV SSL_CERT_DIR=/etc/ssl/certs USER=nobody

EXPOSE 8443

ENTRYPOINT ["/app"]
