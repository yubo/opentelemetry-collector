# GeoIP Processor

Supported pipeline types: metrics, traces, logs

The resource processor can be used to apply changes on resource attributes.
Please refer to [config.go](./config.go) for the config spec.


Examples:

```yaml
processors:
  geoip:
    # The field to get the ip address from for the geographical lookup.
    field:
      - ip
    # The database filename referring to a database the module ships with (GeoLite2-City.mmdb, GeoLite2-Country.mmdb, or GeoLite2-ASN.mmdb) or a custom database in the ingest-geoip config directory.
    database_file: GeoLite2-City.mmdb
    # The field that will hold the geographical information looked up from the MaxMind database.
    target_field: geoip
    # If true and field does not exist, the processor quietly exits without modifying the document
    ignore_missing: false
    # If true only first found geoip data will be returned, even if field contains array
    first_only: true
    # Controls what properties are added to the target_field based on the geoip lookup.
    properties:
      - continent_name
      - country_iso_code
      - country_name
      - region_iso_code
      - region_name
      - city_name
      - location
```

Refer to [config.yaml](./testdata/config.yaml) for detailed
examples on using the processor.

- https://www.elastic.co/guide/en/elasticsearch/reference/current/geoip-processor.html


