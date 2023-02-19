# syntax=docker/dockerfile:1

## Build
FROM golang:1.20-bullseye AS build

WORKDIR /app
COPY . /app

RUN go mod tidy && go build -o /ddns cmd/main.go

## Deploy
FROM debian:stable-slim

WORKDIR /app

COPY --from=build /ddns /app/ddns

EXPOSE 80

ENTRYPOINT ["./ddns"]