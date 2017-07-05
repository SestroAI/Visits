package routing

import (
	"github.com/emicklei/go-restful"
	"strings"
	"errors"
	"github.com/google/logger"
	"github.com/SestroAI/shared/models/auth"
	"github.com/SestroAI/shared/config"
	"net/http"
	"time"
	"encoding/json"
	"github.com/SestroAI/shared"
	"github.com/SestroAI/shared/dao"
)

var (
	ErrNotAuthenticated = errors.New("No authorisation token found.")
)

const (
	RequestUser = "user"
	RequestToken = "token"
	RequestDiner = "diner"
)

func GetJWTFromRequest(req *restful.Request) (string, error) {
	header := req.Request.Header.Get("Authorization")
	if header == "" {
		header = req.Request.URL.Query().Get("autorization")
	}
	if header == ""{
		return "", ErrNotAuthenticated
	}
	return strings.TrimPrefix(header, "Bearer "), nil
}


func AuthorisationFilter(req *restful.Request, res *restful.Response, chain *restful.FilterChain){
	token, err := GetJWTFromRequest(req)
	if err != nil {
		logger.Infof("User not authenticated. Anonymoud request! Error: ", err.Error())
		chain.ProcessFilter(req, res)
		return
	}
	user, err := shared.VerifyIDToken(token, config.GetGoogleProjectID())
	if err != nil {
		logger.Infof("In-valid ID Token with error : %s. Cannot continue", err.Error())
		res.WriteErrorString(http.StatusUnauthorized, "In-valid ID Token")
		return
	}

	req.SetAttribute(RequestUser, user)
	req.SetAttribute(RequestToken, token)
	chain.ProcessFilter(req, res)
	return
}

func LoggedInFilter(req *restful.Request, res *restful.Response, chain *restful.FilterChain)  {
	user, _ := req.Attribute(RequestUser).(*auth.User)
	token, _ := req.Attribute(RequestToken).(string)
	if user == nil || token == ""{
		logger.Infof("Anonymoud User not allowed")
		res.WriteErrorString(http.StatusUnauthorized, "Un-authorized user. Log-In required")
		return
	}
	chain.ProcessFilter(req, res)
}

func DinerFilter(req *restful.Request, res *restful.Response, chain *restful.FilterChain){
	user, _ := req.Attribute(RequestUser).(*auth.User)
	token, _ := req.Attribute(RequestToken).(string)
	if user == nil || token == ""{
		logger.Infof("Anonymoud User is not a Diner")
		res.WriteErrorString(http.StatusUnauthorized, "Un-authorized user")
		return
	}
	userDao := dao.NewUserDao(token)
	uid := user.ID
	diner, err := userDao.GetDiner(uid)
	if err != nil {
		res.WriteErrorString(http.StatusUnauthorized, "Not a valid diner.")
		return
	}
	req.SetAttribute(RequestDiner, diner)
	chain.ProcessFilter(req, res)
}

type AccessEntry struct {
	Type string
	Path string
	Start time.Time
	End time.Time
	User string
}

func LoggingFilter(req *restful.Request, res *restful.Response, chain *restful.FilterChain){
	entry := AccessEntry{}
	entry.Type = "User Access Log"
	entry.Start = time.Now()
	entry.Path = req.Request.URL.Path
	user, _ := req.Attribute(RequestUser).(*auth.User)
	if user != nil {
		entry.User = user.ID
	}
	chain.ProcessFilter(req, res)
	entry.End = time.Now()

	log, err := json.Marshal(&entry)
	if err != nil {
		logger.Errorf("Unable to log the access entry!")
		return
	}
	logger.Infof(string(log))
	return
}