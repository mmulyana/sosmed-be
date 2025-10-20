package main

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func formatValidationError(err error) map[string]string {
	errorsMap := make(map[string]string)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errorsMap[field] = fmt.Sprintf("%s is required", e.Field())
			case "max":
				errorsMap[field] = fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
			default:
				errorsMap[field] = fmt.Sprintf("%s is not valid", e.Field())
			}
		}
	} else {
		errorsMap["error"] = err.Error()
	}
	return errorsMap
}
