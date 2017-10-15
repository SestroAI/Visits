package config

import (
	"os"
	"strconv"
)

//Restful Request Attributes
const (
	RequestUser  = "user"
	RequestToken = "token"
	RequestId    = "rId"
)

var (
	AppScheme   = os.Getenv("APP_HOST_SCHEME")
	AppHost     = os.Getenv("APP_HOST")
	AppPort, _  = strconv.Atoi(os.Getenv("APP_HOST_PORT"))
	ServiceName = os.Getenv("SERVICE_NAME")
//	ServiceAccountKeyPath = os.Getenv("SERVICE_ACCOUNT_KEY_PATH")
	ServiceAccountKeyPath = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	AuthorizedClientList = os.Getenv("AUTHORIZED_CLIENT_LIST")

	FirebaseAdminAccessScopes = []string{
		"https://www.googleapis.com/auth/firebase",
		"https://www.googleapis.com/auth/firebase.readonly",
		"https://www.googleapis.com/auth/userinfo.email",
	}
)

const (
	DefaultUserRole = "sestro_guest"
	DefaultServiceRole = "sestro_admin_service"

	ClientCertURL = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"
)

func GetFirebaseDBAPIKey() string {
	return os.Getenv("FIREBASE_API_KEY")
}

func GetFirebaseDBURL() string {
	return os.Getenv("FIREBASE_URL")
}

func GetGoogleProjectID() string {
	return os.Getenv("GOOGLE_PROJECT_ID")
}
