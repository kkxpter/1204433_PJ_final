package controller

import "github.com/gin-gonic/gin"

func StartServer() {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	UserController(router)
	ProductController(router)
	router.Run(":8081")
}
