package http

import (
	"time"

	restaurantUsecase "github.com/dhyaniarun1993/foody-catalog-service/restaurant/usecase"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/logger"
	"github.com/dhyaniarun1993/foody-common/middlewares"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type restaurantHandler struct {
	restaurantInteractor restaurantUsecase.Interactor
	logger               *logger.Logger
	schemaDecoder        *schema.Decoder
}

// NewRestaurantHandler initialize restaurant endpoint
func NewRestaurantHandler(restaurantInteractor restaurantUsecase.Interactor, logger *logger.Logger,
	schemaDecoder *schema.Decoder) Handler {

	return &restaurantHandler{
		restaurantInteractor: restaurantInteractor,
		logger:               logger,
		schemaDecoder:        schemaDecoder,
	}
}

func (handler *restaurantHandler) LoadRoutes(router *mux.Router) {
	router.Handle("/v1/catalog/restaurants",
		middlewares.ChainHandlerFuncMiddlewares(handler.createRestaurant,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("POST")

	router.Handle("/v1/catalog/restaurants/{restaurantId}",
		middlewares.ChainHandlerFuncMiddlewares(handler.getRestaurantByID,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("GET")

	router.Handle("/v1/catalog/restaurants/{restaurantId}",
		middlewares.ChainHandlerFuncMiddlewares(handler.deleteRestaurantByID,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("DELETE")

	router.Handle("/v1/catalog/restaurants",
		middlewares.ChainHandlerFuncMiddlewares(handler.getAllRestaurants,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("GET")
}
