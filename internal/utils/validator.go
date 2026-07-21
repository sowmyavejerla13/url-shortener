package utils

import (
	"errors"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) map[string]string {

	errors := make(map[string]string)

	for _, fieldErr := range err.(validator.ValidationErrors) {

		field := strings.ToLower(fieldErr.Field())
		
		switch fieldErr.Tag() {

		case "required":
			errors[field] = fieldErr.Field() + " is required"

		case "email":
			errors[field] = "Invalid email address"

		case "min":
			if field == "password" {
				errors[field] = "Password must be at least 8 characters"
			} else {
				errors[field] = fieldErr.Field() + " is too short"
			}

		case "max":
			errors[field] = fieldErr.Field() + " is too long"

		default:
			errors[field] = "Invalid value"
		}
	}

	return errors
}

func ValidateURL(rawURL string) error {

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return errors.New("invalid url")
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return errors.New("invalid url")
	}

	return nil
}