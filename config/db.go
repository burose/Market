package config

import (
	"log"
	"market/global"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initdb() {
	dsn := Appconfig.Database.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		log.Fatal(err)
	}
	global.DB = db
	
}
