package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Address provides the schema definition for Address
type Address struct {
	Location GeoJSON `bson:"location" json:"location" validate:"required"`
}

// GeoJSON provides the model definition for Geo Location
type GeoJSON struct {
	Type        string    `bson:"type" json:"-"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}

// Restaurant provides the model definition for Restaurant
type Restaurant struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MerchantID       string             `bson:"merchant_id" json:"merchant_id"`
	Name             string             `bson:"name" json:"name"`
	Description      string             `bson:"description" json:"description"`
	ReviewsRatingSum int64              `bson:"reviews_rating_sum" json:"reviews_rating_sum"`
	ReviewsCount     int64              `bson:"reviews_count" json:"reviews_count"`
	Address          Address            `bson:"address" json:"address"`
	Status           string             `bson:"status" json:"status"`
	IsFeatured       bool               `bson:"is_featured" json:"is_featured"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}
