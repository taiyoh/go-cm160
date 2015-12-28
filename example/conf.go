package main

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	Volt     int
	Name     string
	Mackerel MackerelConfig
}

type MackerelConfig struct {
	Apikey string
	Hostid string
}

func LoadConfig() *Config {
	var config Config
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		log.Fatalln(err)
	}
	return &config
}
