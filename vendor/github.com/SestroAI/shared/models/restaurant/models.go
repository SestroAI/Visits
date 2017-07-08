package restaurant

import (
	"time"
)

type Restaurant struct {
	ID      string
	Name    string
	Address string
	CuisineId string `json:"cuisineId"`//Cuisine ID
	URL string `json:"url"`
	YelpLink string `json:"yelpLink"`
}

type Cuisine struct {
	ID string
	Name string
}

type Table struct {
	ID string `json:"id"`
	RestaurantId string `json:"restaurantId"`
	OngoingVisitId string `json:"ongoingVisitId"`
	Name string `json:"name"`
}

type Item struct {
	ID string
	RestaurantId string `json:"restaurantId"`
	PriceInDollar float32 `json:"priceInDollar"`
	AverageRating float32 `json:"averageRating"`
}

type Rating struct{
	ID string
	DinerId string `json:"dinerId"`
	Rating int
	Time time.Time
	Comments string
}

type ItemRating struct {
	Rating
	ItemId string `json:"itemId"`
}
