package conf

import (
	"tcp/pkg/configuration"
)

type Config struct {
	ServerAddr  string
	ServerPort  string
	DatasetFile string
}

func ReadConfig() (Config, error) {
	reader := configuration.NewReader()
	cfg := Config{
		ServerAddr: reader.Optional.String("CONFIG_SERVER_ADDR", "pow_server"),
		ServerPort: reader.Optional.String("CONFIG_SERVER_PORT", ":7777"),
		DatasetFile: reader.Mandatory.String("CONFIG_SERVER_DATASET"),
	}
	return cfg, reader.Error()
}
