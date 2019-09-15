ARG GO_VERSION=1.13

FROM golang:${GO_VERSION}-alpine AS builder

ENV CGO_ENABLED=0 GOFLAGS=-mod=vendor

WORKDIR /src

ADD ./ ./

RUN go build ./cmd/store

FROM scratch AS final
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/store /store

ENTRYPOINT ["/store"]