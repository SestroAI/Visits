package merchant

import (
	"time"
	"github.com/SestroAI/shared/models/merchant/menu"
)

type Merchant struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Cuisine  string `json:"cuisine"` //Cuisine ID
	URL      string `json:"url"`
	YelpLink string `json:"yelpLink"`
	Menu   menu.Menu `json:"menu",mapstructure:"menu,squash"`
	isPaymentSetup bool `json:"isPaymentSetup"`
	Tables 	map[string]bool `json:"tables"`
}

type MerchantStripeInfo struct {
	MerchantId string `json:"merchantId"`
	AccountID string `json:"accountId"`
	PublishableKey string `json:"publishableKey"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type Cuisine struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Table struct {
	ID             string `json:"id"`
	MerchantId     string `json:"merchantId"`
	OngoingVisitId string `json:"ongoingVisitId"`
	Name           string `json:"name"`
}