package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"gopkg.in/go-playground/validator.v9"

	"github.com/dhyaniarun1993/foody-catalog-service/acl"
	"github.com/dhyaniarun1993/foody-catalog-service/cmd/catalog-server/config"
	"github.com/dhyaniarun1993/foody-catalog-service/controllers"
	repositories "github.com/dhyaniarun1993/foody-catalog-service/repositories/mongo"
	"github.com/dhyaniarun1993/foody-catalog-service/services"
	"github.com/dhyaniarun1993/foody-common/datastore/mongo"
	"github.com/dhyaniarun1993/foody-common/logger"
	"github.com/dhyaniarun1993/foody-common/tracer"
)

func main() {
	config := config.InitConfiguration()
	logger := logger.CreateLogger(config.Log)
	t, closer := tracer.InitJaeger(config.Jaeger)
	defer closer.Close()

	mongoClient := mongo.CreateMongoDBPool(config.Mongo, t)
	validate := validator.New()
	schemaDecoder := schema.NewDecoder()
	rbac := acl.New()

	healthRepository := repositories.NewHealthRepository(mongoClient)
	restaurantRepository := repositories.NewRestaurantRepository(mongoClient, config.Mongo.Database)
	productRepository := repositories.NewProductRepository(mongoClient, config.Mongo.Database)

	healthService := services.NewHealthService(healthRepository, logger)
	restaurantService := services.NewRestaurantService(restaurantRepository, logger, rbac)
	productService := services.NewProductService(productRepository, restaurantService, logger, rbac)

	router := mux.NewRouter()
	ignoredURLs := []string{"/health1"}
	ignoredMethods := []string{"OPTION"}

	router.Use(tracer.TraceRequest(t, ignoredURLs, ignoredMethods))
	healthController := controllers.NewHealthController(healthService, logger)
	restaurantController := controllers.NewRestaurantController(restaurantService, productService,
		logger, validate, schemaDecoder)

	healthController.LoadRoutes(router)
	restaurantController.LoadRoutes(router)
	serverAddress := ":" + fmt.Sprint(config.Port)
	srv := &http.Server{
		Handler:      router,
		Addr:         serverAddress,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}

	logger.Info("Starting Http server at " + serverAddress)
	serverError := srv.ListenAndServe()
	if serverError != http.ErrServerClosed {
		logger.Error("Http server stopped unexpected")
	} else {
		logger.Info("Http server stopped")
	}
}
