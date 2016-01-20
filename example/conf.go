package main

import (
	"github.com/BurntSushi/toml"
	"log"
)

// Config : config object
type Config struct {
	Volt int
	M2x  M2xConfig
}

// M2xConfig : m2x section in Config
type M2xConfig struct {
	Apikey   string
	Deviceid string
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
