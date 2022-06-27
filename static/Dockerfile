FROM golang:1.18-buster as builder
WORKDIR /go/src/app
ADD . /go/src/app
RUN go build -ldflags "-linkmode external -extldflags -static" -o /go/bin/app
RUN apt update && apt install -y ca-certificates && cat /etc/ssl/certs/* > /tmp/ca-certificates.crt

ARG TIMEZONE="Europe/Chisinau"

FROM scratch

COPY --from=builder /tmp/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/$TIMEZONE /usr/share/zoneinfo/$TIMEZONE
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY server.crt /
COPY server.key /

USER 65534:65524

COPY --from=builder --chown=65534:65534 /go/bin/app /

ENV SSL_CERT_DIR=/etc/ssl/certs TZ=$TIMEZONE USER=nobody

ENTRYPOINT ["/app"]
