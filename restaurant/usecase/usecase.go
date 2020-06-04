package usecase

import (
	"context"

	"gopkg.in/go-playground/validator.v9"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-catalog-service/restaurant"

	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-common/logger"
)

type restaurantRepository interface {
	Create(context.Context, restaurant.Restaurant) (restaurant.Restaurant, errors.AppError)
	GetByID(context.Context, string) (restaurant.Restaurant, errors.AppError)
	DeleteByID(context.Context, string) errors.AppError
	GetAllRestaurants(context.Context, GetAllRestaurantsRequest,
		int64) ([]restaurant.Restaurant, errors.AppError)
	GetAllRestaurantsTotalCount(context.Context, GetAllRestaurantsRequest,
		int64) (int64, errors.AppError)
}

type categoryRespository interface {
	DeleteByRestaurantID(context.Context, string) errors.AppError
}

type productRepository interface {
	DeleteByRestaurantID(context.Context, string) errors.AppError
}

// Interactor provides interface for restaurant interactor
type Interactor interface {
	Create(ctx context.Context, auth authentication.Auth,
		restaurant restaurant.Restaurant) (restaurant.Restaurant, errors.AppError)
	GetByID(ctx context.Context, auth authentication.Auth,
		restaurantID string) (restaurant.Restaurant, errors.AppError)
	DeleteByID(ctx context.Context, auth authentication.Auth, restaurantID string) errors.AppError
	GetAllRestaurants(ctx context.Context, auth authentication.Auth,
		request GetAllRestaurantsRequest) (GetAllRestaurantsResponse, errors.AppError)
}

type restaurantInteractor struct {
	restaurantRepository restaurantRepository
	categoryRespository  categoryRespository
	productRepository    productRepository
	logger               *logger.Logger
	rbac                 acl.RBAC
	validator            *validator.Validate
}

// NewRestaurantInteractor creates and return restaurant Interactor
func NewRestaurantInteractor(restaurantRepository restaurantRepository, categoryRespository categoryRespository,
	productRepository productRepository, logger *logger.Logger, rbac acl.RBAC,
	validator *validator.Validate) Interactor {
	return &restaurantInteractor{
		restaurantRepository: restaurantRepository,
		categoryRespository:  categoryRespository,
		productRepository:    productRepository,
		logger:               logger,
		rbac:                 rbac,
		validator:            validator,
	}
}
