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
	Amount   float64 `bson:"amount" json:"amount" validate:"required"`
	Currency string  `bson:"currency" json:"currency" validate:"required"`
}

// VariantDao provides the schema definition for Products Variant data to be stored in mongodb
type VariantDao struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID   primitive.ObjectID `bson:"product_id" json:"product_id"`
	Name        string             `bson:"name" json:"name" validate:"required,min=6,max=30"`
	Description string             `bson:"description" json:"description" validate:"max=120"`
	Price       PriceDao           `bson:"price" json:"price" validate:"required,dive"`
	InStock     *bool              `bson:"in_stock"  json:"in_stock" validate:"required"`
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
	if variant.ProductID != "" {
		productObjectID, err := primitive.ObjectIDFromHex(variant.ProductID)
		if err != nil {
			return VariantDao{}, errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
		}
		variantDao.ProductID = productObjectID
	}
	return variantDao, nil
}

// ProductDao provides the model definition for product data to be stored in mongodb
type ProductDao struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RestaurantID primitive.ObjectID `bson:"restaurant_id" json:"restaurant_id" validate:"required"`
	Name         string             `bson:"name" json:"name" validate:"required,min=6,max=30"`
	Description  string             `bson:"description" json:"description" validate:"max=120"`
	IsVeg        *bool              `bson:"is_veg" json:"is_veg" validate:"required"`
	InStock      *bool              `bson:"in_stock"  json:"in_stock" validate:"required"`
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
	if product.RestaurantID != "" {
		restaurantObjectID, err := primitive.ObjectIDFromHex(product.RestaurantID)
		if err != nil {
			return ProductDao{}, errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
		}
		productDao.RestaurantID = restaurantObjectID
	}

	return productDao, nil
}
