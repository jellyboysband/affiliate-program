ARG GO_VERSION=1.13

FROM golang:${GO_VERSION}-alpine AS builder

ENV CGO_ENABLED=0 GOFLAGS=-mod=vendor

WORKDIR /src

ADD ./ ./

RUN go build ./cmd/crawler

FROM scratch AS final
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/crawler /crawler

ENTRYPOINT ["/crawler"]