package health

import (
	"context"

	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-common/logger"
)

// healthRepository provides interface for health repositories
type healthRepository interface {
	HealthCheck(context.Context) errors.AppError
}

// Interactor provides interface for health Interactor
type Interactor interface {
	HealthCheck(context.Context) errors.AppError
}

type healthInteractor struct {
	healthRepository healthRepository
	logger           *logger.Logger
}

// NewHealthInteractor creates and return health service object
func NewHealthInteractor(healthRepository healthRepository,
	logger *logger.Logger) Interactor {
	return &healthInteractor{healthRepository, logger}
}
