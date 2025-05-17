// ErrorEnhanced: 2025-05-17
// Package config handles loading, parsing, and validating application configuration.
// It defines the structure for configuration settings, provides default values,
// loads settings from files (e.g., YAML), and applies overrides from environment variables.
// file: internal/config/config.go
package config

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	// Ensure this import path is correct for your hello-tool-base project structure
	"github.com/dkoosis/hello-tool-base/internal/logging"
	"gopkg.in/yaml.v3"
)

// ServerConfig contains settings specific to the server component.
type ServerConfig struct {
	Name            string        `yaml:"name"`
	Port            int           `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	IdleTimeout     time.Duration `yaml:"idleTimeout"`
	GracefulTimeout time.Duration `yaml:"gracefulTimeout"`
}

// Config is the root configuration structure for the application.
type Config struct {
	Server ServerConfig `yaml:"server"`
}

// DefaultConfig returns a configuration populated with default values.
func DefaultConfig() *Config {
	cfg := &Config{
		Server: ServerConfig{
			Name:            "HelloToolBase Service",
			Port:            8080,
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			IdleTimeout:     60 * time.Second,
			GracefulTimeout: 15 * time.Second,
		},
	}
	applyEnvironmentOverrides(cfg, logging.GetLogger("config_default"))
	return cfg
}

// LoadFromFile loads configuration from the specified YAML file path.
// It expands '~' to the user's home directory, reads the file, unmarshals YAML,
// and applies environment variable overrides.
func LoadFromFile(path string) (*Config, error) {
	logger := logging.GetLogger("config_load")

	if len(path) > 0 && path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Added function context to the error message
			return nil, errors.Wrap(err, "LoadFromFile: failed to get home directory to expand path")
		}
		path = filepath.Join(homeDir, path[1:])
		logger.Debug("Expanded config path from '~'", "originalPath", path, "expandedPath", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info("Configuration file not found, using defaults and environment variables.", "path", path)
			config := DefaultConfig()
			return config, nil
		}
		// Added function context to the error message
		return nil, errors.Wrapf(err, "LoadFromFile: failed to read config file: %s", path)
	}
	logger.Info("Successfully read configuration file.", "path", path)

	config := DefaultConfig()

	if err := yaml.Unmarshal(data, config); err != nil {
		// Added function context to the error message
		return nil, errors.Wrapf(err, "LoadFromFile: failed to parse config file YAML: %s", path)
	}
	logger.Info("Successfully unmarshalled YAML from configuration file.", "path", path)

	applyEnvironmentOverrides(config, logger)
	return config, nil
}

// applyEnvironmentOverrides applies configuration overrides from environment variables.
// It logs any overrides applied or any issues encountered during parsing of env vars.
func applyEnvironmentOverrides(config *Config, logger logging.Logger) {
	// Server Port
	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil && port > 0 && port < 65536 { // err scoped to this if/else
			logger.Debug("Overriding server port from environment.", "envVar", "SERVER_PORT", "oldValue", config.Server.Port, "newValue", port)
			config.Server.Port = port
		} else { // err is accessible here
			logger.Warn("Invalid SERVER_PORT environment variable ignored.", "value", portStr, "error", err)
		}
	} else {
		logger.Debug("No SERVER_PORT environment override found, using default/config file value.", "value", config.Server.Port)
	}

	// Server Name
	if serverName := os.Getenv("SERVER_NAME"); serverName != "" {
		logger.Debug("Overriding server name from environment.", "envVar", "SERVER_NAME", "oldValue", config.Server.Name, "newValue", serverName)
		config.Server.Name = serverName
	} else {
		logger.Debug("No SERVER_NAME environment override found, using default/config file value.", "value", config.Server.Name)
	}

	// Helper for parsing duration from environment variable
	getDurationEnv := func(envVar string, currentVal time.Duration, varNameHuman string) time.Duration {
		envValStr := os.Getenv(envVar)
		if envValStr != "" {
			duration, err := time.ParseDuration(envValStr) // Declare duration and err here
			if err == nil {                                // Check err
				if duration != currentVal {
					logger.Debug("Overriding "+varNameHuman+" from environment.", "envVar", envVar, "oldValue", currentVal, "newValue", duration)
				} else {
					logger.Debug("Environment variable "+envVar+" for "+varNameHuman+" matches current value.", "value", duration)
				}
				return duration
			}
			// Now 'err' is in scope for this Warn log
			logger.Warn("Invalid "+envVar+" environment variable for "+varNameHuman+" ignored.", "value", envValStr, "error", err)
		} else {
			logger.Debug("No "+envVar+" environment override found for "+varNameHuman+", using current value.", "value", currentVal)
		}
		return currentVal
	}

	config.Server.ReadTimeout = getDurationEnv("SERVER_READ_TIMEOUT", config.Server.ReadTimeout, "server read timeout")
	config.Server.WriteTimeout = getDurationEnv("SERVER_WRITE_TIMEOUT", config.Server.WriteTimeout, "server write timeout")
	config.Server.IdleTimeout = getDurationEnv("SERVER_IDLE_TIMEOUT", config.Server.IdleTimeout, "server idle timeout")
	config.Server.GracefulTimeout = getDurationEnv("SERVER_GRACEFUL_TIMEOUT", config.Server.GracefulTimeout, "server graceful timeout")
}
