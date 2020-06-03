package http

import (
	"github.com/gorilla/mux"

	"github.com/dhyaniarun1993/foody-catalog-service/health"
	"github.com/dhyaniarun1993/foody-common/logger"
)

type healthHandler struct {
	healthInteractor health.Interactor
	logger           *logger.Logger
}

// NewHealthHandler initialize health endpoint
func NewHealthHandler(healthService health.Interactor,
	logger *logger.Logger) Handler {
	return &healthHandler{
		healthInteractor: healthService,
		logger:           logger,
	}
}

func (handler *healthHandler) LoadRoutes(router *mux.Router) {
	router.HandleFunc("/health", handler.healthCheck).Methods("GET")
}
