package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MinecraftServer struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	IP           string             `bson:"ip,omitempty"`
	Owner        string             `bson:"owner,omitempty"`
	Username     string             `bson:"username,omitempty"`
	InstanceType string             `bson:"instance_type,omitempty"`
	Region       string             `bson:"region,omitempty"`
	Label        string             `bson:"label,omitempty"`
	CreatedAt    time.Time          `bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `bson:"updated_at,omitempty"`
	Status       string             `bson:"status,omitempty"`
	Metadata     map[string]string  `bson:"metadata,omitempty"`
}
