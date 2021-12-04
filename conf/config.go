package conf

import "github.com/caarlos0/env"

type AppConfig struct {
	Environment string `env:"ENVIRONMENT"`
	EnvDev      string `env:"ENV_DEV"`

	Port string `env:"PORT" envDefault:"8082"`

	RBHost   string `env:"RB_HOST" envDefault:"localhost"`
	RBPort   string `env:"RB_PORT" envDefault:"5672"`
	RBUser   string `env:"RB_USER" envDefault:"guest"`
	RBPass   string `env:"RB_PASS" envDefault:"guest"`
	RBPortUI string `env:"RB_PORT_UI" envDefault:"15672"`

	QueueName     string `env:"QUEUE_NAME" envDefault:"cr-product"`
	NumberWorkers string `env:"NUMBER_WORKER" envDefault:"5"`
}

var config AppConfig

func SetEnv() {
	_ = env.Parse(&config)
}

func LoadEnv() AppConfig {
	return config
}
