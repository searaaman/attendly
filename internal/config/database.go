package config

import (
	"fmt"
	"log"

	"github.com/searaaman/attendly/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := Cfg.DatabaseURL
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			Cfg.DBHost, Cfg.DBUser, Cfg.DBPassword, Cfg.DBName, Cfg.DBPort, Cfg.DBSSLMode,
		)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to db: ", err)
	}
	DB = db
	err = DB.AutoMigrate(
		&domain.Session{},
		&domain.Member{},
		&domain.TimeSlot{},
		&domain.Payment{},
		&domain.User{},
	)
	if err != nil {
		log.Fatal("Automigrate was not succesfull: ", err)
	}
	fmt.Println("Succesfully connected to PostgreSQL")
}
