package restaurant

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

// Fees provides the schema definition for Fees
type Fees struct {
	Fee  Price  `bson:"fee" json:"fee" validate:"required,dive"`
	Name string `bson:"name" json:"name" validate:"required"`
}

// Address provides the schema definition for Address
type Address struct {
	Street   string  `bson:"street" json:"street" validate:"required"`
	City     string  `bson:"city" json:"city" validate:"required"`
	State    string  `bson:"state" json:"state" validate:"required"`
	Country  string  `bson:"country" json:"country" validate:"required"`
	Pincode  string  `bson:"pincode" json:"pincode" validate:"required"`
	Location GeoJSON `bson:"location" json:"location" validate:"required,dive"`
}

// GeoJSON provides the model definition for Geo Location
type GeoJSON struct {
	Type        string    `bson:"type" json:"-"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates" validate:"required,min=2,max=2"`
}

// Restaurant provides the model definition for Restaurant
type Restaurant struct {
	ID               string    `bson:"_id,omitempty" json:"id"`
	MerchantID       string    `bson:"merchant_id" json:"merchant_id" validate:"required"`
	Name             string    `bson:"name" json:"name" validate:"required,min=2,max=30"`
	Description      string    `bson:"description" json:"description" validate:"max=120"`
	ReviewsRatingSum int64     `bson:"reviews_rating_sum" json:"reviews_rating_sum"`
	ReviewsCount     int64     `bson:"reviews_count" json:"reviews_count"`
	Address          Address   `bson:"address" json:"address" validate:"required,dive"`
	RestaurantFees   Fees      `bson:"restaurant_fees" json:"restaurant_fees" validate:"required,dive"`
	IsOpen           bool      `bson:"is_open" json:"is_open"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at" json:"updated_at"`
}

// Validate validates Restaurant schema
func (restaurant Restaurant) Validate(validate *validator.Validate) errors.AppError {
	var errMessage string
	// validate struct data
	err := validate.Struct(restaurant)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMessage = fmt.Sprintf("Invalid value for field '%s'", err.Field())
			break
		}
		return errors.NewAppError(errMessage, http.StatusBadRequest, err)
	}

	// validate longitude
	if restaurant.Address.Location.Coordinates[0] < -180 ||
		restaurant.Address.Location.Coordinates[0] > 180 {

		return errors.NewAppError("Invalid value for field longitude", http.StatusBadRequest, err)
	}
	// validate latitude
	if restaurant.Address.Location.Coordinates[1] < -90 ||
		restaurant.Address.Location.Coordinates[1] > 90 {

		return errors.NewAppError("Invalid value for field latitude", http.StatusBadRequest, err)
	}
	return nil
}
