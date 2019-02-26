package jsonvalidate

import (
	"encoding/json"
	"net/url"

	"github.com/json-validate/json-pointer-go"
)

type Validator struct {
	maxErrors int
	maxDepth  int
	registry  map[url.URL]*Schema
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
	registry := map[url.URL]*Schema{}
	for i, schema := range schemas {
		var uri url.URL

		if schema.ID != nil {
			parsedURI, err := url.Parse(*schema.ID)
			if err != nil {
				return Validator{}, err
			}

			uri = *parsedURI
		}

		schemas[i].baseURI = uri
		registry[uri] = &schemas[i]

		for name, def := range schema.Definitions {
			defURI := uri
			defURI.Fragment = name

			registry[defURI] = def
		}
	}

	missingURIs := []url.URL{}
	for _, schema := range registry {
		if err := populateRefs(&missingURIs, registry, &schema.baseURI, schema); err != nil {
			return Validator{}, err
		}
	}

	if len(missingURIs) > 0 {
		return Validator{}, ErrMissingSchemas{URIs: missingURIs}
	}

	return Validator{
		maxDepth:  config.MaxDepth,
		maxErrors: config.MaxErrors,
		registry:  registry,
	}, nil
}

func populateRefs(missingURIs *[]url.URL, registry map[url.URL]*Schema, baseURI *url.URL, schema *Schema) error {
	if schema == nil {
		return nil
	}

	schema.baseURI = *baseURI

	// First, populate the ref on this schema, if any
	if schema.Ref != nil {
		refURI, err := baseURI.Parse(*schema.Ref)
		if err != nil {
			return err
		}

		// Verify that there exists a schema with the given name
		refOk := false
		refBaseURI := *refURI
		refBaseURI.Fragment = ""

		if refSchema, ok := registry[refBaseURI]; ok {
			if refURI.Fragment == "" {
				refOk = true
			} else {
				_, refOk = refSchema.Definitions[refURI.Fragment]
			}
		}

		if refOk {
			schema.refURI = *refURI
		} else {
			*missingURIs = append(*missingURIs, *refURI)
		}
	}

	// Next, walk all sub-schemas
	for _, val := range schema.Definitions {
		if err := populateRefs(missingURIs, registry, baseURI, val); err != nil {
			return err
		}
	}

	if err := populateRefs(missingURIs, registry, baseURI, schema.Elements); err != nil {
		return err
	}

	for _, val := range schema.Properties {
		if err := populateRefs(missingURIs, registry, baseURI, val); err != nil {
			return err
		}
	}

	for _, val := range schema.OptionalProperties {
		if err := populateRefs(missingURIs, registry, baseURI, val); err != nil {
			return err
		}
	}

	if err := populateRefs(missingURIs, registry, baseURI, schema.Values); err != nil {
		return err
	}

	if schema.Disciminator != nil {
		for _, val := range schema.Disciminator.Mapping {
			if err := populateRefs(missingURIs, registry, baseURI, val); err != nil {
				return err
			}
		}
	}

	for _, val := range schema.AnyOf {
		if err := populateRefs(missingURIs, registry, baseURI, val); err != nil {
			return err
		}
	}

	return nil
}

func (v Validator) Validate(instance interface{}) (ValidationResult, error) {
	return v.ValidateURI(url.URL{}, instance)
}

func (v Validator) ValidateURI(uri url.URL, instance interface{}) (ValidationResult, error) {
	schema, ok := v.registry[uri]
	if !ok {
		return ValidationResult{}, ErrNoSuchSchema{URI: uri}
	}

	vm := vm{
		maxErrors:      v.maxErrors,
		maxDepth:       v.maxDepth,
		registry:       v.registry,
		instanceTokens: []string{},
		schemas: []schemaStack{
			schemaStack{
				uri:    uri,
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
