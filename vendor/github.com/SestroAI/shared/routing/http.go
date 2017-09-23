package routing

import (
	"github.com/emicklei/go-restful"
	"fmt"
	"github.com/SestroAI/shared/config"
	"github.com/emicklei/go-restful-swagger12"
	"os"
)

var defaultAllowedDomains = []string{
	"localhost:8080",
	"https://sestro.io",
	"https://stage.sestro.io",
	"https://dashboard.sestro.io",
	"https://dashboard.stage.sestro.io",
}

func GetCorsConfig(allowedDomains []string, wsContainer *restful.Container) restful.CrossOriginResourceSharing {

	if len(allowedDomains) == 0 {
		allowedDomains = defaultAllowedDomains
	}

	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-Set-Authorization"},
		AllowedHeaders: []string{"Content-Type", "Accept", "Authorization", "Access-Control-Allow-Origin"},
		AllowedMethods: []string{"GET", "POST", "PUT", "OPTIONS", "DELETE"},
		CookiesAllowed: true,
		AllowedDomains: allowedDomains,
		Container:      wsContainer}
	return cors
}

func AddSwaggerConfig(container *restful.Container)  {
	wd, err := os.Getwd()
	if err != nil {
		wd = "./"
	}

	swaggerConfig := swagger.Config{
		WebServices:    container.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: fmt.Sprintf("%s://%s:%d", config.AppScheme, config.AppHost, config.AppPort),
		ApiPath:        "/v1/" + config.ServiceName + "/apidocs.json",

		// Optionally, specify where the UI is located
		SwaggerPath:     "/v1/" + config.ServiceName + "/apidocs/",
		SwaggerFilePath: wd + "/static/dist/",
	}
	//swagger.InstallSwaggerService(swaggerConfig)
	swagger.RegisterSwaggerService(swaggerConfig, container)
}