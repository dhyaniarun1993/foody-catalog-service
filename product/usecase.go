package product

import (
	"context"

	"github.com/dhyaniarun1993/foody-common/authentication"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-catalog-service/restaurant"
	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-common/logger"
	"gopkg.in/go-playground/validator.v9"
)

type productRepository interface {
	Create(ctx context.Context, product Product) (Product, errors.AppError)
	GetByID(ctx context.Context, productID string) (Product, errors.AppError)
	DeleteByID(ctx context.Context, productID string) errors.AppError
}

// Interactor provides interface for product interactor
type Interactor interface {
	Create(ctx context.Context, auth authentication.Auth, product Product) (Product, errors.AppError)
	GetByID(ctx context.Context, auth authentication.Auth, productID string) (Product, errors.AppError)
	DeleteByID(ctx context.Context, auth authentication.Auth, productID string) errors.AppError
}

type productInteractor struct {
	productRepository    productRepository
	restaurantInteractor restaurant.Interactor
	logger               *logger.Logger
	rbac                 acl.RBAC
	validator            *validator.Validate
}

// NewProductInteractor creates and return product Interactor
func NewProductInteractor(productRepository productRepository,
	restaurantInteractor restaurant.Interactor, logger *logger.Logger,
	rbac acl.RBAC, validator *validator.Validate) Interactor {
	return &productInteractor{
		productRepository:    productRepository,
		restaurantInteractor: restaurantInteractor,
		logger:               logger,
		rbac:                 rbac,
		validator:            validator,
	}
}
