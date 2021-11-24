package conf

import "github.com/caarlos0/env"

type AppConfig struct {
	Port      string `env:"PORT" envDefault:"8084"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"text"`
}

var config AppConfig

func SetEnv() {
	_ = env.Parse(&config)
}

func LoadEnv() AppConfig {
	return config
}
