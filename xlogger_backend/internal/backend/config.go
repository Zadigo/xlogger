package backend

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type DbConfig struct {
	Provider string
	User     string
	Password string
	Host     string
	Port     int
	Name     string
}

type LogsFolderConfig struct {
	Name string
	// FileRegex string
}

type RedisConfig struct {
	Url string
}

type AllConfig struct {
	// Cron expression for the scheduler to run the job
	Analysis   string
	Db         DbConfig
	Redis      RedisConfig
	LogsFolder LogsFolderConfig
}

type ServerConfig struct {
	Config AllConfig
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func GetServerConfig() *ServerConfig {
	currentPath, err := os.Getwd()

	if err != nil {
		panic("🔴 Could not get current path")
	}

	content, err := os.ReadFile(currentPath + "/../config.yaml")
	if err != nil {
		panic(fmt.Sprintf("🔴 Could not read config file %s", err))
	}

	config := NewServerConfig()
	err = yaml.Unmarshal(content, config)
	if err != nil {
		panic(fmt.Sprintf("🔴 Could not parse config file %s", err))
	}

	return config
}
