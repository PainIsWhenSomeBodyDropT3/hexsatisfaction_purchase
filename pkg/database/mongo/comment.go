package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comment represents a comment model.
type Comment struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"userID"`
	PurchaseID primitive.ObjectID `bson:"purchaseID"`
	Date       time.Time          `bson:"date"`
	Text       string             `bson:"text"`
}
