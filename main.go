package main

import (
	"github.com/emicklei/go-restful"
	"net/http"
	"github.com/google/logger"

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
	wsContainer.Filter(routing.DinerFilter)
	wsContainer.Filter(wsContainer.OPTIONSFilter)

	cors := routing.GetCorsConfig([]string{
		"localhost:",
		"sestro.io",
		"www.sestro.io",
		"stage-dev.sestro.io",
	}, wsContainer)

	wsContainer.Filter(cors.Filter)

	logger.Init("SestroVisitService", false, false, os.Stderr)

	logger.Infof("Sestro Visit API Server: Start listening on port 8080")
	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	logger.Fatal(server.ListenAndServe())
}
