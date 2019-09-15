package dto

import (
	"fmt"

	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-catalog-service/constants"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/models"
	"gopkg.in/go-playground/validator.v9"
)

// CreateProductRequestBody provides the schema definition for create product api request body
type CreateProductRequestBody struct {
	Name         string  `json:"name" validate:"required,min=6,max=30"`
	Description  string  `json:"description" validate:"max=60"`
	Price        float64 `json:"price" validate:"required"`
	DiscountType string  `json:"discount_type" validate:"required"`
	Discount     float64 `json:"discount" validate:"required"`
}

// CreateProductRequestParam provides the schema definition for create product api request params
type CreateProductRequestParam struct {
	RestaurantID string `json:"restaurantId" validate:"required"`
}

// CreateProductRequest provides the schema definition for create product api request
type CreateProductRequest struct {
	UserID   string                    `json:"-"`
	UserRole string                    `json:"-"`
	AppID    string                    `json:"-"`
	Param    CreateProductRequestParam `json:"param" validate:"required,dive"`
	Body     CreateProductRequestBody  `json:"body" validate:"required,dive"`
}

// Validate validates CreateProductRequest
func (dto CreateProductRequest) Validate(validate *validator.Validate) errors.AppError {
	var errMsg string
	err := validate.Struct(dto)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMsg = fmt.Sprintf("validation for field '%s' failed on '%s'", err.Field(), err.Tag())
			break
		}
		return errors.NewAppError(errMsg, errors.StatusBadRequest, err)
	}

	if dto.Body.DiscountType != constants.ProductDiscountTypeFlat &&
		dto.Body.DiscountType != constants.ProductDiscountTypePercentage {
		return errors.NewAppError("Invalid value for `discount_type`", errors.StatusBadRequest, nil)
	}
	return nil
}

// GetProductRequestParam provides the schema definition for get a product api request params
type GetProductRequestParam struct {
	RestaurantID string `json:"restaurantId" validate:"required"`
	ProductID    string `json:"productId" validate:"required"`
}

// GetProductRequest provides the schema definition for get a product api request
type GetProductRequest struct {
	UserID   string                 `json:"-"`
	UserRole string                 `json:"-"`
	AppID    string                 `json:"-"`
	Param    GetProductRequestParam `json:"param" validate:"required,dive"`
}

// Validate validates GetProductRequest
func (dto GetProductRequest) Validate(validate *validator.Validate) errors.AppError {
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

// DeleteProductRequestParam provides the schema definition for delete a product api request params
type DeleteProductRequestParam struct {
	RestaurantID string `json:"restaurantId" validate:"required"`
	ProductID    string `json:"productId" validate:"required"`
}

// DeleteProductRequest provides the schema definition for get a product api request
type DeleteProductRequest struct {
	UserID   string                    `json:"-"`
	UserRole string                    `json:"-"`
	AppID    string                    `json:"-"`
	Param    DeleteProductRequestParam `json:"param" validate:"required,dive"`
}

// Validate validates DeleteProductRequest
func (dto DeleteProductRequest) Validate(validate *validator.Validate) errors.AppError {
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

// GetAllProductsRequestParam provides the schema definition for create product api request params
type GetAllProductsRequestParam struct {
	RestaurantID string `json:"restaurantId" validate:"required"`
}

// GetAllProductsRequestQuery provides the schema definition for get all Products Api request query
type GetAllProductsRequestQuery struct {
	ID         []string
	PageNumber int64 `schema:"pageNumber" json:"pageNumber" validate:"gte=0"`
	PageSize   int64 `schema:"pageSize" json:"pageSize" validate:"lte=100"`
}

// GetAllProductsRequest provides the schema definition for create product api request
type GetAllProductsRequest struct {
	UserID   string                     `json:"-"`
	UserRole string                     `json:"-"`
	AppID    string                     `json:"-"`
	Param    GetAllProductsRequestParam `json:"param" validate:"required,dive"`
	Query    GetAllProductsRequestQuery `json:"query" validate:"required,dive"`
}

// Validate validates GetAllProductsRequest
func (dto GetAllProductsRequest) Validate(validate *validator.Validate) errors.AppError {
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

// GetAllProductsResponse provides the schema definition for create products Api request query
type GetAllProductsResponse struct {
	Total      int64            `json:"total"`
	PageNumber int64            `json:"page_number"`
	PageSize   int64            `json:"page_size"`
	TotalPages int64            `json:"total_pages"`
	Products   []models.Product `json:"products"`
}
