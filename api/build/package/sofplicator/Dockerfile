FROM golang:1.19.4-alpine3.17 AS build

WORKDIR /app

COPY . .

RUN go mod tidy

RUN GOOS=linux go build -o /app/bin/sofplicator cmd/sofplicator/*

FROM bellsoft/alpaquita-linux-base:stream-musl-221207

COPY --from=build /app/bin/sofplicator /usr/local/bin/

RUN apk add --no-cache tini

# Tini is now available at /sbin/tini
ENTRYPOINT ["/sbin/tini", "--", "sofplicator"]