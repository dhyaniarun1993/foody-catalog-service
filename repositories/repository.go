package repositories

import (
	"context"

	"github.com/dhyaniarun1993/foody-catalog-service/product"
	"github.com/dhyaniarun1993/foody-catalog-service/restaurant"
	"github.com/dhyaniarun1993/foody-common/errors"
)

// HealthRepository provides interface for Health repositories
type HealthRepository interface {
	HealthCheck(context.Context) errors.AppError
}

// RestaurantRepository provides interface for Restaurant repository
type RestaurantRepository interface {
	Create(ctx context.Context, restaurant restaurant.Restaurant) (restaurant.Restaurant, errors.AppError)
	GetByID(ctx context.Context, restaurantID string) (restaurant.Restaurant, errors.AppError)
	DeleteByID(ctx context.Context, restaurantID string) errors.AppError
	GetAllRestaurants(context.Context, restaurant.GetAllRestaurantsRequest,
		int64) ([]restaurant.Restaurant, errors.AppError)
	GetAllRestaurantsTotalCount(context.Context, restaurant.GetAllRestaurantsRequest,
		int64) (int64, errors.AppError)
}

// ProductRepository provides interface for Product repository
type ProductRepository interface {
	Create(ctx context.Context, product product.Product) (product.Product, errors.AppError)
	GetByID(ctx context.Context, productID string) (product.Product, errors.AppError)
	DeleteByID(ctx context.Context, productID string) errors.AppError
}
