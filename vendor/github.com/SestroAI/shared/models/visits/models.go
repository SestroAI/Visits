package visits

import (
	"github.com/SestroAI/shared/utils"
	"time"
	"github.com/SestroAI/shared/models/orders"
)

type Rating struct {
	ReviewerId string `json:"reviewerId"`
	Value int `json:"value"`
	Comments string `json:"comments"`
}

type MerchantVisit struct {
	ID           string            `json:"id"`
	Diners       map[string]string `json:"diners"`
	StartTime    time.Time         `json:"startTime"`
	EndTime      time.Time         `json:"endTime"`
	TableId      string            `json:"tableId"`
	Bills 		 []string          `json:"bills"`
	IsComplete   bool              `json:"isComplete"`
	IsOpenForAll bool              `json:"isOpenForAll"`
	Payer 		 string 		   `json:"payer"`
	GuestRating  *Rating 		   `json:"guestRating"`
	MerchantId 	 string 		   `json:"merchantId"`
	IsPaid		 bool 		       `json:"isPaid"`
}

func NewMerchantVisit(id string) *MerchantVisit {
	visit := MerchantVisit{
		IsComplete:   false,
		StartTime:    time.Now(),
		IsOpenForAll: false,
		Payer: "",
		GuestRating:&Rating{},
	}
	if id == "" {
		id = utils.GenerateUUID()
	}
	visit.ID = id
	visit.Diners = map[string]string{}
	visit.Bills = make([]string, 0)
	visit.IsPaid = false
	visit.StartTime = time.Now()
	return &visit
}

type VisitDinerSession struct {
	ID           string   		  `json:"id"`
	DinerId      string   		  `json:"dinerId"`
	Orders 		 map[string]*orders.Order   `json:"orders",mapstructure:"orders,squash"`
	List 		 []string 		  `json:"list",mapstructure:",squash"`
	MerchantRating *Rating		  `json:"merchantRating",mapstructure:"merchantRating,squash"`
	Payer		 string 		  `json:"payer"`
	VisitId 	 string 		  `json:"visitId"`
	IsPaid		 bool 			  `json:"isPaid"`
	BillId       string           `json:"billId"`
}

func NewVisitDinerSession() *VisitDinerSession {
	sess := VisitDinerSession{}
	sess.ID = utils.GenerateUUID()
	sess.Orders = make(map[string]*orders.Order, 0)
	sess.Payer = ""
	sess.MerchantRating = &Rating{}
	sess.DinerId = ""
	sess.VisitId = ""
	sess.IsPaid = false
	return &sess
}