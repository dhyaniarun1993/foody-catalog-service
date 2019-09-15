package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// HealthController provides interface for health controller
type HealthController interface {
	LoadRoutes(*mux.Router)
	HealthCheck(http.ResponseWriter, *http.Request)
}

// RestaurantController provides interface for restaurant controller
type RestaurantController interface {
	LoadRoutes(*mux.Router)
	CreateRestaurant(http.ResponseWriter, *http.Request)
	GetRestaurant(http.ResponseWriter, *http.Request)
	DeleteRestaurant(http.ResponseWriter, *http.Request)
	GetAllRestaurants(http.ResponseWriter, *http.Request)
	CreateProduct(http.ResponseWriter, *http.Request)
	GetProduct(http.ResponseWriter, *http.Request)
	DeleteProduct(http.ResponseWriter, *http.Request)
	GetAllProducts(http.ResponseWriter, *http.Request)
}
