package repositories

import (
	"context"

	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/dto"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/models"
)

// HealthRepository provides interface for Health repositories
type HealthRepository interface {
	HealthCheck(context.Context) errors.AppError
}

// RestaurantRepository provides interface for Restaurant repository
type RestaurantRepository interface {
	Create(context.Context, models.Restaurant) (models.Restaurant, errors.AppError)
	Get(context.Context, string) (models.Restaurant, errors.AppError)
	Delete(context.Context, string) errors.AppError
	GetAllRestaurants(context.Context, dto.GetAllRestaurantsRequestQuery,
		int64) ([]models.Restaurant, errors.AppError)
	GetAllRestaurantsTotalCount(context.Context, dto.GetAllRestaurantsRequestQuery,
		int64) (int64, errors.AppError)
}

// ProductRepository provides interface for Product repository
type ProductRepository interface {
	Create(context.Context, models.Product) (models.Product, errors.AppError)
	Get(context.Context, string, string) (models.Product, errors.AppError)
	Delete(context.Context, string, string) errors.AppError
	GetProductsByRestaurantID(context.Context, string, dto.GetAllProductsRequestQuery) ([]models.Product, errors.AppError)
	GetProductsByRestaurantTotalCount(context.Context, string, dto.GetAllProductsRequestQuery) (int64, errors.AppError)
}
