package usecase

import (
	"context"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-catalog-service/category"
	restaurantUsecase "github.com/dhyaniarun1993/foody-catalog-service/restaurant/usecase"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/logger"
	"gopkg.in/go-playground/validator.v9"

	"github.com/dhyaniarun1993/foody-common/errors"
)

type categoryRepository interface {
	Create(ctx context.Context, category category.Category) (category.Category, errors.AppError)
	GetByID(ctx context.Context, categoryID string) (category.Category, errors.AppError)
	DeleteByID(ctx context.Context, categoryID string) errors.AppError
}

// Interactor provides interface for category interactor
type Interactor interface {
	Create(ctx context.Context, auth authentication.Auth,
		category category.Category) (category.Category, errors.AppError)
	GetByID(ctx context.Context, auth authentication.Auth,
		categoryID string) (category.Category, errors.AppError)
	DeleteByID(ctx context.Context, auth authentication.Auth, categoryID string) errors.AppError
}

type categoryInteractor struct {
	categoryRepository   categoryRepository
	restaurantInteractor restaurantUsecase.Interactor
	logger               *logger.Logger
	validator            *validator.Validate
	rbac                 acl.RBAC
}

// NewCategoryInteractor creates and return category Interactor
func NewCategoryInteractor(categoryRepository categoryRepository,
	restaurantInteractor restaurantUsecase.Interactor,
	logger *logger.Logger, rbac acl.RBAC, validator *validator.Validate) Interactor {

	return &categoryInteractor{
		categoryRepository:   categoryRepository,
		restaurantInteractor: restaurantInteractor,
		logger:               logger,
		validator:            validator,
		rbac:                 rbac,
	}
}
