package main

import (
	"errors"
	"regexp"
	"strings"
)

// Validator function type
type ValidatorFunc func(input string) error

// Common validators map
var Validators = map[string]ValidatorFunc{
	"required": func(input string) error {
		if strings.TrimSpace(input) == "" {
			return errors.New("this field is required")
		}
		return nil
	},
	"email": func(input string) error {
		emailRegex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
		if match, _ := regexp.MatchString(emailRegex, input); !match {
			return errors.New("invalid email format")
		}
		return nil
	},
	"no_numbers": func(input string) error {
		if strings.ContainsAny(input, "0123456789") {
			return errors.New("numbers are not allowed")
		}
		return nil
	},
	"alpha_only": func(input string) error {
		alphaRegex := `^[a-zA-Z]+$`
		if match, _ := regexp.MatchString(alphaRegex, input); !match {
			return errors.New("only alphabets are allowed")
		}
		return nil
	},
	"no_special_chars": func(input string) error {
		if strings.ContainsAny(input, "!@#$%^&*(){}[]<>?/~`|\\:;\"'") {
			return errors.New("special characters are not allowed")
		}
		return nil
	},
}
