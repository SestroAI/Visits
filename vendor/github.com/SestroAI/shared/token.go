package shared

import(
	"net/http"
	"fmt"
	"strings"
	"encoding/pem"
	"crypto/x509"
	"crypto/rsa"
	"errors"
	"encoding/json"

	"gopkg.in/dgrijalva/jwt-go.v3"
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

const (
	clientCertURL = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"
)

func VerifyIDToken(idToken string, googleProjectID string) (string, error) {
	/*
		Returns User ID from a validated JWT token. Returns error if token is not valid/expired
	*/
	keys, err := fetchPublicKeys()

	if err != nil {
		return "", err
	}

	parsedToken, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		kid := token.Header["kid"]

		certPEM := string(*keys[kid.(string)])
		certPEM = strings.Replace(certPEM, "\\n", "\n", -1)
		certPEM = strings.Replace(certPEM, "\"", "", -1)
		block, _ := pem.Decode([]byte(certPEM))
		var cert *x509.Certificate
		cert, _ = x509.ParseCertificate(block.Bytes)
		rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)

		return rsaPublicKey, nil
	})

	if err != nil {
		return "", err
	}

	errMessage := ""

	claims := parsedToken.Claims.(jwt.MapClaims)

	if claims["aud"].(string) != googleProjectID {
		errMessage = "Firebase Auth ID token has incorrect 'aud' claim: " + claims["aud"].(string)
	} else if claims["iss"].(string) != "https://securetoken.google.com/"+googleProjectID {
		errMessage = "Firebase Auth ID token has incorrect 'iss' claim"
	} else if claims["sub"].(string) == "" || len(claims["sub"].(string)) > 128 {
		errMessage = "Firebase Auth ID token has invalid 'sub' claim"
	}
	if errMessage != "" {
		return "", errors.New(errMessage)
	}

	uid := string(claims["user_id"].(string))


	return uid, nil
}

func fetchPublicKeys() (map[string]*json.RawMessage, error) {
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