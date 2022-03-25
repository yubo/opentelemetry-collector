package geoipprocessor

import (
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
	IgnoreMissing bool `mapstructure:"ignore_missing"`
	// If true only first found geoip data will be returned, even if field contains array
	FirstOnly bool `mapstructure:"first_only"`
	// Controls what properties are added to the target_field based on the geoip lookup.
	Properties []string `mapstructure:"properties"`
}

var _ config.Processor = (*Config)(nil)

// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() error {
	return nil
}
