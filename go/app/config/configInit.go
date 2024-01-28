package config

import "log"

var MainConfig ConfigMain

func init() {
	tmp, err := LoadConfigMain("config/config.yml")
	if err != nil {
		log.Fatalln(err)
	}
	MainConfig = *tmp
}
