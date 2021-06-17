package model

import (
	"time"

	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/database/mongo"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comments represents a slice of a comment model.
type Comments []Comment

// CommentsDTO represents a slice of a dto comment model.
type CommentsDTO []CommentDTO

// Comment represents a comment model.
type Comment mongo.Comment

// CommentDTO represents dto of a comment model.
type CommentDTO struct {
	ID         string    `json:"id,omitempty"`
	UserID     int       `json:"userID"`
	PurchaseID string    `json:"purchaseID"`
	Date       time.Time `json:"date"`
	Text       string    `json:"text"`
}

// Entity converts CommentDTO to Comment.
func (c CommentDTO) Entity() (*Comment, error) {
	comment := Comment{
		UserID: c.UserID,
		Date:   c.Date,
		Text:   c.Text,
	}
	var err error
	if c.ID != "" {
		comment.ID, err = primitive.ObjectIDFromHex(c.ID)
		if err != nil {
			return nil, errors.Wrap(err, "invalid id")
		}
	}
	if c.PurchaseID != "" {
		comment.PurchaseID, err = primitive.ObjectIDFromHex(c.PurchaseID)
		if err != nil {
			return nil, errors.Wrap(err, "invalid purchase id")
		}
	}

	return &comment, nil
}

// DTO converts Comment to CommentDTO.
func (c Comment) DTO() *CommentDTO {
	comment := CommentDTO{
		ID:         c.ID.Hex(),
		UserID:     c.UserID,
		PurchaseID: c.PurchaseID.Hex(),
		Date:       c.Date,
		Text:       c.Text,
	}

	return &comment
}

// Entity converts CommentsDTO to Comments.
func (c CommentsDTO) Entity() (Comments, error) {
	var comments Comments
	for _, comment := range c {
		entityComment, err := comment.Entity()
		if err != nil {
			return nil, errors.Wrap(err, "couldn't convert comment dto to comment")
		}
		comments = append(comments, *entityComment)
	}
	return comments, nil
}

// DTO converts Comments to CommentsDTO
func (c Comments) DTO() CommentsDTO {
	var comments CommentsDTO
	for _, comment := range c {
		comments = append(comments, *comment.DTO())
	}
	return comments
}
