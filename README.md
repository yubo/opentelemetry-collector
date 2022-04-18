# Opentelemetry-collector

## Feature

- receiver
  * [prometheusremotewrite](./receiver/prometheusremotewritereceiver)
- processor
  * [geoip](./processor/geoipprocessor)

## Installing
```sh
$ go install github.com/yubo/opentelemetry-collector/cmd/otelcol@latest
```

## Docker Build
```sh
$ make docker-otelcol
```
