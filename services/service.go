package services

import (
	"context"

	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/dto"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/models"
)

// HealthService provides interface for health service
type HealthService interface {
	HealthCheck(context.Context) errors.AppError
}

// RestaurantService provides interface for restaurant service
type RestaurantService interface {
	Create(context.Context, dto.CreateRestaurantRequest) (models.Restaurant, errors.AppError)
	Get(context.Context, dto.GetRestaurantRequest) (models.Restaurant, errors.AppError)
	Delete(context.Context, dto.DeleteRestaurantRequest) errors.AppError
	GetAllRestaurants(context.Context,
		dto.GetAllRestaurantsRequest) (dto.GetAllRestaurantsResponse, errors.AppError)
}

// ProductService provides interface for product service
type ProductService interface {
	Create(context.Context, dto.CreateProductRequest) (models.Product, errors.AppError)
	Get(context.Context, dto.GetProductRequest) (models.Product, errors.AppError)
	Delete(context.Context, dto.DeleteProductRequest) errors.AppError
	GetAllProducts(context.Context,
		dto.GetAllProductsRequest) (dto.GetAllProductsResponse, errors.AppError)
}
