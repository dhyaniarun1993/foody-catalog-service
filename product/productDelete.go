package product

import (
	"context"
	"net/http"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/errors"
)

func (interactor *productInteractor) DeleteByID(ctx context.Context, auth authentication.Auth,
	productID string) errors.AppError {

	product, getProductError := interactor.GetByID(ctx, auth, productID)
	if getProductError != nil {
		return getProductError
	}

	// user should have permission to get restaurant
	restaurant, getRestaurantError := interactor.restaurantInteractor.GetByID(ctx,
		auth, product.RestaurantID)
	if getRestaurantError != nil {
		return getRestaurantError
	}

	if (restaurant.MerchantID == auth.GetUserID() &&
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteOwn)) ||
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteAny) {

		repositoryError := interactor.productRepository.DeleteByID(ctx, productID)
		return repositoryError
	}
	return errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}
