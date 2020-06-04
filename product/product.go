package product

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dhyaniarun1993/foody-common/errors"
	"gopkg.in/go-playground/validator.v9"
)

// Price provides the schema definition for Price
type Price struct {
	Amount   float64 `bson:"amount" json:"amount" validate:"required"`
	Currency string  `bson:"currency" json:"currency" validate:"required"`
}

// Variant provides the schema definition for Products Variant
type Variant struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	ProductID   string    `bson:"product_id" json:"product_id"`
	Name        string    `bson:"name" json:"name" validate:"required,min=6,max=30"`
	Description string    `bson:"description" json:"description" validate:"max=120"`
	Price       Price     `bson:"price" json:"price" validate:"required,dive"`
	InStock     *bool     `bson:"in_stock"  json:"in_stock" validate:"required"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

// Validate validates Variant schema
func (variant Variant) Validate(validate *validator.Validate) errors.AppError {
	var errMessage string
	// validate struct data
	err := validate.Struct(variant)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMessage = fmt.Sprintf("Invalid value for field '%s'", err.Field())
			break
		}
		return errors.NewAppError(errMessage, http.StatusBadRequest, err)
	}
	return nil
}

// Product provides the model definition for Product
type Product struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	RestaurantID string    `bson:"restaurant_id" json:"restaurant_id" validate:"required"`
	CategoryID   string    `bson:"category_id" json:"category_id" validate:"required"`
	Name         string    `bson:"name" json:"name" validate:"required,min=6,max=30"`
	Description  string    `bson:"description" json:"description" validate:"max=120"`
	IsVeg        bool      `bson:"is_veg" json:"is_veg"`
	InStock      bool      `bson:"in_stock"  json:"in_stock" validate:"required"`
	Variants     []Variant `bson:"variants" json:"variants,omitempty" validate:"required,dive"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
}

// Validate validates Product schema
func (product Product) Validate(validate *validator.Validate) errors.AppError {
	var errMessage string
	// validate struct data
	err := validate.Struct(product)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMessage = fmt.Sprintf("Invalid value for field '%s'", err.Field())
			break
		}
		return errors.NewAppError(errMessage, http.StatusBadRequest, err)
	}
	return nil
}
