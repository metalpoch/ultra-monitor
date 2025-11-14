package mongodb

import (
	"context"
	"log"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(mongoURI string) *MongoDB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	// Test the connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("failed to ping MongoDB: %v", err)
	}

	// Extract database name from URI
	databaseName := extractDatabaseNameFromURI(mongoURI)
	if databaseName == "" {
		log.Fatal("failed to extract database name from MongoDB URI")
	}

	log.Printf("Successfully connected to MongoDB database: %s", databaseName)

	return &MongoDB{
		Client:   client,
		Database: client.Database(databaseName),
	}
}

// extractDatabaseNameFromURI extracts the database name from MongoDB URI
func extractDatabaseNameFromURI(uri string) string {
	// Parse the URI
	parsedURI, err := url.Parse(uri)
	if err != nil {
		log.Printf("Error parsing MongoDB URI: %v", err)
		return ""
	}

	// Extract database name from path
	path := strings.TrimPrefix(parsedURI.Path, "/")
	if path == "" {
		log.Printf("No database name found in MongoDB URI path: %s", uri)
		return ""
	}

	return path
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

