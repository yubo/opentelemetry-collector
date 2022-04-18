package geoipprocessor

import (
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
	"go.opentelemetry.io/collector/config"
)

// Config defines configuration for Resource processor.
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct

	// The field to get the ip address from for the geographical lookup.
	Field []string `mapstructure:"field"`
	// The database filename referring to a database the module ships with (GeoLite2-City.mmdb, GeoLite2-Country.mmdb, or GeoLite2-ASN.mmdb) or a custom database in the ingest-geoip config directory.
	DatabaseFile string `mapstructure:"database_file"`
	// The field that will hold the geographical information looked up from the MaxMind database.
	TargetField string `mapstructure:"target_field"`
	// If true and field does not exist, the processor quietly exits without modifying the document
	// IgnoreMissing bool `mapstructure:"ignore_missing"`
	// If true only first found geoip data will be returned, even if field contains array
	// FirstOnly bool `mapstructure:"first_only"`
	// Controls what properties are added to the target_field based on the geoip lookup.
	// "continent_name", "country_iso_code", "country_name", "city_name", "location", "geohash"
	Properties []string `mapstructure:"properties"`

	HashPrecision uint `mapstructure:"hash_precision"`

	reader dbReader
}

type dbReader interface {
	City(ipAddress net.IP) (*geoip2.City, error)
}

var (
	defProperties = []string{"continent_name", "country_iso_code", "country_name", "city_name", "location"}
	mockDatabase  = "mock-database"
)

var _ config.Processor = (*Config)(nil)

// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() (err error) {
	if len(cfg.Properties) == 0 {
		cfg.Properties = defProperties
	}
	if cfg.HashPrecision <= 0 || cfg.HashPrecision > 12 {
		cfg.HashPrecision = 12
	}

	if cfg.DatabaseFile == "" {
		cfg.DatabaseFile = os.Getenv("GEOIP_DB_FILE")
	}

	if cfg.reader, err = geoip2.Open(cfg.DatabaseFile); err != nil {
		return err
	}

	return nil
}
