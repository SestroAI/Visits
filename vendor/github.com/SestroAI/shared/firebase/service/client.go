package service

import (
	"strings"
	"errors"
	"github.com/SestroAI/shared/config"
)

func GetAuthorizedServiceAccountEmails() []string{
	result := make([]string, 0)

	service_account_names := strings.Split(config.AuthorizedClientList, ",")

	for _, name := range service_account_names {
		email := name + "@" + config.GetGoogleProjectID() + ".iam.gserviceaccount.com"
		result = append(result, email)
	}
	return result
}

func IsServiceAccountValid(service_account_email string) bool {
	valid_emails := GetAuthorizedServiceAccountEmails()
	for _, authorized_service := range valid_emails {
		if service_account_email == authorized_service {
			return true
		}
	}
	return false
}

func GetServiceAccountPublicCertificateURL(service_account_email string) (string, error) {
	if !IsServiceAccountValid(service_account_email) {
		return "", errors.New("Unauthorized service account, not in allowed service list")
	}

	public_cert_url := "https://www.googleapis.com/robot/v1/metadata/x509/" + service_account_email
	/*
	TODO: URL Encode public cert url
	 */
	return public_cert_url, nil
}