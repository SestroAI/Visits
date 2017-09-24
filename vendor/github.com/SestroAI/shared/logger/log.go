package logger

import (
	"encoding/json"
	"fmt"
	"github.com/SestroAI/shared/config"
	"github.com/SestroAI/shared/models/auth"
	"github.com/emicklei/go-restful"
	"os"
	"runtime"
	"time"
)

type LogFormat struct {
	Timestamp       time.Time `json:"timestamp"`
	Service         string    `json:"service"`
	GoogleProjectId string    `json:"googleProjectId"`
	Level           string    `json:"level"`
	Message         string    `json:"message"`
	RequestId       string    `json:"requestId"`
	UserId          string    `json:"userId"`
	File            string    `json:"file"`
	Line            int       `json:"line"`
}

func NewLog(req *restful.Request, message string, level string) {

	logObject := LogFormat{}

	if req != nil {
		user, _ := req.Attribute(config.RequestUser).(*auth.User)
		requestID, _ := req.Attribute(config.RequestId).(string)
		if user != nil {
			logObject.UserId = user.ID
		}
		logObject.RequestId = requestID
	}

	_, file, line, ok := runtime.Caller(2)
	if ok {
		logObject.File = file
		logObject.Line = line
		logObject.Message = message
	}

	logObject.Timestamp = time.Now()
	logObject.Level = level
	logObject.Service = config.ServiceName
	logObject.GoogleProjectId = config.GetGoogleProjectID()

	data, err := json.Marshal(logObject)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to marshal log object")
		return
	}

	fmt.Fprintf(os.Stderr, string(data))
	fmt.Fprintf(os.Stderr, "\n")
}

func ReqErrorf(req *restful.Request, format string, v ...interface{}) {
	message := ""
	if len(v) != 0 {
		message = fmt.Sprintf(format, v...)
	} else {
		message = fmt.Sprintf(format)
	}
	NewLog(req, message, "Error")
}

func Errorf(format string, v ...interface{}) {
	message := ""
	if len(v) != 0 {
		message = fmt.Sprintf(format, v...)
	} else {
		message = fmt.Sprintf(format)
	}
	NewLog(nil, message, "Error")
}

func ReqInfof(req *restful.Request, format string, v ...interface{}) {
	message := ""
	if len(v) != 0 {
		message = fmt.Sprintf(format, v...)
	} else {
		message = fmt.Sprintf(format)
	}
	NewLog(req, message, "Info")
}

func Infof(format string, v ...interface{}) {
	message := ""
	if len(v) != 0 {
		message = fmt.Sprintf(format, v...)
	} else {
		message = fmt.Sprintf(format)
	}
	NewLog(nil, message, "Info")
}
