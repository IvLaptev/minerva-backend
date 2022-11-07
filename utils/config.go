package utils

import (
	"fmt"
	"minerva/types"
	"os"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Service struct {
		Host   string   `yaml:"host"`
		Port   string   `yaml:"port"`
		Master bool     `yaml:"master"`
		Slaves []string `yaml:"slaves"`
	} `yaml:"service"`
	Actions []types.Action `yaml:"actions"`
}

var Config Configuration

// Чтение файла конфигурации
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

	Config.Service.Host = Config.Service.Host + ":" + Config.Service.Port
	fmt.Println(Config)

	return nil
}

func GetConfig() Configuration {
	return Config
}
