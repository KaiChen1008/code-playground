package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig              `mapstructure:"server"`
	Languages map[string]LanguageConfig `mapstructure:"languages"`
}

type ServerConfig struct {
	Port                 int           `mapstructure:"port"`
	DataDir              string        `mapstructure:"data_dir"`
	RateLimit            int           `mapstructure:"rate_limit"`
	MaxTotalSubmissions  int           `mapstructure:"max_total_submissions"`
	MaxCodeChars         int           `mapstructure:"max_code_chars"`
	MaxConcurrentRunners int           `mapstructure:"max_concurrent_runners"`
	ExecutionTimeout     time.Duration `mapstructure:"execution_timeout"`
}

type LanguageConfig struct {
	Image   string `mapstructure:"image"`
	Version string `mapstructure:"version"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Support environment variables override
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("Config file not found, using defaults")
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if port := os.Getenv("PORT"); port != "" {
		fmt.Sscanf(port, "%d", &config.Server.Port)
	}
	if dataDir := os.Getenv("DATA_DIR"); dataDir != "" {
		config.Server.DataDir = dataDir
	}

	return &config, nil
}
