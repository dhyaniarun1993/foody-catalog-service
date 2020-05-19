package mongo

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/dhyaniarun1993/foody-catalog-service/repositories"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/dto"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/models"
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
	restaurant models.Restaurant) (models.Restaurant, errors.AppError) {

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
		return models.Restaurant{}, err
	}

	id, _ := insertResult.InsertedID.(primitive.ObjectID)
	restaurant.ID = id
	return restaurant, nil
}

func (db *restaurantRepository) Get(ctx context.Context,
	restaurantID string) (models.Restaurant, errors.AppError) {

	var restaurant models.Restaurant
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
		return restaurant, errors.NewAppError("Unable to get data from DB",
			http.StatusInternalServerError, findError)
	}

	cursorCtx, cursorCancel := context.WithCancel(ctx)
	defer cursorCancel()

	if cursor.Next(cursorCtx) {
		decodeError := cursor.Decode(&restaurant)
		if decodeError != nil {
			return models.Restaurant{}, errors.NewAppError("Unable to decode categories", http.StatusInternalServerError, decodeError)
		}
	}
	return restaurant, nil
}

func (db *restaurantRepository) Delete(ctx context.Context,
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

func (db *restaurantRepository) GetAllRestaurants(ctx context.Context,
	query dto.GetAllRestaurantsRequestQuery, maxDistance int64) ([]models.Restaurant, errors.AppError) {

	restaurants := []models.Restaurant{}
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

	idFilterValue := bson.A{}
	for _, id := range query.ID {
		objectID, _ := primitive.ObjectIDFromHex(id)
		idFilterValue = append(idFilterValue, objectID)
	}

	if len(idFilterValue) > 0 {
		idFilter := bson.E{
			Key: "_id",
			Value: bson.D{
				{
					Key:   "$in",
					Value: idFilterValue,
				},
			},
		}
		filter = append(filter, idFilter)
	}

	if query.MerchantID != "" {
		merchantFilter := bson.E{
			Key:   "merchant_id",
			Value: query.MerchantID,
		}
		filter = append(filter, merchantFilter)
	}

	collection := db.Database(db.database).Collection(restaurantCollection)

	cursor, findError := collection.Find(findCtx, filter, findOptions)
	if findError != nil {
		return restaurants, errors.NewAppError("Unable to get data from DB", http.StatusInternalServerError, findError)
	}

	cursorCtx, cursorCancel := context.WithCancel(ctx)
	defer cursorCancel()
	for cursor.Next(cursorCtx) {
		var restaurant models.Restaurant
		decodeError := cursor.Decode(&restaurant)
		if decodeError != nil {
			return restaurants, errors.NewAppError("Unable to decode categories", http.StatusInternalServerError, decodeError)
		}
		restaurants = append(restaurants, restaurant)
	}
	return restaurants, nil
}

func (db *restaurantRepository) GetAllRestaurantsTotalCount(ctx context.Context,
	query dto.GetAllRestaurantsRequestQuery, maxDistance int64) (int64, errors.AppError) {

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

	idFilterValue := bson.A{}
	for _, id := range query.ID {
		objectID, _ := primitive.ObjectIDFromHex(id)
		idFilterValue = append(idFilterValue, objectID)
	}

	if len(idFilterValue) > 0 {
		idFilter := bson.E{
			Key: "_id",
			Value: bson.D{
				{
					Key:   "$in",
					Value: idFilterValue,
				},
			},
		}
		filter = append(filter, idFilter)
	}

	collection := db.Database(db.database).Collection(restaurantCollection)

	totalCount, findError := collection.CountDocuments(countCtx, filter)
	if findError != nil {
		return totalCount, errors.NewAppError("Unable to get data from DB", http.StatusInternalServerError, findError)
	}

	return totalCount, nil
}
