package pkg

import (
	"fmt"
	"sync"
)

const dsn = ""

var (
	db     *gorm.DB
	dbonce sync.Once
)

func NewDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func GetDB() *gorm.DB {
	dbonce.Do(func() {
		var err error
		if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
			panic(fmt.Errorf("failed to connect database, err: %w", err))
		}
	})
	return db
}