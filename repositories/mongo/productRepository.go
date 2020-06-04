package mongo

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"

	"github.com/dhyaniarun1993/foody-catalog-service/product"
	"github.com/dhyaniarun1993/foody-catalog-service/repositories"
	"github.com/dhyaniarun1993/foody-catalog-service/repositories/mongo/dao"
	"github.com/dhyaniarun1993/foody-common/datastore/mongo"
	"github.com/dhyaniarun1993/foody-common/errors"
)

const (
	productCollection = "product"
	variantCollection = "variant"
)

type productRepository struct {
	*mongo.Client
	database string
}

// NewProductRepository creates and return product repository
func NewProductRepository(mongoClient *mongo.Client, database string) repositories.ProductRepository {
	return &productRepository{mongoClient, database}
}

func (db *productRepository) Create(ctx context.Context, product product.Product) (product.Product, errors.AppError) {

	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	productDao, daoErr := dao.GetProductDao(product)
	if daoErr != nil {
		return product, daoErr
	}

	insertProductCtx, insertProductCancel := context.WithTimeout(ctx, 1*time.Second)
	defer insertProductCancel()
	// Todo: insert product and variant in transaction
	// insert product data in datastore
	productCollection := db.Database(db.database).Collection(productCollection)
	insertProductResult, insertProductError := productCollection.InsertOne(insertProductCtx, productDao)
	if insertProductError != nil {
		return product, errors.NewAppError("Something went wrong",
			http.StatusServiceUnavailable, insertProductError)
	}

	// extract id  of the product inserted
	productObjectID, _ := insertProductResult.InsertedID.(primitive.ObjectID)
	product.ID = productObjectID.Hex()

	// create variant as list of interface to support InsertMany
	variants := make([]interface{}, len(product.Variants))
	for i := range product.Variants {
		// Add product id to product variant data
		product.Variants[i].ProductID = productObjectID.Hex()
		// convert product variant model to product variant dao
		variantDao, err := dao.GetVariantDao(product.Variants[i])
		if err != nil {
			return product, err
		}
		// append the updated dao variant object to the variant interface list
		variants[i] = variantDao
	}

	insertVariantCtx, insertVariantCancel := context.WithTimeout(ctx, 1*time.Second)
	defer insertVariantCancel()
	// insert product variant data in datastore
	variantCollection := db.Database(db.database).Collection(variantCollection)
	insertVariantResult, insertVariantError := variantCollection.InsertMany(insertVariantCtx, variants)
	if insertVariantError != nil {
		return product, errors.NewAppError("Something went wrong",
			http.StatusServiceUnavailable, insertVariantError)
	}

	// extract product variant id and update it in the product data
	for i, variantInsertID := range insertVariantResult.InsertedIDs {
		variantObjectID, _ := variantInsertID.(primitive.ObjectID)
		product.Variants[i].ID = variantObjectID.Hex()
	}

	return product, nil
}

func (db *productRepository) GetByID(ctx context.Context, productID string) (product.Product, errors.AppError) {

	var productObj product.Product
	findCtx, findCancel := context.WithTimeout(ctx, 1*time.Second)
	defer findCancel()

	productObjectID, _ := primitive.ObjectIDFromHex(productID)
	match := bson.D{
		{
			Key: "$match",
			Value: bson.D{
				{Key: "_id", Value: productObjectID},
			},
		},
	}
	lookupVariants := bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				{Key: "localField", Value: "_id"},
				{Key: "from", Value: "variant"},
				{Key: "foreignField", Value: "product_id"},
				{Key: "as", Value: "variants"},
			},
		},
	}

	collection := db.Database(db.database).Collection(productCollection)
	cursor, findError := collection.Aggregate(findCtx, mongoDriver.Pipeline{match, lookupVariants})
	if findError != nil {
		return product.Product{}, errors.NewAppError("Something went wrong",
			http.StatusInternalServerError, findError)
	}

	cursorCtx, cursorCancel := context.WithTimeout(ctx, 1*time.Second)
	defer cursorCancel()

	if cursor.Next(cursorCtx) {
		decodeError := cursor.Decode(&productObj)
		if decodeError != nil {
			return product.Product{}, errors.NewAppError("Something went wrong",
				http.StatusInternalServerError, decodeError)
		}
	}
	return productObj, nil
}

func (db *productRepository) DeleteByID(ctx context.Context, productID string) errors.AppError {

	deleteCtx, deleteCancel := context.WithTimeout(ctx, 1*time.Second)
	defer deleteCancel()

	productObjectID, _ := primitive.ObjectIDFromHex(productID)
	filter := bson.D{
		{
			Key:   "_id",
			Value: productObjectID,
		},
	}

	collection := db.Database(db.database).Collection(productCollection)
	_, deleteError := collection.DeleteOne(deleteCtx, filter)
	if deleteError != nil {
		return errors.NewAppError("Something went wrong",
			http.StatusInternalServerError, deleteError)
	}

	return nil
}

func (db *productRepository) DeleteByRestaurantID(ctx context.Context, restaurantID string) errors.AppError {

	// Todo: Perform all the operations in transaction
	findProductCtx, findProductCancel := context.WithTimeout(ctx, 1*time.Second)
	defer findProductCancel()

	restaurantObjectID, _ := primitive.ObjectIDFromHex(restaurantID)
	productFilter := bson.D{
		{
			Key:   "restaurant_id",
			Value: restaurantObjectID,
		},
	}
	productCollection := db.Database(db.database).Collection(productCollection)

	// find all the products with provided restaurant ID
	cursor, findError := productCollection.Find(findProductCtx, productFilter)
	if findError != nil {
		return errors.NewAppError("Something went wrong", http.StatusInternalServerError, findError)
	}

	// create product ID list
	productIDs := []primitive.ObjectID{}
	cursorCtx, cursorCancel := context.WithCancel(ctx)
	defer cursorCancel()
	for cursor.Next(cursorCtx) {
		var productDao dao.ProductDao
		decodeError := cursor.Decode(&productDao)
		if decodeError != nil {
			return errors.NewAppError("Something went wrong", http.StatusInternalServerError, decodeError)
		}
		productIDs = append(productIDs, productDao.ID)
	}

	deleteVariantFilter := bson.D{
		{
			Key: "product_id",
			Value: bson.D{
				{
					Key:   "$in",
					Value: productIDs,
				},
			},
		},
	}

	deleteVariantCtx, deleteVariantCancel := context.WithTimeout(ctx, 1*time.Second)
	defer deleteVariantCancel()

	// delete all the variants that belong to the above product list
	variantCollection := db.Database(db.database).Collection(variantCollection)
	_, deleteVariantError := variantCollection.DeleteMany(deleteVariantCtx, deleteVariantFilter)
	if deleteVariantError != nil {
		return errors.NewAppError("Something went wrong", http.StatusInternalServerError, deleteVariantError)
	}

	deleteProductCtx, deleteProductCancel := context.WithTimeout(ctx, 1*time.Second)
	defer deleteProductCancel()

	// delete all the products that belong to the restaurant
	_, deleteProductError := productCollection.DeleteMany(deleteProductCtx, productFilter)
	if deleteProductError != nil {
		return errors.NewAppError("Something went wrong", http.StatusInternalServerError, deleteProductError)
	}

	return nil
}

func (db *productRepository) DeleteByCategoryID(ctx context.Context, categoryID string) errors.AppError {

	// Todo: Perform all the operations in transaction
	findProductCtx, findProductCancel := context.WithTimeout(ctx, 1*time.Second)
	defer findProductCancel()

	categoryObjectID, _ := primitive.ObjectIDFromHex(categoryID)
	productFilter := bson.D{
		{
			Key:   "category_id",
			Value: categoryObjectID,
		},
	}
	productCollection := db.Database(db.database).Collection(productCollection)

	// find all the products with provided category ID
	cursor, findError := productCollection.Find(findProductCtx, productFilter)
	if findError != nil {
		return errors.NewAppError("Something went wrong", http.StatusInternalServerError, findError)
	}

	// create product ID list
	productIDs := []primitive.ObjectID{}
	cursorCtx, cursorCancel := context.WithCancel(ctx)
	defer cursorCancel()
	for cursor.Next(cursorCtx) {
		var productDao dao.ProductDao
		decodeError := cursor.Decode(&productDao)
		if decodeError != nil {
			return errors.NewAppError("Something went wrong", http.StatusInternalServerError, decodeError)
		}
		productIDs = append(productIDs, productDao.ID)
	}

	deleteVariantFilter := bson.D{
		{
			Key: "product_id",
			Value: bson.D{
				{
					Key:   "$in",
					Value: productIDs,
				},
			},
		},
	}

	deleteVariantCtx, deleteVariantCancel := context.WithTimeout(ctx, 1*time.Second)
	defer deleteVariantCancel()

	// delete all the variants that belong to the above product list
	variantCollection := db.Database(db.database).Collection(variantCollection)
	_, deleteVariantError := variantCollection.DeleteMany(deleteVariantCtx, deleteVariantFilter)
	if deleteVariantError != nil {
		return errors.NewAppError("Something went wrong", http.StatusInternalServerError, deleteVariantError)
	}

	deleteProductCtx, deleteProductCancel := context.WithTimeout(ctx, 1*time.Second)
	defer deleteProductCancel()

	// delete all the products that belong to the category provided
	_, deleteProductError := productCollection.DeleteMany(deleteProductCtx, productFilter)
	if deleteProductError != nil {
		return errors.NewAppError("Something went wrong", http.StatusInternalServerError, deleteProductError)
	}

	return nil
}
