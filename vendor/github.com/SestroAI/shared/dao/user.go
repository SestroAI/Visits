package dao

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SestroAI/shared/models/auth"
	"net/http"

	"github.com/SestroAI/shared/logger"
	"github.com/SestroAI/shared/models/visits"
)

const (
	FIREBASE_PROFILE_ENDPOINT = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/getAccountInfo?key="
	USER_BASE_PATH            = "/users"
)

type FireBaseUserQuery struct {
	ID string `json:"idToken"`
}

type FirebaseUserResponse struct {
	Kind  string
	Users []*auth.FirebaseUser
}

type UserDao struct {
	Dao
}

func NewUserDao(token string) *UserDao {
	return &UserDao{
		Dao: *NewDao(token),
	}
}

func (ref *UserDao) GetFirebaseUser() (*auth.FirebaseUser, error) {
	/*
		This will not work with the user generated token. It has to be admin token or firbase api key
	*/
	queryData := FireBaseUserQuery{ID: ref.Token}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(queryData)

	endpoint := FIREBASE_PROFILE_ENDPOINT + ref.APIKey
	res, err := http.Post(endpoint, "application/json; charset=utf-8", b)
	if err != nil {
		return nil, err
	}
	if res.Body == nil {
		return nil, errors.New("Unable to get User for token = " + ref.Token)
	}
	var responseObject FirebaseUserResponse
	err = json.NewDecoder(res.Body).Decode(&responseObject)
	if err != nil {
		return nil, err
	}

	if len(responseObject.Users) == 0 {
		return nil, errors.New("No User exists with token = " + ref.Token)
	}

	user := responseObject.Users[0]

	//Update user IDs for all accounts
	for _, account := range user.Accounts {
		account.UID = user.ID
	}

	return user, nil
}

func (ref *UserDao) RegisterFirebaseUser(userId string, roles []*auth.Role) (*auth.User, error) {
	if len(roles) == 0 {
		roles = append(roles, auth.DefaultRole())
	}

	firebaseUser, err := ref.GetFirebaseUser()
	if err != nil {
		return nil, err
	}

	user := auth.User{FirebaseUser: *firebaseUser}
	user.Roles = roles
	user.CustomerProfile = &auth.UserCustomerProfile{}
	err = ref.SaveUser(user.ID, &user)
	if err != nil {
		logger.Infof("Unable to save user with ID = %s", user.ID)
		return nil, err
	}
	return &user, nil
}

func (ref *UserDao) SaveUser(id string, diner *auth.User) error {
	err := ref.SaveObjectById(id, diner, USER_BASE_PATH)

	if err != nil {
		logger.Errorf("Unable to save Diner object with Id = %s", id)
		return err
	}

	return nil
}

func (ref *UserDao) GetUser(id string) (*auth.User, error) {
	object, err := ref.GetObjectById(id, USER_BASE_PATH)
	if object == nil {
		return nil, errors.New("Unable to get diner with id = " + id)
	}
	if err != nil {
		return nil, err
	}
	user := auth.User{}
	_ = MapToStruct(object.(map[string]interface{}), &user)

	user.ID = object.(map[string]interface{})["localId"].(string)
	/*Todo
	Handle the map struct errorin a better way
	*/
	return &user, nil
}

func (ref *UserDao) UpdateDinerOngoingVisit(userId string, visit *visits.MerchantVisit) error {
	user, err := ref.GetUser(userId)
	if err != nil {
		return err
	}

	user.CustomerProfile.OngoingVisitId = visit.ID
	err = ref.SaveUser(user.ID, user)
	return err
}

/*
Resources:
1. https://firebase.google.com/docs/reference/rest/auth/
*/
