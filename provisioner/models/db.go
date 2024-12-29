package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MinecraftServer struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	IP           string             `json:"ip,omitempty" bson:"ip,omitempty"`
	Owner        string             `json:"owner,omitempty" bson:"owner,omitempty"`
	Username     string             `json:"username,omitempty" bson:"username,omitempty"`
	InstanceType string             `json:"instance_type,omitempty" bson:"instance_type,omitempty"`
	Region       string             `json:"region,omitempty" bson:"region,omitempty"`
	Label        string             `json:"label,omitempty" bson:"label,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Status       string             `json:"status,omitempty" bson:"status,omitempty"`
	Metadata     map[string]string  `json:"metadata,omitempty" bson:"metadata,omitempty"`
}
