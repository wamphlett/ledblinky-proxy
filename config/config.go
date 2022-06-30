package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds all the information required to run the proxy
type Config struct {
	LEDBlinkyPath string           `yaml:"ledblinkyPath"`
	Receivers     *ReceiversConfig `yaml:"receivers"`
}

// ReceiversConfig holds information about the receivers which should be published to
type ReceiversConfig struct {
	Executables []string `yaml:"executables"`
	Webhooks    []string `yaml:"webhooks"`
}

// NewFromFile creates a new Config from the config file in the same directory as the
// main executable
func NewFromFile() (*Config, error) {
	buf, err := ioutil.ReadFile(filepath.Join(filepath.Dir(os.Args[0]), "ledblinky-proxy.yaml"))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to read config file: %s", err.Error()))
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to unmarshal config file: %s", err.Error()))
	}

	return c, nil
}
