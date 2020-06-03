package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dhyaniarun1993/foody-catalog-service/restaurant"
	"github.com/dhyaniarun1993/foody-common/authentication"
)

func (handler *restaurantHandler) createRestaurant(w http.ResponseWriter, r *http.Request) {
	var restaurant restaurant.Restaurant
	ctx := r.Context()
	auth, _ := authentication.GetAuthFromContext(ctx)
	logger := handler.logger.WithContext(ctx)

	decodingError := json.NewDecoder(r.Body).Decode(&restaurant)
	if decodingError != nil {
		errorMsg := "Invalid request"
		logger.WithError(decodingError).Error(errorMsg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": %q}`, errorMsg)
		return
	}

	result, serviceError := handler.restaurantInteractor.Create(ctx, auth, restaurant)
	if serviceError != nil {
		logger.WithError(serviceError).Error("Got Error from Service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(serviceError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, serviceError.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}
