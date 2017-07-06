package dao

import(
	"github.com/SestroAI/shared/models/auth"
	"net/http"
	"bytes"
	"encoding/json"
	"errors"

	"github.com/google/logger"
	"github.com/SestroAI/shared/models/visits"
)

const(
	FIREBASE_PROFILE_ENDPOINT = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/getAccountInfo?key="
	DINER_RELATIVE_PATH = "/diners"
	USER_BASE_PATH = "/users"
)

type FireBaseUserQuery struct {
	ID string `json:"idToken"`
}

type FirebaseUserResponse struct {
	Kind string
	Users []*auth.User
}

type UserDao struct {
	Dao
	BasePath string
}

func NewUserDao(token string) *UserDao {
	return &UserDao{
		Dao: *NewDao(token),
		BasePath:USER_BASE_PATH,
	}
}

func (ref *UserDao) GetUser() (*auth.User, error) {
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

	if len(responseObject.Users) == 0{
		return nil, errors.New("No User exists with token = " + ref.Token)
	}

	user := responseObject.Users[0]

	//Update user IDs for all accounts
	for _, account := range user.Accounts {
		account.UID = user.ID
	}

	return user, nil
}

func (ref *UserDao) SaveDiner(id string, diner *auth.Diner) error {
	err := ref.SaveObjectById(id, diner, ref.BasePath + DINER_RELATIVE_PATH)

	if err != nil {
		logger.Errorf("Unable to save Diner object with Id = %s", id)
		return err
	}

	return nil
}

func (ref *UserDao) GetDiner(id string) (*auth.Diner, error) {
	object, _ := ref.GetObjectById(id, ref.BasePath + DINER_RELATIVE_PATH)
	if object == nil {
		return nil, errors.New("Unable to get diner with id = " + id)
	}
	diner := auth.Diner{}
	MapToStruct(object.(map[string]interface{}), &diner)
	return &diner, nil
}

func (ref *UserDao) UpdateDinerOngoingVisit(dinerId string, visit *visits.RestaurantVisit) error {
	diner, err := ref.GetDiner(dinerId)
	if err != nil {
		return err
	}

	diner.OngoingVisitId = visit.ID
	err = ref.SaveDiner(diner.ID, diner)
	return err
}

/*
Resources:
1. https://firebase.google.com/docs/reference/rest/auth/
 */