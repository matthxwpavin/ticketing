package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Metadata struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}

func (s *Metadata) SetUpdatedAt(t time.Time) {
	s.UpdatedAt = t
}

func NewMetadata() *Metadata {
	return &Metadata{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
