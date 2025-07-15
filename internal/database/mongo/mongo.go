package mongo_connection

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/soa-team-11/auth-service/config"
)

var (
	db *mongo.Database
)

func init() {
	cfg := config.LoadConfig()
	uri := cfg.MongoURI

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	log.Printf("Connected to MongoDB on %s", uri)

	db = client.Database(cfg.MongoDB)
}

func GetDatabase() *mongo.Database {
	return db
}
