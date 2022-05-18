# Opentelemetry-collector

## Feature

- receiver
  * [prometheusremotewrite](./receiver/prometheusremotewritereceiver)
- processor
  * [geoip](./processor/geoipprocessor)

## build
```sh
$ git clone https://github.com/yubo/opentelemetry-collector.git
$ cd opentelemetry-collector
$ make
```

## Docker Build
```sh
$ make docker-otelcol
```

## Docker run
```
docker run --rm -i -t \            
	-w / \                     
	-v /etc/otel/config.yaml:/etc/otel/config.yaml \
	ybbbbasdf/otelcol:latest \ 
	--config /etc/otel/config.yaml
```
