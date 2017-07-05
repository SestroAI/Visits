package config

import(
	"os"
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