package docs

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Doc struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string `bson:"title" json:"title"`
	Content   string `bson:"content" json:"content"`
	CreatedBy primitive.ObjectID `bson:"created_by" json:"created_by"  `
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}