package auth

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
type User struct {
	ID string `json:"localId"`
	Email string
	EmailVerified bool `json:"emailVerified"`
	DisplayName string `json:"displayName"`
	Accounts []*UserAccount `json:"providerUserInfo"`
	PhotoUrl string `json:"photoUrl"`
	ValidSince string `json:"validSince"`
	Disabled bool
	LastLoginAt string `json:"lastLoginAt"`
	CreatedAt string `json:"createdAt"`
	CustomeAuth bool `json:"customeAuth"`
}

type Diner struct{
	ID string
	Visits map[string]bool
	OngoingVisitId string `json:"ongoingVisitId"`
}

//Firebase User Provider account
type UserAccount struct {
	UID string `json:"uid"`
	ProviderId string `json:"providerId"`
	DisaplyName string `json:"disaplyName"`
	PhotoUrl string `json:"photoUrl"`
	Email string
	RawId string `json:"rawId"`
	FederatedId string `json:"federatedId"`
	ScreenName string `json:"screenName"`
}
/*
Resources:
1. https://firebase.google.com/docs/reference/rest/auth/
 */