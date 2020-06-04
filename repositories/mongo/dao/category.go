package dao

import (
	"net/http"
	"time"

	"github.com/dhyaniarun1993/foody-catalog-service/category"
	"github.com/dhyaniarun1993/foody-common/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CategoryDao provides the model definition for category data to be stored in mongodb
type CategoryDao struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RestaurantID primitive.ObjectID `bson:"restaurant_id,omitempty" json:"restaurant_id"`
	Name         string             `bson:"name" json:"name"`
	Description  string             `bson:"description" json:"description"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// GetCategoryDao converts and returns category Dao object from category schema
func GetCategoryDao(category category.Category) (CategoryDao, errors.AppError) {
	categoryDao := CategoryDao{
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	if category.ID != "" {
		categoryObjectID, err := primitive.ObjectIDFromHex(category.ID)
		if err != nil {
			return CategoryDao{}, errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
		}
		categoryDao.ID = categoryObjectID
	}

	// restaurant id is required
	restaurantObjectID, err := primitive.ObjectIDFromHex(category.RestaurantID)
	if err != nil {
		return CategoryDao{}, errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
	}
	categoryDao.RestaurantID = restaurantObjectID

	return categoryDao, nil
}
