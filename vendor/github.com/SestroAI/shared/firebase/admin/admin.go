package admin

import (
	"golang.org/x/net/context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"

	"google.golang.org/api/option"
	"github.com/SestroAI/shared/config"
	"github.com/SestroAI/shared/logger"
	"errors"
	"io/ioutil"
	"encoding/json"
)

type ServiceAccount struct {
	ClientEmail string `json:"client_email"`
	ClientId string `json:"client_id"`
	ClientCertURL string `json:"client_x509_cert_url"`
	ProjectID string `json:"project_id"`
	PrivateKeyId string `json:"private_key_id"`
	PrivateKey string `json:"private_key"`
}

func GetJsonKey() ([]byte, error) {
	if config.ServiceAccountKeyPath == "" {
		return nil, errors.New("No Key file found")
	}
	body, err := ioutil.ReadFile(config.ServiceAccountKeyPath)

	return body, err
}

func GetServiceAccount() (*ServiceAccount, error) {
	if config.ServiceAccountKeyPath == "" {
		return nil, errors.New("No Key file found")
	}
	body, err := ioutil.ReadFile(config.ServiceAccountKeyPath)
	if err != nil {
		return nil, err
	}
	var serviceAcc ServiceAccount
	err = json.Unmarshal(body, &serviceAcc)
	return &serviceAcc, err
}

func GetFirebaseAdminApp() (*firebase.App, error) {
	opt := option.WithCredentialsFile(config.ServiceAccountKeyPath)
	if config.ServiceAccountKeyPath == "" {
		logger.Errorf("No Service Account Key path was availble. Firebase admin app withh be without app ")
		opt = nil
	}
	app, err := firebase.NewApp(context.Background(), nil, opt)

	return app, err
}

func GetFirebaseAdminAuthClient() (*auth.Client, error){
	app, err := GetFirebaseAdminApp()
	if err != nil {
		return nil, err
	}
	client, err := app.Auth()
	if err != nil {
		return nil, err
	} else if client == nil {
		return nil, errors.New("Got nil firebase auth client")
	}
	return client, nil
}