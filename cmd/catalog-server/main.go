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
	httpHandler "github.com/dhyaniarun1993/foody-catalog-service/handlers/http"
	"github.com/dhyaniarun1993/foody-catalog-service/health"
	"github.com/dhyaniarun1993/foody-catalog-service/product"
	repositories "github.com/dhyaniarun1993/foody-catalog-service/repositories/mongo"
	"github.com/dhyaniarun1993/foody-catalog-service/restaurant"
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

	healthInteractor := health.NewHealthInteractor(healthRepository, logger)
	restaurantInteractor := restaurant.NewRestaurantInteractor(restaurantRepository, logger, rbac, validate)
	productInteractor := product.NewProductInteractor(productRepository, restaurantInteractor, logger, rbac, validate)

	router := mux.NewRouter()
	ignoredURLs := []string{"/health"}
	ignoredMethods := []string{"OPTION"}

	router.Use(tracer.TraceRequest(t, ignoredURLs, ignoredMethods))
	healthHandler := httpHandler.NewHealthHandler(healthInteractor, logger)
	restaurantHandler := httpHandler.NewRestaurantHandler(restaurantInteractor,
		logger, validate, schemaDecoder)
	productHandler := httpHandler.NewProductHandler(productInteractor, logger,
		validate, schemaDecoder)

	healthHandler.LoadRoutes(router)
	restaurantHandler.LoadRoutes(router)
	productHandler.LoadRoutes(router)
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
