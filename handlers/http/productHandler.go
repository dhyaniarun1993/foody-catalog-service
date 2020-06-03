package http

import (
	"time"

	"github.com/dhyaniarun1993/foody-catalog-service/product"
	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/logger"
	"github.com/dhyaniarun1993/foody-common/middlewares"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"gopkg.in/go-playground/validator.v9"
)

type productHandler struct {
	productInteractor product.Interactor
	logger            *logger.Logger
	validate          *validator.Validate
	schemaDecoder     *schema.Decoder
}

// NewProductHandler initialize product endpoint
func NewProductHandler(productInteractor product.Interactor, logger *logger.Logger,
	validate *validator.Validate, schemaDecoder *schema.Decoder) Handler {

	return &productHandler{
		productInteractor: productInteractor,
		logger:            logger,
		validate:          validate,
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

}
