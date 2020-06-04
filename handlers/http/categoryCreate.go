package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dhyaniarun1993/foody-common/authentication"

	"github.com/dhyaniarun1993/foody-catalog-service/category"
)

func (handler *categoryHandler) create(w http.ResponseWriter, r *http.Request) {
	var category category.Category
	ctx := r.Context()
	auth, _ := authentication.GetAuthFromContext(ctx)
	logger := handler.logger.WithContext(ctx)

	decodeError := json.NewDecoder(r.Body).Decode(&category)
	if decodeError != nil {
		logger.WithError(decodeError).Error("Invalid request body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": %q }`, decodeError.Error())
		return
	}

	result, serviceError := handler.categoryInteractor.Create(ctx, auth, category)
	if serviceError != nil {
		logger.WithError(serviceError).Error("Got Error from service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(serviceError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, serviceError.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}
