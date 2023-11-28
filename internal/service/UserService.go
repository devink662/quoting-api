package service

import (
	"context"
	"encoding/json"
	"freight-quote-api/internal/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// Authentication Service:
//
// Depends on:
//
//	HTTP Client (for API requests)
//
// Interactions:
//
//			Sends user registration and login requests to the API.
//	   	Manages JWT tokens for authentication.
//
// User Service:
// User Registration:
//
// Endpoint: POST /api/users/register
// Description: Registers a new user with provided name, email, and password.
// User Login:
//
// Endpoint: POST /api/users/login
// Description: Authenticates a user with provided email and password, returning an authentication token upon successful login.
// Get User Profile:
//
// Endpoint: GET /api/users/profile
// Description: Retrieves the profile information of the authenticated user.
// Update User Profile:
//
// Endpoint: PUT /api/users/profile
// Description: Allows the user to update their profile information (name, email, etc.).
// main.go
var client *mongo.Client
var db *mongo.Database

type UserService struct {
	Client *mongo.Client
	Db     *mongo.Database
}

func newUserService() *UserService {
	return &UserService{
		Client: client,
		Db:     db,
	}
}
func (us *UserService) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	json.NewDecoder(r.Body).Decode(&user)

	result, err := us.Db.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(result)
	w.Write([]byte("User registered successfully"))
}

func (us *UserService) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	json.NewDecoder(r.Body).Decode(&user)

	var result model.User
	err := db.Collection("users").FindOne(context.TODO(), bson.M{"email": user.Email, "password": user.Password}).Decode(&result)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (us *UserService) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var result model.User
	err = db.Collection("users").FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&result)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (us *UserService) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user model.User
	json.NewDecoder(r.Body).Decode(&user)

	update := bson.M{
		"$set": bson.M{
			"name":  user.Name,
			"email": user.Email,
		},
	}

	_, err = db.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": userID}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("User updated successfully")
}
