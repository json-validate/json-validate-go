package jsonvalidate

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/json-validate/json-pointer-go"
)

type Validator struct {
	MaxErrors int
	MaxDepth  int
	Registry  Registry
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

func (v Validator) Validate(instance interface{}) (ValidationResult, error) {
	return v.ValidateURI(url.URL{}, instance)
}

func (v Validator) ValidateURI(uri url.URL, instance interface{}) (ValidationResult, error) {
	schema, ok := v.Registry.Schemas[uri]
	if !ok {
		return ValidationResult{}, fmt.Errorf("no schema with uri: %s", uri.String())
	}

	vm := vm{
		maxErrors:      v.MaxErrors,
		maxDepth:       v.MaxDepth,
		registry:       v.Registry,
		instanceTokens: []string{},
		schemas: []schemaStack{
			schemaStack{
				uri:    &uri,
				tokens: []string{},
			},
		},
		errors: []ValidationError{},
	}

	if err := vm.eval(schema, instance); err != nil {
		if err != errMaxErrors {
			return ValidationResult{}, err
		}
	}

	return ValidationResult{Errors: vm.errors}, nil
}
