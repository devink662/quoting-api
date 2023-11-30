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
//type Quote struct {
//	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
//	OriginLocation      string             `json:"originLocation,omitempty" bson:"originLocation,omitempty"`
//	DestinationLocation string             `json:"destinationLocation,omitempty" bson:"destinationLocation,omitempty"`
//	CargoType           string             `json:"cargoType,omitempty" bson:"cargoType,omitempty"`
//	Quantity            int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
//	SpecialRequirements string             `json:"specialRequirements,omitempty" bson:"specialRequirements,omitempty"`
//	Status              string             `json:"status,omitempty" bson:"status,omitempty"`
//	Timestamp           time.Time          `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
//}

type Quote struct {
	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Quantity            int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
	SpecialRequirements string             `json:"specialRequirements,omitempty" bson:"specialRequirements,omitempty"`
	Status              string             `json:"status,omitempty" bson:"status,omitempty"`
	Timestamp           time.Time          `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	Origin              string             `json:"origin" bson:"origin"`
	Destination         string             `json:"destination" bson:"destination"`
	CargoType           string             `json:"cargoType" bson:"cargoType"`
	Weight              float64            `json:"weight" bson:"weight"`
	Dimensions          []float64          `json:"dimensions" bson:"dimensions"`
	Units               int                `json:"units" bson:"units"`
	Packaging           string             `json:"packaging" bson:"packaging"`
	Hazardous           bool               `json:"hazardous" bson:"hazardous"`
	Mode                string             `json:"mode" bson:"mode"`
	TransitTime         string             `json:"transitTime" bson:"transitTime"`
	SpecialHandling     string             `json:"specialHandling" bson:"specialHandling"`
	Temperature         string             `json:"temperature" bson:"temperature"`
	CustomsInfo         string             `json:"customsInfo" bson:"customsInfo"`
	PickupDate          string             `json:"pickupDate" bson:"pickupDate"`
	DeliveryDate        string             `json:"deliveryDate" bson:"deliveryDate"`
	Accessorials        []string           `json:"accessorials" bson:"accessorials"`
	Insurance           bool               `json:"insurance" bson:"insurance"`
	InsuranceAmount     float64            `json:"insuranceAmount" bson:"insuranceAmount"`
	Incoterms           string             `json:"incoterms" bson:"incoterms"`
	PaymentTerms        string             `json:"paymentTerms" bson:"paymentTerms"`
	Carrier             string             `json:"carrier" bson:"carrier"`
	ShipperContact      string             `json:"shipperContact" bson:"shipperContact"`
	ConsigneeContact    string             `json:"consigneeContact" bson:"consigneeContact"`
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
			"originLocation":      quote.Origin,
			"destinationLocation": quote.Destination,
			"cargoType":           quote.CargoType,
			"quantity":            quote.Quantity,
			"specialRequirements": quote.SpecialRequirements,
			"weight":              quote.Weight,
			"dimensions":          quote.Dimensions,
			"units":               quote.Units,
			"packaging":           quote.Packaging,
			"hazardous":           quote.Hazardous,
			"mode":                quote.Mode,
			"transitTime":         quote.TransitTime,
			"specialHandling":     quote.SpecialHandling,
			"temperature":         quote.Temperature,
			"customsInfo":         quote.CustomsInfo,
			"pickupDate":          quote.PickupDate,
			"deliveryDate":        quote.DeliveryDate,
			"accessorials":        quote.Accessorials,
			"insurance":           quote.Insurance,
			"insuranceAmount":     quote.InsuranceAmount,
			"incoterms":           quote.Incoterms,
			"paymentTerms":        quote.PaymentTerms,
			"carrier":             quote.Carrier,
			"shipperContact":      quote.ShipperContact,
			"consigneeContact":    quote.ConsigneeContact,
			"timestamp":           quote.Timestamp,
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
