package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App     App     `mapstructure:"app"`
	Metrics Metrics `mapstructure:"metrics"`
	S3      S3      `mapstructure:"s3"`
}

type App struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Secret   string `mapstructure:"secret"`
	LogLevel string `mapstructure:"loglevel"`
}

type Metrics struct {
	Port string `mapstructure:"port"`
}

type S3 struct {
	Bucket string `mapstructure:"bucket"`
}

func LoadConfig(paths ...string) (*Config, error) {
	if len(paths) != 0 {
		for _, path := range paths {
			viper.AddConfigPath(path)
		}
	}

	viper.AddConfigPath("/app/config")
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("app.secret", "APP_SECRET"); err != nil {
		return nil, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
