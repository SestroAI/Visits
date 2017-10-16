package dao

import (
	"errors"
	"fmt"
	"reflect"
	"github.com/SestroAI/shared/logger"
	"github.com/mitchellh/mapstructure"
	"time"
)

func MapToStruct(mapstruct interface{}, structInterface interface{}) error {
	var md mapstructure.Metadata

	stringToDateTimeHook := func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t == reflect.TypeOf(time.Time{}) && f == reflect.TypeOf("") {
			return time.Parse(time.RFC3339, data.(string))
		}

		return data, nil
	}

	config := &mapstructure.DecoderConfig{
		Metadata: &md,
		DecodeHook: stringToDateTimeHook,
		Result:   structInterface,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		logger.Errorf("Unable to get decoder with error = %s", err.Error())
		return err
	}

	err = decoder.Decode(mapstruct)
	if len(md.Unused) != 0 {
		err = errors.New(fmt.Sprintf("Cannot completely convert map to struct with %d fields unsed : %s", len(md.Unused), md.Unused))
	}
	return err
}
