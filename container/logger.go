package container

import (
	"golearn-api-template/config"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func ProvideLogger(config config.SuperConfig) *logrus.Logger {
	log := logrus.New()
	cfg := config.LoggerConfig
	if cfg.Filename != "" {
		log.Out = &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
	}
	return log
}
