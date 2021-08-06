package main

import (
	// "encoding/json"
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	if os.Getenv("ENV") != "deployment" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	/////////////////////////////////////////////////////////////////////////////
	//*************************************************************************//
	//********************** Firestore Authentication *************************//
	//*************************************************************************//
	/////////////////////////////////////////////////////////////////////////////
	credentialsJson := []byte(os.Getenv("FIRESTORE_JSON"))
	ctx := context.Background()
	sa := option.WithCredentialsJSON(credentialsJson)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(client)

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./web", true)))
	api := router.Group("/api")
	api.GET("/")
	api.GET("/test", func(context *gin.Context) {
		message := os.Getenv("TEST")
		context.JSON(200, gin.H{
			"message": message,
		})
	})

	router.Run()
}
