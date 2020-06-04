package usecase

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"reflect"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-catalog-service/restaurant"
	"github.com/dhyaniarun1993/foody-common/async"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/errors"
	"gopkg.in/go-playground/validator.v9"
)

func (interactor *restaurantInteractor) GetByID(ctx context.Context, auth authentication.Auth,
	restaurantID string) (restaurant.Restaurant, errors.AppError) {

	restaurantObj, repositoryError := interactor.restaurantRepository.GetByID(ctx, restaurantID)
	if repositoryError != nil {
		return restaurant.Restaurant{}, repositoryError
	}

	if reflect.DeepEqual(restaurantObj, restaurant.Restaurant{}) {
		return restaurant.Restaurant{}, errors.NewAppError("Unable to find restaurant", http.StatusNotFound, nil)
	}

	if (auth.GetUserID() == restaurantObj.MerchantID &&
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogReadOwn)) ||
		interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogReadAny) {
		return restaurantObj, nil
	}

	return restaurant.Restaurant{}, errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}

func (interactor *restaurantInteractor) GetAllRestaurants(ctx context.Context, auth authentication.Auth,
	request GetAllRestaurantsRequest) (GetAllRestaurantsResponse, errors.AppError) {

	var restaurants []restaurant.Restaurant
	var totalCount, maxDistance int64
	var restaurantResponse GetAllRestaurantsResponse

	maxDistance = 10000
	async, asyncCtx := async.WithContext(ctx)

	if request.PageNumber == 0 {
		request.PageNumber = 1
	}
	if request.PageSize == 0 {
		request.PageSize = 50
	}

	GetAllRestaurants := func() errors.AppError {
		var repositoryError errors.AppError
		restaurants, repositoryError = interactor.restaurantRepository.GetAllRestaurants(asyncCtx,
			request, maxDistance)
		return repositoryError
	}

	getTotalCount := func() errors.AppError {
		var repositoryError errors.AppError
		totalCount, repositoryError = interactor.restaurantRepository.GetAllRestaurantsTotalCount(asyncCtx,
			request, maxDistance)
		return repositoryError
	}

	if interactor.rbac.Can(auth.GetUserRole(), acl.PermissionCatalogReadAny) {

		async.Go(GetAllRestaurants)
		async.Go(getTotalCount)
		err := async.Wait()
		if err != nil {
			return restaurantResponse, err
		}

		restaurantResponse = GetAllRestaurantsResponse{
			Total:       totalCount,
			PageNumber:  request.PageNumber,
			PageSize:    request.PageSize,
			TotalPages:  int64(math.Ceil(float64(totalCount) / float64(request.PageSize))),
			Restaurants: restaurants,
		}
		return restaurantResponse, nil
	}
	return restaurantResponse, errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}

// GetAllRestaurantsRequest provides the schema definition for get all restaurant request
type GetAllRestaurantsRequest struct {
	PageNumber int64   `schema:"pageNumber" json:"pageNumber" validate:"gte=0"`
	PageSize   int64   `schema:"pageSize" json:"pageSize" validate:"lte=100"`
	Latitude   float64 `schema:"latitude" json:"latitude" validate:"required,latitude"`
	Longitude  float64 `schema:"longitude" json:"longitude" validate:"required,longitude"`
}

// Validate validates GetAllRestaurantsRequest
func (request GetAllRestaurantsRequest) Validate(validate *validator.Validate) errors.AppError {
	var errMessage string
	err := validate.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMessage = fmt.Sprintf("validation for field '%s' failed on '%s'", err.Field(), err.Tag())
			break
		}
		return errors.NewAppError(errMessage, http.StatusBadRequest, err)
	}
	return nil
}

// GetAllRestaurantsResponse provides the schema definition for get all restaurant response
type GetAllRestaurantsResponse struct {
	Total       int64                   `json:"total"`
	PageNumber  int64                   `json:"page_number"`
	PageSize    int64                   `json:"page_size"`
	TotalPages  int64                   `json:"total_pages"`
	Restaurants []restaurant.Restaurant `json:"restaurants"`
}
