package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config holds constant parameters for the app
type Config struct {
	HTTPClient struct {
		UserAgent string `yaml:"user-agent"`
		Timeout   int    `yaml:"timeout"`
	} `yaml:"http-client"`

	HTTPServer struct {
		Port string `yaml:"port"`
	} `yaml:"http-server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
}

// LoadFile opens and parses the YAML config file
func LoadFile(path string) (*Config, error) {
	f, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("opening config.yaml: %v", err)
	}

	cfg := Config{}
	err = yaml.Unmarshal([]byte(f), &cfg)
	if err != nil {
		return nil, fmt.Errorf("parsing config.yaml: %v", err)
	}
	return &cfg, nil
}
