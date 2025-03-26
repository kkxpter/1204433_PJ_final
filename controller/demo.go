package controller

import (
	"go-final/dbconn"
	"go-final/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UserController(router *gin.Engine) {
	router.GET("/ping", ping)
	router.GET("/customers", getCustomers)
	router.POST("/login", login) // เพิ่มเส้นทาง /login
	router.POST("/register", register)
	router.POST("/changepass", changePassword)
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pingpong brooootttt",
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

func changePassword(c *gin.Context) {
	var passwordData struct {
		Email       string `json:"email"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	// อ่านข้อมูลจาก request
	if err := c.ShouldBindJSON(&passwordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// ค้นหาลูกค้าจาก email
	var customer model.Customer
	if err := dbconn.DB.Where("email = ?", passwordData.Email).First(&customer).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// ตรวจสอบ old_password กับ hashed password ในฐานข้อมูล
	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(passwordData.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Old password is incorrect"})
		return
	}

	// แฮชรหัสผ่านใหม่
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(passwordData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	// อัปเดตรหัสผ่านใหม่ในฐานข้อมูล
	customer.Password = string(hashedNewPassword)
	if err := dbconn.DB.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func getCustomers(c *gin.Context) {
	var customers []model.Customer
	if err := dbconn.DB.Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve customers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"customers": customers})
}
func ProductController(router *gin.Engine) {
	router.GET("/search-products", searchProducts) // ค้นหาสินค้า
	router.POST("/add-to-cart", addToCart)         // เพิ่มสินค้าลงในรถเข็น
	router.GET("/view-carts", viewCarts)           // เพิ่มเส้นทางดูรถเข็นทั้งหมดของลูกค้า
}

// ฟังก์ชันสำหรับค้นหาสินค้าตามรายละเอียดและช่วงราคา
func searchProducts(c *gin.Context) {
	description := c.DefaultQuery("description", "")
	minPriceStr := c.DefaultQuery("min_price", "0")
	maxPriceStr := c.DefaultQuery("max_price", "10000")

	// แปลงราคาจาก string เป็น float64
	minPrice, err := strconv.Atoi(minPriceStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min_price"})
		return
	}

	maxPrice, err := strconv.Atoi(maxPriceStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max_price"})
		return
	}

	// ค้นหาสินค้าจากชื่อและช่วงราคา
	var products []model.Product
	if err := dbconn.DB.Where("product_name LIKE ? AND price BETWEEN ? AND ?", "%"+description+"%", minPrice, maxPrice).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func addToCart(c *gin.Context) {
	var cartData struct {
		CustomerID int    `json:"customer_id"`
		CartName   string `json:"cart_name"`
		ProductID  int    `json:"product_id"`
		Quantity   int    `json:"quantity"`
	}

	// อ่านข้อมูลจาก request
	if err := c.ShouldBindJSON(&cartData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// ถ้า CartName ไม่มีการระบุ ให้กำหนดเป็น default หรือชื่ออื่นๆ ที่เหมาะสม
	if cartData.CartName == "" {
		cartData.CartName = "default" // ใช้ชื่อ default หากไม่ได้กรอกชื่อรถเข็น
	}

	// ค้นหารถเข็นที่มีอยู่หรือสร้างใหม่
	var cart model.Cart
	if err := dbconn.DB.Where("customer_id = ? AND cart_name = ?", cartData.CustomerID, cartData.CartName).First(&cart).Error; err != nil {
		// ถ้าไม่พบรถเข็น, สร้างรถเข็นใหม่
		cart = model.Cart{
			CustomerID: cartData.CustomerID,
			CartName:   cartData.CartName,
		}
		if err := dbconn.DB.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}
	}

	// ตรวจสอบว่าในรถเข็นมีสินค้านี้อยู่แล้วหรือไม่
	var cartItem model.CartItem
	if err := dbconn.DB.Where("cart_id = ? AND product_id = ?", cart.CartID, cartData.ProductID).First(&cartItem).Error; err == nil {
		// ถ้ามีอยู่แล้ว เพิ่มจำนวน
		cartItem.Quantity += cartData.Quantity
		if err := dbconn.DB.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product quantity updated", "cart_item": cartItem})
		return
	}

	// ถ้ายังไม่มีในรถเข็น, เพิ่มสินค้าลงไป
	newCartItem := model.CartItem{
		CartID:    cart.CartID,
		ProductID: cartData.ProductID,
		Quantity:  cartData.Quantity,
	}
	if err := dbconn.DB.Create(&newCartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart", "cart_item": newCartItem})
}

// ฟังก์ชันสำหรับดูรถเข็นทั้งหมดของลูกค้า
func viewCarts(c *gin.Context) {
	// รับ customer_id จากพารามิเตอร์
	customerIDStr := c.DefaultQuery("customer_id", "")
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil || customerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	// ค้นหารถเข็นทั้งหมดของลูกค้าคนนั้น
	var carts []model.Cart
	if err := dbconn.DB.Where("customer_id = ?", customerID).Find(&carts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve carts"})
		return
	}

	// เตรียมผลลัพธ์ที่จะส่งกลับ
	var cartDetails []gin.H

	// ลูปผ่านรถเข็นทั้งหมด
	for _, cart := range carts {
		// ค้นหาสินค้าในรถเข็นนี้
		var cartItems []model.CartItem
		if err := dbconn.DB.Where("cart_id = ?", cart.CartID).Find(&cartItems).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve cart items"})
			return
		}

		// เตรียมข้อมูลสำหรับการแสดงผล
		var itemsDetails []gin.H

		// ลูปผ่านแต่ละรายการในรถเข็น
		for _, cartItem := range cartItems {
			// ค้นหาข้อมูลสินค้าจาก product_id
			var product model.Product
			if err := dbconn.DB.Where("product_id = ?", cartItem.ProductID).First(&product).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve product"})
				return
			}

			// เพิ่มข้อมูลสินค้าในรถเข็น
			itemsDetails = append(itemsDetails, gin.H{
				"product_id":     product.ProductID,
				"product_name":   product.ProductName,
				"quantity":       cartItem.Quantity,
				"price_per_unit": product.Price,
			})
		}

		// เพิ่มข้อมูลรถเข็นพร้อมรายละเอียดสินค้า
		cartDetails = append(cartDetails, gin.H{
			"cart_id":   cart.CartID,
			"cart_name": cart.CartName,
			"items":     itemsDetails,
		})
	}

	// ส่งข้อมูลรถเข็นทั้งหมดกลับไป
	c.JSON(http.StatusOK, gin.H{
		"customer_id": customerID,
		"carts":       cartDetails,
	})
}
