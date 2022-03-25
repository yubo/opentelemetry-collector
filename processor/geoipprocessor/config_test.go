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

	assert.Equal(t, cfg.Processors[config.NewComponentID(typeStr)], &Config{
		ProcessorSettings: config.NewProcessorSettings(config.NewComponentID(typeStr)),
		Field:             []string{"ip"},
		DatabaseFile:      "./testdata/GeoLite2-City.mmdb",
		TargetField:       "geoip",
		IgnoreMissing:     false,
		FirstOnly:         true,
		Properties: []string{
			"continent_name",
			"country_iso_code",
			"country_name",
			"region_iso_code",
			"region_name",
			"city_name",
			"location",
		},
	})

	assert.Equal(t, cfg.Processors[config.NewComponentIDWithName(typeStr, "invalid")], &Config{
		ProcessorSettings: config.NewProcessorSettings(config.NewComponentIDWithName(typeStr, "invalid")),
	})
}
