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

// PurchaseRepo is a purchase repository.
type PurchaseRepo struct {
	collection *mongo.Collection
}

// NewPurchaseRepo is a PurchaseRepo constructor.
func NewPurchaseRepo(db *mongo.Database) *PurchaseRepo {
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

	return &PurchaseRepo{collection: c}
}

// Create creates new purchase and returns userID.
func (p PurchaseRepo) Create(ctx context.Context, purchase model.PurchaseDTO) (string, error) {
	purchaseEntity, err := purchase.Entity()
	if err != nil {
		return "", err
	}

	res, err := p.collection.InsertOne(ctx, purchaseEntity)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Delete deletes purchase and returns deleted userID.
func (p PurchaseRepo) Delete(ctx context.Context, id string) (string, error) {
	opts := options.FindOneAndDelete().SetProjection(bson.D{{"_id", 1}})
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	query := bson.M{
		"_id": objID,
	}
	var purchase model.Purchase
	err = p.collection.FindOneAndDelete(ctx, query, opts).Decode(&purchase)
	if err != nil {
		return "", err
	}

	return purchase.ID.Hex(), nil
}

// DeleteByFileID deletes purchase by purchase userID and returns purchase userID.
func (p PurchaseRepo) DeleteByFileID(ctx context.Context, id string) (string, error) {
	opts := options.FindOneAndDelete().SetProjection(bson.D{{"fileID", 1}})
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	query := bson.M{
		"fileID": objID,
	}
	var purchase model.Purchase
	err = p.collection.FindOneAndDelete(ctx, query, opts).Decode(&purchase)
	if err != nil {
		return "", err
	}

	return purchase.FileID.Hex(), nil
}

// FindByID finds purchase by userID.
func (p PurchaseRepo) FindByID(ctx context.Context, id string) (*model.PurchaseDTO, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{
		"_id": objID,
	}
	var purchase model.Purchase
	err = p.collection.FindOne(ctx, query).Decode(&purchase)
	if err != nil {
		return nil, err
	}

	return purchase.DTO(), nil
}

// FindLastByUserID finds last purchase by user userID.
func (p PurchaseRepo) FindLastByUserID(ctx context.Context, id int) (*model.PurchaseDTO, error) {
	opts := options.FindOne().SetSort(bson.M{"$natural": -1})
	query := bson.M{
		"userID": id,
	}
	var purchase model.Purchase
	err := p.collection.FindOne(ctx, query, opts).Decode(&purchase)
	if err != nil {
		return nil, err
	}

	return purchase.DTO(), nil
}

// FindAllByUserID finds purchases by user userID.
func (p PurchaseRepo) FindAllByUserID(ctx context.Context, id int) ([]model.PurchaseDTO, error) {
	query := bson.M{
		"userID": id,
	}
	var purchases model.Purchases
	cursor, err := p.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &purchases)
	if err != nil {
		return nil, err
	}

	return purchases.DTO(), nil
}

// FindByUserIDAndPeriod finds purchases by user userID and date period.
func (p PurchaseRepo) FindByUserIDAndPeriod(ctx context.Context, id int, start, end time.Time) ([]model.PurchaseDTO, error) {
	query := bson.M{
		"userID": id,
		"date":   bson.M{"$gte": start, "$lte": end},
	}
	var purchases model.Purchases
	cursor, err := p.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &purchases)
	if err != nil {
		return nil, err
	}

	return purchases.DTO(), nil
}

// FindByUserIDAfterDate finds purchases by user userID and after date.
func (p PurchaseRepo) FindByUserIDAfterDate(ctx context.Context, id int, start time.Time) ([]model.PurchaseDTO, error) {
	query := bson.M{
		"userID": id,
		"date":   bson.M{"$gte": start},
	}
	var purchases model.Purchases
	cursor, err := p.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &purchases)
	if err != nil {
		return nil, err
	}

	return purchases.DTO(), nil
}

// FindByUserIDBeforeDate finds purchases by user userID and before date.
func (p PurchaseRepo) FindByUserIDBeforeDate(ctx context.Context, id int, end time.Time) ([]model.PurchaseDTO, error) {
	query := bson.M{
		"userID": id,
		"date":   bson.M{"$lte": end},
	}
	var purchases model.Purchases
	cursor, err := p.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &purchases)
	if err != nil {
		return nil, err
	}

	return purchases.DTO(), nil
}

// FindByUserIDAndFileID finds purchases by user userID and purchase userID.
func (p PurchaseRepo) FindByUserIDAndFileID(ctx context.Context, userID int, fileID string) ([]model.PurchaseDTO, error) {
	objFileID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return nil, err
	}

	query := bson.M{
		"userID": userID,
		"fileID": objFileID,
	}
	var purchases model.Purchases
	cursor, err := p.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &purchases)
	if err != nil {
		return nil, err
	}

	return purchases.DTO(), nil
}

// FindLast finds last purchase.
func (p PurchaseRepo) FindLast(ctx context.Context) (*model.PurchaseDTO, error) {
	opts := options.FindOne().SetSort(bson.M{"$natural": -1})
	query := bson.M{}
	var purchase model.Purchase
	err := p.collection.FindOne(ctx, query, opts).Decode(&purchase)
	if err != nil {
		return nil, err
	}

	return purchase.DTO(), nil
}

// FindAll finds purchases.
func (p PurchaseRepo) FindAll(ctx context.Context) ([]model.PurchaseDTO, error) {
	query := bson.M{}
	var purchases model.Purchases
	cursor, err := p.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &purchases)
	if err != nil {
		return nil, err
	}

	return purchases.DTO(), nil
}

// FindByPeriod finds purchases by date period.
func (p PurchaseRepo) FindByPeriod(ctx context.Context, start, end time.Time) ([]model.PurchaseDTO, error) {
	query := bson.M{
		"date": bson.M{"$gte": start, "$lte": end},
	}
	var purchases model.Purchases
	cursor, err := p.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &purchases)
	if err != nil {
		return nil, err
	}

	return purchases.DTO(), nil
}

// FindAfterDate finds purchases after date.
func (p PurchaseRepo) FindAfterDate(ctx context.Context, start time.Time) ([]model.PurchaseDTO, error) {
	query := bson.M{
		"date": bson.M{"$gte": start},
	}
	var purchases model.Purchases
	cursor, err := p.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &purchases)
	if err != nil {
		return nil, err
	}

	return purchases.DTO(), nil
}

// FindBeforeDate finds purchases before date.
func (p PurchaseRepo) FindBeforeDate(ctx context.Context, end time.Time) ([]model.PurchaseDTO, error) {
	query := bson.M{
		"date": bson.M{"$lte": end},
	}
	var purchases model.Purchases
	cursor, err := p.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &purchases)
	if err != nil {
		return nil, err
	}

	return purchases.DTO(), nil
}

// FindByFileID finds purchases by purchase userID.
func (p PurchaseRepo) FindByFileID(ctx context.Context, id string) ([]model.PurchaseDTO, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{
		"fileID": objID,
	}
	var purchases model.Purchases
	cursor, err := p.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &purchases)
	if err != nil {
		return nil, err
	}

	return purchases.DTO(), nil
}
