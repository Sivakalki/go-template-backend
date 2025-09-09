package mongodb

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(ctx context.Context) (*mongo.Database, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		fmt.Println("empty uri so this is calling ")
		uri = "mongodb://admin:admin@mongo:27017"
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		fmt.Println("caling the database")
		dbName = "myappdb"
	}

	opts := options.Client().ApplyURI(uri)

	var client *mongo.Client
	var err error
	for i := 0; i < 10; i++ { // try 10 times
		client, err = mongo.Connect(opts)
		if err == nil {
			err = client.Ping(ctx, nil)
		}

		if err == nil {
			db := client.Database(dbName)
			return db, nil
		}

		fmt.Printf("MongoDB not ready yet (%d/10): %v\n", i+1, err)
		time.Sleep(2 * time.Second) 
	}

	return nil, fmt.Errorf("failed to connect to MongoDB after retries: %w", err)
}
