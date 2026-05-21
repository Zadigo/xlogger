package server

import (
	"github.com/Zadigo/goxlogger/internal/models"
)

func LoadConfig(rootDir string) *models.ServerConfig {
	return &models.ServerConfig{
		RootDir: rootDir,
		YamlConfig: models.YamlConfig{
			LogServer: models.LogServerConfig{
				Interval: "* * * * * ",
				Logs: struct {
					Folder string `yaml:"folder"`
				}{
					Folder: "data",
				},
			},
		},
	}
}
