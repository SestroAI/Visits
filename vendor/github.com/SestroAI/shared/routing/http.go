package routing

import "github.com/emicklei/go-restful"

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
