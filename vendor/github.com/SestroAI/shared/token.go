package shared

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/dgrijalva/jwt-go.v3"
	"github.com/SestroAI/shared/config"
	"github.com/SestroAI/shared/firebase/service"
)

/*
Sample Firebase standard ID Token
{
  "iss": "https://securetoken.google.com/sestro-test",
  "name": "Anshul Agrawal",
  "picture": "https://lh3.googleusercontent.com/-5Rs_p2fggwQ/AAAAAAAAAAI/AAAAAAAAR6o/eKbjoKfmzpI/photo.jpg",
  "aud": "sestro-test",
  "auth_time": 1498964190,
  "user_id": "UsNNct2PA1WHSeMQrR0TvF45AY12",
  "sub": "UsNNct2PA1WHSeMQrR0TvF45AY12",
  "iat": 1498964190,
  "exp": 1498967790,
  "email": "anshulagrawal35@gmail.com",
  "email_verified": true,
  "firebase": {
    "identities": {
      "google.com": [
        "103398345827386146432"
      ],
      "email": [
        "anshulagrawal35@gmail.com"
      ]
    },
    "sign_in_provider": "google.com"
  }
}
*/

func VerifyToken(idToken, clientCertURL, aud, pkId string) (map[string]interface{}, error) {
	/*
		Returns User ID from a validated JWT token. Returns error if token is not valid/expired
	*/
	parsedToken, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		if clientCertURL == "" {
			claims, ok := token.Claims.(jwt.StandardClaims)
			if !ok {
				return nil, fmt.Errorf("Invalid Claims")
			}
			clientCertURL = GetClientCertUrl(claims.Issuer)
		}

		keys, err := fetchPublicKeys(clientCertURL)
		if err != nil {
			return nil, err
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			kid = pkId
		}

		certPEM := string(*keys[kid])
		certPEM = strings.Replace(certPEM, "\\n", "\n", -1)
		certPEM = strings.Replace(certPEM, "\"", "", -1)
		block, _ := pem.Decode([]byte(certPEM))
		var cert *x509.Certificate
		cert, _ = x509.ParseCertificate(block.Bytes)
		rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)

		return rsaPublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims := parsedToken.Claims.(jwt.MapClaims)

	if claims["aud"].(string) != aud {
		return nil, errors.New("ID token has incorrect 'aud' claim: " + claims["aud"].(string))
	}

	return claims, nil
}

func VerifyUserIDToken(token, googleProjectId string) (string, error) {
	claims, err := VerifyToken(token,
		config.ClientCertURL,
		googleProjectId,
		"")
	if err != nil {
		return "", err
	}

	uid, ok := claims["user_id"].(string)
	if !ok || uid == ""{
		return "", errors.New("Unable to get user_id from token")
	}

	return uid, nil
}

func fetchPublicKeys(clientCertURL string) (map[string]*json.RawMessage, error) {
	resp, err := http.Get(clientCertURL)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var objmap map[string]*json.RawMessage
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&objmap)

	return objmap, err
}

func GetClientCertUrl(iss string) string {
	switch iss{
	case "https://securetoken.google.com/"+ config.GetGoogleProjectID():
		return config.ClientCertURL
	default:
		url, err := service.GetServiceAccountPublicCertificateURL(iss)
		if err != nil {
			return ""
		}
		return url
	}
}