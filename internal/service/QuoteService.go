package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

// Define the Quote struct
type Quote struct {
	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OriginLocation      string             `json:"originLocation,omitempty" bson:"originLocation,omitempty"`
	DestinationLocation string             `json:"destinationLocation,omitempty" bson:"destinationLocation,omitempty"`
	CargoType           string             `json:"cargoType,omitempty" bson:"cargoType,omitempty"`
	Quantity            int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
	SpecialRequirements string             `json:"specialRequirements,omitempty" bson:"specialRequirements,omitempty"`
	Status              string             `json:"status,omitempty" bson:"status,omitempty"`
	Timestamp           time.Time          `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

func initDB() error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	db = client.Database("spot_freight_quotes")
	fmt.Println("Connected to MongoDB")

	return nil
}

type QuoteService struct {
	Client *mongo.Client
	Db     *mongo.Database
}

func NewQuoteService() *QuoteService {
	return &QuoteService{
		Client: client,
		Db:     db,
	}
}

func (qs *QuoteService) SubmitSpotFreightQuoteRequest(w http.ResponseWriter, r *http.Request) {
	var quote Quote
	json.NewDecoder(r.Body).Decode(&quote)

	quote.Status = "submitted"
	quote.Timestamp = time.Now()

	result, err := qs.Db.Collection("quotes").InsertOne(context.TODO(), quote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (qs *QuoteService) GetQuoteDetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	quoteID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid quote ID", http.StatusBadRequest)
		return
	}

	var result Quote
	err = db.Collection("quotes").FindOne(context.TODO(), bson.M{"_id": quoteID}).Decode(&result)
	if err != nil {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (qs *QuoteService) UpdateQuoteDetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	quoteID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid quote ID", http.StatusBadRequest)
		return
	}

	var quote Quote
	json.NewDecoder(r.Body).Decode(&quote)

	update := bson.M{
		"$set": bson.M{
			"originLocation":      quote.OriginLocation,
			"destinationLocation": quote.DestinationLocation,
			"cargoType":           quote.CargoType,
			"quantity":            quote.Quantity,
			"specialRequirements": quote.SpecialRequirements,
			"status":              "updated",
		},
	}

	_, err = db.Collection("quotes").UpdateOne(context.TODO(), bson.M{"_id": quoteID}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Quote updated successfully")
}

func (qs *QuoteService) GetQuotes(w http.ResponseWriter, r *http.Request) {
	var quotes []Quote
	cursor, err := qs.Db.Collection("quotes").Find(context.Background(), bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = cursor.All(context.Background(), &quotes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(quotes)
}
