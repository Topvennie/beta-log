// Package dto forms the bridge between the api data and the internal models
package dto

import "github.com/go-playground/validator/v10"

var Validate *validator.Validate

func InitValidator() error {
	Validate = validator.New(validator.WithRequiredStructEnabled())

	// Register custom validators

	return nil
}
