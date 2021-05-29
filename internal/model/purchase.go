package model

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Purchase represents a purchase model.
type Purchase struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"userID"`
	Date   time.Time          `bson:"date"`
	FileID primitive.ObjectID `bson:"fileID"`
}

// PurchaseDTO represents dto of a purchase model.
type PurchaseDTO struct {
	ID     string    `json:"id,omitempty"`
	UserID string    `json:"userID"`
	Date   time.Time `json:"date"`
	FileID string    `json:"fileID"`
}

// Entity converts PurchaseDTO to Purchase.
func (p PurchaseDTO) Entity() (*Purchase, error) {
	purchase := Purchase{
		Date: p.Date,
	}
	var err error
	if p.ID != "" {
		purchase.ID, err = primitive.ObjectIDFromHex(p.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid id : %v", err)
		}
	}
	if p.UserID != "" {
		purchase.UserID, err = primitive.ObjectIDFromHex(p.UserID)
		if err != nil {
			return nil, fmt.Errorf("invalid user id : %v", err)
		}
	}
	if p.FileID != "" {
		purchase.FileID, err = primitive.ObjectIDFromHex(p.FileID)
		if err != nil {
			return nil, fmt.Errorf("invalid file id : %v", err)
		}
	}

	return &purchase, nil
}

// Entity converts Purchase to PurchaseDTO.
func (p Purchase) DTO() *PurchaseDTO {
	purchase := PurchaseDTO{
		ID:     p.ID.Hex(),
		UserID: p.UserID.Hex(),
		Date:   time.Time{},
		FileID: p.FileID.Hex(),
	}

	return &purchase
}
