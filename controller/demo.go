package controller

import (
	"go-final/dbconn"
	"go-final/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UserController(router *gin.Engine) {
	router.GET("/ping", ping)
	router.GET("/customers", getCustomers)
	router.POST("/login", login) // เพิ่มเส้นทาง /login
	router.POST("/register", register)
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pingpong broooo",
	})
}

// login - ฟังก์ชันสำหรับตรวจสอบการเข้าสู่ระบบด้วยอีเมลและรหัสผ่าน
func login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// อ่านข้อมูลจาก request
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// ค้นหา Email ในฐานข้อมูล
	var customer model.Customer
	if err := dbconn.DB.Where("email = ?", loginData.Email).First(&customer).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// ตรวจสอบรหัสผ่านที่กรอกกับรหัสที่แฮชในฐานข้อมูล
	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(loginData.Password)); err != nil {
		// ถ้ารหัสผ่านไม่ตรงกัน
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// ลบรหัสผ่านออกจากข้อมูลก่อนส่งกลับ
	customer.Password = "" // กำหนดให้เป็นค่าว่าง

	// ถ้ารหัสผ่านถูกต้อง
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"customer": gin.H{
			"customer_id":  customer.CustomerID,
			"first_name":   customer.FirstName,
			"last_name":    customer.LastName,
			"email":        customer.Email,
			"phone_number": customer.PhoneNumber,
			"address":      customer.Address,
			"created_at":   customer.CreatedAt,
			"updated_at":   customer.UpdatedAt,
		},
	})
}

func register(c *gin.Context) {
	var registerData struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
		Password    string `json:"password"`
	}

	// อ่านข้อมูลจาก request
	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// ตรวจสอบว่า email ซ้ำหรือไม่
	var existingCustomer model.Customer
	if err := dbconn.DB.Where("email = ?", registerData.Email).First(&existingCustomer).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// แฮชรหัสผ่าน
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// สร้าง customer ใหม่
	newCustomer := model.Customer{
		FirstName:   registerData.FirstName,
		LastName:    registerData.LastName,
		Email:       registerData.Email,
		PhoneNumber: registerData.PhoneNumber,
		Address:     registerData.Address,
		Password:    string(hashedPassword), // แฮชรหัสผ่าน
	}

	// บันทึกข้อมูลลูกค้าใหม่ลงฐานข้อมูล
	if err := dbconn.DB.Create(&newCustomer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful", "customer": newCustomer})
}

func getCustomers(c *gin.Context) {
	var customers []model.Customer
	if err := dbconn.DB.Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve customers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"customers": customers})
}
