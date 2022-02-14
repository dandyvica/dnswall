package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type Config struct {
	Resolver string
	Domains []string
}

// Load configuration file
func (config *Config) getConfig(configFile string) *Config {

    yamlFile, err := ioutil.ReadFile(configFile)
    if err != nil {
        log.Fatalf("unable to open config file: <%s>, error: <%v> ", configFile, err)
    }
    err = yaml.Unmarshal(yamlFile, config)
    if err != nil {
        log.Fatalf("unmarshal error: <%v> for config file: <%s>", err, configFile)
    }

    return config
}