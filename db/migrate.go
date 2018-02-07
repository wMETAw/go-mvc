package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"models"
)

func main() {
	db, err := gorm.Open("mysql", "root@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.CreateTable(&models.User{})
	db.Create(&models.User{Name: "Yamada", Age: 20})
	db.Create(&models.User{Name: "Suzuki", Age: 25})
	db.Create(&models.User{Name: "Sato", Age: 29})
}
