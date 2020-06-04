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

	if (auth.GetUserID() == restaurantObj.MerchantID &&
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteOwn)) ||
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteAny) {

		deleteError := interactor.restaurantRepository.DeleteByID(ctx, restaurantID)
		return deleteError
	}
	return errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}
