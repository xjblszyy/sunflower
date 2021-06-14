package config

import (
	"os"

	"github.com/jinzhu/configor"
	"go.uber.org/zap"
)

type Config struct {
	Debug       bool   `yaml:"debug,omitempty" default:"false"`
	ServiceName string `yaml:"service_name,omitempty" default:"sunflower"`

	Logger   LoggerConfig   `yaml:"logger,omitempty"`
	Database DatabaseConfig `yaml:"database,omitempty"`
	Server   ServerConfig   `yaml:"server,omitempty"`
	File     FileConfig     `yaml:"file,omitempty"`
}

type DatabaseConfig struct {
	// 仅支持 mysql
	DSN          string `yaml:"dsn"`
	MaxIdleConns int    `yaml:"max_idle_conns" default:"10"`
	MaxOpenConns int    `yaml:"max_open_conns" default:"100"`
	// format: https://golang.org/pkg/time/#ParseDuration
	ConnMaxLifetime string `yaml:"conn_max_lifetime" default:"1h"`
}

type LoggerConfig struct {
	Level string `yaml:"level,omitempty" default:"debug"`
	// json or text
	Format string `yaml:"format,omitempty" default:"json"`
	// file
	Output string `yaml:"output,omitempty" default:""`
}

type ServerConfig struct {
	HttpAddr string `yaml:"http_addr,omitempty" default:":8080"`
	GrpcAddr string `yaml:"grpc_addr,omitempty" default:":8082"`
}

type FileConfig struct {
	// 允许上传的文件大小 单位：byte,10485760 = 10MB
	MaxSize  int64 `yaml:"max_size" default:"10485760"`
	MaxCount int   `yaml:"max_count" default:"5"`
}

var C *Config

func Init(cfgFile string) {
	_ = os.Setenv("GUARD_ENV_PREFIX", "-")

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	if cfgFile != "" {
		if err := configor.New(&configor.Config{AutoReload: true}).Load(C, cfgFile); err != nil {
			zap.L().Panic("init config fail", zap.Error(err))
		}
	} else {
		if err := configor.New(&configor.Config{AutoReload: true}).Load(C); err != nil {
			zap.L().Panic("init config fail", zap.Error(err))
		}
	}

	zap.L().Debug("loaded config")
}

func init() {
	C = &Config{}
}
