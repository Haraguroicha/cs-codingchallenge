package Configs

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config use defaults and need initialized by loadConfigFromFile
var Config *Configs

// Configs is the configuration which loaded from config.yaml
type Configs struct {
	MaximumTopicLength int `yaml:"maxTopicLength"`
}

// NewConfig is initialized an empty config set
func NewConfig(filename string) *Configs {
	config := Configs{}
	config.loadConfigFromFile(filename)
	return &config
}

// loadConfigFromFile typicly is load config.yaml
func (config *Configs) loadConfigFromFile(filename string) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(raw, config)
	if err != nil {
		panic(err)
	}
}
