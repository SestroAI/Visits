package billing

import (
	"github.com/SestroAI/shared/models/orders"
	"github.com/stripe/stripe-go"
	"github.com/SestroAI/shared/utils"
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
	ID string
	BilledToUserId string `json:"billedToUserId"`
	Amount float32 `json:"amount"`
	OrderMap map[string][]*orders.Order `json:"orderMap"`
	AssociatedVisitId string `json:"associatedVisitId"`
	StripeCharge *stripe.Charge `json:"stripeCharge",mapstructure:"stripeCharge,squash"`
	IsPaid bool `json:"isPaid"`
	IsRefunded bool `json:"isRefunded"`
}

func NewUserBill() *UserBill {
	var ub UserBill
	ub.ID = utils.GenerateUUID()
	ub.OrderMap = make(map[string][]*orders.Order, 0)
	ub.StripeCharge = &stripe.Charge{}
	ub.StripeCharge = nil
	ub.IsPaid = false
	ub.IsRefunded = false

	return &ub
}