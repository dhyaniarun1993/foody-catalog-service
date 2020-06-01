package controllers

import (
	"github.com/gorilla/mux"
)

// HealthController provides interface for health controller
type HealthController interface {
	LoadRoutes(*mux.Router)
}

// RestaurantController provides interface for restaurant controller
type RestaurantController interface {
	LoadRoutes(*mux.Router)
}
