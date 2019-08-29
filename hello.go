package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbURI = os.Getenv("MONGO_URI")
var port = os.Getenv("PORT")
var client *mongo.Client

// Product struct
type Product struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Brand *Brand             `json:"brand" bson:"brand,omitempty"`
	Desc  string             `json:"desc,omitempty" bson:"desc,omitempty"`
	Price int                `json:"price,omitempty" bson:"price,omitempty"`
	Ts    int64              `json:"ts,omitempty" bson:"ts,omitempty"`
	Port  string             `json:"port,omitempty" bson:"port,omitempty"`
}

// Brand Struct (Model)
type Brand struct {
	Name string `json:"name" bson:"name,omitempty"`
	Cod  int    `json:"cod" bson:"cod,omitempty"`
}

// Response Struct (Model)
type Response struct {
	Date   int      `json:"name" bson:"name,omitempty"`
	Result *Product `json:"result" bson:"result,omitempty"`
}

func getItem(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var product Product
	query := mux.Vars(req)
	product = Product{Desc: query["key"], Price: 1500, Ts: time.Now().UnixNano(), Port: port, Brand: &Brand{Name: "Acetaminofen", Cod: 2000}}
	// _ = json.NewDecoder(req.Body).Decode(&product)
	collection := client.Database("minimax").Collection("sampleData")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, product)
	userSent := Response{
		Date:   123,
		Result: &product,
	}
	fmt.Println(result)
	json.NewEncoder(res).Encode(userSent)
	defer cancel()
}

func createItem(res http.ResponseWriter, req *http.Request) {
	// res.Header().Set("Content-Type", "application/json")
	// var product Product
	// _ = json.NewDecoder(req.Body).Decode(&product)
	// collection := client.Database(`testGo`).Collection(`example`)
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// result, _ := collection.InsertOne(ctx, product)
	// json.NewEncoder(res).Encode(result)
}
func main() {
	if port == "" {
		port = "5000"
	}
	fmt.Println("Starting server on PORT", port)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	defer cancel()
	if client != nil {
		fmt.Println("Client Done")
		r := mux.NewRouter()
		// Route handler / Endpoints
		r.HandleFunc("/products/{key}", getItem).Methods("GET")
		r.HandleFunc("/products/{key}", createItem).Methods("POST")
		log.Fatal(http.ListenAndServe(":"+port, r))
	}
}
