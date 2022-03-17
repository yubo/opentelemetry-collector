all: otelcol

otelcol:
	go build -o $@ ./cmd/otelcol