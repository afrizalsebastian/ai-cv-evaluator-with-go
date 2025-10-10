package helper

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateParams(ctx context.Context, req interface{}) []string {
	err := validate.Struct(req)
	if err == nil {
		return nil
	}

	var fieldErrors []string
	for _, f := range err.(validator.ValidationErrors) {
		fieldErrors = append(fieldErrors, fmt.Sprintf("%s invalid value", f.Field()))
	}
	return fieldErrors
}
