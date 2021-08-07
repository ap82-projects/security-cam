package main

import (
	// "encoding/json"
	"context"
	// "fmt"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	if os.Getenv("ENV") != "deployment" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file: ", err.Error())
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
		log.Fatalln(err.Error())
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err.Error())
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
			log.Println("Failed adding user: ", err.Error())
		}
	})

	api.GET("/user", func(c *gin.Context) {
		log.Println("Request on /api/user type: GET")
		userDocId := c.Query("id")
		log.Println("Retrieving user with Document ID", userDocId)
		query, errQ := client.Collection("users").Doc(userDocId).Get(ctx)
		if errQ != nil {
			log.Println("Error retrieving user data: ", err.Error())
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

	api.DELETE("/user", func(c *gin.Context) {
		log.Println("Request on /api/user type: DELETE")
		userId := c.Query("id")

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

	api.GET("/user/google", func(c *gin.Context) {
		log.Println("Request on /api/user/google type: GET")
		userGoogleId := c.Query("id")
		log.Println("Retrieving user with Google ID", userGoogleId)
		query := client.Collection("users").Where("googleid", "==", userGoogleId).Documents(ctx)
		for {
			doc, err := query.Next()
			if err == iterator.Done {
				break
			}

			id, err := json.Marshal(doc.Ref.ID)
			if err != nil {
				log.Println("Error:", err.Error())
			}
			docId := string(id)
			log.Println("User found. Sending response.")
			log.Println("Document id", docId)
			// fmt.Fprintln(w, "{ \"id\":", docId, "}")
			c.String(http.StatusOK, "{ \"id\":", docId, "}")
		}
	})

	api.PUT("/user/incident", func(c *gin.Context) {
		log.Println("Request on user/incident type: PUT")

		var newIncident Incident
		err := json.NewDecoder(c.Request.Body).Decode(&newIncident)
		if err != nil {
			// http.Error(w, err.Error(), http.StatusBadRequest)
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		log.Println("Received New Incident:")
		// log.Println(newIncident)

		log.Println("Retrieving User Data")
		userId := c.Query("id")
		query, errQ := client.Collection("users").Doc(userId).Get(ctx)
		if errQ != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		var currentUser User
		mapstructure.Decode(query.Data(), &currentUser)
		// log.Println(currentUser)
		log.Println("Updating User Data")
		currentUser.Incidents = append(currentUser.Incidents, newIncident)

		// log.Println(currentUser)

		_, err = client.Collection("users").Doc(userId).Update(ctx, []firestore.Update{
			{
				Path:  "incidents",
				Value: currentUser.Incidents,
			},
		})
		if err != nil {
			log.Println("An error has occurred:", err.Error())
		}

		// Send notification to user
		from := mail.NewEmail(os.Getenv("AUTH_EMAIL_NAME"), os.Getenv("AUTH_EMAIL_ADDR"))
		subject := "Notification from Security Cam"
		to := mail.NewEmail(currentUser.Name, currentUser.Email)
		plainTextContent := "Movement has been detected.  Please log in to check status."
		// htmlContent := "<img src=" + newIncident.Image + "alt=\"img\" />"
		htmlContent := "<strong>Movement has been detected.  Please log in to check status.</strong>"
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
		response, err := client.Send(message)
		if err != nil {
			log.Println("THIS IS AN ERROR")
			log.Println(err.Error())
		} else {
			log.Println("SUCCESS")
			log.Println(response.StatusCode)
			log.Println(response.Body)
			log.Println(response.Headers)
		}
		c.String(http.StatusOK, "{ \"time\":", newIncident.Time, "}")
	})

	api.DELETE("/user/incident", func(c *gin.Context) {
		log.Println("Request on user/incident type: DELETE")
		userId := c.Query("id")
		incidentTime := c.Query("time")

		log.Println("Retrieving data for user ID", userId)
		query, errQ := client.Collection("users").Doc(userId).Get(ctx)
		if errQ != nil {

		}

		var currentUser User
		mapstructure.Decode(query.Data(), &currentUser)
		// log.Println(currentUser)
		log.Println("Updating User Data")

		i := 0
		for _, incident := range currentUser.Incidents {
			if incident.Time != incidentTime {
				currentUser.Incidents[i] = incident
				i++
			}
		}
		currentUser.Incidents = currentUser.Incidents[:i]

		_, err = client.Collection("users").Doc(userId).Update(ctx, []firestore.Update{
			{
				Path:  "incidents",
				Value: currentUser.Incidents,
			},
		})
		if err != nil {
			log.Println("An error has occurred:", err.Error())
		}
	})

	router.Run()
}
