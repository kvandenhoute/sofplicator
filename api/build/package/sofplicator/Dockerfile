FROM golang:1.19.2-alpine3.16 AS build

WORKDIR /app

COPY . .

RUN go mod tidy

RUN GOOS=linux go build -o /app/bin/sofplicator cmd/sofplicator/*

FROM alpine:3.16.0

COPY --from=build /app/bin/sofplicator /usr/local/bin/

RUN apk add --no-cache tini

# Tini is now available at /sbin/tini
ENTRYPOINT ["/sbin/tini", "--", "sofplicator"]