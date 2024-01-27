package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

//	type ConfigSql struct {
//		DriverName string `yaml:"driverName"`
//		User       string `yaml:"userName"`
//		Pass       string `yaml:"password"`
//		Address    string `yaml:"address"`
//		Port       string `yaml:"port"`
//		Database   string `yaml:"database"`
//	}

type ConfigMain struct {
	AppKeyQiita string `yaml:"qiita_key"`
	AppKeyTSAPI string `yaml:"text_summarization_key"`
}

func _loadConfig(confObj any, filepath *string) error {
	bytes, err := os.ReadFile(*filepath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, confObj)
	if err != nil {
		return err
	}
	return nil
}

func (conf *ConfigMain) LoadConfig(filepath string) error {
	return _loadConfig(&conf, &filepath)
}

func LoadConfigMain(filepath string) (config *ConfigMain, err error) {
	result := new(ConfigMain)
	err = result.LoadConfig(filepath)
	if err != nil {
		return nil, err
	}
	return result, nil
}
