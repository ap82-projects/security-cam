# SecurityCam
This app can be used to repurpose devices with webcams to be used as security cameras.  All that is necessary other than the camera is a browser and connection to the internet.  You can see a live version by clicking the link below.

[SecurityCam on Heroku](https://security-cam-go.herokuapp.com/)

## Running SecurityCam Locally
Clone the repository locally and create a .env file with the following environment variables defined
### API_KEY, AUTH_DOMAIN, PROJECT_ID, STORAGE_BUCKET, MESSAGING_SENDER_ID, APP_ID
These are associated with your Firebase account and can be found in the bottom of
"Project Settings" -> "General"
of your associated Firebase project

### FIRESTORE_JSON
This is the JSON data created by going to 
"Project Settings" -> "Service Accounts" -> "Firebase Admin SDK" -> "Generate new private key"
of your associated Firebase project

### SENDGRID_API_KEY
This key can be created by going to
"Settings" -> "API Keys" -> "Create API Key"
from your SendGrid dashboard

### AUTH_EMAIL_ADDR
The email address that the notifications are to be sent from.  It must be verified in SendGrid by going to
"Settings" -> "Sender Authentication" -> "Verify a Single Sender"
from your SendGrid dashboard

### AUTH_EMAIL_NAME
The name to be associated with the email address above.  Best to use the name registered with SendGrid

## Running server and client independently
### Starting the server
1. Copy the .env file to the root of the server directory
2. From a new terminal instance, navigate into the server directory and run "go run main.go"

### Starting the client
1. From a new terminal instance, navigate into the client directory and run "npm install"
2. To start the client in a development environment with hot-reloading, run "npm start" and access the client from localhost:3000
3. To build the client and have it served by the server, run "npm run build-local" and access the client from localhost:8080

## Running server and client within a Docker container
1. Copy the .env file to the root of the cloned repository
2. From a new terminal instance, navigate to the root of the cloned repository
3. Make sure your docker daemon is running and execute "docker build -t security-cam ."
4. If the container was built with no errors, then execute "docker run --env-file .env -p 3000:8080 -d security-cam" and access the client from localhost:3000

## API Endpoints

The following are the API endpoints used on the server

### GET /api/test
For testing connection to server.  Returns the following in the body if successful
\{
  "message": "passed"
\}

### POST /api/user
Creates a new user.  The following format is passed to the body
\{
  "email": email address as string
  "googleid: google id as string
  "incidents": \[\]
  "name": name as string
  "phone": phone number as string
  "watching": false
\}

It returns the newly created user's document id in the following format
\{
   "id": document id as string
\}

### GET /api/user?id=<"user document id as string">
Returns specified user's data in the following format
\{
  "email": email address as string
  "googleid": google id as string
  "incidents": \[\{
                 "time": time of incident as string
                 "image": image taken as string
               \}\],
  "name": name as string
  "phone": phone number as string
  "watching": status of whether the user is watching as a boolean
\}

### DELETE /api/user?id=<"user id as string">
Deletes specified user from database

### GET /api/user/google?id=<"google id as string">
If a user with the specified google id exists, document id
is returned in the following format
\{
  "id": document id as string
\}

### PUT /api/user/incident?id=<"user id as string">
Used to add incidents to the user's data as well as send notification email via SendGrid.  Takes and object in the following format
\{
  "time": time of incident as string
  "image": image taken as string
\}

### DELETE /api/user/incident?id=<"user id as string">&time=<"time of incident as string">
Deletes specified incident from users data

### PUT /api/user/watching?id=<"user id as string">
Sets whether the user is currently trying to watch the live camera feed.  Can be set with the following.
\{
  "watching": boolean representing whether the user is watching or not
\}

### GET /api/firebase
Used for interacting with Firebase

### GET /api/twiliodata
Used for interacting with Twilio

## Technologies Used
### Golang
Golang was used to write the server backend.  It was used for creating the API endpoints as well as handling communication with the database via Firestore and user notification using SendGrid

### Firebase
Firebase was used for user login and authentication as well as for storing incident images with timestamps and user data via Firestore

### SendGrid
SendGrid is used to notify the user via email that the security camera app detected motion and has recorded an image of the incident

### React
React is used for creating the frontend for this app

### Docker
Docker is used as a container for the backend and frontend allowing for a consistent environment for both development and deployment

### CircleCI
CircleCI is used for CI/CD.  While Heroku can be used to automatically build and deploy the container directly from GitHub, the project is large enough that it runs out of memory during the build.  Therefore CircleCI is used to automatically build the container and then deploy it in it's built state to Heroku.

### Heroku
Heroku is used for hosting the app on the web.  Originally it was also used for CI/CD but the build process took up too much memory for the free tier.

## Future Features
### [Realtime communication to trigger live video feed](https://github.com/ap82-projects/security-cam/issues/1)

### [Implement live video feed](https://github.com/ap82-projects/security-cam/issues/2)

### [Improve UI](https://github.com/ap82-projects/security-cam/issues/3)

