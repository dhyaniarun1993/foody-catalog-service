package restaurant

import (
	"context"
	"net/http"
	"reflect"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/errors"
)

func (interactor *restaurantInteractor) DeleteByID(ctx context.Context, auth authentication.Auth,
	restaurantID string) errors.AppError {

	restaurant, getError := interactor.restaurantRepository.GetByID(ctx, restaurantID)
	if getError != nil {
		return getError
	}

	if reflect.DeepEqual(restaurant, Restaurant{}) {
		return errors.NewAppError("Resource not found", http.StatusNotFound, nil)
	}

	if (auth.GetUserID() == restaurant.MerchantID &&
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteOwn)) ||
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteAny) {

		deleteError := interactor.restaurantRepository.DeleteByID(ctx, restaurantID)
		return deleteError
	}
	return errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}
