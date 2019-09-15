package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"gopkg.in/go-playground/validator.v9"

	"github.com/dhyaniarun1993/foody-common/authentication"
	"github.com/dhyaniarun1993/foody-common/logger"
	"github.com/dhyaniarun1993/foody-common/middlewares"
	"github.com/dhyaniarun1993/foody-catalog-service/schemas/dto"
	"github.com/dhyaniarun1993/foody-catalog-service/services"
)

type restaurantController struct {
	restaurantService services.RestaurantService
	productService    services.ProductService
	logger            *logger.Logger
	validate          *validator.Validate
	schemaDecoder     *schema.Decoder
}

// NewRestaurantController initialize restaurant endpoint
func NewRestaurantController(restaurantService services.RestaurantService,
	productService services.ProductService, logger *logger.Logger, validate *validator.Validate,
	schemaDecoder *schema.Decoder) RestaurantController {

	return &restaurantController{
		restaurantService: restaurantService,
		productService:    productService,
		logger:            logger,
		validate:          validate,
		schemaDecoder:     schemaDecoder,
	}
}

func (controller *restaurantController) LoadRoutes(router *mux.Router) {
	router.Handle("/v1/catalog/restaurants",
		middlewares.ChainHandlerFuncMiddlewares(controller.CreateRestaurant,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("POST")

	router.Handle("/v1/catalog/restaurants/{restaurantId}",
		middlewares.ChainHandlerFuncMiddlewares(controller.GetRestaurant,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("GET")

	router.Handle("/v1/catalog/restaurants/{restaurantId}",
		middlewares.ChainHandlerFuncMiddlewares(controller.DeleteRestaurant,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("DELETE")

	router.Handle("/v1/catalog/restaurants",
		middlewares.ChainHandlerFuncMiddlewares(controller.GetAllRestaurants,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("GET")

	router.Handle("/v1/catalog/restaurants/{restaurantId}/products",
		middlewares.ChainHandlerFuncMiddlewares(controller.CreateProduct,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("POST")

	router.Handle("/v1/catalog/restaurants/{restaurantId}/products/{productId}",
		middlewares.ChainHandlerFuncMiddlewares(controller.GetProduct,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("GET")

	router.Handle("/v1/catalog/restaurants/{restaurantId}/products/{productId}",
		middlewares.ChainHandlerFuncMiddlewares(controller.DeleteProduct,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("DELETE")

	router.Handle("/v1/catalog/restaurants/{restaurantId}/products",
		middlewares.ChainHandlerFuncMiddlewares(controller.GetAllProducts,
			authentication.AuthHandler(), middlewares.TimeoutHandler(2*time.Second))).Methods("GET")
}

func (controller *restaurantController) CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var requestBody dto.CreateRestaurantRequestBody
	var request dto.CreateRestaurantRequest
	ctx := r.Context()
	request.UserID, _ = authentication.GetUserID(ctx)
	request.UserRole, _ = authentication.GetUserRole(ctx)
	request.AppID, _ = authentication.GetAppID(ctx)
	logger := controller.logger.WithContext(ctx)

	decodingError := json.NewDecoder(r.Body).Decode(&requestBody)
	if decodingError != nil {
		errorMsg := "Invalid request"
		logger.WithError(decodingError).Error(errorMsg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": %q}`, errorMsg)
		return
	}

	request.Body = requestBody
	validationError := request.Validate(controller.validate)
	if validationError != nil {
		logger.WithError(validationError).Error("Invalid request body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(validationError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, validationError.Error())
		return
	}

	result, serviceError := controller.restaurantService.Create(ctx, request)
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

func (controller *restaurantController) GetRestaurant(w http.ResponseWriter, r *http.Request) {
	var request dto.GetRestaurantRequest
	ctx := r.Context()
	request.UserID, _ = authentication.GetUserID(ctx)
	request.UserRole, _ = authentication.GetUserRole(ctx)
	request.AppID, _ = authentication.GetAppID(ctx)
	logger := controller.logger.WithContext(ctx)
	params := mux.Vars(r)
	request.Param.RestaurantID = params["restaurantId"]

	validationError := request.Validate(controller.validate)
	if validationError != nil {
		logger.WithError(validationError).Error("Invalid request query Params")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(validationError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, validationError.Error())
		return
	}

	result, serviceError := controller.restaurantService.Get(ctx, request)
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

func (controller *restaurantController) DeleteRestaurant(w http.ResponseWriter, r *http.Request) {
	var request dto.DeleteRestaurantRequest
	ctx := r.Context()
	request.UserID, _ = authentication.GetUserID(ctx)
	request.UserRole, _ = authentication.GetUserRole(ctx)
	request.AppID, _ = authentication.GetAppID(ctx)
	logger := controller.logger.WithContext(ctx)
	params := mux.Vars(r)
	request.Param.RestaurantID = params["restaurantId"]

	validationError := request.Validate(controller.validate)
	if validationError != nil {
		logger.WithError(validationError).Error("Invalid request query Params")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(validationError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, validationError.Error())
		return
	}

	serviceError := controller.restaurantService.Delete(ctx, request)
	if serviceError != nil {
		logger.WithError(serviceError).Error("Got Error from Service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(serviceError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, serviceError.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (controller *restaurantController) GetAllRestaurants(w http.ResponseWriter, r *http.Request) {
	var queryParams dto.GetAllRestaurantsRequestQuery
	var request dto.GetAllRestaurantsRequest
	ctx := r.Context()
	request.UserID, _ = authentication.GetUserID(ctx)
	request.UserRole, _ = authentication.GetUserRole(ctx)
	request.AppID, _ = authentication.GetAppID(ctx)
	logger := controller.logger.WithContext(ctx)
	queryParamsData := r.URL.Query()

	decodeError := controller.schemaDecoder.Decode(&queryParams, queryParamsData)
	if decodeError != nil {
		errorMsg := "Invalid request query Params"
		logger.WithError(decodeError).Error(errorMsg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": %q}`, errorMsg)
		return
	}

	request.Query = queryParams
	validationError := request.Validate(controller.validate)
	if validationError != nil {
		logger.WithError(validationError).Error("Invalid request query Params")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(validationError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, validationError.Error())
		return
	}

	result, serviceError := controller.restaurantService.GetAllRestaurants(ctx, request)
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

func (controller *restaurantController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var requestBody dto.CreateProductRequestBody
	var request dto.CreateProductRequest
	ctx := r.Context()
	request.UserID, _ = authentication.GetUserID(ctx)
	request.UserRole, _ = authentication.GetUserRole(ctx)
	request.AppID, _ = authentication.GetAppID(ctx)
	logger := controller.logger.WithContext(ctx)
	params := mux.Vars(r)
	request.Param.RestaurantID = params["restaurantId"]

	decodeError := json.NewDecoder(r.Body).Decode(&requestBody)
	if decodeError != nil {
		logger.WithError(decodeError).Error("Invalid request body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": %q}`, decodeError.Error())
		return
	}

	request.Body = requestBody
	validationError := request.Validate(controller.validate)
	if validationError != nil {
		logger.WithError(validationError).Error("Invalid request body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(validationError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, validationError.Error())
		return
	}

	result, serviceError := controller.productService.Create(ctx, request)
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

func (controller *restaurantController) GetProduct(w http.ResponseWriter, r *http.Request) {
	var request dto.GetProductRequest
	ctx := r.Context()
	request.UserID, _ = authentication.GetUserID(ctx)
	request.UserRole, _ = authentication.GetUserRole(ctx)
	request.AppID, _ = authentication.GetAppID(ctx)
	logger := controller.logger.WithContext(ctx)
	params := mux.Vars(r)
	request.Param.ProductID = params["productId"]
	request.Param.RestaurantID = params["restaurantId"]

	validationError := request.Validate(controller.validate)
	if validationError != nil {
		logger.WithError(validationError).Error("Invalid request query Params")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(validationError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, validationError.Error())
		return
	}

	result, serviceError := controller.productService.Get(ctx, request)
	if serviceError != nil {
		logger.WithError(serviceError).Error("Got Error from service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(serviceError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, serviceError.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (controller *restaurantController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var request dto.DeleteProductRequest
	ctx := r.Context()
	request.UserID, _ = authentication.GetUserID(ctx)
	request.UserRole, _ = authentication.GetUserRole(ctx)
	request.AppID, _ = authentication.GetAppID(ctx)
	logger := controller.logger.WithContext(ctx)
	params := mux.Vars(r)
	request.Param.ProductID = params["productId"]
	request.Param.RestaurantID = params["restaurantId"]

	validationError := request.Validate(controller.validate)
	if validationError != nil {
		logger.WithError(validationError).Error("Invalid request query Params")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(validationError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, validationError.Error())
		return
	}

	serviceError := controller.productService.Delete(ctx, request)
	if serviceError != nil {
		logger.WithError(serviceError).Error("Got Error from service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(serviceError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, serviceError.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (controller *restaurantController) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	var queryParams dto.GetAllProductsRequestQuery
	var request dto.GetAllProductsRequest
	ctx := r.Context()
	request.UserID, _ = authentication.GetUserID(ctx)
	request.UserRole, _ = authentication.GetUserRole(ctx)
	request.AppID, _ = authentication.GetAppID(ctx)
	params := mux.Vars(r)
	request.Param.RestaurantID = params["restaurantId"]
	logger := controller.logger.WithContext(ctx)
	queryParamsData := r.URL.Query()

	decodeError := controller.schemaDecoder.Decode(&queryParams, queryParamsData)
	if decodeError != nil {
		errorMsg := "Invalid request query Params"
		logger.WithError(decodeError).Error(errorMsg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": %q}`, errorMsg)
		return
	}

	request.Query = queryParams
	validationError := request.Validate(controller.validate)
	if validationError != nil {
		logger.WithError(validationError).Error("Invalid request query Params")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(validationError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, validationError.Error())
		return
	}

	result, serviceError := controller.productService.GetAllProducts(ctx, request)
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
