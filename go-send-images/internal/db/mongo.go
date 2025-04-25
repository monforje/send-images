package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	MongoDBName string
)

// InitMongo подключает MongoDB и проверяет соединение
func InitMongo(uri, dbName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ MongoDB ping failed: %v", err)
	}

	MongoClient = client
	MongoDBName = dbName

	log.Println("✅ Connected to MongoDB")
}

// Collection возвращает коллекцию по имени
func Collection(name string) *mongo.Collection {
	return MongoClient.Database(MongoDBName).Collection(name)
}
