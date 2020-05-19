package services_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang/mock/gomock"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	mockacl "github.com/dhyaniarun1993/foody-catalog-service/acl/mocks"
	"github.com/dhyaniarun1993/foody-catalog-service/constants"
	mockrepositories "github.com/dhyaniarun1993/foody-catalog-service/repositories/mocks"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/dto"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/models"
	"github.com/dhyaniarun1993/foody-catalog-service/services"
	mockservices "github.com/dhyaniarun1993/foody-catalog-service/services/mocks"
	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-common/logger"
)

func toObjectID(id string) primitive.ObjectID {
	objectID, _ := primitive.ObjectIDFromHex(id)
	return objectID
}

func TestProductCreate(t *testing.T) {
	testingTable := []struct {
		name                      string
		request                   dto.CreateProductRequest
		expectedResponse          models.Product
		expectedError             errors.AppError
		createAnyPermissionResult bool
		createOwnPermissionResult bool
		restaurantServiceResponse models.Restaurant
		restaurantServiceError    errors.AppError
		repositoryResponse        models.Product
		repositoryResponseError   errors.AppError
	}{
		{
			name: "Create Product Forbidden Error",
			request: dto.CreateProductRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.CreateProductRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
				},
				Body: dto.CreateProductRequestBody{
					Name:         "Thai Soup",
					Description:  "",
					Price:        150,
					DiscountType: "flat",
					Discount:     50,
				},
			},
			expectedResponse:          models.Product{},
			expectedError:             errors.NewAppError("Forbidden", http.StatusForbidden, nil),
			createAnyPermissionResult: false,
			createOwnPermissionResult: false,
			restaurantServiceResponse: models.Restaurant{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				MerchantID:       "1234",
				Name:             "Pasta place",
				Description:      "",
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Address: models.Address{
					Location: models.GeoJSON{
						Coordinates: []float64{40.714254, -73.961472},
					},
				},
				Status:     constants.RestaurantStatusClosed,
				IsFeatured: false,
				CreatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			restaurantServiceError:  nil,
			repositoryResponse:      models.Product{},
			repositoryResponseError: nil,
		},
		{
			name: "Create Product Error From Restaurant Service",
			request: dto.CreateProductRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.CreateProductRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
				},
				Body: dto.CreateProductRequestBody{
					Name:         "Thai Soup",
					Description:  "",
					Price:        150,
					DiscountType: "flat",
					Discount:     50,
				},
			},
			expectedResponse:          models.Product{},
			expectedError:             errors.NewAppError("Unable to get restaurant data", http.StatusInternalServerError, nil),
			createAnyPermissionResult: false,
			createOwnPermissionResult: false,
			restaurantServiceResponse: models.Restaurant{},
			restaurantServiceError:    errors.NewAppError("Unable to get restaurant data", http.StatusInternalServerError, nil),
			repositoryResponse:        models.Product{},
			repositoryResponseError:   nil,
		},
		{
			name: "Create Product Error From Repository",
			request: dto.CreateProductRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.CreateProductRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
				},
				Body: dto.CreateProductRequestBody{
					Name:         "Thai Soup",
					Description:  "",
					Price:        150,
					DiscountType: "flat",
					Discount:     50,
				},
			},
			expectedResponse:          models.Product{},
			expectedError:             errors.NewAppError("unable to create product", http.StatusInternalServerError, nil),
			createAnyPermissionResult: false,
			createOwnPermissionResult: true,
			restaurantServiceResponse: models.Restaurant{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				MerchantID:       "1234",
				Name:             "Pasta place",
				Description:      "",
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Address: models.Address{
					Location: models.GeoJSON{
						Coordinates: []float64{40.714254, -73.961472},
					},
				},
				Status:     constants.RestaurantStatusClosed,
				IsFeatured: false,
				CreatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			restaurantServiceError:  nil,
			repositoryResponse:      models.Product{},
			repositoryResponseError: errors.NewAppError("unable to create product", http.StatusInternalServerError, nil),
		},
		{
			name: "Create Product Success with Own Permission",
			request: dto.CreateProductRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.CreateProductRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
				},
				Body: dto.CreateProductRequestBody{
					Name:         "Thai Soup",
					Description:  "",
					Price:        150,
					DiscountType: "flat",
					Discount:     50,
				},
			},
			expectedResponse: models.Product{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				Name:             "Some Product",
				Description:      "",
				Price:            150,
				DiscountType:     constants.ProductDiscountTypeFlat,
				Discount:         100,
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Status:           constants.ProductStatusAvailable,
				CreatedAt:        time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:        time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			expectedError:             nil,
			createAnyPermissionResult: false,
			createOwnPermissionResult: true,
			restaurantServiceResponse: models.Restaurant{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				MerchantID:       "1234",
				Name:             "Pasta place",
				Description:      "",
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Address: models.Address{
					Location: models.GeoJSON{
						Coordinates: []float64{40.714254, -73.961472},
					},
				},
				Status:     constants.RestaurantStatusClosed,
				IsFeatured: false,
				CreatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			restaurantServiceError: nil,
			repositoryResponse: models.Product{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				Name:             "Some Product",
				Description:      "",
				Price:            150,
				DiscountType:     constants.ProductDiscountTypeFlat,
				Discount:         100,
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Status:           constants.ProductStatusAvailable,
				CreatedAt:        time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:        time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			repositoryResponseError: nil,
		},
		{
			name: "Create Product Success with Any Permission",
			request: dto.CreateProductRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.CreateProductRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
				},
				Body: dto.CreateProductRequestBody{
					Name:         "Thai Soup",
					Description:  "",
					Price:        150,
					DiscountType: "flat",
					Discount:     50,
				},
			},
			expectedResponse: models.Product{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				Name:             "Some Product",
				Description:      "",
				Price:            150,
				DiscountType:     constants.ProductDiscountTypeFlat,
				Discount:         100,
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Status:           constants.ProductStatusAvailable,
				CreatedAt:        time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:        time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			expectedError:             nil,
			createAnyPermissionResult: true,
			createOwnPermissionResult: false,
			restaurantServiceResponse: models.Restaurant{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				MerchantID:       "1234",
				Name:             "Pasta place",
				Description:      "",
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Address: models.Address{
					Location: models.GeoJSON{
						Coordinates: []float64{40.714254, -73.961472},
					},
				},
				Status:     constants.RestaurantStatusClosed,
				IsFeatured: false,
				CreatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			restaurantServiceError: nil,
			repositoryResponse: models.Product{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				Name:             "Some Product",
				Description:      "",
				Price:            150,
				DiscountType:     constants.ProductDiscountTypeFlat,
				Discount:         100,
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Status:           constants.ProductStatusAvailable,
				CreatedAt:        time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:        time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			repositoryResponseError: nil,
		},
	}

	logger := logger.CreateLogger(logger.Configuration{
		Level:  "INFO",
		Format: "json",
	})

	for _, testData := range testingTable {
		t.Run(testData.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockProductRepository := mockrepositories.NewMockProductRepository(ctrl)
			mockRestaurantService := mockservices.NewMockRestaurantService(ctrl)
			mockRBAC := mockacl.NewMockRBAC(ctrl)

			service := services.NewProductService(mockProductRepository, mockRestaurantService,
				logger, mockRBAC)

			getRestaurantRequest := dto.GetRestaurantRequest{
				UserID:   testData.request.UserID,
				UserRole: testData.request.UserRole,
				AppID:    testData.request.AppID,
				Param: dto.GetRestaurantRequestParam{
					RestaurantID: testData.request.Param.RestaurantID,
				},
			}
			mockRestaurantService.EXPECT().Get(context.TODO(), getRestaurantRequest).Return(testData.restaurantServiceResponse, testData.restaurantServiceError)
			if testData.restaurantServiceError == nil {
				if testData.restaurantServiceResponse.MerchantID == testData.request.UserID {
					first := mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionCreateProductOwn).Return(testData.createOwnPermissionResult)
					if !testData.createOwnPermissionResult {
						second := mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionCreateProductAny).Return(testData.createAnyPermissionResult)
						gomock.InOrder(
							first,
							second,
						)
					}
				} else {
					mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionCreateProductAny).Return(testData.createAnyPermissionResult)
				}
				if testData.createAnyPermissionResult || testData.createOwnPermissionResult {
					repositoryRequest := models.Product{
						Name:         testData.request.Body.Name,
						RestaurantID: testData.restaurantServiceResponse.ID,
						Description:  testData.request.Body.Description,
						Price:        testData.request.Body.Price,
						DiscountType: testData.request.Body.DiscountType,
						Discount:     testData.request.Body.Discount,
						Status:       constants.ProductStatusAvailable,
					}
					mockProductRepository.EXPECT().Create(context.TODO(), repositoryRequest).Return(testData.repositoryResponse, testData.repositoryResponseError)
				}
			}
			response, err := service.Create(context.TODO(), testData.request)
			assert.Equal(t, testData.expectedResponse, response, "they should be equal")
			if err != nil || testData.expectedError != nil {
				assert.Equal(t, testData.expectedError.StatusCode(), err.StatusCode(), "they should be equal")
				assert.Equal(t, testData.expectedError.Error(), err.Error(), "they should be equal")
			}
		})
	}
}

func TestProductGet(t *testing.T) {
	testingTable := []struct {
		name                      string
		request                   dto.GetProductRequest
		expectedResponse          models.Product
		expectedError             errors.AppError
		getAnyPermissionResult    bool
		getOwnPermissionResult    bool
		restaurantServiceResponse models.Restaurant
		restaurantServiceError    errors.AppError
		repositoryResponse        models.Product
		repositoryResponseError   errors.AppError
	}{
		{
			name: "Get Product Forbidden Error",
			request: dto.GetProductRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.GetProductRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
					ProductID:    "5d78d4975eff2a81dda94819",
				},
			},
			expectedResponse:       models.Product{},
			expectedError:          errors.NewAppError("Forbidden", http.StatusForbidden, nil),
			getAnyPermissionResult: false,
			getOwnPermissionResult: false,
			restaurantServiceResponse: models.Restaurant{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				MerchantID:       "123",
				Name:             "Pasta place",
				Description:      "",
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Address: models.Address{
					Location: models.GeoJSON{
						Coordinates: []float64{40.714254, -73.961472},
					},
				},
				Status:     constants.RestaurantStatusClosed,
				IsFeatured: false,
				CreatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			restaurantServiceError:  nil,
			repositoryResponse:      models.Product{},
			repositoryResponseError: nil,
		},
		{
			name: "Get Product Error From Restaurant Service",
			request: dto.GetProductRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.GetProductRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
					ProductID:    "5d78d4975eff2a81dda94819",
				},
			},
			expectedResponse:          models.Product{},
			expectedError:             errors.NewAppError("Unable to get restaurant data", http.StatusInternalServerError, nil),
			getAnyPermissionResult:    false,
			getOwnPermissionResult:    false,
			restaurantServiceResponse: models.Restaurant{},
			restaurantServiceError:    errors.NewAppError("Unable to get restaurant data", http.StatusInternalServerError, nil),
			repositoryResponse:        models.Product{},
			repositoryResponseError:   nil,
		},
		{
			name: "Get Product Error From Repository",
			request: dto.GetProductRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.GetProductRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
					ProductID:    "5d78d4975eff2a81dda94819",
				},
			},
			expectedResponse:       models.Product{},
			expectedError:          errors.NewAppError("Unable to get product data", http.StatusInternalServerError, nil),
			getAnyPermissionResult: false,
			getOwnPermissionResult: true,
			restaurantServiceResponse: models.Restaurant{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				MerchantID:       "1234",
				Name:             "Pasta place",
				Description:      "",
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Address: models.Address{
					Location: models.GeoJSON{
						Coordinates: []float64{40.714254, -73.961472},
					},
				},
				Status:     constants.RestaurantStatusClosed,
				IsFeatured: false,
				CreatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			restaurantServiceError:  nil,
			repositoryResponse:      models.Product{},
			repositoryResponseError: errors.NewAppError("Unable to get product data", http.StatusInternalServerError, nil),
		},
		{
			name: "Get Product Error Not Found",
			request: dto.GetProductRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.GetProductRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
					ProductID:    "5d78d4975eff2a81dda94819",
				},
			},
			expectedResponse:       models.Product{},
			expectedError:          errors.NewAppError("Resource not found", http.StatusNotFound, nil),
			getAnyPermissionResult: false,
			getOwnPermissionResult: true,
			restaurantServiceResponse: models.Restaurant{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				MerchantID:       "1234",
				Name:             "Pasta place",
				Description:      "",
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Address: models.Address{
					Location: models.GeoJSON{
						Coordinates: []float64{40.714254, -73.961472},
					},
				},
				Status:     constants.RestaurantStatusClosed,
				IsFeatured: false,
				CreatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			restaurantServiceError:  nil,
			repositoryResponse:      models.Product{},
			repositoryResponseError: nil,
		},
		{
			name: "Get Product Success With Own Permission",
			request: dto.GetProductRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.GetProductRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
					ProductID:    "5d78d4975eff2a81dda94819",
				},
			},
			expectedResponse: models.Product{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				Name:             "Some Product",
				Description:      "",
				Price:            150,
				DiscountType:     constants.ProductDiscountTypeFlat,
				Discount:         100,
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Status:           constants.ProductStatusAvailable,
				CreatedAt:        time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:        time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			expectedError:          nil,
			getAnyPermissionResult: false,
			getOwnPermissionResult: true,
			restaurantServiceResponse: models.Restaurant{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				MerchantID:       "1234",
				Name:             "Pasta place",
				Description:      "",
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Address: models.Address{
					Location: models.GeoJSON{
						Coordinates: []float64{40.714254, -73.961472},
					},
				},
				Status:     constants.RestaurantStatusClosed,
				IsFeatured: false,
				CreatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			restaurantServiceError: nil,
			repositoryResponse: models.Product{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				Name:             "Some Product",
				Description:      "",
				Price:            150,
				DiscountType:     constants.ProductDiscountTypeFlat,
				Discount:         100,
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Status:           constants.ProductStatusAvailable,
				CreatedAt:        time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:        time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			repositoryResponseError: nil,
		},
		{
			name: "Get Product Success With Any Permission",
			request: dto.GetProductRequest{
				UserID:   "123",
				UserRole: "customer",
				AppID:    "com.foody.customer",
				Param: dto.GetProductRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
					ProductID:    "5d78d4975eff2a81dda94819",
				},
			},
			expectedResponse: models.Product{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				Name:             "Some Product",
				Description:      "",
				Price:            150,
				DiscountType:     constants.ProductDiscountTypeFlat,
				Discount:         100,
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Status:           constants.ProductStatusAvailable,
				CreatedAt:        time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:        time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			expectedError:          nil,
			getAnyPermissionResult: true,
			getOwnPermissionResult: false,
			restaurantServiceResponse: models.Restaurant{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				MerchantID:       "1234",
				Name:             "Pasta place",
				Description:      "",
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Address: models.Address{
					Location: models.GeoJSON{
						Coordinates: []float64{40.714254, -73.961472},
					},
				},
				Status:     constants.RestaurantStatusClosed,
				IsFeatured: false,
				CreatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			restaurantServiceError: nil,
			repositoryResponse: models.Product{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				Name:             "Some Product",
				Description:      "",
				Price:            150,
				DiscountType:     constants.ProductDiscountTypeFlat,
				Discount:         100,
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Status:           constants.ProductStatusAvailable,
				CreatedAt:        time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:        time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			repositoryResponseError: nil,
		},
	}

	logger := logger.CreateLogger(logger.Configuration{
		Level:  "INFO",
		Format: "json",
	})

	for _, testData := range testingTable {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockProductRepository := mockrepositories.NewMockProductRepository(ctrl)
		mockRestaurantService := mockservices.NewMockRestaurantService(ctrl)
		mockRBAC := mockacl.NewMockRBAC(ctrl)

		service := services.NewProductService(mockProductRepository, mockRestaurantService,
			logger, mockRBAC)

		getRestaurantRequest := dto.GetRestaurantRequest{
			UserID:   testData.request.UserID,
			UserRole: testData.request.UserRole,
			AppID:    testData.request.AppID,
			Param: dto.GetRestaurantRequestParam{
				RestaurantID: testData.request.Param.RestaurantID,
			},
		}
		mockRestaurantService.EXPECT().Get(context.TODO(), getRestaurantRequest).Return(testData.restaurantServiceResponse, testData.restaurantServiceError)
		if testData.restaurantServiceError == nil {
			if testData.restaurantServiceResponse.MerchantID == testData.request.UserID {
				first := mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionGetProductOwn).Return(testData.getOwnPermissionResult)
				if !testData.getOwnPermissionResult {
					second := mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionGetProductAny).Return(testData.getAnyPermissionResult)
					gomock.InOrder(
						first,
						second,
					)
				}
			} else {
				mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionGetProductAny).Return(testData.getAnyPermissionResult)
			}
			if testData.getAnyPermissionResult || testData.getOwnPermissionResult {
				mockProductRepository.EXPECT().Get(context.TODO(), testData.request.Param.ProductID, testData.request.Param.RestaurantID).Return(testData.repositoryResponse, testData.repositoryResponseError)
			}
		}
		response, err := service.Get(context.TODO(), testData.request)
		assert.Equal(t, testData.expectedResponse, response, "they should be equal")
		if err != nil || testData.expectedError != nil {
			assert.Equal(t, testData.expectedError.StatusCode(), err.StatusCode(), "they should be equal")
			assert.Equal(t, testData.expectedError.Error(), err.Error(), "they should be equal")
		}
	}
}

func TestProductDelete(t *testing.T) {
	testingTable := []struct {
		name                      string
		request                   dto.DeleteProductRequest
		expectedError             errors.AppError
		deleteAnyPermissionResult bool
		deleteOwnPermissionResult bool
		restaurantServiceResponse models.Restaurant
		restaurantServiceError    errors.AppError
		repositoryResponseError   errors.AppError
	}{
		{
			name: "Delete Product Forbidden Error",
			request: dto.DeleteProductRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.DeleteProductRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
					ProductID:    "5d78d4975eff2a81dda94819",
				},
			},
			expectedError:             errors.NewAppError("Forbidden", http.StatusForbidden, nil),
			deleteAnyPermissionResult: false,
			deleteOwnPermissionResult: false,
			restaurantServiceResponse: models.Restaurant{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				MerchantID:       "123",
				Name:             "Pasta place",
				Description:      "",
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Address: models.Address{
					Location: models.GeoJSON{
						Coordinates: []float64{40.714254, -73.961472},
					},
				},
				Status:     constants.RestaurantStatusClosed,
				IsFeatured: false,
				CreatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			restaurantServiceError:  nil,
			repositoryResponseError: nil,
		},
		{
			name: "Delete Product Error From Restaurant Service",
			request: dto.DeleteProductRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.DeleteProductRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
					ProductID:    "5d78d4975eff2a81dda94819",
				},
			},
			expectedError:             errors.NewAppError("Unable to get restaurant data", http.StatusInternalServerError, nil),
			deleteAnyPermissionResult: false,
			deleteOwnPermissionResult: false,
			restaurantServiceResponse: models.Restaurant{},
			restaurantServiceError:    errors.NewAppError("Unable to get restaurant data", http.StatusInternalServerError, nil),
			repositoryResponseError:   nil,
		},
		{
			name: "Delete Product Error From Repository",
			request: dto.DeleteProductRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.DeleteProductRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
					ProductID:    "5d78d4975eff2a81dda94819",
				},
			},
			expectedError:             errors.NewAppError("Unable to get product data", http.StatusInternalServerError, nil),
			deleteAnyPermissionResult: false,
			deleteOwnPermissionResult: true,
			restaurantServiceResponse: models.Restaurant{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				MerchantID:       "1234",
				Name:             "Pasta place",
				Description:      "",
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Address: models.Address{
					Location: models.GeoJSON{
						Coordinates: []float64{40.714254, -73.961472},
					},
				},
				Status:     constants.RestaurantStatusClosed,
				IsFeatured: false,
				CreatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			restaurantServiceError:  nil,
			repositoryResponseError: errors.NewAppError("Unable to get product data", http.StatusInternalServerError, nil),
		},
		{
			name: "Delete Product Success With Own Permission",
			request: dto.DeleteProductRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.DeleteProductRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
					ProductID:    "5d78d4975eff2a81dda94819",
				},
			},
			expectedError:             nil,
			deleteAnyPermissionResult: false,
			deleteOwnPermissionResult: true,
			restaurantServiceResponse: models.Restaurant{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				MerchantID:       "1234",
				Name:             "Pasta place",
				Description:      "",
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Address: models.Address{
					Location: models.GeoJSON{
						Coordinates: []float64{40.714254, -73.961472},
					},
				},
				Status:     constants.RestaurantStatusClosed,
				IsFeatured: false,
				CreatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			restaurantServiceError:  nil,
			repositoryResponseError: nil,
		},
		{
			name: "Delete Product Success With Any Permission",
			request: dto.DeleteProductRequest{
				UserID:   "123",
				UserRole: "customer",
				AppID:    "com.foody.customer",
				Param: dto.DeleteProductRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
					ProductID:    "5d78d4975eff2a81dda94819",
				},
			},
			expectedError:             nil,
			deleteAnyPermissionResult: true,
			deleteOwnPermissionResult: false,
			restaurantServiceResponse: models.Restaurant{
				ID:               toObjectID("5d78d4975eff2a81dda94810"),
				MerchantID:       "1234",
				Name:             "Pasta place",
				Description:      "",
				ReviewsRatingSum: 0,
				ReviewsCount:     0,
				Address: models.Address{
					Location: models.GeoJSON{
						Coordinates: []float64{40.714254, -73.961472},
					},
				},
				Status:     constants.RestaurantStatusClosed,
				IsFeatured: false,
				CreatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
				UpdatedAt:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			restaurantServiceError:  nil,
			repositoryResponseError: nil,
		},
	}

	logger := logger.CreateLogger(logger.Configuration{
		Level:  "INFO",
		Format: "json",
	})

	for _, testData := range testingTable {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockProductRepository := mockrepositories.NewMockProductRepository(ctrl)
		mockRestaurantService := mockservices.NewMockRestaurantService(ctrl)
		mockRBAC := mockacl.NewMockRBAC(ctrl)

		service := services.NewProductService(mockProductRepository, mockRestaurantService,
			logger, mockRBAC)

		getRestaurantRequest := dto.GetRestaurantRequest{
			UserID:   testData.request.UserID,
			UserRole: testData.request.UserRole,
			AppID:    testData.request.AppID,
			Param: dto.GetRestaurantRequestParam{
				RestaurantID: testData.request.Param.RestaurantID,
			},
		}
		mockRestaurantService.EXPECT().Get(context.TODO(), getRestaurantRequest).Return(testData.restaurantServiceResponse, testData.restaurantServiceError)
		if testData.restaurantServiceError == nil {
			if testData.restaurantServiceResponse.MerchantID == testData.request.UserID {
				first := mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionDeleteProductOwn).Return(testData.deleteOwnPermissionResult)
				if !testData.deleteOwnPermissionResult {
					second := mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionDeleteProductAny).Return(testData.deleteAnyPermissionResult)
					gomock.InOrder(
						first,
						second,
					)
				}
			} else {
				mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionDeleteProductAny).Return(testData.deleteAnyPermissionResult)
			}
			if testData.deleteAnyPermissionResult || testData.deleteOwnPermissionResult {
				mockProductRepository.EXPECT().Delete(context.TODO(), testData.request.Param.ProductID, testData.request.Param.RestaurantID).Return(testData.repositoryResponseError)
			}
		}
		err := service.Delete(context.TODO(), testData.request)
		if err != nil || testData.expectedError != nil {
			assert.Equal(t, testData.expectedError.StatusCode(), err.StatusCode(), "they should be equal")
			assert.Equal(t, testData.expectedError.Error(), err.Error(), "they should be equal")
		}
	}
}
