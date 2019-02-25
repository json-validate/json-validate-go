package jsonvalidate

import (
	"net/url"

	"github.com/json-validate/json-pointer-go"
)

type Validator struct {
	Schemas   map[url.URL]Schema
	MaxErrors int
	MaxDepth  int
}

type ValidatorConfig struct {
	MaxErrors int
	MaxDepth  int
}

type ValidationResult struct {
	Errors []ValidationError `json:"errors"`
}

func (r ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

type ValidationError struct {
	InstancePath jsonpointer.Ptr
	SchemaPath   jsonpointer.Ptr
	SchemaURI    url.URL
}

func NewValidator(schemas []Schema) (Validator, error) {
	return NewValidatorWithConfig(ValidatorConfig{}, schemas)
}

func NewValidatorWithConfig(config ValidatorConfig, schemas []Schema) (Validator, error) {
	return Validator{
		MaxDepth:  config.MaxDepth,
		MaxErrors: config.MaxErrors,
	}, nil
}

func (v Validator) Validate(instance interface{}) ValidationResult {
	return v.ValidateURI(url.URL{}, instance)
}

func (v Validator) ValidateURI(uri url.URL, instance interface{}) ValidationResult {
	return ValidationResult{}
}
