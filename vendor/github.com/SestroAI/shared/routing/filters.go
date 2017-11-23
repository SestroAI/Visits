package routing

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SestroAI/shared"
	"github.com/SestroAI/shared/config"
	"github.com/SestroAI/shared/dao"
	"github.com/SestroAI/shared/logger"
	"github.com/SestroAI/shared/models/auth"
	"github.com/SestroAI/shared/utils"
	"github.com/emicklei/go-restful"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	ErrNotAuthenticated = errors.New("No authorisation token found.")
)

func GetJWTFromRequest(req *restful.Request) (string, error) {
	token := req.Request.Header.Get("Authorization")
	if token == "" {
		token = req.Request.URL.Query().Get("authorization")
	}
	if token == "" {
		//Check for query param
		token = req.QueryParameter("auth")
		if token == ""{
			return "", ErrNotAuthenticated
		}
	}
	return strings.TrimPrefix(token, "Bearer "), nil
}

func AuthorisationFilter(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
	token, err := GetJWTFromRequest(req)
	if err != nil {
		//User is not authenticated
		//chain.ProcessFilter(req, res)
		res.WriteErrorString(http.StatusUnauthorized, "No token found")
		return
	}

	uid, err := shared.VerifyUserIDToken(token, config.GetGoogleProjectID())
	if err != nil {
		logger.ReqInfof(req, "Invalid user token with err = %s", err.Error())
		res.WriteErrorString(http.StatusUnauthorized, "Token is Invalid or Expired")
		return
	}

	type TokenValidOutput struct {
		RegistrationRequired bool `json:"registrationRequired"`
	}

	var ref = dao.NewUserDao(token)
	user, err := ref.GetUser(uid)
	if err != nil {
		//New user. Register First
		res.WriteHeaderAndEntity(http.StatusNotFound, TokenValidOutput{true})
		return
	}

	req.SetAttribute(config.RequestUser, user)
	req.SetAttribute(config.RequestToken, token)
	chain.ProcessFilter(req, res)
	return
}

func LoggedInFilter(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
	user, _ := req.Attribute(config.RequestUser).(*auth.User)
	token, _ := req.Attribute(config.RequestToken).(string)
	if user == nil || token == "" {
		res.WriteErrorString(http.StatusUnauthorized, "Log-In required")
		return
	}
	chain.ProcessFilter(req, res)
}

type AccessEntry struct {
	Type      string    `json:"type"`
	RequestID string    `json:"requestId"`
	Path      string    `json:"path"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Token     string    `json:"token"`
	URL       string    `json:"url"`
}

func LoggingFilter(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
	entry := AccessEntry{}
	entry.Type = "User Access Log"
	entry.Start = time.Now()
	entry.Path = req.Request.URL.Path

	entry.RequestID = utils.GenerateUUID()
	entry.URL = req.Request.URL.String()
	req.SetAttribute(config.RequestId, entry.RequestID)

	chain.ProcessFilter(req, res)

	entry.End = time.Now()
	log, err := json.Marshal(&entry)
	if err != nil {
		fmt.Fprint(os.Stderr, "Unable to marshal Log Entry with error = "+err.Error())
		return
	}
	logger.ReqInfof(req, string(log))
	return
}
