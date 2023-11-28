package main

import (
	"context"
	"fmt"
	"freight-quote-api/internal/service"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var db *mongo.Database

func initDB() error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database to verify if connection is successful
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	db = client.Database("spot_freight_quotes")
	fmt.Println("Connected to MongoDB")

	return nil
}

func main() {
	err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	us := service.UserService{
		Client: client,
		Db:     db,
	}

	qs := service.QuoteService{
		Client: client,
		Db:     db,
	}

	// Add CORS handling middleware
	r.Use(corsMiddleware)

	// Define routes
	r.HandleFunc("/api/users/register", us.RegisterUser).Methods("OPTIONS")
	r.HandleFunc("/api/users/login", us.LoginUser).Methods("OPTIONS")
	r.HandleFunc("/api/users/register", us.RegisterUser).Methods("POST")
	r.HandleFunc("/api/users/login", us.LoginUser).Methods("POST")
	r.HandleFunc("/api/users/profile", us.GetUserProfile).Methods("GET")
	r.HandleFunc("/api/users/profile", us.UpdateUserProfile).Methods("PUT")

	r.HandleFunc("/api/quotes/submit", qs.SubmitSpotFreightQuoteRequest).Methods("POST")
	r.HandleFunc("/api/quotes/{id}", qs.GetQuoteDetails).Methods("GET")
	r.HandleFunc("/api/quotes/{id}", qs.UpdateQuoteDetails).Methods("PUT")
	r.HandleFunc("/api/quotes", qs.GetQuotes).Methods("GET")

	// Add other routes for quotes, carriers, negotiations, reports, etc.

	// Start the server
	port := ":3000"
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
