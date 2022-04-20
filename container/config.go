package container

import "golearn-api-template/config"

func ProvideConfig() config.SuperConfig {
	config.Initialize(".", "yaml")
	if config.IsSet() == false {
		panic("Config must be set up")
	}
	return config.Config
}
