package conf

import "github.com/caarlos0/env"

type AppConfig struct {
	Environment string `env:"ENVIRONMENT"`
	EnvDev      string `env:"ENV_DEV"`

	Port string `env:"PORT"`

	RBHost   string `env:"RB_HOST"`
	RBPort   string `env:"RB_PORT"`
	RBUser   string `env:"RB_USER"`
	RBPass   string `env:"RB_PASS"`
	RBPortUI string `env:"RB_PORT_UI"`

	QueueName     string `env:"QUEUE_NAME"`
	NumberWorkers string `env:"NUMBER_WORKER"`
}

var config AppConfig

func SetEnv() {
	_ = env.Parse(&config)
}

func LoadEnv() AppConfig {
	return config
}
