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

func (interactor *categoryInteractor) DeleteByID(ctx context.Context, auth authentication.Auth,
	categoryID string) errors.AppError {

	categoryObj, getCategoryError := interactor.categoryRepository.GetByID(ctx, categoryID)
	if getCategoryError != nil {
		return getCategoryError
	}

	// check if category is empty
	if reflect.DeepEqual(categoryObj, category.Category{}) {
		return errors.NewAppError("Resource not found", http.StatusNotFound, nil)
	}

	// user should have permission to get the restaurant
	restaurant, getRestuarantError := interactor.restaurantInteractor.GetByID(ctx, auth,
		categoryObj.RestaurantID)
	if getRestuarantError != nil {
		return getRestuarantError
	}

	// check if user have permission to delete category
	if restaurant.MerchantID == auth.GetUserID() &&
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteOwn) ||
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteAny) {

		deleteCategoryError := interactor.categoryRepository.DeleteByID(ctx, categoryID)
		return deleteCategoryError
	}
	return errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}
