package auth

import (
	"github.com/SestroAI/shared/config"
	"github.com/stripe/stripe-go"
)

/*
Sample User Account from Firebase
{
	"localId": "UsNNct2PA1WHSeMQrR0TvF45AY12",
	"email": "anshulagrawal35@gmail.com",
	"emailVerified": true,
	"displayName": "Anshul Agrawal",
	"providerUserInfo": [
		{
			"providerId": "google.com",
			"displayName": "Anshul Agrawal",
			"photoUrl": "https://lh3.googleusercontent.com/-5Rs_p2fggwQ/AAAAAAAAAAI/AAAAAAAAR6o/eKbjoKfmzpI/photo.jpg",
			"federatedId": "103398345827386146432",
			"email": "anshulagrawal35@gmail.com",
			"rawId": "103398345827386146432"
		}
	],
	"photoUrl": "https://lh3.googleusercontent.com/-5Rs_p2fggwQ/AAAAAAAAAAI/AAAAAAAAR6o/eKbjoKfmzpI/photo.jpg",
	"validSince": "1498948066",
	"lastLoginAt": "1499034186000",
	"createdAt": "1498948066000"
}
*/

//Firebase user profile
type FirebaseUser struct {
	ID            string                 `json:"localId",mapstructure:"localId"`
	Email         string                 `json:"email",mapstructure:"email"`
	EmailVerified bool                   `json:"emailVerified",mapstructure:"emailVerified"`
	DisplayName   string                 `json:"displayName",mapstructure:"displayName"`
	Accounts      []*FirebaseUserAccount `json:"providerUserInfo",mapstructure:"providerUserInfo,squash"`
	PhotoUrl      string                 `json:"photoUrl",mapstructure:"photoUrl"`
	ValidSince    string                 `json:"validSince",mapstructure:"validSince"`
	Disabled      bool                   `json:"disabled",mapstructure:"disabled"`
	LastLoginAt   string                 `json:"lastLoginAt",mapstructure:"lastLoginAt"`
	CreatedAt     string                 `json:"createdAt",mapstructure:"createdAt"`
	CustomAuth    bool                   `json:"customAuth",mapstructure:"customAuth"`
}

type User struct {
	FirebaseUser    `mapstructure:",squash"`
	Roles           map[string]*Role     `json:"roles",mapstructure:"roles,squash"`
	CustomerProfile *UserCustomerProfile `json:"customerProfile",mapstructure:"customerProfile,squash"`
	MerchantProfiles map[string]*UserMerchantProfile `json:"merchantProfiles",mapstructure:"merchantProfiles,squash"`
}

type UserCustomerProfile struct {
	Visits         map[string]bool `json:"visits",mapstructure:"visits,squash"`
	OngoingVisitId string          `json:"ongoingVisitId",mapstructure:"ongoingVisitId"`
	StripeCustomer *stripe.Customer `json:"stripeCustomer",mapstructure:"stripeCustomer,squash"`
}

type UserMerchantProfile struct {
	IsMerchant bool `json:"isMerchant",mapstructure:"isMerchant"`
	AssociatedMerchantId string `json:"associatedMerchantId",mapstructure:"associatedMerchantId"`
}

type Role struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func DefaultRole() *Role {
	return &Role{
		Name:        config.DefaultUserRole,
		Description: "Default Role assigned by Sestro",
	}
}

//Firebase User Provider account
type FirebaseUserAccount struct {
	UID         string `json:"uid"`
	ProviderId  string `json:"providerId"`
	DisplayName string `json:"disaplyName"`
	PhotoUrl    string `json:"photoUrl"`
	Email       string `json:"email"`
	RawId       string `json:"rawId"`
	FederatedId string `json:"federatedId"`
	ScreenName  string `json:"screenName"`
}

/*
Resources:
1. https://firebase.google.com/docs/reference/rest/auth/
*/