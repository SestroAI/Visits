package config

import(
	"os"
)

var(
	AppScheme = os.Getenv("APP_HOST_SCHEME")
	AppHost = os.Getenv("APP_HOST")
	AppPort = os.Getenv("APP_HOST_PORT")
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