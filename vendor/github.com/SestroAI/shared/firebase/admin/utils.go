package admin

import "github.com/SestroAI/shared/config"

func GenerateServiceUsername() string {
	return "__" + config.ServiceName
}
