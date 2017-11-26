package menu

import "time"

type Menu struct {
	Meta 	 MenuMeta `json:"meta",mapstructure:"meta,squash"`
	ID       string	`json:"id"`
	Sections map[string]*Section `json:"sections",mapstructure:",squash"`
	Items	 map[string]*Item `json:"items",mapstructure:",squash"`
}

type MenuMeta struct {
	Url 		 string `json:"url"`
	Language     string	`json:"language"`
	LanguageCode string `json:"languageCode",mapstructure:"languageCode"`
	LastUpdated  string `json:"lastUpdated",mapstructure:"lastUpdated"`
}

type Section struct {
	//e.g. Appetizers, Desserts etc
	ID          string `json:"id"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Meta 		ItemMeta    `json:"meta",mapstructure:",squash"`
}

type ItemMeta struct {
	//All the fields to make menu smart will come here
	IsOutOfStock bool      `json:"isOutOfStock",mapstructure:"isOutOfStock"`
	Tags         []string `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
}

type ItemTag struct {
	Name     string   `json:"name"`
	Synonyms []string `json:"synonyms"`
}

type Item struct {
	ID            string	`json:"id"`
	Price         float32  `json:"price"`
	Currency      string   `json:"currency"`
	AverageRating float32  `json:"averageRating"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Images        []string `json:"images"`
	Meta 		  ItemMeta `json:"meta",mapstructure:",squash"`
	SectionId     string   `json:"sectionId"`
}

type Rating struct {
	ID       string
	DinerId  string `json:"dinerId"`
	Rating   int
	Time     time.Time
	Comments string
}

type ItemRating struct {
	Rating
	ItemId string `json:"itemId"`
}