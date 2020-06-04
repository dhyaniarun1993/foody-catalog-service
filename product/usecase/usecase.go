package usecase

import (
	"context"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	categoryUsecase "github.com/dhyaniarun1993/foody-catalog-service/category/usecase"
	"github.com/dhyaniarun1993/foody-catalog-service/product"
	restaurantUsecase "github.com/dhyaniarun1993/foody-catalog-service/restaurant/usecase"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-common/logger"
	"gopkg.in/go-playground/validator.v9"
)

type productRepository interface {
	CreateProduct(ctx context.Context, product product.Product) (product.Product, errors.AppError)
	CreateVariant(ctx context.Context, variant product.Variant) (product.Variant, errors.AppError)
	GetProductByID(ctx context.Context, productID string) (product.Product, errors.AppError)
	GetVariantByID(ctx context.Context, variantID string) (product.Variant, errors.AppError)
	DeleteProductByID(ctx context.Context, productID string) errors.AppError
	DeleteVariantByID(ctx context.Context, variantID string) errors.AppError
}

// Interactor provides interface for product interactor
type Interactor interface {
	CreateProduct(ctx context.Context, auth authentication.Auth, product product.Product) (product.Product, errors.AppError)
	AddVariant(ctx context.Context, auth authentication.Auth,
		productID string, variant product.Variant) (product.Variant, errors.AppError)
	GetProductByID(ctx context.Context, auth authentication.Auth, productID string) (product.Product, errors.AppError)
	DeleteProductByID(ctx context.Context, auth authentication.Auth, productID string) errors.AppError
	RemoveVariant(ctx context.Context, auth authentication.Auth, productID string,
		variantID string) errors.AppError
}

type productInteractor struct {
	productRepository    productRepository
	restaurantInteractor restaurantUsecase.Interactor
	categoryInteractor   categoryUsecase.Interactor
	logger               *logger.Logger
	rbac                 acl.RBAC
	validator            *validator.Validate
}

// NewProductInteractor creates and return product Interactor
func NewProductInteractor(productRepository productRepository, restaurantInteractor restaurantUsecase.Interactor,
	categoryInteractor categoryUsecase.Interactor, logger *logger.Logger, rbac acl.RBAC,
	validator *validator.Validate) Interactor {
	return &productInteractor{
		productRepository:    productRepository,
		restaurantInteractor: restaurantInteractor,
		categoryInteractor:   categoryInteractor,
		logger:               logger,
		rbac:                 rbac,
		validator:            validator,
	}
}
