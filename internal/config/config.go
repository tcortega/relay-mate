package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ForwardTo      string `yaml:"forward_to"`
	StarterMessage string `yaml:"starter_message"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return attemptToLoadFromEnv()
	}
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}

func attemptToLoadFromEnv() (Config, error) {
	var cfg Config
	forwardTo := os.Getenv("FORWARD_TO")
	starterMessage := os.Getenv("STARTER_MESSAGE")

	if forwardTo == "" {
		return cfg, fmt.Errorf("FORWARD_TO environment variable is not set")
	}

	if starterMessage == "" {
		return cfg, fmt.Errorf("STARTER_MESSAGE environment variable is not set")
	}

	cfg.ForwardTo = forwardTo
	cfg.StarterMessage = starterMessage

	return cfg, nil
}
