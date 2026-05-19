package server

import "github.com/Zadigo/xlogger_backend/internal/models"

func LoadConfig(rootDir string) *models.ServerConfig {
	return &models.ServerConfig{
		RootDir: rootDir,
		YamlConfig: models.YamlConfig{
			LogServerConfig: models.LogServerConfig{
				Interval:     "0 0 * * *",
				LocationName: "data",
			},
		},
	}
}
