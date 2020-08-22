package main

import (
	"context"
	"fmt"
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func database()	*mongo.Client  {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// koneksi ke mongoDB
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err.Error())
	}

	// cek koneksi
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("CONNECTED TO MONGODB")

	return client
}