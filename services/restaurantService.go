package services

import (
	"context"
	"math"
	"net/http"
	"reflect"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"

	"github.com/dhyaniarun1993/foody-catalog-service/constants"
	"github.com/dhyaniarun1993/foody-catalog-service/repositories"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/dto"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/models"
	"github.com/dhyaniarun1993/foody-common/async"
	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-common/logger"
)

type restaurantService struct {
	restaurantRepository repositories.RestaurantRepository
	logger               *logger.Logger
	rbac                 acl.RBAC
}

// NewRestaurantService creates and return restaurant service
func NewRestaurantService(restaurantRepository repositories.RestaurantRepository,
	logger *logger.Logger, rbac acl.RBAC) RestaurantService {
	return &restaurantService{
		restaurantRepository: restaurantRepository,
		logger:               logger,
		rbac:                 rbac,
	}
}

func (service *restaurantService) Create(ctx context.Context,
	request dto.CreateRestaurantRequest) (models.Restaurant, errors.AppError) {

	var repositoryError errors.AppError
	if (request.UserID == request.Body.MerchantID && service.rbac.Can(request.UserRole, acl.PermissionCreateRestaurantOwn)) ||
		(service.rbac.Can(request.UserRole, acl.PermissionCreateRestaurantAny)) {
		restaurant := models.Restaurant{
			MerchantID:  request.Body.MerchantID,
			Name:        request.Body.Name,
			Description: request.Body.Description,
			Address: models.Address{
				Location: models.GeoJSON{
					Coordinates: request.Body.Address.Location.Coordinates,
				},
			},
			Status: constants.RestaurantStatusClosed,
		}
		restaurant, repositoryError = service.restaurantRepository.Create(ctx, restaurant)
		return restaurant, repositoryError
	}
	return models.Restaurant{}, errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}

func (service *restaurantService) Get(ctx context.Context,
	request dto.GetRestaurantRequest) (models.Restaurant, errors.AppError) {

	restaurant, repositoryError := service.restaurantRepository.Get(ctx, request.Param.RestaurantID)
	if repositoryError != nil {
		return restaurant, repositoryError
	}

	if reflect.DeepEqual(restaurant, models.Restaurant{}) {
		return restaurant, errors.NewAppError("Resource not found", http.StatusNotFound, nil)
	}

	if (request.UserID == restaurant.MerchantID && service.rbac.Can(request.UserRole, acl.PermissionGetRestaurantOwn)) ||
		service.rbac.Can(request.UserRole, acl.PermissionGetRestaurantAny) {
		return restaurant, nil
	}

	return models.Restaurant{}, errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}

func (service *restaurantService) Delete(ctx context.Context,
	request dto.DeleteRestaurantRequest) errors.AppError {

	restaurant, getError := service.restaurantRepository.Get(ctx, request.Param.RestaurantID)
	if getError != nil {
		return getError
	}

	if reflect.DeepEqual(restaurant, models.Restaurant{}) {
		return errors.NewAppError("Resource not found", http.StatusNotFound, nil)
	}

	if (request.UserID == restaurant.MerchantID && service.rbac.Can(request.UserRole, acl.PermissionDeleteRestaurantOwn)) ||
		service.rbac.Can(request.UserRole, acl.PermissionDeleteRestaurantAny) {

		deleteError := service.restaurantRepository.Delete(ctx, request.Param.RestaurantID)
		return deleteError
	}
	return errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}

func (service *restaurantService) GetAllRestaurants(ctx context.Context,
	request dto.GetAllRestaurantsRequest) (dto.GetAllRestaurantsResponse, errors.AppError) {

	var restaurants []models.Restaurant
	var totalCount, maxDistance int64
	var restaurantResponse dto.GetAllRestaurantsResponse

	maxDistance = 10000
	async, asyncCtx := async.WithContext(ctx)

	if request.Query.PageNumber == 0 {
		request.Query.PageNumber = 1
	}
	if request.Query.PageSize == 0 {
		request.Query.PageSize = 10
	}

	GetAllRestaurants := func() errors.AppError {
		var repositoryError errors.AppError
		restaurants, repositoryError = service.restaurantRepository.GetAllRestaurants(asyncCtx,
			request.Query, maxDistance)
		return repositoryError
	}

	getTotalCount := func() errors.AppError {
		var repositoryError errors.AppError
		totalCount, repositoryError = service.restaurantRepository.GetAllRestaurantsTotalCount(asyncCtx,
			request.Query, maxDistance)
		return repositoryError
	}

	if (request.Query.MerchantID == request.UserID && service.rbac.Can(request.UserRole, acl.PermissionGetRestaurantOwn)) ||
		service.rbac.Can(request.UserRole, acl.PermissionGetRestaurantAny) {

		async.Go(GetAllRestaurants)
		async.Go(getTotalCount)
		err := async.Wait()
		if err != nil {
			return restaurantResponse, err
		}

		restaurantResponse = dto.GetAllRestaurantsResponse{
			Total:       totalCount,
			PageNumber:  request.Query.PageNumber,
			PageSize:    request.Query.PageSize,
			TotalPages:  int64(math.Ceil(float64(totalCount) / float64(request.Query.PageSize))),
			Restaurants: restaurants,
		}
		return restaurantResponse, nil
	}
	return restaurantResponse, errors.NewAppError("Forbidden", http.StatusForbidden, nil)
}
