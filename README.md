# SOFPLICATOR: The Sophisticated Replicator
Is Harbor the glue of your company?
Is harbor failing too much?
Is harbor not scaling nicely?
Enter: The Sofplicator

Sofplicator aimes to provide easy scaling from one registry to another.


## Build
`docker build -f ./build/package/sofplicator . -t sofplicator:0.1.0`

## RUN
```
docker run \
-p 8080:8080 \
sofplicator:0.1.2

```