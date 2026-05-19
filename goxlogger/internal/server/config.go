package server

import "github.com/Zadigo/goxlogger/internal/models"

func LoadConfig(rootDir string) *models.ServerConfig {
	return &models.ServerConfig{
		RootDir: rootDir,
		YamlConfig: models.YamlConfig{
			LogServerConfig: models.LogServerConfig{
				Interval:     "* * * * * ",
				LocationName: "data",
			},
		},
	}
}
