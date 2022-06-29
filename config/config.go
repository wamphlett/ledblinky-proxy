package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LEDBlinkyPath string           `yaml:"ledblinkyPath"`
	Receivers     *ReceiversConfig `yaml:"receivers"`
}

type ReceiversConfig struct {
	Executables []string `yaml:"executables"`
	Webhooks    []string `yaml:"webhooks"`
}

func NewFromINI() (*Config, error) {
	buf, err := ioutil.ReadFile(filepath.Join(filepath.Dir(os.Args[0]), "ledblinky-proxy.yaml"))
	if err != nil {

		log.Fatal(err.Error())
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		log.Fatal(err.Error())
		return nil, fmt.Errorf("failed to unmarshal config file")
	}

	return c, nil
}
