FROM golang:1.17 as build

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM busybox

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY otelcol /

ENV OTEL_EXPORTER_OTLP_ENDPOINT=0.0.0.0:4317
EXPOSE 4317 55680 55679

ENTRYPOINT ["/otelcol"]
CMD ["--config", "/etc/otel/config.yaml"]
