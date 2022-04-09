package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Configuration model
type Configuration struct {
	Port             string
	ConnectionString string
}

// GetConfiguration method basically populate configuration information from .env and return Configuration model
func GetConfiguration() Configuration {
	err := godotenv.Load(".env")
	// dotenv := goDotEnvVariable("STRONGEST_AVENGER")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	configuration := Configuration{
		os.Getenv("PORT"),
		os.Getenv("CONNECTION_STRING"),
	}

	return configuration
}

// ConnectDB : This is helper function to connect mongoDB
func ConnectDB() *mongo.Collection {
	//way1 ------------------------------------------------------------------------------------------------------
	// config := GetConfiguration()
	// // Set client options
	// clientOptions := options.Client().ApplyURI(config.ConnectionString)

	// // Connect to MongoDB
	// client, err := mongo.Connect(context.TODO(), clientOptions)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Connected to MongoDB!", config.Port)
	// collection := client.Database("thepolyglotdeveloper").Collection("books")
	// return collection

	//way2---------------------------------------------------------------------------------------------------------------------------
	config := GetConfiguration()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.ConnectionString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	fmt.Println("Connected to MongoDB!", config.Port)
	collection := client.Database("thepolyglotdeveloper").Collection("books")

	return collection

}

// ErrorResponse : This is error model.
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// GetError : This is helper function to prepare error model.
// If you want to export your function. You must to start upper case function name. Otherwise you won't see your function when you import that on other class.
func GetError(err error, w http.ResponseWriter) {

	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}

//you can use package otherwise you have to load .env in every function
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
