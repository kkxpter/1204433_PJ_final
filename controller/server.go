package controller

import "github.com/gin-gonic/gin"

func StartServer() {
	router := gin.Default()

	UserController(router)
	router.Run()
}
