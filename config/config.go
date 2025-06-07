package config

import (
	"github.com/spf13/viper"
)

// Environment types
const (
	DevEnvironment  = "development"
	TestEnvironment = "test"
	ProdEnvironment = "production"
)

type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	Port        string `mapstructure:"PORT"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      string `mapstructure:"DB_PORT"`
	DBUser      string `mapstructure:"DB_USER"`
	DBPass      string `mapstructure:"DB_PASSWORD"`
	DBName      string `mapstructure:"DB_NAME"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
}

func LoadFromFile(file string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(file)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	err := v.Unmarshal(config)
	return config, err
}

func Load() (*Config, error) {
	env := viper.GetString("ENVIRONMENT")
	if env == "" {
		env = DevEnvironment
	}

	configFile := ".env"
	switch env {
	case TestEnvironment:
		configFile = ".env.test"
	case ProdEnvironment:
		configFile = ".env.prod"
	}

	return LoadFromFile(configFile)
}
