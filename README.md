# security-cam
This is a remake of an app I wrote during my time as a student in Code Chrysalis.  The original app was separated into two repositories, one for the Golang server and one for the React client.  It is now hosted in a monorepo and the entire app has been dockerized to make for easier deployment and to have a consistent testing environment.  You can see a live version by clicking the link below.

[Security Cam on Heroku](https://security-cam-go.herokuapp.com/)

## Technologies Used
### Golang
Golang was used to write the server backend.  It was used for creating the API endpoints as well as handling communication with the database via Firestore and user notification using SendGrid

### Firebase
Firebase was used for user login and authentication as well as for storing incident images with timestamps and user data via Firestore

### SendGrid
SendGrid is used to notify the user via email that the security camera app detected motion and has recorded an image of the incident.

### React
React is used for creating the frontend for this app.

### Docker
Docker is used as a container for the backend and frontend allowing for a consistent environment for both development and deployment

### CircleCI
CircleCI is used for CI/CD.  While Heroku can be used to automatically build and deploy the container directly from GitHub, the project is large enough that it runs out of memory during the build.  Therefore CircleCI is used to automatically build the container and then deploy it in it's built state to Heroku.

## Future Features
### Socket.io
This will be used for a future feature where the user can log in and monitor the security camera remotely.  It will tell the security camera to stop detecting motion and transmit the video feed.  It will also be used to let the security camera know when the user has logged out to then resume monitoring.

### Twilio
This will be used to transmit the live video feed when the user wants to watch a live video feed.  It's activation and deactivation will be handled by communication via Socket.io.

### Bootstrap
Bootstrap will be used to improve the look of the front end app.

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

### /api/user/watching?id=<"user id as string">
Sets whether the user is currently trying to watch the live camera feed.  Can be set with the following.
\{
  "watching": boolean representing whether the user is watching or not
\}

### /api/firebase
Used for interacting with firebase

### /api/twiliodata
Used for interacting with twilio
