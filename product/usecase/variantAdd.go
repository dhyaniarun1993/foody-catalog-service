package usecase

import (
	"context"
	"net/http"

	"github.com/dhyaniarun1993/foody-common/authentication"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-catalog-service/product"
	"github.com/dhyaniarun1993/foody-common/errors"
)

func (interactor *productInteractor) AddVariant(ctx context.Context, auth authentication.Auth,
	productID string, variant product.Variant) (product.Variant, errors.AppError) {

	variant.ProductID = productID
	// validate product schema
	validationError := variant.Validate(interactor.validator)
	if validationError != nil {
		return product.Variant{}, validationError
	}

	// check if product exist
	productObj, getProductError := interactor.GetProductByID(ctx, auth, variant.ProductID)
	if getProductError != nil {
		return product.Variant{}, getProductError
	}

	// get restaurant to check if user have permission to add variant
	// user should have permission to get the restaurant
	restaurant, getRestaurantError := interactor.restaurantInteractor.GetByID(ctx, auth,
		productObj.RestaurantID)
	if getRestaurantError != nil {
		return product.Variant{}, getRestaurantError
	}

	// check if user have permission to add variant to product
	if (restaurant.MerchantID == auth.GetUserID() &&
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteOwn)) ||
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteAny) {

		var createVariantError errors.AppError
		variant, createVariantError = interactor.productRepository.CreateVariant(ctx, variant)
		return variant, createVariantError
	}
	return product.Variant{}, errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}
