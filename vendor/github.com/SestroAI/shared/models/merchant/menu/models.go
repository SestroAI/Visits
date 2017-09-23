package menu

import "time"

type Menu struct {
	MenuMeta `json:"meta"`
	ID string
	Sections []*Item `json:"items"`
}

type MenuMeta struct {
	Language string
	LanguageCode string
	LastUpdated string `json:"lastUpdated"`
}

type Section struct {
	//e.g. Appetizers, Desserts etc
	ID string
	Name string `json:"name"`
	Category string `json:"category"`
	Description string 	`json:"description"`
	Items []Item `json:"items"`
	ItemMeta `json:"meta"`
}

type ItemMeta struct {
	//All the fields to make menu smart will come here
	IsOutOfStock bool `json:"isOutOfStock"`
	Tags []ItemTag `json:"tags"`
	Ingredients []string `json:"ingredients"`
}

type ItemTag struct {
	Name string `json:"name"`
	Synonyms []string `json:"synonyms"`
}

type Item struct {
	ID string
	Price float32 `json:"price"`
	Currency string `json:"currency"`
	AverageRating float32 `json:"averageRating"`
	Name string `json:"name"`
	Description string `json:"description"`
	Images []string `json:"images"`
	ItemMeta `json:"meta"`
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