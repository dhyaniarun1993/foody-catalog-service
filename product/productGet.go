package product

import (
	"context"
	"net/http"
	"reflect"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/errors"
)

func (interactor *productInteractor) GetByID(ctx context.Context, auth authentication.Auth,
	productID string) (Product, errors.AppError) {

	// get Product from datastore
	product, repositoryError := interactor.productRepository.GetByID(ctx, productID)
	if repositoryError != nil {
		return Product{}, repositoryError
	}

	// check if product is empty
	if reflect.DeepEqual(product, Product{}) {
		return Product{}, errors.NewAppError("Resource not found", http.StatusNotFound, nil)
	}

	// user should have permission to get the restaurant
	restaurant, getRestaurantError := interactor.restaurantInteractor.GetByID(ctx, auth,
		product.RestaurantID)
	if getRestaurantError != nil {
		return Product{}, getRestaurantError
	}

	if (restaurant.MerchantID == auth.GetUserID() &&
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogReadOwn)) ||
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogReadAny) {

		return product, nil
	}
	return Product{}, errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}
