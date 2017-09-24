package merchant

type Merchant struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Cuisine  string `json:"cuisine"` //Cuisine ID
	URL      string `json:"url"`
	YelpLink string `json:"yelpLink"`
	MenuID   string `json:"menuId"`
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
