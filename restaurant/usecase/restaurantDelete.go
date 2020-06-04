package usecase

import (
	"context"
	"net/http"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/errors"
)

func (interactor *restaurantInteractor) DeleteByID(ctx context.Context, auth authentication.Auth,
	restaurantID string) errors.AppError {

	restaurantObj, getError := interactor.GetByID(ctx, auth, restaurantID)
	if getError != nil {
		return getError
	}

	// check if user have access to delete the restaurant
	if (auth.GetUserID() == restaurantObj.MerchantID &&
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteOwn)) ||
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteAny) {

		// only closed restaurant can be deleted
		if restaurantObj.IsOpen {
			return errors.NewAppError("Restaurant in open state cannot be deleted", http.StatusBadRequest, nil)
		}

		// delete products of the provided restaurant
		deleteProductError := interactor.productRepository.DeleteByRestaurantID(ctx, restaurantID)
		if deleteProductError != nil {
			return deleteProductError
		}

		// delete categories of the provided restaurant
		deleteCategoryError := interactor.categoryRespository.DeleteByRestaurantID(ctx, restaurantID)
		if deleteCategoryError != nil {
			return deleteCategoryError
		}

		// finially delete the restaurant
		deleteError := interactor.restaurantRepository.DeleteByID(ctx, restaurantID)
		return deleteError
	}
	return errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}
