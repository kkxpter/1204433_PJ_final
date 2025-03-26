package dbconn

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBconnect() *gorm.DB {
	viper.SetConfigName("config") // เอาวงเล็บออก
	viper.AddConfigPath(".")      // ค้นหาไฟล์ config ในโฟลเดอร์ปัจจุบัน

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	dsn := viper.GetString("mysql.dsn")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Connection successful")
	DB = db // กำหนดให้ตัวแปร DB (ซึ่งเป็นตัวแปรแพ็กเกจ var DB *gorm.DB) อ้างอิงไปที่ฐานข้อมูลที่เชื่อมต่อไว้
	return db
}
