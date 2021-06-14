package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Strings []string

func (s *Strings) Scan(src interface{}) error {
	switch typ := src.(type) {
	default:
		return fmt.Errorf("%s not supported", typ)
	case []byte:
		return json.Unmarshal(src.([]byte), s)
	}
}

func (s Strings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type BaseModel struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func IsDuplicateError(err error) bool {
	// nolint: errorlint
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok {
		if mysqlErr.Number == 1062 {
			return true
		}
	}

	return false
}
