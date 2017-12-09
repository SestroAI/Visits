package controller

import (
	"github.com/SestroAI/shared/models/auth"
	"github.com/SestroAI/shared/models/visits"
)

func IsUserAllowedToOrder(user *auth.User, visit *visits.MerchantVisit) (bool, string) {
	if visit.IsComplete {
		return false, "visitEnded"
	}

	if user.CustomerProfile.StripeCustomer.Delinquent {
		//User has unpaid payments
		return false, "delinquent"
	}

	if visit.IsOpenForAll {
		//PayForAllOpted
		return true, "ok"
	}

	if user.CustomerProfile.StripeCustomer == nil || user.CustomerProfile.StripeCustomer.Sources.ListMeta.Count == 0{
		//No Stripe Customer
		return false, "noCard"
	}

	if _, ok := visit.Diners[user.ID]; !ok {
		//User not a part of visit
		return false, "userNotPartOfVisit"
	}

	return true, "ok"
}