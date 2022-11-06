package types

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Service struct {
		Port   string   `yaml:"port"`
		Master bool     `yaml:"master"`
		Slaves []string `yaml:"slaves"`
	} `yaml:"service"`
}

var Config Configuration

func ReadConfig() error {
	f, err := os.Open("config.yaml")
	if err != nil {
		return err
	}

	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Config)
	if err != nil {
		return err
	}

	log.Println(Config)

	return nil
}

func GetConfig() Configuration {
	return Config
}
