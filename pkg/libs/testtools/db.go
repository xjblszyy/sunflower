package testtools

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"

	gormUtil "sunflower/pkg/libs/gormutil_v2"
)

func MockedGORMDBForTest(t *testing.T, sqlDB *sql.DB) *gorm.DB {
	gormDB := gormUtil.MockedGORMDBForTest(t, sqlDB)
	return gormDB
}

func DbCnnForTest() *gorm.DB {
	viper.AutomaticEnv()
	dsn := viper.GetString("FOUNDATION_DATABASE_DSN")
	db, err := gormUtil.ConnectWithDSN(dsn, gormUtil.Conf{LogMode: true})
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("init test case db client error: %v", err))
	}
	return db.Debug()

}

// 清空表数据
func CleanTableForTest(db *gorm.DB, tableNames []string) (err error) {
	SQL := "TRUNCATE TABLE "

	for k, v := range tableNames {
		if len(tableNames)-1 == k {
			// 最后一条数据,以分号结尾
			SQL += fmt.Sprintf("%s;", v)
		} else {
			SQL += fmt.Sprintf("%s, ", v)
		}
	}
	if err := db.Exec(SQL).Error; err != nil {
		return err
	}
	return nil
}

func MockedZAPForTest(t *testing.T) *zap.Logger {
	return zap.L().With(zap.String("impl", "mock for test"))
}
