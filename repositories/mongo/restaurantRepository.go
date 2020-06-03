package mongo

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/dhyaniarun1993/foody-catalog-service/repositories"
	"github.com/dhyaniarun1993/foody-catalog-service/restaurant"
	"github.com/dhyaniarun1993/foody-common/datastore/mongo"
	"github.com/dhyaniarun1993/foody-common/errors"
)

const (
	restaurantCollection = "restaurant"
)

type restaurantRepository struct {
	*mongo.Client
	database string
}

// NewRestaurantRepository creates and return restaurant repository
func NewRestaurantRepository(mongoClient *mongo.Client, database string) repositories.RestaurantRepository {
	return &restaurantRepository{mongoClient, database}
}

func (db *restaurantRepository) Create(ctx context.Context,
	restaurant restaurant.Restaurant) (restaurant.Restaurant, errors.AppError) {

	restaurant.CreatedAt = time.Now()
	restaurant.UpdatedAt = time.Now()
	restaurant.Address.Location.Type = "Point"
	insertCtx, insertCancel := context.WithTimeout(ctx, 1*time.Second)
	defer insertCancel()

	collection := db.Database(db.database).Collection(restaurantCollection)

	insertResult, insertError := collection.InsertOne(insertCtx, restaurant)
	if insertError != nil {
		err := errors.NewAppError("Unable to insert restaurant to DB",
			http.StatusServiceUnavailable, insertError)
		return restaurant, err
	}

	id, _ := insertResult.InsertedID.(primitive.ObjectID)
	restaurant.ID = id.Hex()
	return restaurant, nil
}

func (db *restaurantRepository) GetByID(ctx context.Context, restaurantID string) (restaurant.Restaurant, errors.AppError) {

	var restaurantObj restaurant.Restaurant
	findCtx, findCancel := context.WithTimeout(ctx, 1*time.Second)
	defer findCancel()

	objectID, _ := primitive.ObjectIDFromHex(restaurantID)
	filter := bson.D{
		{
			Key:   "_id",
			Value: objectID,
		},
	}

	collection := db.Database(db.database).Collection(restaurantCollection)

	cursor, findError := collection.Find(findCtx, filter)
	if findError != nil {
		return restaurant.Restaurant{}, errors.NewAppError("Unable to get data from DB",
			http.StatusInternalServerError, findError)
	}

	cursorCtx, cursorCancel := context.WithCancel(ctx)
	defer cursorCancel()

	if cursor.Next(cursorCtx) {
		decodeError := cursor.Decode(&restaurantObj)
		if decodeError != nil {
			return restaurant.Restaurant{}, errors.NewAppError("Unable to decode categories", http.StatusInternalServerError, decodeError)
		}
	}
	return restaurantObj, nil
}

func (db *restaurantRepository) DeleteByID(ctx context.Context,
	restaurantID string) errors.AppError {

	deleteCtx, deleteCancel := context.WithTimeout(ctx, 1*time.Second)
	defer deleteCancel()

	objectID, _ := primitive.ObjectIDFromHex(restaurantID)
	filter := bson.D{
		{
			Key:   "_id",
			Value: objectID,
		},
	}

	collection := db.Database(db.database).Collection(restaurantCollection)

	_, deleteErr := collection.DeleteOne(deleteCtx, filter)
	if deleteErr != nil {
		return errors.NewAppError("Unable to delete from DB", http.StatusInternalServerError, deleteErr)
	}
	return nil
}

func (db *restaurantRepository) GetAllRestaurants(ctx context.Context, query restaurant.GetAllRestaurantsRequest,
	maxDistance int64) ([]restaurant.Restaurant, errors.AppError) {

	restaurants := []restaurant.Restaurant{}
	offset := (query.PageNumber - 1) * query.PageSize
	findCtx, findCancel := context.WithTimeout(ctx, 1*time.Second)
	defer findCancel()

	findOptions := &mongoOptions.FindOptions{
		Skip:  &offset,
		Limit: &query.PageSize,
	}

	filter := bson.D{
		{
			Key: "address.location",
			Value: bson.D{
				{
					Key: "$geoWithin",
					Value: bson.D{
						{
							Key: "$centerSphere",
							Value: bson.A{
								bson.A{query.Longitude, query.Latitude},
								maxDistance},
						},
					},
				},
			},
		},
	}

	collection := db.Database(db.database).Collection(restaurantCollection)

	cursor, findError := collection.Find(findCtx, filter, findOptions)
	if findError != nil {
		return restaurants, errors.NewAppError("Unable to get data from DB", http.StatusInternalServerError, findError)
	}

	cursorCtx, cursorCancel := context.WithCancel(ctx)
	defer cursorCancel()
	for cursor.Next(cursorCtx) {
		var restaurantObj restaurant.Restaurant
		decodeError := cursor.Decode(&restaurantObj)
		if decodeError != nil {
			return restaurants, errors.NewAppError("Unable to decode categories", http.StatusInternalServerError, decodeError)
		}
		restaurants = append(restaurants, restaurantObj)
	}
	return restaurants, nil
}

func (db *restaurantRepository) GetAllRestaurantsTotalCount(ctx context.Context,
	query restaurant.GetAllRestaurantsRequest, maxDistance int64) (int64, errors.AppError) {

	countCtx, countCancel := context.WithTimeout(ctx, 1*time.Second)
	defer countCancel()

	filter := bson.D{
		{
			Key: "address.location",
			Value: bson.D{
				{
					Key: "$geoWithin",
					Value: bson.D{
						{
							Key: "$centerSphere",
							Value: bson.A{
								bson.A{query.Longitude, query.Latitude},
								maxDistance},
						},
					},
				},
			},
		},
	}

	collection := db.Database(db.database).Collection(restaurantCollection)

	totalCount, findError := collection.CountDocuments(countCtx, filter)
	if findError != nil {
		return totalCount, errors.NewAppError("Unable to get data from DB", http.StatusInternalServerError, findError)
	}

	return totalCount, nil
}
