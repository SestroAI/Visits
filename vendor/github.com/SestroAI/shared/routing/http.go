package routing

import (
	"github.com/emicklei/go-restful"
	"fmt"
	"github.com/SestroAI/shared/config"
	"github.com/emicklei/go-restful-swagger12"
)

func GetCorsConfig(allowedDomains []string, wsContainer *restful.Container) restful.CrossOriginResourceSharing {
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-Set-Authorization"},
		AllowedHeaders: []string{"Content-Type", "Accept", "Authorization"},
		AllowedMethods: []string{"GET", "POST", "PUT", "OPTIONS"},
		CookiesAllowed: true,
		AllowedDomains: allowedDomains,
		Container:      wsContainer}
	return cors
}

func AddSwaggerConfig(wsContainer *restful.Container)  {
	swaggerConfig := swagger.Config{
		WebServices:    wsContainer.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: fmt.Sprintf("%s://%s:%d", config.AppScheme, config.AppHost, config.AppPort),
		ApiPath:        "/apidocs.json",

		// Optionally, specify where the UI is located
		//SwaggerPath:     "/apidocs/",
		//SwaggerFilePath: config.SwaggerRoot,
	}
	swagger.RegisterSwaggerService(swaggerConfig, wsContainer)
}