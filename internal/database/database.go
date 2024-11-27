package database

import (
	"app/cmd/auth/model"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Parameters struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
}

var DB *gorm.DB

// Инициализация БД
func Init(Params Parameters) {
	conf := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", Params.Host, Params.User, Params.Password, Params.DatabaseName, Params.Port)
	var err error
	DB, err = gorm.Open(postgres.Open(conf), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// Автоматическая миграция
	err = DB.AutoMigrate(&model.AuthRefreshTokens{})
	if err != nil {
		log.Fatalf("Database migration error: %v", err)
	}
}
