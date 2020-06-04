package usecase

import (
	"context"
	"net/http"
	"reflect"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-catalog-service/product"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/errors"
)

func (interactor *productInteractor) GetProductByID(ctx context.Context, auth authentication.Auth,
	productID string) (product.Product, errors.AppError) {

	// get Product from datastore
	productObj, repositoryError := interactor.productRepository.GetProductByID(ctx, productID)
	if repositoryError != nil {
		return product.Product{}, repositoryError
	}

	// check if product is empty
	if reflect.DeepEqual(productObj, product.Product{}) {
		return product.Product{}, errors.NewAppError("Unable to find product", http.StatusNotFound, nil)
	}

	// user should have permission to get the restaurant
	restaurant, getRestaurantError := interactor.restaurantInteractor.GetByID(ctx, auth,
		productObj.RestaurantID)
	if getRestaurantError != nil {
		return product.Product{}, getRestaurantError
	}

	if (restaurant.MerchantID == auth.GetUserID() &&
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogReadOwn)) ||
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogReadAny) {

		return productObj, nil
	}
	return product.Product{}, errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}
