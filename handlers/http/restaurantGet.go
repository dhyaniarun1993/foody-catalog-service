package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dhyaniarun1993/foody-catalog-service/restaurant"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/gorilla/mux"
)

func (handler *restaurantHandler) getRestaurantByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth, _ := authentication.GetAuthFromContext(ctx)
	logger := handler.logger.WithContext(ctx)
	params := mux.Vars(r)
	restaurantID := params["restaurantId"]

	result, serviceError := handler.restaurantInteractor.GetByID(ctx, auth, restaurantID)
	if serviceError != nil {
		logger.WithError(serviceError).Error("Got Error from Service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(serviceError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, serviceError.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (handler *restaurantHandler) getAllRestaurants(w http.ResponseWriter, r *http.Request) {
	var request restaurant.GetAllRestaurantsRequest
	ctx := r.Context()
	auth, _ := authentication.GetAuthFromContext(ctx)
	logger := handler.logger.WithContext(ctx)
	queryParamsData := r.URL.Query()

	decodeError := handler.schemaDecoder.Decode(&request, queryParamsData)
	if decodeError != nil {
		errorMsg := "Invalid request query Params"
		logger.WithError(decodeError).Error(errorMsg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": %q}`, errorMsg)
		return
	}

	result, serviceError := handler.restaurantInteractor.GetAllRestaurants(ctx, auth, request)
	if serviceError != nil {
		logger.WithError(serviceError).Error("Got Error from service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(serviceError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, serviceError.Error())
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
