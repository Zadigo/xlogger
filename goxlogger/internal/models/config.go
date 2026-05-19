package models

type LogServerConfig struct {
	Interval     string `yaml:"interval"`
	LocationName string `yaml:"location"`
}

type YamlConfig struct {
	LogServerConfig LogServerConfig `yaml:"log_server_config"`
}

type ServerConfig struct {
	RootDir    string
	YamlConfig YamlConfig
}
