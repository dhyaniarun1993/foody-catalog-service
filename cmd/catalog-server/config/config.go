package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"

	"github.com/dhyaniarun1993/foody-common/datastore/mongo"
	"github.com/dhyaniarun1993/foody-common/logger"
	"github.com/dhyaniarun1993/foody-common/tracer"
)

// Configuration provides application configuration
type Configuration struct {
	Port   int `required:"true" split_words:"true"`
	Mongo  mongo.Configuration
	Log    logger.Configuration
	Jaeger tracer.Configuration
}

// InitConfiguration initialize the configuration
func InitConfiguration() Configuration {
	var config Configuration
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatalln(err)
	}
	return config
}
