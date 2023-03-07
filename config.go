package main

import (
	"fmt"
	"gopkg.in/yaml.v2"

	"log"
	"os"
)

var ConfigFile = "./polly_config.yaml"

type Db struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	IP       string `yaml:"ip"`
}

type Config struct {
	Db        Db      `yaml:"db"`
	AuthToken string  `yaml:"authToken"`
	Charger   Charger `yaml:"charger"`
}

type Charger struct {
	Name     string `yaml:"name"`
	IP       string `yaml:"ip"`
	Location string `yaml:"location"`
}

func defaultConfig() Config {
	return Config{Db: Db{Name: "admin", Password: "admin123"}, AuthToken: "my-token"}
}

func loadConfig() Config {
	config := defaultConfig()
	if _, err := os.Stat(ConfigFile); err != nil {
		return defaultConfig()
	}

	file, err := os.ReadFile(ConfigFile)
	if err != nil {
		log.Printf("Unable to load file %s, returning default config", err)
		return config
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded config file %s\n", config)
	return config
}
