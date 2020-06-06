package category

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dhyaniarun1993/foody-catalog-service/product"
	"github.com/dhyaniarun1993/foody-common/errors"
	"gopkg.in/go-playground/validator.v9"
)

// Category provides the model definition for Product Category
type Category struct {
	ID           string            `bson:"_id,omitempty" json:"id"`
	RestaurantID string            `bson:"restaurant_id,omitempty" json:"restaurant_id" validate:"required"`
	Name         string            `bson:"name" json:"name" validate:"required,min=2,max=30"`
	Description  string            `bson:"description" json:"description" validate:"max=120"`
	Products     []product.Product `bson:"products" json:"products,omitempty"`
	CreatedAt    time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time         `bson:"updated_at" json:"updated_at"`
}

// Validate validates Category schema
func (category *Category) Validate(validate *validator.Validate) errors.AppError {
	var errMessage string
	// validate struct data
	err := validate.Struct(category)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMessage = fmt.Sprintf("Invalid value for field '%s'", err.Field())
			break
		}
		return errors.NewAppError(errMessage, http.StatusBadRequest, err)
	}
	return nil
}
