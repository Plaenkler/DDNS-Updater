# syntax=docker/dockerfile:1

## Build
FROM golang:1.24-bookworm AS build

WORKDIR /app
COPY . /app

RUN go mod tidy && go build -o /ddns-updater cmd/def/main.go

## Deploy
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=build /ddns-updater /app/ddns-updater

RUN chmod +x /app/ddns-updater

ARG CA_CERTIFICATES_VERSION=20230311        # https://packages.debian.org/bookworm/ca-certificates
ARG CURL_VERSION=7.88.1*                    # https://packages.debian.org/bookworm/curl

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates=${CA_CERTIFICATES_VERSION} curl=${CURL_VERSION}

RUN useradd -m -u 1000 appuser && \
    chown appuser:appuser /app/ddns-updater

HEALTHCHECK CMD curl --fail http://localhost:80

#checkov:skip=CKV_DOCKER_3:Irrelevant

USER appuser

EXPOSE 80

ENTRYPOINT ["./ddns-updater"]
