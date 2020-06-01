package controllers

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	services "github.com/dhyaniarun1993/foody-catalog-service/services/mocks"
	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-common/logger"
	"github.com/stretchr/testify/assert"
)

func TestControllerHealthCheck(t *testing.T) {
	testingTable := []struct {
		name                 string
		expectedResponse     string
		expectedStatusCode   int
		serviceResponseError errors.AppError
	}{
		{name: "HealthCheck OK", expectedResponse: ``, expectedStatusCode: http.StatusOK, serviceResponseError: nil},
		{name: "HealthCheck Error From Service", expectedResponse: `{"message": "database down"}`, expectedStatusCode: http.StatusInternalServerError, serviceResponseError: errors.NewAppError("database down", http.StatusInternalServerError, nil)},
	}

	logger := logger.CreateLogger(logger.Configuration{
		Level:  "INFO",
		Format: "json",
	})

	for _, testData := range testingTable {
		t.Run(testData.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/health", nil)
			if err != nil {
				t.Fatalf("Unable to create request: %v", err)
			}

			recorder := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockHealthService := services.NewMockHealthService(ctrl)
			mockHealthService.EXPECT().HealthCheck(context.TODO()).Return(testData.serviceResponseError)

			healthController := &healthController{mockHealthService, logger}

			healthController.healthCheck(recorder, req)

			res := recorder.Result()
			defer res.Body.Close()

			assert.Equal(t, testData.expectedStatusCode, res.StatusCode, "they should be equal")

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response body: %v", err)
			}

			responseBody := string(bytes.TrimSpace(b))
			if res.StatusCode != http.StatusOK {
				assert.JSONEq(t, testData.expectedResponse, responseBody, "they should be equal")
			} else {
				assert.Equal(t, testData.expectedResponse, responseBody, "they should be equal")
			}
		})
	}
}
