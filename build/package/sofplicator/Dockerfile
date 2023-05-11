FROM docker.io/golang:1.20.4-alpine3.17 AS build

WORKDIR /app

COPY . .

RUN go mod tidy

RUN GOOS=linux go build -o /app/bin/sofplicator cmd/sofplicator/*

FROM docker.io/bellsoft/alpaquita-linux-base:stream-musl-230424

COPY --from=build /app/bin/sofplicator /usr/local/bin/

RUN apk add --no-cache tini

# Tini is now available at /sbin/tini
ENTRYPOINT ["/sbin/tini", "--", "sofplicator"]