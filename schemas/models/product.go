package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product provides the model definition for Product
type Product struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RestaurantID     primitive.ObjectID `bson:"restaurant_id" json:"restaurant_id"`
	Name             string             `bson:"name" json:"name"`
	Description      string             `bson:"description" json:"description"`
	Price            float64            `bson:"price" json:"price"`
	DiscountType     string             `bson:"discount_type" json:"discount_type"`
	Discount         float64            `bson:"discount" json:"discount"`
	ReviewsRatingSum int64              `bson:"reviews_rating_sum" json:"reviews_rating_sum"`
	ReviewsCount     int64              `bson:"reviews_count" json:"reviews_count"`
	Status           string             `bson:"status" json:"status"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}
