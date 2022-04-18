package prometheusremotewritereceiver

import (
	"go.opentelemetry.io/collector/config"
)

// Config defines configuration for StatsD receiver.
type Config struct {
	config.ReceiverSettings `mapstructure:",squash"`
	Endpoint                string `mapstructure:"endpoint"`
	WriteURL                string `mapstructure:"write_url"`
	//Namespace               string `mapstructure:"namespace"`
}

func (c *Config) validate() error {
	if c.Endpoint == "" {
		c.Endpoint = "localhost:9090"
	}
	if c.WriteURL == "" {
		c.WriteURL = PromWriteURL
	}
	return nil
}
