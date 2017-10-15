package state_token

import(
	"github.com/SestroAI/shared/firebase/admin"
	"github.com/SestroAI/shared/config"
	"errors"
	"github.com/SestroAI/shared"
	"github.com/SestroAI/shared/firebase/service"
)

func GenerateOauthStateToken(data string, redirectUrl string) (string, error) {
	return service.GenerateServiceToken(map[string]string{
		"redirectURL" : redirectUrl,
		"data" : data,
	})
}

func VerifyOauthAStateToken(token string) (map[string]interface{}, error) {
	svcAcc, err := admin.GetServiceAccount()
	if err != nil {
		return nil, err
	}

	return shared.VerifyToken(token,
		svcAcc.ClientCertURL,
		config.GetGoogleProjectID(),
		svcAcc.PrivateKeyId,
	)
}

func GetDataFromOauthStateToken(token string) (string, string, error) {
	claims, err := VerifyOauthAStateToken(token)
	if err != nil {
		return "","", err
	}

	claims_data, ok := claims["data"].(map[string]interface{})
	if !ok {
		return "", "", errors.New("Custom claim data not available for")
	}

	data, ok := claims_data["data"].(string)
	if !ok {
		//data does not exist in claims
		return "", "", errors.New("Missing/Invalid data in Oauth state token claims")
	}

	returnUrl, ok := claims_data["redirectURL"].(string)
	if !ok {
		//data does not exist in claims
		return "", "", errors.New("Missing/Invalid redirectURL in Oauth state token claims")
	}

	return data, returnUrl, err
}