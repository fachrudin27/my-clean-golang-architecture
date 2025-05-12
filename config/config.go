package config

import (
	"github.com/spf13/viper"
)

type (
	// Config
	Config struct {
		App  `yaml:"app" mapstructure:"app"`
		HTTP `yaml:"http" mapstructure:"http"`
		Log  `yaml:"logger" mapstructure:"logger"`
		PG   `yaml:"postgres" mapstructure:"postgres"`
		RMQ  `yaml:"rmq" mapstructure:"rmq"`
		JWT  `yaml:"jwt" mapstructure:"jwt"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME" mapstructure:"name"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION" mapstructure:"version"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT" mapstructure:"port"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"level" mapstructure:"level"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" mapstructure:"pool_max"`
		URL     string `env-required:"true" yaml:"pg_url" mapstructure:"pg_url"`
	}

	// RMQ -.
	RMQ struct {
		RmqUrl string `env-required:"true" yaml:"rmq_url" mapstructure:"rmq_url"`
	}

	// JWT -,
	JWT struct {
		SecretKey string `env-required:"true" env:"SECRET_KEY" mapstructure:"secret_key"`
	}
)

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil

}
