package mongo

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/soa-team-11/auth-service/utils"
)

var (
	db *mongo.Database
)

func init() {
	uri := utils.Getenv("MONGO_URI", "mongodb://localhost:27017")
	dbName := utils.Getenv("MONGO_DB", "auth-db")

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

	db = client.Database(dbName)
}

func GetDatabase() *mongo.Database {
	return db
}
