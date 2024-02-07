package gateway

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	DefaultTarget     string
	TerminalsFilePath string
	DBPath            string
	ListenAddress     string
	TerminalFieldName string
	GracefulTimeout   time.Duration
	LogDir            string
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	IdleTimeout       time.Duration
}

func NewConfig(filepath string) (config *Config, err error) {
	return loadConfig(filepath)
}

// loadConfig reads a file using given filepath and decodes it into *Config
func loadConfig(filepath string) (config *Config, err error) {
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	config = &Config{}
	err = json.Unmarshal(fileContent, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return config, nil
}
