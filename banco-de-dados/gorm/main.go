package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Notebook struct {
	ID    int `gorm:"primary_key"`
	Name  string
	Price float64
}

func main() {
	dsn := "root:root@tcp(localhost:3306)/goexpert"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Notebook{})

	db.Create(&Notebook{
		Name:  "MacBook Pro",
		Price: 25000.00,
	})
}
