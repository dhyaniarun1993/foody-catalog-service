package usecase

import (
	"context"
	"net/http"
	"reflect"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-catalog-service/category"

	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/errors"
)

func (interactor *categoryInteractor) GetByID(ctx context.Context, auth authentication.Auth,
	categoryID string) (category.Category, errors.AppError) {

	categoryObj, getCategoryError := interactor.categoryRepository.GetByID(ctx, categoryID)
	if getCategoryError != nil {
		return category.Category{}, getCategoryError
	}

	// check if category is empty
	if reflect.DeepEqual(categoryObj, category.Category{}) {
		return category.Category{}, errors.NewAppError("Unable to find category", http.StatusNotFound, nil)
	}

	// user should have permission to get the restaurant
	restaurant, getRestaurantError := interactor.restaurantInteractor.GetByID(ctx, auth,
		categoryObj.RestaurantID)
	if getRestaurantError != nil {
		return category.Category{}, getRestaurantError
	}

	// check if user have permission to get category
	if (restaurant.MerchantID == auth.GetUserID() &&
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogReadOwn)) ||
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogReadAny) {

		return categoryObj, nil
	}
	return category.Category{}, errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}
