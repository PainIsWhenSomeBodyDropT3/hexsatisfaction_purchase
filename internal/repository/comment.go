package repository

import (
	"context"
	"time"

	"github.com/JesusG2000/hexsatisfaction_purchase/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// CommentRepo is a purchase repository.
type CommentRepo struct {
	collection *mongo.Collection
}

// NewCommentRepo is a CommentRepo constructor.
func NewCommentRepo(db *mongo.Database) *CommentRepo {
	c := db.Collection("purchase")
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{
				Key:   "_id",
				Value: bsonx.Int64(1),
			}},
			Options: options.Index().SetName("userID"),
		},
	}
	_, err := c.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return nil
	}

	return &CommentRepo{collection: c}
}

// Create creates purchase and returns userID.
func (c CommentRepo) Create(context context.Context, comment model.CommentDTO) (string, error) {
	commentEntity, err := comment.Entity()
	if err != nil {
		return "", err
	}

	res, err := c.collection.InsertOne(context, commentEntity)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Update updates purchase and returns userID.
func (c CommentRepo) Update(context context.Context, id string, comment model.CommentDTO) (string, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	commentEntity, err := comment.Entity()
	if err != nil {
		return "", err
	}

	query := bson.M{
		"_id": objID,
	}
	update := bson.M{
		"$set": commentEntity,
	}
	var updateComment model.Comment
	err = c.collection.FindOneAndUpdate(context, query, update).Decode(&updateComment)
	if err != nil {
		return "", err
	}

	return updateComment.ID.Hex(), nil
}

// Delete deletes purchase and returns userID.
func (c CommentRepo) Delete(context context.Context, id string) (string, error) {
	opts := options.FindOneAndDelete().SetProjection(bson.D{{"_id", 1}})
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	query := bson.M{
		"_id": objID,
	}
	var delComment model.Comment
	err = c.collection.FindOneAndDelete(context, query, opts).Decode(&delComment)
	if err != nil {
		return "", err
	}

	return delComment.ID.Hex(), nil
}

// DeleteByPurchaseID deletes purchase by purchase userID and returns purchase userID.
func (c CommentRepo) DeleteByPurchaseID(context context.Context, id string) (string, error) {
	opts := options.FindOneAndDelete().SetProjection(bson.D{{"purchaseID", 1}})
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	query := bson.M{
		"purchaseID": objID,
	}
	var delComment model.Comment
	err = c.collection.FindOneAndDelete(context, query, opts).Decode(&delComment)
	if err != nil {
		return "", err
	}

	return delComment.PurchaseID.Hex(), nil
}

// FindByID finds purchase by userID.
func (c CommentRepo) FindByID(context context.Context, id string) (*model.CommentDTO, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{
		"_id": objID,
	}
	var comment model.Comment
	err = c.collection.FindOne(context, query).Decode(&comment)
	if err != nil {
		return nil, err
	}

	return comment.DTO(), nil
}

// FindAllByUserID finds purchases by user userID.
func (c CommentRepo) FindAllByUserID(context context.Context, id int) ([]model.CommentDTO, error) {
	query := bson.M{
		"userID": id,
	}
	var comments model.Comments
	cursor, err := c.collection.Find(context, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context, &comments)
	if err != nil {
		return nil, err
	}

	return comments.DTO(), nil
}

// FindByPurchaseID finds purchases by purchase userID.
func (c CommentRepo) FindByPurchaseID(context context.Context, id string) ([]model.CommentDTO, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{
		"purchaseID": objID,
	}
	var comments model.Comments
	cursor, err := c.collection.Find(context, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context, &comments)
	if err != nil {
		return nil, err
	}

	return comments.DTO(), nil
}

// FindByUserIDAndPurchaseID finds purchases by purchase and user userID.
func (c CommentRepo) FindByUserIDAndPurchaseID(context context.Context, userID int, purchaseID string) ([]model.CommentDTO, error) {
	objPurchaseID, err := primitive.ObjectIDFromHex(purchaseID)
	if err != nil {
		return nil, err
	}

	query := bson.M{
		"userID":     userID,
		"purchaseID": objPurchaseID,
	}
	var comments model.Comments
	cursor, err := c.collection.Find(context, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context, &comments)
	if err != nil {
		return nil, err
	}

	return comments.DTO(), nil
}

// FindAll finds purchases.
func (c CommentRepo) FindAll(context context.Context) ([]model.CommentDTO, error) {
	query := bson.M{}
	var comments model.Comments
	cursor, err := c.collection.Find(context, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context, &comments)
	if err != nil {
		return nil, err
	}

	return comments.DTO(), nil
}

// FindByText finds purchases by text.
func (c CommentRepo) FindByText(context context.Context, text string) ([]model.CommentDTO, error) {
	query := bson.M{
		"text": bson.M{"$regex": text},
	}
	var comments model.Comments
	cursor, err := c.collection.Find(context, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context, &comments)
	if err != nil {
		return nil, err
	}

	return comments.DTO(), nil
}

// FindByPeriod finds purchases by date period.
func (c CommentRepo) FindByPeriod(context context.Context, start, end time.Time) ([]model.CommentDTO, error) {
	query := bson.M{
		"date": bson.M{"$gte": start, "$lte": end},
	}
	var comments model.Comments
	cursor, err := c.collection.Find(context, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context, &comments)
	if err != nil {
		return nil, err
	}

	return comments.DTO(), nil
}
