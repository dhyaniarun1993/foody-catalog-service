package restaurant

import (
	"context"
	"net/http"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/errors"
)

func (interactor *restaurantInteractor) Create(ctx context.Context, auth authentication.Auth,
	restaurant Restaurant) (Restaurant, errors.AppError) {

	validationError := restaurant.Validate(interactor.validator)
	if validationError != nil {
		return Restaurant{}, validationError
	}

	if (auth.GetUserID() == restaurant.MerchantID &&
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteOwn)) ||
		(interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteAny)) {

		var repositoryError errors.AppError
		restaurant, repositoryError = interactor.restaurantRepository.Create(ctx, restaurant)
		return restaurant, repositoryError
	}
	return Restaurant{}, errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}
