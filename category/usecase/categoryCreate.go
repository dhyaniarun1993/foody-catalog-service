package usecase

import (
	"context"
	"net/http"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-catalog-service/category"

	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/errors"
)

func (interactor *categoryInteractor) Create(ctx context.Context, auth authentication.Auth,
	categoryObj category.Category) (category.Category, errors.AppError) {

	validationError := categoryObj.Validate(interactor.validator)
	if validationError != nil {
		return category.Category{}, validationError
	}

	// user should have permission to get the restaurant
	restaurant, getRestaurantError := interactor.restaurantInteractor.GetByID(ctx, auth,
		categoryObj.RestaurantID)
	if getRestaurantError != nil {
		return category.Category{}, getRestaurantError
	}

	if (restaurant.MerchantID == auth.GetUserID() &&
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteOwn)) ||
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteAny) {

		var createCategoryError errors.AppError
		categoryObj, createCategoryError := interactor.categoryRepository.Create(ctx, categoryObj)
		if createCategoryError != nil {
			return category.Category{}, createCategoryError
		}
		return categoryObj, nil
	}

	return category.Category{}, errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}
