package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// Config is the config struct
type Config struct {
	EqDir            string `yaml:"eq_dir" json:"eq_dir"`
	CurrentExpansion int    `yaml:"current_expansion" json:"current_expansion"`
}

const configFileName = "eq-expansion-switcher-config.yaml"

// Save saves the config struct to the config file
func Save(config Config) error {
	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	fmt.Println(string(yamlData))

	f, err := os.Create(configFileName)
	if err != nil {
		return err
	}

	_, err = f.Write(yamlData)
	if err != nil {
		return err
	}

	_ = f.Close()

	return nil
}

// Get returns the config struct
func Get() Config {
	// check if file exists
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		fmt.Println("Config file does not exist")
		return Config{}
	}

	var config Config
	yamlFile, err := os.ReadFile(configFileName)
	if err != nil {
		return config
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return Config{}
	}

	return config
}
