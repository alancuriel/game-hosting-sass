package db

import (
	"context"
	"fmt"
	"os"
	"time"

	m "github.com/alancuriel/game-hosting-sass/provisioner/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProvisionerDB struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewProvisioner() (*ProvisionerDB, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, fmt.Errorf("MONGODB_URI not found")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	db := client.Database("provisioner_db")

	return &ProvisionerDB{
		client: client,
		db:     db,
	}, nil
}

func (s *ProvisionerDB) SaveServer(server *m.MinecraftServer) (interface{}, error) {
	collection := s.db.Collection("minecraft_servers")
	result, err := collection.InsertOne(context.Background(), server)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (s *ProvisionerDB) UpdateServerStatus(id interface{}, status string) error {
	collection := s.db.Collection("minecraft_servers")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"status":      status,
		"lastUpdated": time.Now(),
	}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProvisionerDB) ListMcServerByOwner(owner string) ([]*m.MinecraftServer, error) {
	collection := s.db.Collection("minecraft_servers")
	filter := bson.M{"owner": owner}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}


	var servers []*m.MinecraftServer
	if err = cursor.All(context.Background(), &servers); err != nil {
		return nil, err
	}

	return servers, nil
}
