package DataBase

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbase *gorm.DB

func Init() {

	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=1221 dbname=Lesson1 port=5432 sslmode=disable"))
	if err != nil {
		fmt.Println("failed to connect database")
	}

	//err = db.AutoMigrate(&models.User{})
	//if err != nil {
	//
	//}
	dbase = db

}

func GetDB() *gorm.DB {
	return dbase
}
