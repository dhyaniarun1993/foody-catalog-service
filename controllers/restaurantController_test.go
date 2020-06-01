package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dhyaniarun1993/foody-catalog-service/constants"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/dto"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/models"
	services "github.com/dhyaniarun1993/foody-catalog-service/services/mocks"
	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-common/logger"
	"github.com/gorilla/schema"
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

func toObjectID(id string) primitive.ObjectID {
	objectID, _ := primitive.ObjectIDFromHex(id)
	return objectID
}

func TestCreateRestaurant(t *testing.T) {
	testingTable := []struct {
		name                 string
		request              dto.CreateRestaurantRequest
		expectedResponse     string
		expectedStatusCode   int
		willServiceBeCalled  bool
		serviceResponse      models.Restaurant
		serviceResponseError errors.AppError
	}{
		{
			name: "Create Restaurant Success",
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
			expectedResponse:    `{ "id": "5d78d4975eff2a81dda94810", "merchant_id": "1234", "name": "Pasta place", "description": "", "reviews_rating_sum": 0, "reviews_count": 0, "address": { "location": { "coordinates": [40.714254, -73.961472] } }, "status": "closed", "is_featured": false, "created_at": "2019-11-17T20:34:58.651387237Z", "updated_at": "2019-11-17T20:34:58.651387237Z"}`,
			expectedStatusCode:  http.StatusCreated,
			willServiceBeCalled: true,
			serviceResponse: models.Restaurant{
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
			serviceResponseError: nil,
		},
		{
			name: "Create Restaurant Validation Error Merchant ID Required",
			request: dto.CreateRestaurantRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Body: dto.CreateRestaurantRequestBody{
					Address: dto.Address{
						Location: dto.GeoJSON{
							Coordinates: []float64{40.714254, -73.961472},
						},
					},
				},
			},
			expectedResponse:     `{ "message": "validation for field 'MerchantID' failed on 'required'" }`,
			expectedStatusCode:   http.StatusBadRequest,
			willServiceBeCalled:  false,
			serviceResponse:      models.Restaurant{},
			serviceResponseError: nil,
		},
		{
			name: "Create Restaurant Validation Error Name Required",
			request: dto.CreateRestaurantRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Body: dto.CreateRestaurantRequestBody{
					MerchantID: "1234",
					Address: dto.Address{
						Location: dto.GeoJSON{
							Coordinates: []float64{40.714254, -73.961472},
						},
					},
				},
			},
			expectedResponse:     `{ "message": "validation for field 'Name' failed on 'required'" }`,
			expectedStatusCode:   http.StatusBadRequest,
			willServiceBeCalled:  false,
			serviceResponse:      models.Restaurant{},
			serviceResponseError: nil,
		},
		{
			name: "Create Restaurant Validation Error Name Length",
			request: dto.CreateRestaurantRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Body: dto.CreateRestaurantRequestBody{
					MerchantID:  "1234",
					Name:        "Very Long name so that validation on name length fails",
					Description: "",
					Address: dto.Address{
						Location: dto.GeoJSON{
							Coordinates: []float64{40.714254, -73.961472},
						},
					},
				},
			},
			expectedResponse:     `{ "message": "validation for field 'Name' failed on 'max'" }`,
			expectedStatusCode:   http.StatusBadRequest,
			willServiceBeCalled:  false,
			serviceResponse:      models.Restaurant{},
			serviceResponseError: nil,
		},
		{
			name: "Create Restaurant Validation Error Coordinates Required",
			request: dto.CreateRestaurantRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Body: dto.CreateRestaurantRequestBody{
					MerchantID:  "1234",
					Name:        "Pasta place",
					Description: "",
					Address: dto.Address{
						Location: dto.GeoJSON{},
					},
				},
			},
			expectedResponse:     `{ "message": "validation for field 'Coordinates' failed on 'required'" }`,
			expectedStatusCode:   http.StatusBadRequest,
			willServiceBeCalled:  false,
			serviceResponse:      models.Restaurant{},
			serviceResponseError: nil,
		},
		{
			name: "Create Restaurant Validation Error Description Length",
			request: dto.CreateRestaurantRequest{
				UserID:   "1234",
				UserRole: "merchant",
				AppID:    "com.foody.merchant",
				Body: dto.CreateRestaurantRequestBody{
					MerchantID:  "1234",
					Name:        "Pasta place",
					Description: "ansflkna,snsam,nsm,snaksnsfa,mnm,xnm,nm,snmxz c,msnsiudkfjhksjfhbsnkjsfhnjsk,mnsjksfmnfskjsfbnsjkfsanfkjnfsjkas",
					Address: dto.Address{
						Location: dto.GeoJSON{
							Coordinates: []float64{40.714254, -73.961472},
						},
					},
				},
			},
			expectedResponse:     `{ "message": "validation for field 'Description' failed on 'max'" }`,
			expectedStatusCode:   http.StatusBadRequest,
			willServiceBeCalled:  false,
			serviceResponse:      models.Restaurant{},
			serviceResponseError: nil,
		},
	}

	logger := logger.CreateLogger(logger.Configuration{
		Level:  "INFO",
		Format: "json",
	})

	for _, testData := range testingTable {
		t.Run(testData.name, func(t *testing.T) {
			requestByte, marshalError := json.Marshal(testData.request.Body)
			if marshalError != nil {
				t.Fatalf("Unable to marshal request: %v", marshalError)
			}
			req, err := http.NewRequest("POST", "/v1/catalog/restaurants", bytes.NewReader(requestByte))
			if err != nil {
				t.Fatalf("Unable to create request: %v", err)
			}

			recorder := httptest.NewRecorder()
			validate := validator.New()
			schemaDecoder := schema.NewDecoder()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockProductService := services.NewMockProductService(ctrl)
			mockRestaurantService := services.NewMockRestaurantService(ctrl)

			if testData.willServiceBeCalled {
				mockRestaurantService.EXPECT().Create(context.TODO(), testData.request).Return(testData.serviceResponse, testData.serviceResponseError)
			}

			restaurantController := &restaurantController{mockRestaurantService,
				mockProductService, logger, validate, schemaDecoder}
			restaurantController.createRestaurant(recorder, req)

			res := recorder.Result()
			defer res.Body.Close()

			assert.Equal(t, testData.expectedStatusCode, res.StatusCode, "they should be equal")

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response body: %v", err)
			}

			responseBody := string(bytes.TrimSpace(b))
			if res.StatusCode == http.StatusCreated {
				assert.JSONEq(t, testData.expectedResponse, responseBody, "they should be equal")
			} else {
				assert.JSONEq(t, testData.expectedResponse, responseBody, "they should be equal")
			}
		})
	}
}
