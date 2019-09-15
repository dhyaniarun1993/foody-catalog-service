package services_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-common/logger"
	mockrepositories "github.com/dhyaniarun1993/foody-catalog-service/repositories/mocks"
	"github.com/dhyaniarun1993/foody-catalog-service/services"
	"github.com/stretchr/testify/assert"
)

func TestServiceHealthCheck(t *testing.T) {
	testingTable := []struct {
		name            string
		repositoryError errors.AppError
	}{
		{name: "Success From Service", repositoryError: nil},
		{name: "Error From Service", repositoryError: errors.NewAppError("error from repository", errors.StatusInternalServerError, nil)},
	}

	logger := logger.CreateLogger(logger.Configuration{
		Level:  "INFO",
		Format: "json",
	})

	for _, testingData := range testingTable {
		t.Run(testingData.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockHealthRepository := mockrepositories.NewMockHealthRepository(ctrl)
			mockHealthRepository.EXPECT().HealthCheck(context.TODO()).Return(testingData.repositoryError)

			service := services.NewHealthService(mockHealthRepository, logger)
			serviceError := service.HealthCheck(context.TODO())

			if serviceError != nil {
				assert.Equal(t, testingData.repositoryError.StatusCode(), serviceError.StatusCode(), "Status code should be equal")
				assert.Equal(t, testingData.repositoryError.Error(), serviceError.Error(), "Error message should be equal")
			} else {
				assert.Nil(t, serviceError, "Error retured from service should be nil")
			}
		})
	}
}
