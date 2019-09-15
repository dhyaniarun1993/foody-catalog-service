package dto

import (
	"fmt"

	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/models"
	"gopkg.in/go-playground/validator.v9"
)

// Address provides the schema definition for Address
type Address struct {
	Location GeoJSON `json:"location" validate:"required,dive"`
}

// GeoJSON provides the schema definition for Geo Location
type GeoJSON struct {
	Coordinates []float64 `json:"coordinates" validate:"required,eq=2"`
}

// CreateRestaurantRequestBody provides the schema definition for create restaurant Api request body
type CreateRestaurantRequestBody struct {
	MerchantID  string  `json:"merchant_id" validate:"required"`
	Name        string  `json:"name" validate:"required,min=6,max=30"`
	Description string  `json:"description" validate:"max=60"`
	Address     Address `json:"address" validate:"required,dive"`
}

// CreateRestaurantRequest provides the schema definition for create restaurant api request
type CreateRestaurantRequest struct {
	UserID   string                      `json:"-"`
	UserRole string                      `json:"-"`
	AppID    string                      `json:"-"`
	Body     CreateRestaurantRequestBody `json:"body" validate:"required,dive"`
}

// Validate validates CreateMerchantRequest
func (dto CreateRestaurantRequest) Validate(validate *validator.Validate) errors.AppError {
	var errMessage string
	err := validate.Struct(dto)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMessage = fmt.Sprintf("validation for field '%s' failed on '%s'", err.Field(), err.Tag())
			break
		}
		return errors.NewAppError(errMessage, errors.StatusBadRequest, err)
	}
	if dto.Body.Address.Location.Coordinates[0] < -180 || dto.Body.Address.Location.Coordinates[0] > 180 {
		return errors.NewAppError("Invalid longitude", errors.StatusBadRequest, err)
	}
	if dto.Body.Address.Location.Coordinates[1] < -90 || dto.Body.Address.Location.Coordinates[1] > 90 {
		return errors.NewAppError("Invalid latitude", errors.StatusBadRequest, err)
	}
	return nil
}

// GetRestaurantRequestParam provides the schema definition for get a restaurant api request params
type GetRestaurantRequestParam struct {
	RestaurantID string `json:"restaurantId" validate:"required"`
}

// GetRestaurantRequest provides the schema definition for get a restaurant api request
type GetRestaurantRequest struct {
	UserID   string                    `json:"-"`
	UserRole string                    `json:"-"`
	AppID    string                    `json:"-"`
	Param    GetRestaurantRequestParam `json:"param" validate:"required,dive"`
}

// Validate validates GetRestaurantRequest
func (dto GetRestaurantRequest) Validate(validate *validator.Validate) errors.AppError {
	var errMessage string
	err := validate.Struct(dto)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMessage = fmt.Sprintf("validation for field '%s' failed on '%s'", err.Field(), err.Tag())
			break
		}
		return errors.NewAppError(errMessage, errors.StatusBadRequest, err)
	}
	return nil
}

// DeleteRestaurantRequestParam provides the schema definition for delete a restaurant api request params
type DeleteRestaurantRequestParam struct {
	RestaurantID string `json:"restaurantId" validate:"required"`
}

// DeleteRestaurantRequest provides the schema definition for get a restaurant api request
type DeleteRestaurantRequest struct {
	UserID   string                       `json:"-"`
	UserRole string                       `json:"-"`
	AppID    string                       `json:"-"`
	Param    DeleteRestaurantRequestParam `json:"param" validate:"required,dive"`
}

// Validate validates DeleteRestaurantRequest
func (dto DeleteRestaurantRequest) Validate(validate *validator.Validate) errors.AppError {
	var errMessage string
	err := validate.Struct(dto)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMessage = fmt.Sprintf("validation for field '%s' failed on '%s'", err.Field(), err.Tag())
			break
		}
		return errors.NewAppError(errMessage, errors.StatusBadRequest, err)
	}
	return nil
}

// GetAllRestaurantsRequestQuery provides the schema definition for get all restaurant Api request
type GetAllRestaurantsRequestQuery struct {
	ID         []string `schema:"id" json:"id"`
	MerchantID string   `schema:"merchantId" json:"merchantId"`
	PageNumber int64    `schema:"pageNumber" json:"pageNumber" validate:"gte=0"`
	PageSize   int64    `schema:"pageSize" json:"pageSize" validate:"lte=50"`
	Latitude   float64  `schema:"latitude" json:"latitude" validate:"required,latitude"`
	Longitude  float64  `schema:"longitude" json:"longitude" validate:"required,longitude"`
}

// GetAllRestaurantsRequest provides the schema definition for create restaurant api request
type GetAllRestaurantsRequest struct {
	UserID   string                        `json:"-"`
	UserRole string                        `json:"-"`
	AppID    string                        `json:"-"`
	Query    GetAllRestaurantsRequestQuery `json:"body" validate:"required,dive"`
}

// Validate validates GetAllRestaurantRequest
func (dto GetAllRestaurantsRequest) Validate(validate *validator.Validate) errors.AppError {
	var errMessage string
	err := validate.Struct(dto)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMessage = fmt.Sprintf("validation for field '%s' failed on '%s'", err.Field(), err.Tag())
			break
		}
		return errors.NewAppError(errMessage, errors.StatusBadRequest, err)
	}
	return nil
}

// GetAllRestaurantsResponse provides the schema definition for create restaurant Api request query
type GetAllRestaurantsResponse struct {
	Total       int64               `json:"total"`
	PageNumber  int64               `json:"page_number"`
	PageSize    int64               `json:"page_size"`
	TotalPages  int64               `json:"total_pages"`
	Restaurants []models.Restaurant `json:"restaurants"`
}
