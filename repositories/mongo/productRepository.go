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
	productCollection = "product"
)

type productRepository struct {
	*mongo.Client
	database string
}

// NewProductRepository creates and return product repository
func NewProductRepository(mongoClient *mongo.Client, database string) repositories.ProductRepository {
	return &productRepository{mongoClient, database}
}

func (db *productRepository) Create(ctx context.Context,
	product models.Product) (models.Product, errors.AppError) {

	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	insertCtx, insertCancel := context.WithTimeout(ctx, 1*time.Second)
	defer insertCancel()

	collection := db.Database(db.database).Collection(productCollection)

	insertResult, insertError := collection.InsertOne(insertCtx, product)
	if insertError != nil {
		return product, errors.NewAppError("Unable to insert restaurant to DB",
			http.StatusServiceUnavailable, insertError)
	}

	id, _ := insertResult.InsertedID.(primitive.ObjectID)
	product.ID = id
	return product, nil
}

func (db *productRepository) Get(ctx context.Context, productID string,
	restaurantID string) (models.Product, errors.AppError) {

	var product models.Product
	findCtx, findCancel := context.WithTimeout(ctx, 1*time.Second)
	defer findCancel()

	productObjectID, _ := primitive.ObjectIDFromHex(productID)
	restaurantObjectID, _ := primitive.ObjectIDFromHex(restaurantID)
	filter := bson.D{
		{
			Key:   "_id",
			Value: productObjectID,
		},
		{
			Key:   "restaurant_id",
			Value: restaurantObjectID,
		},
	}

	collection := db.Database(db.database).Collection(productCollection)
	cursor, findError := collection.Find(findCtx, filter)
	if findError != nil {
		return models.Product{}, errors.NewAppError("Unable to get data from DB",
			http.StatusInternalServerError, findError)
	}

	cursorCtx, cursorCancel := context.WithTimeout(ctx, 1*time.Second)
	defer cursorCancel()

	if cursor.Next(cursorCtx) {
		decodeError := cursor.Decode(&product)
		if decodeError != nil {
			return models.Product{}, errors.NewAppError("Unable to decode categories",
				http.StatusInternalServerError, decodeError)
		}
	}
	return product, nil
}

func (db *productRepository) Delete(ctx context.Context,
	productID string, restaurantID string) errors.AppError {

	deleteCtx, deleteCancel := context.WithTimeout(ctx, 1*time.Second)
	defer deleteCancel()

	productObjectID, _ := primitive.ObjectIDFromHex(productID)
	restaurantObjectID, _ := primitive.ObjectIDFromHex(restaurantID)
	filter := bson.D{
		{
			Key:   "_id",
			Value: productObjectID,
		},
		{
			Key:   "restaurant_id",
			Value: restaurantObjectID,
		},
	}

	collection := db.Database(db.database).Collection(productCollection)
	_, deleteError := collection.DeleteOne(deleteCtx, filter)
	if deleteError != nil {
		return errors.NewAppError("Unable to delete data from DB",
			http.StatusInternalServerError, deleteError)
	}

	return nil
}

func (db *productRepository) GetProductsByRestaurantID(ctx context.Context, restaurantID string,
	query dto.GetAllProductsRequestQuery) ([]models.Product, errors.AppError) {

	products := []models.Product{}
	offset := (query.PageNumber - 1) * query.PageSize

	findCtx, findCancel := context.WithTimeout(ctx, 1*time.Second)
	defer findCancel()

	findOptions := &mongoOptions.FindOptions{
		Skip:  &offset,
		Limit: &query.PageSize,
	}

	restaurantObjectID, _ := primitive.ObjectIDFromHex(restaurantID)
	filter := bson.D{
		{
			Key:   "restaurant_id",
			Value: restaurantObjectID,
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

	collection := db.Database(db.database).Collection(productCollection)

	cursor, findError := collection.Find(findCtx, filter, findOptions)
	if findError != nil {
		return products, errors.NewAppError("Unable to get data from DB", http.StatusInternalServerError, findError)
	}

	cursorCtx, cursorCancel := context.WithCancel(ctx)
	defer cursorCancel()
	for cursor.Next(cursorCtx) {
		var product models.Product
		decodeError := cursor.Decode(&product)
		if decodeError != nil {
			return products, errors.NewAppError("Unable to decode categories", http.StatusInternalServerError, decodeError)
		}
		products = append(products, product)
	}
	return products, nil
}

func (db productRepository) GetProductsByRestaurantTotalCount(ctx context.Context, restaurantID string,
	query dto.GetAllProductsRequestQuery) (int64, errors.AppError) {

	countCtx, countCancel := context.WithTimeout(ctx, 1*time.Second)
	defer countCancel()

	restaurantObjectID, _ := primitive.ObjectIDFromHex(restaurantID)
	filter := bson.D{
		{
			Key:   "restaurant_id",
			Value: restaurantObjectID,
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

	collection := db.Database(db.database).Collection(productCollection)

	totalCount, countError := collection.CountDocuments(countCtx, filter)
	if countError != nil {
		return totalCount, errors.NewAppError("Unable to get data from DB",
			http.StatusInternalServerError, countError)
	}

	return totalCount, nil
}
