package DataBase

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"workWithDB/models"
)

var db *gorm.DB

var TEST = false

func init() {
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=1221 dbname=Lesson1 port=5432 sslmode=disable"))
	if err != nil {
		fmt.Println("failed to connect database")
		return
	}

	if TEST {
		dbName := "test_Lesson1"
		err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
		if err != nil {
			fmt.Println("failed to drop database lesson1_test")
			return
		}
		err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			fmt.Println("failed to create database lesson1_test")
			return
		}

		dbTest, err2 := gorm.Open(postgres.Open("host=localhost user=postgres password=1221 dbname=test_Lesson1 port=5432 sslmode=disable"))
		if err2 != nil {
			fmt.Println("failed to connect database")
			return
		}
		db = dbTest
	}

	db.AutoMigrate(&models.User{})
}

func GetDB() *gorm.DB {
	return db
}
