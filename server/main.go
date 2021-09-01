package main

import (
	// "encoding/json"
	"context"
	"encoding/json"

	// "fmt"
	"log"
	"net/http"

	// "net/url"
	"os"

	// "time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"

	// "github.com/gorilla/websocket"

	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	// twilio "github.com/xaviiic/twilioGo"
	// jose "github.com/dvsekhvalnov/jose2go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func GinSocketIOServerWrapper(server *gosocketio.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.IsWebsocket() {
			server.ServeHTTP(ctx.Writer, ctx.Request)
		} else {
			_, _ = ctx.Writer.WriteString("===not websocket request===")
		}
	}
}

func main() {
	if os.Getenv("NODE_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file: ", err.Error())
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

	// type UserId struct {
	// 	Id string
	// }
	type Incident struct {
		Time  string
		Image string
	}
	type WatchingUpdate struct {
		Watching bool
	}
	// type WatchingMessageToCamera struct {
	// 	Id       string
	// 	Watching bool
	// }

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
	//****************************** Socket.io ********************************//
	//*************************************************************************//
	/////////////////////////////////////////////////////////////////////////////
	type Message struct {
		Text string `json:"text"`
	}

	socket := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	socket.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("!!!!!!!!!!!!! SOCKET !!!!!!!!!!!!!")
		log.Println("New client connected")
		log.Println("id: ", c.Id())
		log.Println("ip: ", c.Ip())
		c.Emit("message", Message{"you connected"})
		log.Println("message sent")
		//join them to room
		c.Join("chat")
	})

	socket.On("send", func(c *gosocketio.Channel, msg Message) string {
		log.Println("!!!!! message received !!!!!")
		//send event to all in room
		c.BroadcastTo("chat", "message", msg)
		return "OK"
	})

	socket.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Println("!!!!!!!!!!!!! SOCKET !!!!!!!!!!!!!")
		log.Println("Client disconnected")
		log.Println("id: ", c.Id())
		log.Println("ip: ", c.Ip())
		//join them to room
		c.Leave("chat")
	})

	/////////////////////////////////////////////////////////////////////////////
	//*************************************************************************//
	//************************** REST API Endpoints ***************************//
	//*************************************************************************//
	/////////////////////////////////////////////////////////////////////////////
	router := gin.Default()

	router.Use(static.Serve("/", static.LocalFile("./web", true)))
	api := router.Group("/api")
	api.GET("/")

	api.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "passed",
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
		// log.Println(string(currentUserData))
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
			// c.String(http.StatusOK, "{ \"id\":", docId, "}")
			c.JSON(http.StatusOK, gin.H{
				"id": doc.Ref.ID,
			})
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
		log.Println("Updating User Data for ", currentUser)

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

	api.PUT("/user/watching", func(c *gin.Context) {
		log.Println("Request on /api/user/watching type: PUT")

		userId := c.Query("id")
		var currently WatchingUpdate
		err := json.NewDecoder(c.Request.Body).Decode(&currently)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		log.Println("User", userId, "current watching status:", currently.Watching)
		log.Println("Updating user data")
		_, err = client.Collection("users").Doc(userId).Update(ctx, []firestore.Update{
			{
				Path:  "watching",
				Value: currently.Watching,
			},
		})
		if err != nil {
			log.Println("An error has occurred:", err.Error())
		}
	})

	api.GET("/twiliodata", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"accoundSid":   os.Getenv("TWILIO_ACCOUNT_SID"),
			"apiKeySid":    os.Getenv("TWILIO_API_KEY_SID"),
			"apiKeySecret": os.Getenv("TWILIO_API_KEY_SECRET"),
		})
	})

	api.GET("/firebase", func(c *gin.Context) {
		log.Println("RemoteIP")
		log.Println(c.ClientIP())
		c.JSON(http.StatusOK, gin.H{
			"apiKey":            os.Getenv("API_KEY"),
			"authDomain":        os.Getenv("AUTH_DOMAIN"),
			"projectId":         os.Getenv("PROJECT_ID"),
			"storageBucket":     os.Getenv("STORAGE_BUCKET"),
			"messagingSenderId": os.Getenv("MESSAGING_SENDER_ID"),
			"appId":             os.Getenv("APP_ID"),
		})
	})

	// api.GET("/socket.io/", gin.WrapH(socket.Handler))
	// socketHandle := http.NewServeMux()
	// socketHandle.Handle("/socket.io/", socket)
	// api.GET("/socket.io/", gin.WrapH(socketHandle))

	// router.GET("/socket.io", gin.WrapH(socket))
	// router.Handle("/socket.io/", gin.WrapH())
	// api.GET("/socket.io", GinSocketIOServerWrapper(socket))
	router.GET("/socket.io/", GinSocketIOServerWrapper(socket))
	router.POST("/socket.io/", GinSocketIOServerWrapper(socket))

	// The below works for serving the socket.io connection
	// serveMux := http.NewServeMux()
	// serveMux.Handle("/socket.io/", socket)
	// http.ListenAndServe(":8080", serveMux)

	router.Use(cors.Default())
	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{os.Getenv("URL")}
	// router.Use(cors.New(config))
	// log.Println("Current URL")
	// log.Println(os.Getenv("URL"))
	router.Run()

	// routerPrivate := gin.Default()
	// apiPrivate := routerPrivate.Group("/api")
	// apiPrivate.GET("/privatetest", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "private",
	// 	})
	// })
	// routerPrivate.Use(cors.Default())
	// routerPrivate.Run(":8088")
}
