package usecase

import (
	"context"
	"net/http"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-catalog-service/product"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/errors"
)

func (interactor *productInteractor) CreateProduct(ctx context.Context, auth authentication.Auth,
	productObj product.Product) (product.Product, errors.AppError) {

	// validate product schema
	validationError := productObj.Validate(interactor.validator)
	if validationError != nil {
		return product.Product{}, validationError
	}

	// check if restaurant exist
	// user should have permission to get the restaurant
	restaurant, getRestaurantError := interactor.restaurantInteractor.GetByID(ctx, auth,
		productObj.RestaurantID)
	if getRestaurantError != nil {
		return product.Product{}, getRestaurantError
	}

	// check if category exist
	// user should have permission to get the category
	category, getCategoryError := interactor.categoryInteractor.GetByID(ctx, auth, productObj.CategoryID)
	if getCategoryError != nil {
		return product.Product{}, getCategoryError
	}

	// check if category belongs to the restaurant
	if category.RestaurantID != restaurant.ID {
		return product.Product{}, errors.NewAppError("Category doesnot belong to the restaurant",
			http.StatusBadRequest, nil)
	}

	// check if user have permission to create product
	if (restaurant.MerchantID == auth.GetUserID() &&
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteOwn)) ||
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogWriteAny) {

		var createProductError errors.AppError
		productObj, createProductError = interactor.productRepository.CreateProduct(ctx, productObj)
		return productObj, createProductError
	}
	return product.Product{}, errors.NewAppError("Forbidden", http.StatusForbidden, nil)

}
