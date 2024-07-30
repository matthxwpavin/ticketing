package service

import (
	"time"

	"github.com/matthxwpavin/ticketing/database/mongo"
)

type Metadata struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func MetadataFrom(data *mongo.Metadata) *Metadata {
	return &Metadata{
		ID:        data.ID.Hex(),
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}
