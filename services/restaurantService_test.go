package services

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/dhyaniarun1993/foody-catalog-service/constants"

	"github.com/stretchr/testify/assert"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"

	"github.com/golang/mock/gomock"

	mockacl "github.com/dhyaniarun1993/foody-catalog-service/acl/mocks"
	mockrepositories "github.com/dhyaniarun1993/foody-catalog-service/repositories/mocks"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/dto"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/models"
	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-common/logger"
)

func TestRestaurantCreate(t *testing.T) {
	testingTable := []struct {
		name                      string
		request                   dto.CreateRestaurantRequest
		expectedResponse          models.Restaurant
		expectedError             errors.AppError
		createAnyPermissionResult bool
		createOwnPermissionResult bool
		repositoryRequest         models.Restaurant
		repositoryResponse        models.Restaurant
		repositoryResponseError   errors.AppError
	}{
		{
			name: "Create Restaurant Forbidden Error",
			request: dto.CreateRestaurantRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Body: dto.CreateRestaurantRequestBody{
					MerchantID: "1234",
					Name:       "Pasta place",
					Address: dto.Address{
						Location: dto.GeoJSON{
							Coordinates: []float64{40.714254, -73.961472},
						},
					},
				},
			},
			expectedResponse:          models.Restaurant{},
			expectedError:             errors.NewAppError("Forbidden", http.StatusForbidden, nil),
			createAnyPermissionResult: false,
			createOwnPermissionResult: false,
			repositoryRequest:         models.Restaurant{},
			repositoryResponse:        models.Restaurant{},
			repositoryResponseError:   nil,
		},
		{
			name: "Create Restaurant Success Using Create Own Permission",
			request: dto.CreateRestaurantRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Body: dto.CreateRestaurantRequestBody{
					MerchantID: "1234",
					Name:       "Pasta place",
					Address: dto.Address{
						Location: dto.GeoJSON{
							Coordinates: []float64{40.714254, -73.961472},
						},
					},
				},
			},
			expectedResponse: models.Restaurant{
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
			expectedError:             nil,
			createAnyPermissionResult: false,
			createOwnPermissionResult: true,
			repositoryRequest: models.Restaurant{
				MerchantID: "1234",
				Name:       "Pasta place",
				Address: models.Address{
					Location: models.GeoJSON{
						Coordinates: []float64{40.714254, -73.961472},
					},
				},
				Status: constants.RestaurantStatusClosed,
			},
			repositoryResponse: models.Restaurant{
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
			repositoryResponseError: nil,
		},
		{
			name: "Create Restaurant Success Using Create Any Permission",
			request: dto.CreateRestaurantRequest{
				UserID:   "1234",
				UserRole: "admin",
				AppID:    "com.foody.merchant",
				Body: dto.CreateRestaurantRequestBody{
					MerchantID: "1234",
					Name:       "Pasta place",
					Address: dto.Address{
						Location: dto.GeoJSON{
							Coordinates: []float64{40.714254, -73.961472},
						},
					},
				},
			},
			expectedResponse: models.Restaurant{
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
			expectedError:             nil,
			createAnyPermissionResult: true,
			createOwnPermissionResult: false,
			repositoryRequest: models.Restaurant{
				MerchantID: "1234",
				Name:       "Pasta place",
				Address: models.Address{
					Location: models.GeoJSON{
						Coordinates: []float64{40.714254, -73.961472},
					},
				},
				Status: constants.RestaurantStatusClosed,
			},
			repositoryResponse: models.Restaurant{
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
			repositoryResponseError: nil,
		},
		{
			name: "Create Restaurant Error From Repository",
			request: dto.CreateRestaurantRequest{
				UserID:   "123",
				UserRole: "admin",
				AppID:    "com.foody.merchant",
				Body: dto.CreateRestaurantRequestBody{
					MerchantID: "1234",
					Name:       "Pasta place",
					Address: dto.Address{
						Location: dto.GeoJSON{
							Coordinates: []float64{40.714254, -73.961472},
						},
					},
				},
			},
			expectedResponse:          models.Restaurant{},
			expectedError:             errors.NewAppError("Unable to insert", http.StatusInternalServerError, nil),
			createAnyPermissionResult: true,
			createOwnPermissionResult: false,
			repositoryRequest: models.Restaurant{
				MerchantID: "1234",
				Name:       "Pasta place",
				Address: models.Address{
					Location: models.GeoJSON{
						Coordinates: []float64{40.714254, -73.961472},
					},
				},
				Status: constants.RestaurantStatusClosed,
			},
			repositoryResponse:      models.Restaurant{},
			repositoryResponseError: errors.NewAppError("Unable to insert", http.StatusInternalServerError, nil),
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
			mockRestaurantRepository := mockrepositories.NewMockRestaurantRepository(ctrl)
			mockRBAC := mockacl.NewMockRBAC(ctrl)

			restaurantService := NewRestaurantService(mockRestaurantRepository,
				logger, mockRBAC)
			if (testData.request.UserID == testData.request.Body.MerchantID && testData.createOwnPermissionResult) || testData.createAnyPermissionResult {
				mockRestaurantRepository.EXPECT().Create(context.TODO(), testData.repositoryRequest).Return(testData.repositoryResponse, testData.repositoryResponseError)
			}
			if testData.request.Body.MerchantID == testData.request.UserID {
				first := mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionCreateRestaurantOwn).Return(testData.createOwnPermissionResult)
				if !testData.createOwnPermissionResult {
					second := mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionCreateRestaurantAny).Return(testData.createAnyPermissionResult)
					gomock.InOrder(
						first,
						second,
					)
				}
			} else {
				mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionCreateRestaurantAny).Return(testData.createAnyPermissionResult)
			}
			response, err := restaurantService.Create(context.TODO(), testData.request)
			assert.Equal(t, testData.expectedResponse, response, "they should be equal")
			if err != nil {
				assert.Equal(t, testData.expectedError.StatusCode(), err.StatusCode(), "they should be equal")
				assert.Equal(t, testData.expectedError.Error(), err.Error(), "they should be equal")
			}
		})
	}
}

func TestRestaurantGet(t *testing.T) {
	testingTable := []struct {
		name                    string
		request                 dto.GetRestaurantRequest
		expectedResponse        models.Restaurant
		expectedError           errors.AppError
		getAnyPermissionResult  bool
		getOwnPermissionResult  bool
		repositoryResponse      models.Restaurant
		repositoryResponseError errors.AppError
	}{
		{
			name: "Get Restaurant Resource Not Found Error",
			request: dto.GetRestaurantRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.GetRestaurantRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
				},
			},
			expectedResponse:        models.Restaurant{},
			expectedError:           errors.NewAppError("Resource not found", http.StatusNotFound, nil),
			getAnyPermissionResult:  false,
			getOwnPermissionResult:  false,
			repositoryResponse:      models.Restaurant{},
			repositoryResponseError: nil,
		},
		{
			name: "Get Restaurant Resource Forbidden Error",
			request: dto.GetRestaurantRequest{
				UserID:   "123",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.GetRestaurantRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
				},
			},
			expectedResponse:       models.Restaurant{},
			expectedError:          errors.NewAppError("Forbidden", http.StatusForbidden, nil),
			getAnyPermissionResult: false,
			getOwnPermissionResult: false,
			repositoryResponse: models.Restaurant{
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
			repositoryResponseError: nil,
		},
		{
			name: "Get Restaurant Resource Success Using Get Any Permission",
			request: dto.GetRestaurantRequest{
				UserID:   "123",
				UserRole: "customer",
				AppID:    "com.foody.customer",
				Param: dto.GetRestaurantRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
				},
			},
			expectedResponse: models.Restaurant{
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
			expectedError:          nil,
			getAnyPermissionResult: true,
			getOwnPermissionResult: false,
			repositoryResponse: models.Restaurant{
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
			repositoryResponseError: nil,
		},
		{
			name: "Get Restaurant Resource Success Using Get Own Permission",
			request: dto.GetRestaurantRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.GetRestaurantRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
				},
			},
			expectedResponse: models.Restaurant{
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
			expectedError:          nil,
			getAnyPermissionResult: false,
			getOwnPermissionResult: true,
			repositoryResponse: models.Restaurant{
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
			repositoryResponseError: nil,
		},
		{
			name: "Get Restaurant Resource Error From Service",
			request: dto.GetRestaurantRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.GetRestaurantRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
				},
			},
			expectedResponse:        models.Restaurant{},
			expectedError:           errors.NewAppError("Unable to get data from DB", http.StatusForbidden, nil),
			getAnyPermissionResult:  false,
			getOwnPermissionResult:  true,
			repositoryResponse:      models.Restaurant{},
			repositoryResponseError: errors.NewAppError("Unable to get data from DB", http.StatusForbidden, nil),
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

			mockRestaurantRepository := mockrepositories.NewMockRestaurantRepository(ctrl)
			mockRBAC := mockacl.NewMockRBAC(ctrl)

			restaurantService := NewRestaurantService(mockRestaurantRepository,
				logger, mockRBAC)

			mockRestaurantRepository.EXPECT().Get(context.TODO(), testData.request.Param.RestaurantID).Return(testData.repositoryResponse, testData.repositoryResponseError)
			if testData.repositoryResponseError == nil && !reflect.DeepEqual(testData.repositoryResponse, models.Restaurant{}) {
				if testData.request.UserID == testData.repositoryResponse.MerchantID {
					first := mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionGetRestaurantOwn).Return(testData.getOwnPermissionResult)
					if !testData.getOwnPermissionResult {
						second := mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionGetRestaurantAny).Return(testData.getAnyPermissionResult)
						gomock.InOrder(
							first,
							second,
						)
					}
				} else {
					mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionGetRestaurantAny).Return(testData.getAnyPermissionResult)
				}
			}
			response, err := restaurantService.Get(context.TODO(), testData.request)
			assert.Equal(t, testData.expectedResponse, response)
			if err != nil || testData.expectedError != nil {
				assert.Equal(t, testData.expectedError.StatusCode(), err.StatusCode(), "they should be equal")
				assert.Equal(t, testData.expectedError.Error(), err.Error(), "they should be equal")
			}
		})
	}
}

func TestRestaurantDelete(t *testing.T) {
	testingTable := []struct {
		name                          string
		request                       dto.DeleteRestaurantRequest
		expectedError                 errors.AppError
		deleteAnyPermissionResult     bool
		deleteOwnPermissionResult     bool
		repositoryGetResponse         models.Restaurant
		repositoryGetResponseError    errors.AppError
		repositoryDeleteResponseError errors.AppError
	}{
		{
			name: "Delete Restaurant Resource Not Found Error",
			request: dto.DeleteRestaurantRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.DeleteRestaurantRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
				},
			},
			expectedError:                 errors.NewAppError("Resource not found", http.StatusNotFound, nil),
			deleteAnyPermissionResult:     false,
			deleteOwnPermissionResult:     false,
			repositoryGetResponse:         models.Restaurant{},
			repositoryGetResponseError:    nil,
			repositoryDeleteResponseError: nil,
		},
		{
			name: "Delete Restaurant Resource Forbidden Error",
			request: dto.DeleteRestaurantRequest{
				UserID:   "123",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.DeleteRestaurantRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
				},
			},
			expectedError:             errors.NewAppError("Forbidden", http.StatusForbidden, nil),
			deleteAnyPermissionResult: false,
			deleteOwnPermissionResult: false,
			repositoryGetResponse: models.Restaurant{
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
			repositoryGetResponseError:    nil,
			repositoryDeleteResponseError: nil,
		},
		{
			name: "Delete Restaurant Resource Success With Own Permission",
			request: dto.DeleteRestaurantRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Param: dto.DeleteRestaurantRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
				},
			},
			expectedError:             nil,
			deleteAnyPermissionResult: false,
			deleteOwnPermissionResult: true,
			repositoryGetResponse: models.Restaurant{
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
			repositoryGetResponseError:    nil,
			repositoryDeleteResponseError: nil,
		},
		{
			name: "Delete Restaurant Resource Success With Any Permission",
			request: dto.DeleteRestaurantRequest{
				UserID:   "123",
				UserRole: "admin",
				AppID:    "com.foody.merchant",
				Param: dto.DeleteRestaurantRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
				},
			},
			expectedError:             nil,
			deleteAnyPermissionResult: true,
			deleteOwnPermissionResult: false,
			repositoryGetResponse: models.Restaurant{
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
			repositoryGetResponseError:    nil,
			repositoryDeleteResponseError: nil,
		},
		{
			name: "Delete Restaurant Resource Error While Getting",
			request: dto.DeleteRestaurantRequest{
				UserID:   "1234",
				UserRole: "admin",
				AppID:    "com.foody.merchant",
				Param: dto.DeleteRestaurantRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
				},
			},
			expectedError:             errors.NewAppError("Unable to get data from database", http.StatusInternalServerError, nil),
			deleteAnyPermissionResult: true,
			deleteOwnPermissionResult: false,
			repositoryGetResponse: models.Restaurant{
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
			repositoryGetResponseError:    nil,
			repositoryDeleteResponseError: errors.NewAppError("Unable to get data from database", http.StatusInternalServerError, nil),
		},
		{
			name: "Delete Restaurant Resource Error While Deleting",
			request: dto.DeleteRestaurantRequest{
				UserID:   "1234",
				UserRole: "admin",
				AppID:    "com.foody.merchant",
				Param: dto.DeleteRestaurantRequestParam{
					RestaurantID: "5d78d4975eff2a81dda94810",
				},
			},
			expectedError:             errors.NewAppError("Unable to delete from database", http.StatusInternalServerError, nil),
			deleteAnyPermissionResult: true,
			deleteOwnPermissionResult: false,
			repositoryGetResponse: models.Restaurant{
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
			repositoryGetResponseError:    nil,
			repositoryDeleteResponseError: errors.NewAppError("Unable to delete from database", http.StatusInternalServerError, nil),
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

			mockRestaurantRepository := mockrepositories.NewMockRestaurantRepository(ctrl)
			mockRBAC := mockacl.NewMockRBAC(ctrl)

			restaurantService := NewRestaurantService(mockRestaurantRepository,
				logger, mockRBAC)

			mockRestaurantRepository.EXPECT().Get(context.TODO(), testData.request.Param.RestaurantID).Return(testData.repositoryGetResponse, testData.repositoryGetResponseError)
			if testData.repositoryGetResponseError != nil || !reflect.DeepEqual(testData.repositoryGetResponse, models.Restaurant{}) {
				if testData.request.UserID == testData.repositoryGetResponse.MerchantID {
					first := mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionDeleteRestaurantOwn).Return(testData.deleteOwnPermissionResult)
					if !testData.deleteOwnPermissionResult {
						second := mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionDeleteRestaurantAny).Return(testData.deleteAnyPermissionResult)
						gomock.InOrder(
							first,
							second,
						)
					}
				} else {
					mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionDeleteRestaurantAny).Return(testData.deleteAnyPermissionResult)
				}
				if testData.deleteAnyPermissionResult || testData.deleteOwnPermissionResult {
					mockRestaurantRepository.EXPECT().Delete(context.TODO(), testData.request.Param.RestaurantID).Return(testData.repositoryDeleteResponseError)
				}
			}

			err := restaurantService.Delete(context.TODO(), testData.request)
			if err != nil || testData.expectedError != nil {
				assert.Equal(t, testData.expectedError.StatusCode(), err.StatusCode(), "they should be equal")
				assert.Equal(t, testData.expectedError.Error(), err.Error(), "they should be equal")
			}
		})
	}
}

func TestGetAllRestaurants(t *testing.T) {
	testingTable := []struct {
		name                                 string
		request                              dto.GetAllRestaurantsRequest
		maxDistance                          int64
		expectedError                        errors.AppError
		expectedResponse                     dto.GetAllRestaurantsResponse
		getAnyPermissionResult               bool
		getOwnPermissionResult               bool
		repositoryGetResponse                []models.Restaurant
		repositoryGetResponseError           errors.AppError
		repositoryGetTotalCountResponse      int64
		repositoryGetTotalCountResponseError errors.AppError
	}{
		{
			name: "Get All Restaurant Resource Forbidden Error",
			request: dto.GetAllRestaurantsRequest{
				UserID:   "123",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Query: dto.GetAllRestaurantsRequestQuery{
					Longitude:  40.714254,
					Latitude:   -73.961472,
					MerchantID: "1234",
				},
			},
			maxDistance:                          10000,
			expectedError:                        errors.NewAppError("Forbidden", http.StatusForbidden, nil),
			expectedResponse:                     dto.GetAllRestaurantsResponse{},
			getAnyPermissionResult:               false,
			getOwnPermissionResult:               false,
			repositoryGetResponse:                []models.Restaurant{},
			repositoryGetResponseError:           nil,
			repositoryGetTotalCountResponse:      0,
			repositoryGetTotalCountResponseError: nil,
		},
		{
			name: "Get All Restaurant Resource Forbidden Error With Merchant ID",
			request: dto.GetAllRestaurantsRequest{
				UserID:   "123",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Query: dto.GetAllRestaurantsRequestQuery{
					Longitude: 40.714254,
					Latitude:  -73.961472,
				},
			},
			maxDistance:                          10000,
			expectedError:                        errors.NewAppError("Forbidden", http.StatusForbidden, nil),
			expectedResponse:                     dto.GetAllRestaurantsResponse{},
			getAnyPermissionResult:               false,
			getOwnPermissionResult:               false,
			repositoryGetResponse:                []models.Restaurant{},
			repositoryGetResponseError:           nil,
			repositoryGetTotalCountResponse:      0,
			repositoryGetTotalCountResponseError: nil,
		},
		{
			name: "Get All Restaurant Resource Success Own Permission",
			request: dto.GetAllRestaurantsRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Query: dto.GetAllRestaurantsRequestQuery{
					Longitude:  40.714254,
					Latitude:   -73.961472,
					MerchantID: "1234",
				},
			},
			maxDistance:   10000,
			expectedError: nil,
			expectedResponse: dto.GetAllRestaurantsResponse{
				Total:      1,
				PageNumber: 1,
				PageSize:   10,
				TotalPages: 1,
				Restaurants: []models.Restaurant{
					{
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
				},
			},
			getAnyPermissionResult: false,
			getOwnPermissionResult: true,
			repositoryGetResponse: []models.Restaurant{
				models.Restaurant{
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
			},
			repositoryGetResponseError:           nil,
			repositoryGetTotalCountResponse:      1,
			repositoryGetTotalCountResponseError: nil,
		},
		{
			name: "Get All Restaurant Resource Success Any Permission",
			request: dto.GetAllRestaurantsRequest{
				UserID:   "1234",
				UserRole: "customer",
				AppID:    "com.foody.customer",
				Query: dto.GetAllRestaurantsRequestQuery{
					Longitude: 40.714254,
					Latitude:  -73.961472,
				},
			},
			maxDistance:   10000,
			expectedError: nil,
			expectedResponse: dto.GetAllRestaurantsResponse{
				Total:      1,
				PageNumber: 1,
				PageSize:   10,
				TotalPages: 1,
				Restaurants: []models.Restaurant{
					{
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
				},
			},
			getAnyPermissionResult: true,
			getOwnPermissionResult: false,
			repositoryGetResponse: []models.Restaurant{
				models.Restaurant{
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
			},
			repositoryGetResponseError:           nil,
			repositoryGetTotalCountResponse:      1,
			repositoryGetTotalCountResponseError: nil,
		},
		{
			name: "Get All Restaurant Resource Success With Empty Result",
			request: dto.GetAllRestaurantsRequest{
				UserID:   "1234",
				UserRole: "customer",
				AppID:    "com.foody.customer",
				Query: dto.GetAllRestaurantsRequestQuery{
					Longitude: 40.714254,
					Latitude:  -73.961472,
				},
			},
			maxDistance:   10000,
			expectedError: nil,
			expectedResponse: dto.GetAllRestaurantsResponse{
				Total:       0,
				PageNumber:  1,
				PageSize:    10,
				TotalPages:  0,
				Restaurants: []models.Restaurant{},
			},
			getAnyPermissionResult:               true,
			getOwnPermissionResult:               false,
			repositoryGetResponse:                []models.Restaurant{},
			repositoryGetResponseError:           nil,
			repositoryGetTotalCountResponse:      0,
			repositoryGetTotalCountResponseError: nil,
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

			mockRestaurantRepository := mockrepositories.NewMockRestaurantRepository(ctrl)
			mockRBAC := mockacl.NewMockRBAC(ctrl)

			restaurantService := NewRestaurantService(mockRestaurantRepository,
				logger, mockRBAC)

			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()
			if testData.request.UserID == testData.request.Query.MerchantID {
				first := mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionGetRestaurantOwn).Return(testData.getOwnPermissionResult)
				if !testData.getOwnPermissionResult {
					second := mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionGetRestaurantAny).Return(testData.getAnyPermissionResult)
					gomock.InOrder(
						first,
						second,
					)
				}
			} else {
				mockRBAC.EXPECT().Can(testData.request.UserRole, acl.PermissionGetRestaurantAny).Return(testData.getAnyPermissionResult)
			}
			if (testData.request.UserID == testData.request.Query.MerchantID && testData.getOwnPermissionResult) || testData.getAnyPermissionResult {
				request := testData.request
				if request.Query.PageNumber == 0 {
					request.Query.PageNumber = 1
				}
				if request.Query.PageSize == 0 {
					request.Query.PageSize = 10
				}
				mockRestaurantRepository.EXPECT().GetAllRestaurants(ctx, request.Query, testData.maxDistance).Return(testData.repositoryGetResponse, testData.repositoryGetResponseError)
				mockRestaurantRepository.EXPECT().GetAllRestaurantsTotalCount(ctx, request.Query, testData.maxDistance).Return(testData.repositoryGetTotalCountResponse, testData.repositoryGetTotalCountResponseError)
			}

			response, err := restaurantService.GetAllRestaurants(context.TODO(), testData.request)
			assert.Equal(t, testData.expectedResponse, response)
			if err != nil || testData.expectedError != nil {
				assert.Equal(t, testData.expectedError.StatusCode(), err.StatusCode(), "they should be equal")
				assert.Equal(t, testData.expectedError.Error(), err.Error(), "they should be equal")
			}
		})
	}
}
