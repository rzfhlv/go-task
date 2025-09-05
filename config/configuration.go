package config

import (
	"context"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Configuration struct {
	App      AppConfiguration      `mapstructure:"app"`
	Database DatabaseConfiguration `mapstructure:"database"`
	Redis    RedisConfiguration    `mapstructure:"redis"`
	JWT      JWTConfiguration      `mapstructure:"jwt"`
}

type AppConfiguration struct {
	Env      string `mapstructure:"env"`
	Name     string `mapstructure:"name"`
	Port     string `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
}

type DatabaseConfiguration struct {
	Driver   string `mapstructure:"driver"`
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type RedisConfiguration struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type JWTConfiguration struct {
	Secret    string        `mapstructure:"secret"`
	ExpiresIn time.Duration `mapstructure:"expires_in"`
}

var (
	configuration *Configuration
	once          sync.Once
)

func All() *Configuration {
	once.Do(func() {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yml")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()

		err := viper.ReadInConfig()
		if err != nil {
			slog.ErrorContext(context.Background(), "failed to initiate config", slog.String("error", err.Error()))
		}

		viper.Unmarshal(&configuration)
	})

	return configuration
}

func Get() *Configuration {
	return configuration
}
