package main

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	Mackerel MackerelConfig
	App      AppConfig
}

type MackerelConfig struct {
	Apikey string "toml:apikey"
	Hostid string "toml:hostid"
}

type AppConfig struct {
	Volt int    "toml:volt"
	Name string "toml:name"
}

func LoadConfig() *Config {
	var config Config
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		log.Fatalln(err)
	}
	return &config
}
