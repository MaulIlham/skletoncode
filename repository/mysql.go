package repository

import (
	"ELKExample/models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rs/zerolog"
	"os"
)

type Databases struct {
	Conn *gorm.DB
	Logger zerolog.Logger
}


func InitDB(keyUser, keyPassword, keyHost, keyPort, keySchema string, logger zerolog.Logger) (Databases, error) {
	db := Databases{}
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		os.Getenv(keyUser),
		os.Getenv(keyPassword),
		os.Getenv(keyHost),
		os.Getenv(keyPort),
		os.Getenv(keySchema),
	)

	conn, err := gorm.Open("mysql", dataSource)
	if err != nil {
		return db, err
	}

	db.Conn = conn
	db.Logger = logger

	return db, nil
}

func InitMigration(db *gorm.DB)  {
	db.AutoMigrate(&models.Posts{})
	db.AutoMigrate(&models.PostLog{})
}

