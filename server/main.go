package main

import (
	// "encoding/json"
	"context"
	// "fmt"
	"encoding/json"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
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

	/////////////////////////////////////////////////////////////////////////////
	//*************************************************************************//
	//******************** Structs for interacting with DB ********************//
	//*************************************************************************//
	/////////////////////////////////////////////////////////////////////////////

	type UserId struct {
		Id string
	}
	type Incident struct {
		Time  string
		Image string
	}
	type WatchingUpdate struct {
		Watching bool
	}
	type WatchingMessageToCamera struct {
		Id       string
		Watching bool
	}
	type User struct {
		Name      string
		GoogleId  string
		Email     string
		Phone     string
		Incidents []Incident
		Watching  bool
	}

	/////////////////////////////////////////////////////////////////////////////
	//*************************************************************************//
	//************************** REST API Endpoints ***************************//
	//*************************************************************************//
	/////////////////////////////////////////////////////////////////////////////

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

	api.POST("/user", func(c *gin.Context) {
		log.Println("Request on /api/user type: POST")

		var newUser User
		err := json.NewDecoder(c.Request.Body).Decode(&newUser)
		if err != nil {
			// http.Error(w, err.Error(), http.StatusBadRequest)
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		log.Println("Adding user:")
		log.Println(newUser)

		fsDocRef, fsWriteResult, err := client.Collection("users").Add(ctx, map[string]interface{}{
			"name":      newUser.Name,
			"googleid":  newUser.GoogleId,
			"email":     newUser.Email,
			"phone":     newUser.Phone,
			"incidents": newUser.Incidents,
			"watching":  newUser.Watching,
		})
		log.Println("New user id", fsDocRef.ID, "created at", fsWriteResult)
		// fmt.Fprintln(w, "{ \"id\":", fsDocRef.ID, "}")
		c.String(http.StatusOK, "{ \"id\":", fsDocRef.ID, "}")
		if err != nil {
			log.Fatalf("Failed adding user: %v", err.Error())
		}
	})

	api.GET("/user/:id", func(c *gin.Context) {
		log.Println("Request on /api/user type: GET")
		userDocId := c.Param("id")
		log.Println("Retrieving user with Document ID", userDocId)
		query, errQ := client.Collection("users").Doc(userDocId).Get(ctx)
		if errQ != nil {
			log.Fatal("Error retrieving user data")
		}
		var currentUser User
		mapstructure.Decode(query.Data(), &currentUser)
		currentUserData, err := json.Marshal(currentUser)
		if err != nil {
			log.Println("Error:", err.Error())
		}
		log.Println("Sending user data")
		log.Println(string(currentUserData))
		c.String(http.StatusOK, string(currentUserData))
	})

	api.DELETE("/user/:id", func(c *gin.Context) {
		log.Println("Request on /api/user type: DELETE")
		userId := c.Param("id")

		log.Println("Deleting user ID", userId)
		fsDeleteTime, err := client.Collection("users").Doc(userId).Delete(ctx)
		if err != nil {
			log.Println("An error has occurred:", err.Error())
			c.String(http.StatusBadRequest, "An error has occurred:", err.Error())
		} else {
			log.Println("User", userId, "deleted at", fsDeleteTime)
			c.String(http.StatusOK, "User", userId, "deleted at", fsDeleteTime)
		}
	})

	api.GET("/user/google/:id", func(c *gin.Context) {
		log.Println("Request type: GET")
		userGoogleId := c.Param("id")
		log.Println("Retrieving user with Google ID", userGoogleId)
		query := client.Collection("users").Where("googleid", "==", userGoogleId).Documents(ctx)
		for {
			doc, err := query.Next()
			if err == iterator.Done {
				break
			}

			id, err := json.Marshal(doc.Ref.ID)
			if err != nil {
				log.Println("Error:", err)
			}
			docId := string(id)
			log.Println("User found. Sending response.")
			log.Println("Document id", docId)
			// fmt.Fprintln(w, "{ \"id\":", docId, "}")
			c.String(http.StatusOK, "{ \"id\":", docId, "}")
		}
	})

	router.Run()
}
