package token

import(
	"github.com/SestroAI/shared/firebase/admin"
	"github.com/SestroAI/shared/models/auth"
	"github.com/SestroAI/shared/config"
	"errors"
	"github.com/SestroAI/shared"
)

func getAdminRoles() []*auth.Role {
	return []*auth.Role{
		{
			Name:config.DefaultServiceRole,
		},
	}
}

func GenerateServiceToken() (string, error) {
	uid := admin.GenerateServiceUsername()
	client, err := admin.GetFirebaseAdminAuthClient()
	if err != nil {
		return "", err
	}
	claims := map[string]interface{}{
		"roles": getAdminRoles(),
	}
	return client.CustomTokenWithClaims(uid, claims)
}

func GenerateOauthStateToken(data string, redirectUrl string) (string, error) {
	client, err := admin.GetFirebaseAdminAuthClient()
	if err != nil {
		return "", err
	}
	claims := map[string]interface{}{
		"data" : data,
		"redirectURL" : redirectUrl,
	}
	return client.CustomTokenWithClaims(data, claims)
}

//Return Claims if verified
func VerifyCustomToken(token string) (map[string]interface{}, error) {

	svcAcc, err := admin.GetServiceAccount()
	if err != nil {
		return nil, err
	}

	return shared.VerifyToken(token,
		svcAcc.ClientCertURL,
		"https://identitytoolkit.googleapis.com/google.identity.identitytoolkit.v1.IdentityToolkit",
		svcAcc.ClientEmail,
		svcAcc.PrivateKeyId,
	)
}

func GetDataFromOauthStateToken(token string) (string, string, error) {
	tokenData, err := VerifyCustomToken(token)
	if err != nil {
		return "","", err
	}
	claims, ok := tokenData["claims"].(map[string]interface{})
	if !ok {
		return "", "", errors.New("No Claims in state")
	}


	data, ok := claims["data"].(string)
	if !ok {
		//data does not exist in claims
		return "", "", errors.New("Missing/Invalid data in Oauth state token claims")
	}

	returnUrl, ok := claims["redirectURL"].(string)
	if !ok {
		//data does not exist in claims
		return "", "", errors.New("Missing/Invalid redirectURL in Oauth state token claims")
	}

	return data, returnUrl, err
}