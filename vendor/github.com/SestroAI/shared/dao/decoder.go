package dao

import (
	"github.com/mitchellh/mapstructure"
	"github.com/SestroAI/shared/logger"
)

func MapToStruct(mapstruct map[string]interface{}, structInterface interface{}) error {
	err := mapstructure.Decode(mapstruct, structInterface)
	if err != nil {
		logger.Errorf("Unable to convert map struct to struct with err = %s", err)
	}
	return err
}