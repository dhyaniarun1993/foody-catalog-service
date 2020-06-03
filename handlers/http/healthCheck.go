package http

import (
	"fmt"
	"net/http"
)

func (handler *healthHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := handler.logger.WithContext(ctx)

	serviceError := handler.healthInteractor.HealthCheck(ctx)
	if serviceError != nil {
		logger.WithError(serviceError).Error("Error occurred")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(serviceError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, serviceError.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
