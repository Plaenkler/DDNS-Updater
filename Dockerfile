# syntax=docker/dockerfile:1

## Build
FROM golang:1.20-bullseye AS build

WORKDIR /app
COPY . /app

RUN go mod tidy && go build -o /ddns-updater cmd/main.go

## Deploy
FROM debian:stable-slim

WORKDIR /app

COPY --from=build /ddns-updater /app/ddns-updater

RUN apt-get update && \
    apt-get install -y ca-certificates

EXPOSE 80

ENTRYPOINT ["./ddns-updater"]