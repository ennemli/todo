package db

import (
	"fmt"

	"github.com/ennemli/todo/user/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	configs := configs.GetConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		configs.DBConfig.DB_HOST,
		configs.DBConfig.DB_USER,
		configs.DBConfig.DB_PASSWORD,
		configs.DBConfig.DB_NAME,
		configs.DBConfig.DB_PORT,
	)
	if db == nil {
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	}
	return db
}
