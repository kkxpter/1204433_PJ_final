package main

import (
	"fmt"
	"go-final/controller"
	"go-final/dbconn"
	"go-final/model"
)

func main() {
	db := dbconn.DBconnect()
	Customer := []model.Customer{}
	result := db.Find(&Customer)
	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Print(Customer)

	controller.StartServer()
}
