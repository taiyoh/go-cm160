package main

import (
	"github.com/BurntSushi/toml"
	"log"
)

// Config : config object
type Config struct {
	Volt     int
	Name     string
	Mackerel MackerelConfig
}

// MackerelConfig : mackerel section in Config
type MackerelConfig struct {
	Apikey string
	Hostid string
}

// LoadConfig : load from toml file
func LoadConfig(file string) *Config {
	var config Config
	_, err := toml.DecodeFile(file, &config)
	if err != nil {
		log.Fatalln(err)
	}
	return &config
}
