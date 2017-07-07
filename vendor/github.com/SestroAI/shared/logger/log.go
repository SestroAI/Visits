package logger

import(
	"github.com/google/logger"
	"fmt"
	"github.com/emicklei/go-restful"
	"runtime"
	"github.com/SestroAI/shared/config"
	"github.com/SestroAI/shared/models/auth"
	"encoding/json"
)

type LogFormat struct {
	Message string
	RequestId string `json:"requestId"`
	UserId string `json:"userId"`
	File string
	Line int
}

func NewLog(req *restful.Request, message string) string {

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

	data, err := json.Marshal(logObject)
	if err != nil {
		logger.Errorf("Unable to marshal log object")
	}
	return string(data)
}

func ReqErrorf(req *restful.Request, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v)
	log := NewLog(req, message)
	logger.Errorf(log)
}

func Errorf(format string, v ...interface{}){
	message := fmt.Sprintf(format, v)
	log := NewLog(nil, message)
	logger.Errorf(log)
}

func ReqInfof(req *restful.Request, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v)
	log := NewLog(req, message)
	logger.Infof(log)
}

func Infof(format string, v ...interface{}){
	message := fmt.Sprintf(format, v)
	log := NewLog(nil, message)
	logger.Infof(log)
}