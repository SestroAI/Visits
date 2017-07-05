package dao

import (
	"github.com/SestroAI/shared/models/restaurant"
	"github.com/google/logger"
	"errors"
)

type RestaurantDao struct {
	Dao
	BasePath string
}

const RESTAURANT_BASE_PATH = "/restaurants"


func NewRestaurantDao(token string) *RestaurantDao {
	return &RestaurantDao{
		Dao: *NewDao(token),
		BasePath:RESTAURANT_BASE_PATH,
	}
}

func (ref *RestaurantDao) SaveRestaurant(id string, restro restaurant.Restaurant) error {
	err := ref.SaveObjectById(id, restro, ref.BasePath)

	if err != nil {
		logger.Errorf("Unable to save RESTAURANT object with Id = %s", id)
		return err
	}

	return nil
}

func (ref *RestaurantDao) GetRestaurant(id string) (*restaurant.Restaurant, error) {
	object, _ := ref.GetObjectById(id, ref.BasePath)
	if object == nil {
		return nil, errors.New("Unable to get RESTAURANT with id = " + id)
	}

	restro := restaurant.Restaurant{}
	MapToStruct(object.(map[string]interface{}), &restro)
	return &restro, nil
}