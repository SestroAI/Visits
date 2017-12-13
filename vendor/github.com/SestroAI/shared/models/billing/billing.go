package billing

import (
	"github.com/SestroAI/shared/models/orders"
	"github.com/SestroAI/shared/utils"
	"time"
)

/*
	orderMap: {
		"payer1" : {
			"dinerID1" : [order1, order2, ...],
			...
		},
		...
	}
	 */

type UserBill struct {
	ID string `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	ChargedAt time.Time `json:"chargedAt"`
	BilledToUserId string `json:"billedToUserId"`
	TaxAmount float32 `json:"taxAmount"`
	DiscountAmount float32 `json:"discountAmount"`
	TipAmount float32 `json:"tipAmount"`
	Amount float32 `json:"amount"`
	OrderMap map[string][]*orders.Order `json:"orderMap"`
	AssociatedVisitId string `json:"associatedVisitId"`
	StripeChargeId string `json:"stripeChargeId",mapstructure:"stripeChargeId,"`
	IsPaid bool `json:"isPaid"`
	IsRefunded bool `json:"isRefunded"`
	MerchantName string `json:"merchantName"`
	MerchantEmail string `json:"merchantEmail"`
	MerchantPhone string `json:"merchantPhone"`
	MerchantContactPreference string `json:"merchantContactPreference"`
}

func NewUserBill() *UserBill {
	var ub UserBill
	ub.ID = utils.GenerateUUID()
	ub.CreatedAt = time.Now()
	ub.OrderMap = make(map[string][]*orders.Order, 0)
	ub.StripeChargeId = ""
	ub.IsPaid = false
	ub.IsRefunded = false

	return &ub
}