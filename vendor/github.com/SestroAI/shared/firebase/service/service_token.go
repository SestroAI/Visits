package service

import (
	"fmt"
	"time"

	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/dgrijalva/jwt-go"

	"github.com/SestroAI/shared/firebase/admin"
	"github.com/SestroAI/shared/config"

	"golang.org/x/oauth2/google"

	"os"
	"golang.org/x/net/context"
)

type ServiceJwtClaim struct {
	jwt.StandardClaims
	Data map[string]string `json:"data"`
}

func GetServiceAccessKey(scopes ...string) (string, error) {
	//Make sure env var for default creds is setup
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.ServiceAccountKeyPath)
	oauthTokenSource, err := google.DefaultTokenSource(context.Background(),
		scopes...)
	if err != nil {
		return "", err
	}
	token, err := oauthTokenSource.Token()
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func GetServiceFirebaseAccessKey() (string, error) {
	return GetServiceAccessKey(config.FirebaseAdminAccessScopes...)
}

func GenerateServiceToken(data map[string]string) (string, error){
	serviceAccount, err := admin.GetServiceAccount()
	if err != nil {
		return "", fmt.Errorf("Could not read service account: %v", err)
	}

	rsaKey, err := parseKey([]byte(serviceAccount.PrivateKey))
	if err != nil {
		return "", fmt.Errorf("Could not get RSA key: %v", err)
	}

	iat := time.Now()
	exp := iat.Add(time.Minute * 10)

	if data == nil {
		data = make(map[string]string, 0)
	}

	claims := ServiceJwtClaim{
		jwt.StandardClaims{
			Issuer:serviceAccount.ClientEmail,
			Audience:config.GetGoogleProjectID(),
			IssuedAt:iat.Unix(),
			ExpiresAt:exp.Unix(),
			Subject:"",
		},
		data,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = serviceAccount.PrivateKeyId
	return token.SignedString(rsaKey)
}

func parseKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block != nil {
		key = block.Bytes
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("private key should be a PEM or plain PKSC1 or PKCS8; parse error: %v", err)
		}
	}
	parsed, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is invalid")
	}
	return parsed, nil
}