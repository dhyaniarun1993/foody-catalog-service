package health

import (
	"context"

	"github.com/dhyaniarun1993/foody-common/errors"
)

func (interactor *healthInteractor) HealthCheck(ctx context.Context) errors.AppError {
	repositoryError := interactor.healthRepository.HealthCheck(ctx)
	return repositoryError
}
