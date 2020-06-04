package dao

import (
	"net/http"
	"time"

	"github.com/dhyaniarun1993/foody-catalog-service/product"
	"github.com/dhyaniarun1993/foody-common/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PriceDao provides the schema definition for Price
type PriceDao struct {
	Amount   float64 `bson:"amount" json:"amount"`
	Currency string  `bson:"currency" json:"currency"`
}

// VariantDao provides the schema definition for Products Variant data to be stored in mongodb
type VariantDao struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID   primitive.ObjectID `bson:"product_id" json:"product_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Price       PriceDao           `bson:"price" json:"price"`
	InStock     *bool              `bson:"in_stock"  json:"in_stock"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// GetVariantDao converts and returns product variant Dao object from product variant schema
func GetVariantDao(variant product.Variant) (VariantDao, errors.AppError) {
	variantDao := VariantDao{
		Name:        variant.Name,
		Description: variant.Description,
		Price: PriceDao{
			Amount:   variant.Price.Amount,
			Currency: variant.Price.Currency,
		},
		InStock:   variant.InStock,
		CreatedAt: variant.CreatedAt,
		UpdatedAt: variant.UpdatedAt,
	}

	if variant.ID != "" {
		variantObjectID, err := primitive.ObjectIDFromHex(variant.ID)
		if err != nil {
			return VariantDao{}, errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
		}
		variantDao.ID = variantObjectID
	}

	// product id is required
	productObjectID, err := primitive.ObjectIDFromHex(variant.ProductID)
	if err != nil {
		return VariantDao{}, errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
	}
	variantDao.ProductID = productObjectID

	return variantDao, nil
}

// ProductDao provides the model definition for product data to be stored in mongodb
type ProductDao struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RestaurantID primitive.ObjectID `bson:"restaurant_id" json:"restaurant_id"`
	CategoryID   primitive.ObjectID `bson:"category_id" json:"category_id"`
	Name         string             `bson:"name" json:"name"`
	Description  string             `bson:"description" json:"description"`
	IsVeg        bool               `bson:"is_veg" json:"is_veg"`
	InStock      bool               `bson:"in_stock"  json:"in_stock"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// GetProductDao converts and returns product Dao object from product schema
func GetProductDao(product product.Product) (ProductDao, errors.AppError) {
	productDao := ProductDao{
		Name:        product.Name,
		Description: product.Description,
		IsVeg:       product.IsVeg,
		InStock:     product.InStock,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	if product.ID != "" {
		productObjectID, err := primitive.ObjectIDFromHex(product.ID)
		if err != nil {
			return ProductDao{}, errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
		}
		productDao.ID = productObjectID
	}

	// restaurant id is required
	restaurantObjectID, err := primitive.ObjectIDFromHex(product.RestaurantID)
	if err != nil {
		return ProductDao{}, errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
	}
	productDao.RestaurantID = restaurantObjectID

	categoryObjectID, err := primitive.ObjectIDFromHex(product.CategoryID)
	if err != nil {
		return ProductDao{}, errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
	}
	productDao.CategoryID = categoryObjectID

	return productDao, nil
}
