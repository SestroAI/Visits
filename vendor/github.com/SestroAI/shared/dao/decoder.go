package dao

import (
	"errors"
	"fmt"
	"github.com/SestroAI/shared/logger"
	"github.com/mitchellh/mapstructure"
)

func MapToStruct(mapstruct interface{}, structInterface interface{}) error {
	var md mapstructure.Metadata
	config := &mapstructure.DecoderConfig{
		Metadata: &md,
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
