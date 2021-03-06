package repositories

import (
	"context"

	"github.com/dhyaniarun1993/foody-catalog-service/category"
	"github.com/dhyaniarun1993/foody-catalog-service/product"
	"github.com/dhyaniarun1993/foody-catalog-service/restaurant"
	restaurantUsecase "github.com/dhyaniarun1993/foody-catalog-service/restaurant/usecase"
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
	GetAllRestaurants(context.Context, restaurantUsecase.GetAllRestaurantsRequest,
		int64) ([]restaurant.Restaurant, errors.AppError)
	GetAllRestaurantsTotalCount(context.Context, restaurantUsecase.GetAllRestaurantsRequest,
		int64) (int64, errors.AppError)
}

// ProductRepository provides interface for Product repository
type ProductRepository interface {
	CreateProduct(ctx context.Context, product product.Product) (product.Product, errors.AppError)
	CreateVariant(ctx context.Context, variant product.Variant) (product.Variant, errors.AppError)
	GetProductByID(ctx context.Context, productID string) (product.Product, errors.AppError)
	GetVariantByID(ctx context.Context, variantID string) (product.Variant, errors.AppError)
	DeleteProductByID(ctx context.Context, productID string) errors.AppError
	DeleteVariantByID(ctx context.Context, variantID string) errors.AppError
	DeleteProductByRestaurantID(ctx context.Context, restaurantID string) errors.AppError
	DeleteProductByCategoryID(ctx context.Context, categoryID string) errors.AppError
}

// CategoryRepository provides interface for Category repository
type CategoryRepository interface {
	Create(ctx context.Context, category category.Category) (category.Category, errors.AppError)
	GetByID(ctx context.Context, categoryID string) (category.Category, errors.AppError)
	DeleteByID(ctx context.Context, categoryID string) errors.AppError
	DeleteByRestaurantID(ctx context.Context, restaurantID string) errors.AppError
}
