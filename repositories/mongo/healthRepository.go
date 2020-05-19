package mongo

import (
	"context"
	"net/http"
	"time"

	"github.com/dhyaniarun1993/foody-catalog-service/repositories"
	"github.com/dhyaniarun1993/foody-common/datastore/mongo"
	"github.com/dhyaniarun1993/foody-common/errors"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type healthRepository struct {
	*mongo.Client
}

// NewHealthRepository creates and return health repository
func NewHealthRepository(mongoClient *mongo.Client) repositories.HealthRepository {
	return &healthRepository{mongoClient}
}

func (mongo *healthRepository) HealthCheck(ctx context.Context) errors.AppError {
	timedCtx, pingCancel := context.WithTimeout(ctx, 1*time.Second)
	defer pingCancel()
	pingError := mongo.Ping(timedCtx, readpref.Primary())
	if pingError != nil {
		return errors.NewAppError("Unable to connect to MongoDB", http.StatusServiceUnavailable, pingError)
	}
	return nil
}
