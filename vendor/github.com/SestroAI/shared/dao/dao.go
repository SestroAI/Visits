package dao

import (
	"bytes"
	"encoding/json"
	"github.com/SestroAI/shared/logger"
	"net/http"

	"errors"
	"github.com/SestroAI/shared/config"
	"io/ioutil"
	"strconv"
)

type Dao struct {
	APIKey      string
	Token       string
	FireBaseURL string
}

func NewDao(token string) *Dao {
	dao := Dao{
		Token:       token,
		APIKey:      config.GetFirebaseDBAPIKey(),
		FireBaseURL: config.GetFirebaseDBURL(),
	}
	return &dao
}

func (ref *Dao) PrepareURL(url string) string {
	if ref.Token != "" {
		url += "?auth=" + ref.Token
	}
	return url
}

func (ref *Dao) PrepareRequest(req *http.Request) {
	q := req.URL.Query()
	q.Add("auth", ref.Token)
	req.URL.RawQuery = q.Encode()
}

func (ref *Dao) GetObjectById(id string, objectPath string) (interface{}, error) {
	var objectInstance interface{}
	url := ref.FireBaseURL + objectPath + "/" + id + ".json"
	url = ref.PrepareURL(url)
	res, err := http.Get(url)
	if err != nil {
		logger.Errorf("Unable to get object at url %s from firebase with err = %s", url, err.Error())
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Unable to get object at url " + url + " with firebase response status " + res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(&objectInstance)
	if err != nil {
		logger.Errorf("Unable to get object with ID = %s and Error: %s", id, err.Error())
		return nil, err
	}

	return objectInstance, nil
}

func (ref *Dao) SaveObjectById(id string, object interface{}, objectPath string) error {
	if id == "" {
		return errors.New("Cannot save object with empty ID")
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(object)

	url := ref.FireBaseURL + objectPath + "/" + id + ".json"
	url = ref.PrepareURL(url)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, url, b)
	if err != nil {
		logger.Errorf("Unable to create PUT request to save object with Error = %s", err.Error())
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("Unable to make req to save object with ID = %s and Error: %s", id, err.Error())
		return err
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		b, _ := ioutil.ReadAll(res.Body)
		logger.Errorf("Unable to save object with ID = %s and firebase response: %s", id, string(b))
		return errors.New("Unable to save object. Firebase returned code: " + strconv.Itoa(res.StatusCode))
	}
	return nil
}
