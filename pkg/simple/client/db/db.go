package db

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"sunflower/config"
)

var DB *gorm.DB

type Scope = func(db *gorm.DB) *gorm.DB

func InitDB(cfg *config.Config) {
	zap.L().Debug("connect db", zap.String("dsn", cfg.Database.DSN))

	var dbLogger logger.Interface

	if cfg.Debug {
		dbLogger = logger.Default.LogMode(logger.Info)
	} else {
		dbLogger = logger.Default.LogMode(logger.Warn)
	}

	var err error
	DB, err = gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{
		Logger: dbLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
		},
	})
	if err != nil {
		zap.L().Panic("connect db failed", zap.Error(err))
	}

	sqlDB, err := DB.DB()
	if err != nil {
		zap.L().Panic("connect db failed", zap.Error(err))
	}

	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	if cfg.Database.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	}

	// SetMaxOpenCons 设置数据库的最大连接数量。
	if cfg.Database.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	}

	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	if cfg.Database.ConnMaxLifetime != "" {
		maxLifetime, err := time.ParseDuration(cfg.Database.ConnMaxLifetime)
		if err != nil {
			zap.L().Panic("db ConnMaxLifetime parse failed", zap.Error(err))
		}

		sqlDB.SetConnMaxLifetime(maxLifetime)
	}

	if err := sqlDB.Ping(); err != nil {
		zap.L().Panic("ping db failed", zap.Error(err))
	}
}

func CloseDB() {
}

func AutoMigrateDB() {
	query := DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	if err := query.AutoMigrate(); err != nil {
		zap.L().Panic("migrate db fail", zap.Error(err))
	}
}
