package main

import (
	"github.com/BurntSushi/toml"
	"log"
)

// Config for this app
type Config struct {
	Volt     int
	Name     string
	Mackerel MackerelConfig
}

// MackerelConfig for mackerel section in Config
type MackerelConfig struct {
	Apikey string
	Hostid string
}

// LoadConfig is loading config from toml
func LoadConfig(file string) *Config {
	var config Config
	_, err := toml.DecodeFile(file, &config)
	if err != nil {
		log.Fatalln(err)
	}
	return &config
}
