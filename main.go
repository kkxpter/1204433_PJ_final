package main

import (
	"fmt"
	"go-final/controller"
	"go-final/dbconn"
	"go-final/model"
)

func main() {
	// เชื่อมต่อกับฐานข้อมูล
	db := dbconn.DBconnect()

	// ดึงข้อมูลลูกค้าจากฐานข้อมูล
	Customer := []model.Customer{}
	result := db.Find(&Customer)
	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Println("Customers:", Customer)

	// ดึงข้อมูลสินค้าจากฐานข้อมูล
	products := []model.Product{}
	resultp := db.Find(&products)
	if resultp.Error != nil {
		panic(resultp.Error)
	}
	fmt.Println("Products:", products)

	// เริ่มต้น server
	controller.StartServer()
}
