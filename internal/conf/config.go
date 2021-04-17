package conf

import (
	"encoding/json"
	"os"
	"path"
)

type Config struct {
	Host                string
	Token               string
	DefaultActivityDesc string
}

func getConfDir() (string, error) {
	userConfDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	confDir := path.Join(userConfDir, "overtime-cli")
	return confDir, os.MkdirAll(confDir, 0765)
}

func WriteConfig(c Config) error {
	confDir, err := getConfDir()
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(&c, "", " ")

	return os.WriteFile(path.Join(confDir, "conf.json"), bytes, os.ModePerm)
}

func LoadConfig() (*Config, error) {
	confDir, err := getConfDir()
	if err != nil {
		return nil, err
	}
	bytes, err := os.ReadFile(path.Join(confDir, "conf.json"))
	if err != nil {
		return nil, err
	}
	var c Config
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
