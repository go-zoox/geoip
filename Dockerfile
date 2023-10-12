# Builder
FROM whatwewant/builder-go:v1.20-1 as builder

WORKDIR /build

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
  go build \
  -trimpath \
  -ldflags '-w -s -buildid=' \
  -v -o geoip ./cmd/geoip

# Server
FROM whatwewant/alpine:v3.17-1

LABEL MAINTAINER="Zero<tobewhatwewant@gmail.com>"

LABEL org.opencontainers.image.source="https://github.com/go-zoox/geoip"

ARG VERSION=latest

ENV GEOIP_VERSION=${VERSION}

COPY --from=builder /build/geoip /bin

RUN geoip --version

CMD geoip server
