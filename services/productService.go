package services

import (
	"context"
	"math"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dhyaniarun1993/foody-common/async"
	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-common/logger"
	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-catalog-service/constants"
	"github.com/dhyaniarun1993/foody-catalog-service/repositories"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/dto"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/models"
)

type productService struct {
	productRepository repositories.ProductRepository
	restaurantService RestaurantService
	logger            *logger.Logger
	rbac              acl.RBAC
}

// NewProductService creates and return product service
func NewProductService(productRepository repositories.ProductRepository,
	restaurantService RestaurantService, logger *logger.Logger,
	rbac acl.RBAC) ProductService {
	return &productService{
		productRepository: productRepository,
		restaurantService: restaurantService,
		logger:            logger,
		rbac:              rbac,
	}
}

func (service *productService) Create(ctx context.Context,
	request dto.CreateProductRequest) (models.Product, errors.AppError) {

	getRestaurantRequest := dto.GetRestaurantRequest{
		UserID:   request.UserID,
		UserRole: request.UserRole,
		AppID:    request.AppID,
		Param: dto.GetRestaurantRequestParam{
			RestaurantID: request.Param.RestaurantID,
		},
	}
	restaurant, getRestaurantError := service.restaurantService.Get(ctx, getRestaurantRequest)
	if getRestaurantError != nil {
		return models.Product{}, getRestaurantError
	}

	if (restaurant.MerchantID == request.UserID && service.rbac.Can(request.UserRole, acl.PermissionCreateProductOwn)) ||
		service.rbac.Can(request.UserRole, acl.PermissionCreateProductAny) {
		restaurantObjectID, _ := primitive.ObjectIDFromHex(request.Param.RestaurantID)

		product := models.Product{
			Name:         request.Body.Name,
			RestaurantID: restaurantObjectID,
			Description:  request.Body.Description,
			Price:        request.Body.Price,
			DiscountType: request.Body.DiscountType,
			Discount:     request.Body.Discount,
			Status:       constants.ProductStatusAvailable,
		}

		var createProductError errors.AppError
		product, createProductError = service.productRepository.Create(ctx, product)
		return product, createProductError
	}
	return models.Product{}, errors.NewAppError("Forbidden", errors.StatusForbidden, nil)

}

func (service *productService) Get(ctx context.Context,
	request dto.GetProductRequest) (models.Product, errors.AppError) {

	getRestaurantRequest := dto.GetRestaurantRequest{
		UserID:   request.UserID,
		UserRole: request.UserRole,
		AppID:    request.AppID,
		Param: dto.GetRestaurantRequestParam{
			RestaurantID: request.Param.RestaurantID,
		},
	}
	restaurant, getRestaurantError := service.restaurantService.Get(ctx, getRestaurantRequest)
	if getRestaurantError != nil {
		return models.Product{}, getRestaurantError
	}

	if (restaurant.MerchantID == request.UserID && service.rbac.Can(request.UserRole, acl.PermissionGetProductOwn)) ||
		service.rbac.Can(request.UserRole, acl.PermissionGetProductAny) {

		product, repositoryError := service.productRepository.Get(ctx,
			request.Param.ProductID, request.Param.RestaurantID)
		if repositoryError != nil {
			return models.Product{}, repositoryError
		}

		if reflect.DeepEqual(product, models.Product{}) {
			return models.Product{}, errors.NewAppError("Resource not found", errors.StatusNotFound, nil)
		}
		return product, nil
	}
	return models.Product{}, errors.NewAppError("Forbidden", errors.StatusForbidden, nil)
}

func (service *productService) Delete(ctx context.Context,
	request dto.DeleteProductRequest) errors.AppError {

	getRestaurantRequest := dto.GetRestaurantRequest{
		UserID:   request.UserID,
		UserRole: request.UserRole,
		AppID:    request.AppID,
		Param: dto.GetRestaurantRequestParam{
			RestaurantID: request.Param.RestaurantID,
		},
	}
	restaurant, getRestaurantError := service.restaurantService.Get(ctx, getRestaurantRequest)
	if getRestaurantError != nil {
		return getRestaurantError
	}

	if (restaurant.MerchantID == request.UserID && service.rbac.Can(request.UserRole, acl.PermissionDeleteProductOwn)) ||
		service.rbac.Can(request.UserRole, acl.PermissionDeleteProductAny) {

		repositoryError := service.productRepository.Delete(ctx,
			request.Param.ProductID, request.Param.RestaurantID)
		return repositoryError
	}
	return errors.NewAppError("Forbidden", errors.StatusForbidden, nil)
}

func (service *productService) GetAllProducts(ctx context.Context,
	request dto.GetAllProductsRequest) (dto.GetAllProductsResponse, errors.AppError) {

	var products []models.Product
	var totalCount int64
	var productResponse dto.GetAllProductsResponse

	async, asyncCtx := async.WithContext(ctx)

	if request.Query.PageNumber == 0 {
		request.Query.PageNumber = 1
	}
	if request.Query.PageSize == 0 {
		request.Query.PageSize = 50
	}

	getRestaurantRequest := dto.GetRestaurantRequest{
		UserID:   request.UserID,
		UserRole: request.UserRole,
		AppID:    request.AppID,
		Param: dto.GetRestaurantRequestParam{
			RestaurantID: request.Param.RestaurantID,
		},
	}
	restaurant, getRestaurantError := service.restaurantService.Get(ctx, getRestaurantRequest)
	if getRestaurantError != nil {
		return productResponse, getRestaurantError
	}

	getAllProducts := func() errors.AppError {
		var repositoryError errors.AppError
		products, repositoryError = service.productRepository.GetProductsByRestaurantID(asyncCtx,
			request.Param.RestaurantID, request.Query)
		return repositoryError
	}

	getTotalCount := func() errors.AppError {
		var repositoryError errors.AppError
		totalCount, repositoryError = service.productRepository.GetProductsByRestaurantTotalCount(asyncCtx,
			request.Param.RestaurantID, request.Query)
		return repositoryError
	}

	if (restaurant.MerchantID == request.UserID && service.rbac.Can(request.UserRole, acl.PermissionGetProductOwn)) ||
		service.rbac.Can(request.UserRole, acl.PermissionGetProductAny) {

		async.Go(getAllProducts)
		async.Go(getTotalCount)
		err := async.Wait()
		if err != nil {
			return productResponse, err
		}

		productResponse = dto.GetAllProductsResponse{
			Total:      totalCount,
			PageNumber: request.Query.PageNumber,
			PageSize:   request.Query.PageSize,
			TotalPages: int64(math.Ceil(float64(totalCount) / float64(request.Query.PageSize))),
			Products:   products,
		}
		return productResponse, nil
	}
	return productResponse, errors.NewAppError("Forbidden", errors.StatusForbidden, nil)
}
