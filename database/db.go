package database

import (
	"fmt"
	"go-jwt-challenge/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "/"
	dbPort   = 5432
	dbName   = "simple-api"
	db       *gorm.DB
	err      error
)

func StartDB() {

	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, dbName, dbPort)
	dsn := config
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to database :", err)
	}

	fmt.Println("Successfully connected to database")
	db.AutoMigrate(models.User{}, models.Product{})
}

func GetDB() *gorm.DB {
	return db
}
