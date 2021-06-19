package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// File represents a file model.
type File struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Size        int                `bson:"size"`
	Path        string             `bson:"path"`
	AddDate     time.Time          `bson:"addDate"`
	UpdateDate  time.Time          `bson:"updateDate"`
	Actual      bool               `bson:"actual"`
	AuthorID    int                `bson:"authorID"`
}
