package config

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Configuration struct {
	Name    string
	HTTP    HTTP
	Metrics Metrics
	Logger  Logger
}

type Logger struct {
	Level    zapcore.Level
	Encoding string
}

type HTTP struct {
	Port      int
	Profiling Profiling
}

type Metrics struct {
	Namespace string
}

type Profiling struct {
	Enabled bool
}

func Read(prefix, configPath string) (Configuration, error) {
	path, extension := filepath.Dir(configPath), filepath.Ext(configPath)
	file := strings.TrimSuffix(filepath.Base(configPath), extension)

	viper.AddConfigPath(path)
	viper.SetConfigName(file)
	viper.SetConfigType(extension[1:])
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetEnvPrefix(prefix)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Configuration{}, fmt.Errorf("error reading config file: %w", err)
	}

	var config Configuration
	err := viper.Unmarshal(&config)
	if err != nil {
		return Configuration{}, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return config, nil
}
