package controller

import "github.com/gin-gonic/gin"

func UserController(router *gin.Engine) {
	router.GET("/ping", ping)
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pingpong broooo",
	})
}
