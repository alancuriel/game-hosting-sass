package db

import (
    "context"
    "fmt"
    "os"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    m "github.com/alancuriel/game-hosting-sass/provisioner/models"
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

func (s *ProvisionerDB) SaveServer(server *m.MinecraftServer) error {
    collection := s.db.Collection("minecraft_servers")
    _, err := collection.InsertOne(context.Background(), server)
    return err
}
