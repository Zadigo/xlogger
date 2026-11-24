package utils

import (
	"context"
	"os"
	"path"

	"github.com/goccy/go-yaml"
	"github.com/redis/go-redis/v9"
)

type LogServerConfig struct {
	Interval string   `yaml:"interval"`
	Logs     []string `yaml:"logs"` // Directory for the logs
}

type ServerOptions struct {
	LogServer LogServerConfig `yaml:"log_server"`
	Redis     *redis.Client   `yaml:"-"`
}

func (c *ServerOptions) Load(ctx context.Context) error {
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
