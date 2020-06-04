package http

import (
	"fmt"
	"net/http"

	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/gorilla/mux"
)

func (handler *categoryHandler) deleteByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth, _ := authentication.GetAuthFromContext(ctx)
	logger := handler.logger.WithContext(ctx)

	params := mux.Vars(r)
	categoryID := params["categoryId"]

	serviceError := handler.categoryInteractor.DeleteByID(ctx, auth, categoryID)
	if serviceError != nil {
		logger.WithError(serviceError).Error("Got error from service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(serviceError.StatusCode())
		fmt.Fprintf(w, `{"message": %q }`, serviceError.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
