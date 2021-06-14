package cmd

import (
	cfg "sunflower/config"
	"sunflower/pkg/simple/client/db"

	"go.uber.org/zap"
)

// 初始化db
func setupDB() error {
	db.InitDB(cfg.C)
	db.AutoMigrateDB()

	return nil
}

func setupLogger() error {
	var conf zap.Config
	if cfg.C.Debug {
		conf = zap.NewDevelopmentConfig()
	} else {
		conf = zap.NewProductionConfig()
	}

	var zapLevel = zap.NewAtomicLevel()
	if err := zapLevel.UnmarshalText([]byte(cfg.C.Logger.Level)); err != nil {
		zap.L().Panic("set logger level fail",
			zap.Strings("only", []string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}),
			zap.Error(err),
		)
	}

	conf.Level = zapLevel
	conf.Encoding = "json"

	if cfg.C.Logger.Output != "" {
		conf.OutputPaths = []string{cfg.C.Logger.Output}
		conf.ErrorOutputPaths = []string{cfg.C.Logger.Output}
	}

	logger, _ := conf.Build()

	zap.RedirectStdLog(logger)
	zap.ReplaceGlobals(logger)

	return nil
}
