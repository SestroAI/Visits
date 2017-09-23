package routing

import (
	"github.com/emicklei/go-restful"
	"strings"
	"errors"
	"github.com/SestroAI/shared/logger"
	"github.com/SestroAI/shared/models/auth"
	"github.com/SestroAI/shared/config"
	"net/http"
	"time"
	"encoding/json"
	"github.com/SestroAI/shared"
	"github.com/SestroAI/shared/dao"
	"github.com/SestroAI/shared/utils"
)

var (
	ErrNotAuthenticated = errors.New("No authorisation token found.")
)

func GetJWTFromRequest(req *restful.Request) (string, error) {
	header := req.Request.Header.Get("Authorization")
	if header == "" {
		header = req.Request.URL.Query().Get("authorization")
	}
	if header == ""{
		return "", ErrNotAuthenticated
	}
	return strings.TrimPrefix(header, "Bearer "), nil
}


func AuthorisationFilter(req *restful.Request, res *restful.Response, chain *restful.FilterChain){
	token, err := GetJWTFromRequest(req)
	if err != nil {
		//User is not authenticated
		chain.ProcessFilter(req, res)
		return
	}

	uid, err := shared.VerifyIDToken(token, config.GetGoogleProjectID())
	if err != nil {
		logger.ReqInfof(req, "Invalid user token with err = %s", err.Error())
		res.WriteErrorString(http.StatusUnauthorized, "Token is Invalid or Expired")
		return
	}

	var ref = dao.NewUserDao(token)
	user, err := ref.GetUser(uid)
	if err != nil {
		//New user. Register
		user, err = ref.RegisterFirebaseUser(uid, nil) //nil perms means default
		if err != nil {
			logger.ReqErrorf(req, "Unable to register the new user with ID = %s", uid)
			res.WriteErrorString(http.StatusInternalServerError, "Unable to register the user.")
			return
		}
	}

	req.SetAttribute(config.RequestUser, user)
	req.SetAttribute(config.RequestToken, token)
	chain.ProcessFilter(req, res)
	return
}

func LoggedInFilter(req *restful.Request, res *restful.Response, chain *restful.FilterChain)  {
	user, _ := req.Attribute(config.RequestUser).(*auth.User)
	token, _ := req.Attribute(config.RequestToken).(string)
	if user == nil || token == ""{
		res.WriteErrorString(http.StatusUnauthorized, "Log-In required")
		return
	}
	chain.ProcessFilter(req, res)
}

type AccessEntry struct {
	Type string
	RequestID string `json:"requestId"`
	Path string
	Start time.Time
	End time.Time
	Token string
	Request http.Request
}

func LoggingFilter(req *restful.Request, res *restful.Response, chain *restful.FilterChain){
	entry := AccessEntry{}
	entry.Type = "User Access Log"
	entry.Start = time.Now()
	entry.Path = req.Request.URL.Path

	entry.RequestID = utils.GenerateUUID()
	entry.Request = *req.Request
	req.SetAttribute(config.RequestId, entry.RequestID)

	chain.ProcessFilter(req, res)

	entry.End = time.Now()

	log, err := json.Marshal(&entry)
	if err != nil {
		return
	}
	logger.ReqInfof(req, string(log))
	return
}