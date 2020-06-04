package http

import (
	"time"

	productUsecase "github.com/dhyaniarun1993/foody-catalog-service/product/usecase"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/logger"
	"github.com/dhyaniarun1993/foody-common/middlewares"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type productHandler struct {
	productInteractor productUsecase.Interactor
	logger            *logger.Logger
	schemaDecoder     *schema.Decoder
}

// NewProductHandler initialize product endpoint
func NewProductHandler(productInteractor productUsecase.Interactor, logger *logger.Logger,
	schemaDecoder *schema.Decoder) Handler {

	return &productHandler{
		productInteractor: productInteractor,
		logger:            logger,
		schemaDecoder:     schemaDecoder,
	}
}

func (handler *productHandler) LoadRoutes(router *mux.Router) {
	router.Handle("/v1/catalog/products",
		middlewares.ChainHandlerFuncMiddlewares(handler.createProduct,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("POST")

	router.Handle("/v1/catalog/products/{productId}",
		middlewares.ChainHandlerFuncMiddlewares(handler.getProductByID,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("GET")

	router.Handle("/v1/catalog/products/{productId}",
		middlewares.ChainHandlerFuncMiddlewares(handler.deleteProductByID,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("DELETE")

	router.Handle("/v1/catalog/products/{productId}/variants",
		middlewares.ChainHandlerFuncMiddlewares(handler.AddVariant,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("POST")

	router.Handle("/v1/catalog/products/{productId}/variants/{variantId}",
		middlewares.ChainHandlerFuncMiddlewares(handler.RemoveVariant,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("DELETE")
}
