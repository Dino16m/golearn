package container

import (
	"golearn-api-template/config"

	gwa "github.com/gobuffalo/gocraft-work-adapter"

	redigo "github.com/gomodule/redigo/redis"
)

func provideRedisPool(cfg config.SuperConfig) *redigo.Pool {
	options := cfg.RedisStoreOptions
	dialPasswordOption := redigo.DialPassword(options.Password)
	return &redigo.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial(
				options.Network, options.Address, dialPasswordOption)
		},
	}
}

func ProvideWorker(cfg config.SuperConfig) *gwa.Adapter {
	return gwa.New(gwa.Options{
		Pool:           provideRedisPool(cfg),
		Name:           cfg.AppName,
		MaxConcurrency: 25,
	})
}
