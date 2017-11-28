package main

import (
	"github.com/SestroAI/shared/logger"
	"github.com/emicklei/go-restful"
	"net/http"

	"github.com/SestroAI/Visits/services/visits/endpoints"
	sessionEP "github.com/SestroAI/Visits/services/sessions/endpoints"
	"github.com/SestroAI/shared/routing"
	"github.com/SestroAI/shared/utils"
	"os"
)

type APIService struct {
}

func main() {
	wsContainer := restful.NewContainer()

	u := endpoints.VisitResource{}
	u.Register(wsContainer, utils.GetServicePrefix())

	u2 := sessionEP.SessionResource{}
	u2.Register(wsContainer, utils.GetServicePrefix())

	cors := routing.GetCorsConfig([]string{}, wsContainer)
	wsContainer.Filter(cors.Filter)
	wsContainer.Filter(wsContainer.OPTIONSFilter)

	wsContainer.Filter(routing.LoggingFilter)

	wsContainer.Filter(routing.AuthorisationFilter)
	wsContainer.Filter(routing.LoggedInFilter)

	routing.AddSwaggerConfig(wsContainer)

	logger.Infof("Sestro Visit API Server: Start listening on port 8080")
	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	server.ListenAndServe()
	os.Exit(-1)
}
