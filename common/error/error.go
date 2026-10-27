package error

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

// handle error jika ada request error
type ValidationResponse struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

var ErrValidator = map[string]string{}

func ErrValidationResponse(err error) (validationResponse []ValidationResponse) {
	var fieldError validator.ValidationErrors
	if errors.As(err, &fieldError) {
		for _, err := range fieldError {
			switch err.Tag() {
			case "required":
				validationResponse = append(validationResponse, ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("Field %s is required", err.Field()),
				})
			case "email":
				validationResponse = append(validationResponse, ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("Email %s is invalid", err.Field()),
				})
			default:
				errValidator, ok := ErrValidator[err.Tag()]
				if ok {
					count := strings.Count(errValidator, "%s")
					if count == 1 {
						validationResponse = append(validationResponse, ValidationResponse{
							Field:   err.Field(),
							Message: fmt.Sprintf(errValidator, err.Field()),
						})
					} else {
						validationResponse = append(validationResponse, ValidationResponse{
							Field:   err.Field(),
							Message: fmt.Sprintf(errValidator, err.Field(), err.Param()),
						})
					}
				} else {
					validationResponse = append(validationResponse, ValidationResponse{
						Field:   err.Field(),
						Message: fmt.Sprintf("something went wrong on %s; %s", err.Field(), err.Tag()),
					})
				}

			}
		}
	}
	return validationResponse
}

func WrapError(err error) error {
	logrus.Errorf("error: %v", err)
	return err
}
