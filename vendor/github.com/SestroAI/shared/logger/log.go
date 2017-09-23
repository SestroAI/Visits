package logger

import(
	"fmt"
	"github.com/emicklei/go-restful"
	"runtime"
	"github.com/SestroAI/shared/config"
	"github.com/SestroAI/shared/models/auth"
	"encoding/json"
	"os"
	"time"
)

type LogFormat struct {
	Timestamp time.Time
	Service string
	GoogleProjectId string `json:"googleProjectId"`
	Level string
	Message string
	RequestId string `json:"requestId"`
	UserId string `json:"userId"`
	File string
	Line int
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
	logObject.Service = os.Getenv("SERVICE_NAME")
	logObject.GoogleProjectId = os.Getenv("GOOGLE_PROJECT_ID")

	data, err := json.Marshal(logObject)
	if err != nil {
		fmt.Fprintf(os.Stderr,"Error: Unable to marshal log object")
		return
	}

	fmt.Fprintf(os.Stderr, string(data))
	fmt.Fprintf(os.Stderr, "\n")
}

func ReqErrorf(req *restful.Request, format string, v ...interface{}) {
	message := ""
	if len(v) != 0 {
		message = fmt.Sprintf(format, v)
	} else {
		message = fmt.Sprintf(format)
	}
	NewLog(req, message, "Error")
}

func Errorf(format string, v ...interface{}){
	message := ""
	if len(v) != 0 {
		message = fmt.Sprintf(format, v)
	} else {
		message = fmt.Sprintf(format)
	}
	NewLog(nil, message, "Error")
}

func ReqInfof(req *restful.Request, format string, v ...interface{}) {
	message := ""
	if len(v) != 0 {
		message = fmt.Sprintf(format, v)
	} else {
		message = fmt.Sprintf(format)
	}
	NewLog(req, message, "Info")
}

func Infof(format string, v ...interface{}){
	message := ""
	if len(v) != 0 {
		message = fmt.Sprintf(format, v)
	} else {
		message = fmt.Sprintf(format)
	}
	NewLog(nil, message, "Info")
}