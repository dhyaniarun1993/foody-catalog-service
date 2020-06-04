package mongo

import (
	"context"
	"net/http"
	"time"

	"github.com/dhyaniarun1993/foody-catalog-service/category"
	"github.com/dhyaniarun1993/foody-catalog-service/repositories"
	"github.com/dhyaniarun1993/foody-catalog-service/repositories/mongo/dao"
	"github.com/dhyaniarun1993/foody-common/datastore/mongo"
	"github.com/dhyaniarun1993/foody-common/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	categoryCollection = "category"
)

type categoryRepository struct {
	*mongo.Client
	database string
}

// NewCategoryRepository creates and return category repository
func NewCategoryRepository(mongoClient *mongo.Client, database string) repositories.CategoryRepository {
	return &categoryRepository{mongoClient, database}
}

func (db *categoryRepository) Create(ctx context.Context,
	category category.Category) (category.Category, errors.AppError) {

	category.ID = ""
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	categoryDao, daoErr := dao.GetCategoryDao(category)
	if daoErr != nil {
		return category, daoErr
	}

	insertCtx, insertCancel := context.WithTimeout(ctx, 1*time.Second)
	defer insertCancel()

	collection := db.Database(db.database).Collection(categoryCollection)

	insertResult, insertError := collection.InsertOne(insertCtx, categoryDao)
	if insertError != nil {
		err := errors.NewAppError("Something went wrong",
			http.StatusServiceUnavailable, insertError)
		return category, err
	}

	id, _ := insertResult.InsertedID.(primitive.ObjectID)
	category.ID = id.Hex()
	return category, nil
}

func (db *categoryRepository) GetByID(ctx context.Context, categoryID string) (category.Category, errors.AppError) {

	var categoryObj category.Category
	findCtx, findCancel := context.WithTimeout(ctx, 1*time.Second)
	defer findCancel()

	objectID, _ := primitive.ObjectIDFromHex(categoryID)
	filter := bson.D{
		{
			Key:   "_id",
			Value: objectID,
		},
	}

	collection := db.Database(db.database).Collection(categoryCollection)

	cursor, findError := collection.Find(findCtx, filter)
	if findError != nil {
		return category.Category{}, errors.NewAppError("Something went wrong",
			http.StatusInternalServerError, findError)
	}

	cursorCtx, cursorCancel := context.WithCancel(ctx)
	defer cursorCancel()

	if cursor.Next(cursorCtx) {
		decodeError := cursor.Decode(&categoryObj)
		if decodeError != nil {
			return category.Category{}, errors.NewAppError("Something went wrong",
				http.StatusInternalServerError, decodeError)
		}
	}
	return categoryObj, nil
}

func (db *categoryRepository) DeleteByID(ctx context.Context, categoryID string) errors.AppError {

	deleteCtx, deleteCancel := context.WithTimeout(ctx, 1*time.Second)
	defer deleteCancel()

	objectID, _ := primitive.ObjectIDFromHex(categoryID)
	filter := bson.D{
		{
			Key:   "_id",
			Value: objectID,
		},
	}

	collection := db.Database(db.database).Collection(categoryCollection)

	_, deleteErr := collection.DeleteOne(deleteCtx, filter)
	if deleteErr != nil {
		return errors.NewAppError("Something went wrong", http.StatusInternalServerError, deleteErr)
	}
	return nil
}

func (db *categoryRepository) DeleteByRestaurantID(ctx context.Context, restaurantID string) errors.AppError {

	deleteCtx, deleteCancel := context.WithTimeout(ctx, 1*time.Second)
	defer deleteCancel()

	objectID, _ := primitive.ObjectIDFromHex(restaurantID)
	filter := bson.D{
		{
			Key:   "restaurant_id",
			Value: objectID,
		},
	}

	collection := db.Database(db.database).Collection(categoryCollection)

	_, deleteErr := collection.DeleteMany(deleteCtx, filter)
	if deleteErr != nil {
		return errors.NewAppError("Something went wrong", http.StatusInternalServerError, deleteErr)
	}
	return nil
}
