package config

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App      App            `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
}

type App struct {
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Admin  Addmin       `yaml:"admin"`
}

type Addmin struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Timezone string `yaml:"timezone"`
}

func LoadConfig(configFilePath string) Config {
	file, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("Error opening config file: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Error unmarshaling config: %s", err)
	}

	return config
}
