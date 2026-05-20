package models

import (
	"os"

	"github.com/goccy/go-yaml"
)

type LogServerConfig struct {
	Interval string `yaml:"interval"`
	Logs     struct {
		Folder string `yaml:"folder"`
	} `yaml:"logs"`
}

type YamlConfig struct {
	LogServer LogServerConfig `yaml:"log_server"`
	Redis     struct {
		Addr string `yaml:"addr"`
	}
}

func (c *YamlConfig) Load(rootDir string) error {
	filePath := rootDir + "/config.yaml"
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, c)
	if err != nil {
		return err
	}
	return nil
}

type ServerConfig struct {
	RootDir    string
	YamlConfig YamlConfig
}
