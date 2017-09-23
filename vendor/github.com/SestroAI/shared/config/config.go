package config

import(
	"os"
	"strconv"
)

//Restful Request Attributes
const (
	RequestUser = "user"
	RequestToken = "token"
	RequestDiner = "diner"
	RequestId = "rId"
)

var(
	AppScheme = os.Getenv("APP_HOST_SCHEME")
	AppHost = os.Getenv("APP_HOST")
	AppPort, _ = strconv.Atoi(os.Getenv("APP_HOST_PORT"))
	ServiceName = os.Getenv("SERVICE_NAME")
)

const(
	DefaultUserRole = "sestro_guest"
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