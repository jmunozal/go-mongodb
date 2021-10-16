package main

import (

	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

)

var (
	client 				*mongo.Client
	ctx					context.Context 		= context.Background()
	examplesDatabase	*mongo.Database
	moviesCollection 	*mongo.Collection

	WarningLogger 		*log.Logger
	InfoLogger    		*log.Logger
	ErrorLogger   		*log.Logger
)

func main() {

	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	mongoDbHost := os.Getenv("mongodb.host")
	mongoDBPort := os.Getenv("mongodb.port")

	var err error

	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://" + mongoDbHost + ":" + mongoDBPort + ""))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		ErrorLogger.Println(err)
	}

	examplesDatabase = client.Database("examples")
	moviesCollection = examplesDatabase.Collection("movies")

	defer client.Disconnect(ctx)
	handleRequests()

}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/movies/year/{year}", getMoviesByYear)
	myRouter.HandleFunc("/movies", getAllMovies)

	InfoLogger.Println("Handling requests in port 8080")
	http.ListenAndServe(":8080", myRouter)

}

func getMoviesByYear(responseWriter http.ResponseWriter, request *http.Request) {


	InfoLogger.Println("getMoviesByYear method called")

	// parse path parameters
	vars := mux.Vars(request)
	year, _ := strconv.Atoi(vars["year"])

	moviesCursor, err := moviesCollection.Find(ctx, bson.M{"year": year})
	if err != nil {
		ErrorLogger.Println(err)
	}

	executeCursor(moviesCursor, responseWriter, request)

}

func getAllMovies(responseWriter http.ResponseWriter, request *http.Request) {

	InfoLogger.Println("getAllMovies method called")
	moviesCursor, err := moviesCollection.Find(ctx, bson.M{})
	if err != nil {
		ErrorLogger.Println(err)
	}

	executeCursor(moviesCursor, responseWriter, request)

}

func executeCursor(moviesCursor *mongo.Cursor, responseWriter http.ResponseWriter, request *http.Request) {

	var movies []bson.M
	err := moviesCursor.All(ctx, &movies)
	if err != nil {
		ErrorLogger.Println(err.Error())
	}

	jsonResp, err := json.Marshal(movies)
	if err != nil {
		ErrorLogger.Println("Error marshalling data")
	}

	responseWriter.Header().Add("Content-Type", "application-json")
	responseWriter.Write(jsonResp)

}

