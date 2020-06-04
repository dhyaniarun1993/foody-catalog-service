package usecase

import (
	"context"
	"net/http"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/errors"
)

func (interactor *productInteractor) RemoveVariant(ctx context.Context, auth authentication.Auth,
	productID string, variantID string) errors.AppError {

	// check if product exist
	productObj, getProductError := interactor.GetProductByID(ctx, auth, productID)
	if getProductError != nil {
		return getProductError
	}

	// get restaurant to check if user have permission to remove variant
	// user should have permission to get the restaurant
	restaurant, getRestaurantError := interactor.restaurantInteractor.GetByID(ctx, auth,
		productObj.RestaurantID)
	if getRestaurantError != nil {
		return getRestaurantError
	}

	// check if user have permission to remove variant from product
	if (restaurant.MerchantID == auth.GetUserID() &&
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteOwn)) ||
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteAny) {

		variant, getVariantError := interactor.productRepository.GetVariantByID(ctx, variantID)
		if getVariantError != nil {
			return getVariantError
		}

		// check if variant belong to the product
		if variant.ProductID != productObj.ID {
			return errors.NewAppError("Variant is not part of the provided product", http.StatusBadRequest, nil)
		}

		// delete variant
		var deleteVariantError errors.AppError
		deleteVariantError = interactor.productRepository.DeleteVariantByID(ctx, variantID)
		return deleteVariantError
	}
	return errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}
