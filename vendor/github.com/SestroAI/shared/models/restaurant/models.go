package restaurant

import (
	"time"
)

type Restaurant struct {
	ID      string
	Name    string
	Address string
	CuisineId string //Cuisine ID
}

type Cuisine struct {
	ID string
	Name string
}

type Table struct {
	ID string
	RestaurantId string
}

type Item struct {
	ID string
	RestaurantId string
	PriceInDollar float32
	AverageRating float32
}

type Rating struct{
	ID string
	DinerId string
	Rating int
	Time time.Time
	Comments string
}

type ItemRating struct {
	Rating
	ItemId string
}
