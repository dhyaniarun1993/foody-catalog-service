package product

import (
	"context"
	"net/http"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/errors"
)

func (interactor *productInteractor) Create(ctx context.Context, auth authentication.Auth,
	product Product) (Product, errors.AppError) {

	// validate product schema
	validationError := product.Validate(interactor.validator)
	if validationError != nil {
		return Product{}, validationError
	}

	// user should have permission to get the restaurant
	restaurant, getRestaurantError := interactor.restaurantInteractor.GetByID(ctx, auth,
		product.RestaurantID)
	if getRestaurantError != nil {
		return Product{}, getRestaurantError
	}

	// check if user have permission to create product
	if (restaurant.MerchantID == auth.GetUserID() &&
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteOwn)) ||
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteAny) {

		var createProductError errors.AppError
		product, createProductError = interactor.productRepository.Create(ctx, product)
		return product, createProductError
	}
	return Product{}, errors.NewAppError("Forbidden", http.StatusForbidden, nil)

}
