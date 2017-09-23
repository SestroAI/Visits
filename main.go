package main

import (
	"github.com/emicklei/go-restful"
	"net/http"
	"github.com/SestroAI/shared/logger"

	"github.com/SestroAI/shared/routing"
	"github.com/SestroAI/Visits/services/visits/endpoints"
	"os"
)


type APIService struct {

}

func main()  {
	wsContainer := restful.NewContainer()

	u := endpoints.VisitResource{}

	u.Register(wsContainer)

	wsContainer.Filter(routing.AuthorisationFilter)
	wsContainer.Filter(routing.LoggingFilter)
	wsContainer.Filter(wsContainer.OPTIONSFilter)

	cors := routing.GetCorsConfig([]string{}, wsContainer)

	wsContainer.Filter(cors.Filter)
	routing.AddSwaggerConfig(wsContainer)

	logger.Infof("Sestro Visit API Server: Start listening on port 8080")
	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	server.ListenAndServe().Error()
	os.Exit(-1)
}
