package main

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("../client", true)))
	api := router.Group("/api")
	api.GET("/")
	api.GET("/test", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "passed",
		})
	})

	router.Run()
}
