package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	HomeDir     string
	DataPort    uint16
	LogLevel    string
	Moniker     string
	NodeKeyFile string
	StateFile   string
}

// NewConfig creates a new Config instance
func NewConfig(homeDir string) *Config {
	return &Config{
		HomeDir:     homeDir,
		DataPort:    6900,
		LogLevel:    "info",
		StateFile:   filepath.Join(homeDir, "dseq.bin"),
		NodeKeyFile: filepath.Join(homeDir, "config", "node_key.json"),
	}
}

// Load loads configuration from file
func (c *Config) Load() error {
	viper.SetConfigFile(filepath.Join(c.HomeDir, "config", "config.toml"))
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	c.Moniker = viper.GetString("moniker")
	c.LogLevel = viper.GetString("log_level")
	if port := viper.GetUint("data_port"); port > 0 {
		c.DataPort = uint16(port)
	}

	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.HomeDir == "" {
		return fmt.Errorf("home directory cannot be empty")
	}
	if c.DataPort == 0 {
		return fmt.Errorf("data port cannot be 0")
	}
	if c.Moniker == "" {
		return fmt.Errorf("moniker cannot be empty")
	}
	return nil
}
