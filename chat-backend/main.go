package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	//"net/http/httputil"
	//"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `json:"username"`
	Password string             `json:"password"`
}

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Content   string             `json:"content"`
	Timestamp time.Time          `json:"timestamp"`
	Username  string             `json:"username"`
	Upvotes   int                `json:"upvotes"`
	Downvotes int                `json:"downvotes"`
	Votes     map[string]int     `json:"votes"`
}

var (
	store = sessions.NewCookieStore([]byte("secret")) // Change "secret" to preferred session secret
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the username already exists in the database
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Access the "users" collection
	collection := client.Database("chatapp").Collection("users")

	// Check if a user with the same username already exists
	existingUser := User{}
	err = collection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		// User with the same username already exists
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User already exists. Please login!"))
		return
	} else if err != mongo.ErrNoDocuments {
		// An error occurred while querying the database
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Store user in the database or perform any necessary operations
	log.Println("User registered:", user)

	w.Write([]byte("Redirecting to login page...."))

	// Insert the user document into the collection
	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	// Get the inserted document ID
	insertedID := result.InsertedID.(primitive.ObjectID)
	user.ID = insertedID

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Access the "users" collection
	collection := client.Database("chatapp").Collection("users")

	// Find the user by username
	filter := bson.M{"username": user.Username}
	result := collection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		// User not found
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid username or password"))
		return
	}

	// Retrieve the hashed password from the database
	var dbUser User
	err = result.Decode(&dbUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Compare the hashed password with the input password
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		// Passwords don't match
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid username or password"))
		return
	}

	// Passwords match, user is authenticated
	w.WriteHeader(http.StatusOK)
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var message Message
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&message)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := message.Username
	// Attach the username to the message
	message.Username = username

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Access the "messages" collection
	collection := client.Database("chatapp").Collection("messages")

	// Set the current timestamp
	message.Timestamp = time.Now()
	message.Votes = make(map[string]int)
	message.Upvotes = 0
	message.Downvotes = 0

	// Insert the message document into the collection
	result, err := collection.InsertOne(context.Background(), message)
	if err != nil {
		log.Fatal(err)
	}

	// Get the inserted document ID
	insertedID := result.InsertedID.(primitive.ObjectID)
	message.ID = insertedID

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func FetchHistoryMessagesHandler(w http.ResponseWriter, r *http.Request) {
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Access the "messages" collection
	collection := client.Database("chatapp").Collection("messages")

	// Find all messages in the collection
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var messages []Message
	for cursor.Next(context.Background()) {
		var message Message
		err := cursor.Decode(&message)
		if err != nil {
			log.Fatal(err)
		}

		// Fetch the username associated with the message
		userCollection := client.Database("chatapp").Collection("users")
		filter := bson.M{"username": message.Username}
		userResult := userCollection.FindOne(context.Background(), filter)
		if userResult.Err() == nil {
			var user User
			err := userResult.Decode(&user)
			if err == nil {
				// Append the username to the message content
				message.Content = fmt.Sprintf("%s : %s", user.Username, message.Content)
			}
		}

		messages = append(messages, message)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}

func UpvoteMessageHandler(w http.ResponseWriter, r *http.Request) {

	//Get token
	token := r.Header.Get("Authorization")

	// Extract the username from the token
	var tokenData struct {
		Username string `json:"username"`
	}

	err := json.NewDecoder(strings.NewReader(strings.TrimPrefix(token, "Bearer "))).Decode(&tokenData)

	if err != nil {
		// Handle error (failed to decode token or missing username field)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Use the extracted username
	username := tokenData.Username

	fmt.Println("Upvoter's username:", username)

	// Get the message ID from the URL parameter
	params := mux.Vars(r)
	//fmt.Println(params)
	messageID, err := primitive.ObjectIDFromHex(params["messageId"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Attempting to connect to mongoDB")
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Access the "messages" collection
	collection := client.Database("chatapp").Collection("messages")

	//fmt.Println("Connected")

	// Find the message by ID
	filter := bson.M{"_id": messageID}
	result := collection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		// Message not found
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("error!")
		return
	}

	// Retrieve the message document
	var message Message
	err = result.Decode(&message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	//fmt.Println("Printing message")
	//fmt.Println(result)
	//fmt.Println(err)

	// Get the userID from the message
	userID := username

	//fmt.Println(userID)

	// Check if the user has already voted for this message
	if _, ok := message.Votes[userID]; ok {
		// User has already voted, disallow voting again
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User has already voted for this message"))
		return
	}

	fmt.Println("Attempt to update message!")
	//fmt.Println("Filter:", filter)
	//fmt.Println("Message ID:",message.ID)
	// Update the message with an upvote
	update := bson.M{
		"$inc": bson.M{"upvotes": 1},
		"$set": bson.M{"votes." + userID: 1},
	}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DownvoteMessageHandler(w http.ResponseWriter, r *http.Request) {
	//Get token
	token := r.Header.Get("Authorization")

	// Extract the username from the token
	var tokenData struct {
		Username string `json:"username"`
	}

	err := json.NewDecoder(strings.NewReader(strings.TrimPrefix(token, "Bearer "))).Decode(&tokenData)

	if err != nil {
		// Handle error (failed to decode token or missing username field)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Use the extracted username
	username := tokenData.Username

	fmt.Println("Upvoter's username:", username)

	// Get the message ID from the URL parameter
	params := mux.Vars(r)
	fmt.Println(params)
	messageID, err := primitive.ObjectIDFromHex(params["messageId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Access the "messages" collection
	collection := client.Database("chatapp").Collection("messages")

	// Find the message by ID
	filter := bson.M{"_id": messageID}
	result := collection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		// Message not found
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Retrieve the message document
	var message Message
	err = result.Decode(&message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the userID from the message
	userID := username

	// Check if the user has already voted for this message
	if _, ok := message.Votes[userID]; ok {
		// User has already voted, disallow voting again
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User has already voted for this message"))
		return
	}

	// Update the message with a downvote
	update := bson.M{
		"$inc": bson.M{"downvotes": 1},
		"$set": bson.M{"votes." + userID: 1},
	}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session for the current request
	session, err := store.Get(r, "session-name")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("Deleting session:", session)

	// Clear the session values
	session.Values = nil

	// Save the session to apply the changes
	err = session.Save(r, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Logout successful"))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/signup", SignUpHandler).Methods("POST")
	router.HandleFunc("/api/login", LoginHandler).Methods("POST")
	router.HandleFunc("/api/messages", SendMessageHandler).Methods("POST")
	router.HandleFunc("/api/messages/history", FetchHistoryMessagesHandler).Methods("GET")
	router.HandleFunc("/api/messages/{messageId}/upvote", UpvoteMessageHandler).Methods("POST")
	router.HandleFunc("/api/messages/{messageId}/downvote", DownvoteMessageHandler).Methods("POST")
	router.HandleFunc("/api/logout", LogoutHandler).Methods("POST")

	// // Proxy requests to the React development server
	// reactURL, _ := url.Parse("http://localhost:4000")
	// proxy := httputil.NewSingleHostReverseProxy(reactURL)
	// router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	proxy.ServeHTTP(w, r)
	// })

	// Set up CORS middleware
	corsMiddleware := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				return
			}
			handler.ServeHTTP(w, r)
		})
	}

	// Start the front-end server on port 4000
	go func() {
		log.Fatal(http.ListenAndServe(":4000", corsMiddleware(http.FileServer(http.Dir("./build")))))
	}()

	// Start the back-end server on port 8000
	log.Println("Server listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", corsMiddleware(router)))

	//log.Fatal(http.ListenAndServe(":4000", router))
}
