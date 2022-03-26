package geoipprocessor

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/service/servicetest"
)

func TestLoadConfig(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.NoError(t, err)

	factories.Processors[typeStr] = NewFactory()

	cfg, err := servicetest.LoadConfigAndValidate(filepath.Join("testdata", "config.yaml"), factories)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	expected := &Config{
		ProcessorSettings: config.NewProcessorSettings(config.NewComponentID(typeStr)),
		Field:             []string{"ip"},
		DatabaseFile:      "./testdata/GeoIP2-City-Test.mmdb",
		TargetField:       "geoip",
		Properties: []string{
			"continent_name",
			"country_iso_code",
			"country_name",
			"region_iso_code",
			"region_name",
			"city_name",
			"location",
			"geohash",
		},
		HashPrecision: 3,
	}

	expected.Validate()

	assert.Equal(t, cfg.Processors[config.NewComponentID(typeStr)], expected)
}
