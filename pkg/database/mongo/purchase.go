package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Purchase represents a purchase model.
type Purchase struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID int                `bson:"userID"`
	Date   time.Time          `bson:"date"`
	FileID primitive.ObjectID `bson:"fileID"`
}
