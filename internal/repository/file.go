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

// FileRepo is a purchase repository.
type FileRepo struct {
	collection *mongo.Collection
}

// NewFileRepo is a FileRepo constructor.
func NewFileRepo(db *mongo.Database) *FileRepo {
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

	return &FileRepo{collection: c}
}

// Create creates new purchase and returns userID.
func (f FileRepo) Create(context context.Context, file model.FileDTO) (string, error) {
	fileEntity, err := file.Entity()
	if err != nil {
		return "", err
	}

	res, err := f.collection.InsertOne(context, fileEntity)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Update updates purchase and returns userID.
func (f FileRepo) Update(context context.Context, id string, file model.FileDTO) (string, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	fileEntity, err := file.Entity()
	if err != nil {
		return "", err
	}

	query := bson.M{
		"_id": objID,
	}
	update := bson.M{
		"$set": fileEntity,
	}
	var updateFile model.File
	err = f.collection.FindOneAndUpdate(context, query, update).Decode(&updateFile)
	if err != nil {
		return "", err
	}

	return updateFile.ID.Hex(), nil
}

// Delete deletes purchase and returns deleted userID.
func (f FileRepo) Delete(context context.Context, id string) (string, error) {
	opts := options.FindOneAndDelete().SetProjection(bson.D{{"_id", 1}})
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	query := bson.M{
		"_id": objID,
	}
	var delFile model.File
	err = f.collection.FindOneAndDelete(context, query, opts).Decode(&delFile)
	if err != nil {
		return "", err
	}

	return delFile.ID.Hex(), nil
}

// DeleteByAuthorID deletes purchase by authorID and returns authorID.
func (f FileRepo) DeleteByAuthorID(context context.Context, id string) (string, error) {
	opts := options.FindOneAndDelete().SetProjection(bson.D{{"authorID", 1}})
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	query := bson.M{
		"authorID": objID,
	}
	var delFile model.File
	err = f.collection.FindOneAndDelete(context, query, opts).Decode(&delFile)
	if err != nil {
		return "", err
	}

	return delFile.AuthorID.Hex(), nil
}

// FindByID finds purchase by userID.
func (f FileRepo) FindByID(context context.Context, id string) (*model.FileDTO, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{
		"_id": objID,
	}
	var file model.File
	err = f.collection.FindOne(context, query).Decode(&file)
	if err != nil {
		return nil, err
	}

	return file.DTO(), nil
}

// FindByName finds purchases by name.
func (f FileRepo) FindByName(context context.Context, name string) ([]model.FileDTO, error) {
	query := bson.M{
		"name": name,
	}
	var files model.Files
	cursor, err := f.collection.Find(context, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context, &files)
	if err != nil {
		return nil, err
	}

	return files.DTO(), nil
}

// FindAll finds purchases.
func (f FileRepo) FindAll(context context.Context) ([]model.FileDTO, error) {
	query := bson.M{}
	var files model.Files
	cursor, err := f.collection.Find(context, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context, &files)
	if err != nil {
		return nil, err
	}

	return files.DTO(), nil
}

// FindByAuthorID finds purchases by author userID.
func (f FileRepo) FindByAuthorID(context context.Context, id string) ([]model.FileDTO, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{
		"authorID": objID,
	}
	var files model.Files
	cursor, err := f.collection.Find(context, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context, &files)
	if err != nil {
		return nil, err
	}

	return files.DTO(), nil
}

// FindNotActual finds not actual purchases.
func (f FileRepo) FindNotActual(context context.Context) ([]model.FileDTO, error) {
	query := bson.M{
		"actual": false,
	}
	var files model.Files
	cursor, err := f.collection.Find(context, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context, &files)
	if err != nil {
		return nil, err
	}

	return files.DTO(), nil
}

// FindActual finds actual purchases.
func (f FileRepo) FindActual(context context.Context) ([]model.FileDTO, error) {
	query := bson.M{
		"actual": true,
	}
	var files model.Files
	cursor, err := f.collection.Find(context, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context, &files)
	if err != nil {
		return nil, err
	}

	return files.DTO(), nil
}

// FindAddedByPeriod finds added purchases by date period.
func (f FileRepo) FindAddedByPeriod(context context.Context, start, end time.Time) ([]model.FileDTO, error) {
	query := bson.M{
		"addDate": bson.M{"$gte": start, "$lte": end},
	}
	var files model.Files
	cursor, err := f.collection.Find(context, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context, &files)
	if err != nil {
		return nil, err
	}

	return files.DTO(), nil
}

// FindUpdatedByPeriod finds updated purchases by date period.
func (f FileRepo) FindUpdatedByPeriod(context context.Context, start, end time.Time) ([]model.FileDTO, error) {
	query := bson.M{
		"updateDate": bson.M{"$gte": start, "$lte": end},
	}
	var files model.Files
	cursor, err := f.collection.Find(context, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context, &files)
	if err != nil {
		return nil, err
	}

	return files.DTO(), nil
}
