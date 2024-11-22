package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// Initialize the Firebase app
var app *firebase.App

func init() {
	// Replace with the path to your service account key JSON file
	opt := option.WithCredentialsFile("C:/Users/janku/Downloads/myappauth-a6302-firebase-adminsdk-1s913-ccdfb1d37b.json")
	var err error
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}
}

// UserRequest represents the request body for user sign-in/registration
type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Handler for user sign-in or registration
func signInOrRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var reqBody UserRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		http.Error(w, "Failed to initialize Firebase Auth client", http.StatusInternalServerError)
		return
	}

	// Check if the user is already registered
	userRecord, err := client.GetUserByEmail(context.Background(), reqBody.Email)
	if err != nil {
		// User is not registered, so register them
		params := (&auth.UserToCreate{}).
			Email(reqBody.Email).
			Password(reqBody.Password)
		userRecord, err = client.CreateUser(context.Background(), params)
		if err != nil {
			http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"message": "User registered successfully",
			"uid":     userRecord.UID,
			"email":   userRecord.Email,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// User is already registered
	// Note: Firebase Admin SDK doesn't allow password verification on the server
	response := map[string]string{
		"message": "User already registered, signed in successfully",
		"uid":     userRecord.UID,
		"email":   userRecord.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/sign-in-or-register", signInOrRegisterHandler)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
