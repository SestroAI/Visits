package dao

import (
	"github.com/mitchellh/mapstructure"
	"github.com/google/logger"
)

func MapToStruct(mapstruct map[string]interface{}, structInterface interface{}){
	err := mapstructure.Decode(mapstruct, structInterface)
	if err != nil {
		logger.Errorf("Unable to convert map struct to struct with err = %s", err)
		panic(err)
	}
	return
}