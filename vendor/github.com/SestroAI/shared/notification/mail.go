package notification

import (
	"io/ioutil"
	"html/template"
	"bytes"
)

const EmailNotificationEventname = "email_notification"

type NotifyUserEmailEvent struct {
	UserId string `json:"userId"`
	EmailHtmlContent string `json:"emailHtmlContent"`
	Subject string `json:"subject"`
	Tags []string `json:"tags"`
}

func GenerateHTMLFromTemplate(templatePath string, data interface{}) (string, error) {

	content, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	t, err := template.New("email").Parse(string(content))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}