package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Redis *RedisConfig `mapstructure:"redis"`

	Postgres *PostgresConfig `mapstructure:"postgres"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	Timeout  int    `mapstructure:"timeout"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       string `mapstructure:"db"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

var conf *AppConfig
var once *sync.Once = &sync.Once{}

func ParseConfig(path string) (*AppConfig, error) {
	var (
		err error
	)

	once.Do(func() {
		conf = &AppConfig{}
		viper.SetConfigName(path)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")

		if err = viper.ReadInConfig(); err != nil {
			err = fmt.Errorf("failed to read config: %w", err)
			return
		}

		if err = viper.Unmarshal(conf); err != nil {
			err = fmt.Errorf("failed to unmarshal config: %w", err)
			return
		}
	})
	if err != nil {
		return nil, err
	}

	return conf, nil
}
