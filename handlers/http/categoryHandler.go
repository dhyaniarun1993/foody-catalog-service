package http

import (
	"time"

	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/middlewares"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	categoryUsecase "github.com/dhyaniarun1993/foody-catalog-service/category/usecase"
	"github.com/dhyaniarun1993/foody-common/logger"
)

type categoryHandler struct {
	categoryInteractor categoryUsecase.Interactor
	logger             *logger.Logger
	schemaDecoder      *schema.Decoder
}

// NewCategoryHandler initialize category endpoint
func NewCategoryHandler(categoryInteractor categoryUsecase.Interactor,
	logger *logger.Logger, schemaDecoder *schema.Decoder) Handler {
	return &categoryHandler{
		categoryInteractor: categoryInteractor,
		logger:             logger,
		schemaDecoder:      schemaDecoder,
	}
}

func (handler *categoryHandler) LoadRoutes(router *mux.Router) {
	router.Handle("/v1/catalog/categories",
		middlewares.ChainHandlerFuncMiddlewares(handler.create,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("POST")

	router.Handle("/v1/catalog/categories/{categoryId}",
		middlewares.ChainHandlerFuncMiddlewares(handler.getByID,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("GET")

	router.Handle("/v1/catalog/categories/{categoryId}",
		middlewares.ChainHandlerFuncMiddlewares(handler.deleteByID,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("DELETE")
}
