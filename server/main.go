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

	router.Run()
}
