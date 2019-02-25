package jsonvalidate

import (
	"encoding/json"
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

func (e *ValidationError) UnmarshalJSON(data []byte) error {
	var strings map[string]string
	if err := json.Unmarshal(data, &strings); err != nil {
		return err
	}

	if val, ok := strings["instancePath"]; ok {
		ptr, err := jsonpointer.New(val)
		if err != nil {
			return err
		}

		e.InstancePath = ptr
	}

	if val, ok := strings["schemaPath"]; ok {
		ptr, err := jsonpointer.New(val)
		if err != nil {
			return err
		}

		e.SchemaPath = ptr
	}

	if val, ok := strings["schemaURI"]; ok {
		uri, err := url.Parse(val)
		if err != nil {
			return err
		}

		e.SchemaURI = *uri
	}

	return nil
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
