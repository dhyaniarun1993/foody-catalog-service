package usecase

import (
	"context"
	"net/http"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-catalog-service/restaurant"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/errors"
)

func (interactor *restaurantInteractor) Create(ctx context.Context, auth authentication.Auth,
	restaurantObj restaurant.Restaurant) (restaurant.Restaurant, errors.AppError) {

	validationError := restaurantObj.Validate(interactor.validator)
	if validationError != nil {
		return restaurant.Restaurant{}, validationError
	}

	if (auth.GetUserID() == restaurantObj.MerchantID &&
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteOwn)) ||
		(interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteAny)) {

		var repositoryError errors.AppError
		restaurantObj, repositoryError = interactor.restaurantRepository.Create(ctx, restaurantObj)
		return restaurantObj, repositoryError
	}
	return restaurant.Restaurant{}, errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}
