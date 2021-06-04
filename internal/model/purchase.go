package model

import (
	"time"

	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/database/mongo"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Purchases represents a slice of a purchase model.
type Purchases []Purchase

// PurchasesDTO represents a slice of a dto purchase model.
type PurchasesDTO []PurchaseDTO

// Purchase represents a purchase model.
type Purchase mongo.Purchase

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
			return nil, errors.Wrap(err, "invalid id")
		}
	}
	if p.UserID != "" {
		purchase.UserID, err = primitive.ObjectIDFromHex(p.UserID)
		if err != nil {
			return nil, errors.Wrap(err, "invalid user id")
		}
	}
	if p.FileID != "" {
		purchase.FileID, err = primitive.ObjectIDFromHex(p.FileID)
		if err != nil {
			return nil, errors.Wrap(err, "invalid file id")
		}
	}

	return &purchase, nil
}

// Entity converts Purchase to PurchaseDTO.
func (p Purchase) DTO() *PurchaseDTO {
	purchase := PurchaseDTO{
		ID:     p.ID.Hex(),
		UserID: p.UserID.Hex(),
		Date:   p.Date,
		FileID: p.FileID.Hex(),
	}

	return &purchase
}

// Entity converts PurchasesDTO to Purchases.
func (p PurchasesDTO) Entity() (Purchases, error) {
	var purchases Purchases
	for _, purchase := range p {
		entityComment, err := purchase.Entity()
		if err != nil {
			return nil, errors.Wrap(err, "couldn't convert purchase dto to purchase")
		}
		purchases = append(purchases, *entityComment)
	}
	return purchases, nil
}

// DTO converts Purchases to PurchasesDTO
func (p Purchases) DTO() PurchasesDTO {
	var purchases PurchasesDTO
	for _, purchase := range p {
		purchases = append(purchases, *purchase.DTO())
	}
	return purchases
}
