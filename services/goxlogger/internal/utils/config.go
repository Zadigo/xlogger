package utils

import (
	"context"
	"os"
	"path"

	"github.com/goccy/go-yaml"
)

type LogServerConfig struct {
	Interval string `yaml:"interval"`
	Logs     struct {
		Folder string `yaml:"folder"`
	} `yaml:"logs"`
}

type ServerConfig struct {
	LogServer LogServerConfig `yaml:"log_server"`
	Redis     struct {
		Addr string `yaml:"addr"`
	} `yaml:"redis"`
}

func (c *ServerConfig) Load(ctx context.Context) error {
	rootDir := ctx.Value("rootDir").(string)
	filePath := path.Join(rootDir, "config.yaml")

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
