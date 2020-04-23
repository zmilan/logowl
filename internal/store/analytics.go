package store

import (
	"context"

	"github.com/jz222/loggy/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type interfaceAnalytics interface {
	InsertOne(models.Analytics) (primitive.ObjectID, error)
	DeleteMany(bson.M) (int64, error)
	FindOne(bson.M) (models.Analytics, error)
	FindOneAndUpdate(bson.M, bson.M) (models.Analytics, error)
}

type analytics struct {
	db *mongo.Database
}

func (a *analytics) InsertOne(analyticsDocument models.Analytics) (primitive.ObjectID, error) {
	collection := a.db.Collection(CollectionAnalytics)
	res, err := collection.InsertOne(context.TODO(), analyticsDocument)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (a *analytics) DeleteMany(filter bson.M) (int64, error) {
	collection := a.db.Collection(CollectionAnalytics)

	res, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return 0, nil
	}

	return res.DeletedCount, nil
}

func (a *analytics) FindOne(filter bson.M) (models.Analytics, error) {
	var analyticsDocument models.Analytics

	collection := a.db.Collection(CollectionAnalytics)

	queryResult := collection.FindOne(context.TODO(), filter)
	if queryResult.Err() != nil {
		return models.Analytics{}, queryResult.Err()
	}

	err := queryResult.Decode(&analyticsDocument)
	if err != nil {
		return models.Analytics{}, err
	}

	return analyticsDocument, nil
}

func (a *analytics) FindOneAndUpdate(filter, update bson.M) (models.Analytics, error) {
	collection := a.db.Collection(CollectionAnalytics)

	res := collection.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		options.MergeFindOneAndUpdateOptions().SetUpsert(true),
		options.MergeFindOneAndUpdateOptions().SetReturnDocument(options.After),
	)
	if res.Err() != nil {
		return models.Analytics{}, res.Err()
	}

	var analyticsDocument models.Analytics

	err := res.Decode(&analyticsDocument)
	if err != nil {
		return models.Analytics{}, err
	}

	return analyticsDocument, nil
}
