package db

import (
	"context"
	"fmt"
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Collection {
	
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://0.0.0.0:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	fmt.Println("Connected to MongoDB!")

	//collection := client.Database("samples").Collection("employees")
	collection := client.Database("samples").Collection("employees_docker")

	return collection
}



