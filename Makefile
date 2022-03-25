all: otelcol

otelcol:
	go build -o $@ ./cmd/otelcol

.PHONY: clean
clean:
	rm -f otelcol
