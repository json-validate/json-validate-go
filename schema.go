package jsonvalidate

import (
	"encoding/json"
	"net/url"
)

// SchemaStruct is a JSON- and YAML-friendly representation of a JSON Validate
// schema.
//
// Just because data is a SchemaStruct, doesn't mean it's a valid JSON Validate
// schema. For example, some keywords of JSON Schema aren't allowed to be used
// alongside each other, the "type" keyword only has a specific set of valid
// values it can take on, the "id" has to be a valid URI, et cetera.
//
// If you want to make sure your schema is meaningful, you'll need to use
// NewRegistry to construct a registry using your schema. That function handles
// checking all the rules about what makes a schema valid.
type SchemaStruct struct {
	ID                 *string                    `json:"id,omitempty"`
	Ref                *string                    `json:"ref,omitempty"`
	Definitions        *map[string]SchemaStruct   `json:"definitions,omitempty"`
	Type               *string                    `json:"type,omitempty"`
	Elements           *SchemaStruct              `json:"elements,omitempty"`
	Properties         *map[string]SchemaStruct   `json:"properties,omitempty"`
	OptionalProperties *map[string]SchemaStruct   `json:"optionalProperties,omitempty"`
	Values             *SchemaStruct              `json:"values,omitempty"`
	Discriminator      *SchemaStructDiscriminator `json:"discriminator,omitempty"`

	// Extra stores data that's in a schema, but isn't part of the formal spec.
	//
	// If a schema contains non-formalized data like `title` or `description`,
	// those data will be present in this field.
	Extra map[string]interface{} `json:"-"`
}

// SchemaStructDiscriminator represents the "discriminator" keyword value of a
// SchemaStruct.
type SchemaStructDiscriminator struct {
	PropertyName string                  `json:"propertyName"`
	Mapping      map[string]SchemaStruct `json:"mapping"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (s *SchemaStruct) UnmarshalJSON(data []byte) error {
	type schemaStruct SchemaStruct
	var raw schemaStruct
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var extra map[string]interface{}
	if err := json.Unmarshal(data, &extra); err != nil {
		return err
	}

	s.ID = raw.ID
	delete(extra, "id")

	s.Ref = raw.Ref
	delete(extra, "ref")

	s.Definitions = raw.Definitions
	delete(extra, "definitions")

	s.Type = raw.Type
	delete(extra, "type")

	s.Elements = raw.Elements
	delete(extra, "elements")

	s.Properties = raw.Properties
	delete(extra, "properties")

	s.OptionalProperties = raw.OptionalProperties
	delete(extra, "optionalProperties")

	s.Values = raw.Values
	delete(extra, "values")

	s.Discriminator = raw.Discriminator
	delete(extra, "discriminator")

	s.Extra = extra
	return nil
}

// MarshalJSON satisfies the json.Marshaler interface.
func (s SchemaStruct) MarshalJSON() ([]byte, error) {
	out := make(map[string]interface{}, len(s.Extra)+2)
	for k, v := range s.Extra {
		out[k] = v
	}

	if s.ID != nil {
		out["id"] = s.ID
	}

	if s.Ref != nil {
		out["ref"] = s.Ref
	}

	if s.Definitions != nil {
		out["definitions"] = s.Definitions
	}

	if s.Type != nil {
		out["type"] = s.Type
	}

	if s.Elements != nil {
		out["elements"] = s.Elements
	}

	if s.Properties != nil {
		out["properties"] = s.Properties
	}

	if s.OptionalProperties != nil {
		out["optionalProperties"] = s.OptionalProperties
	}

	if s.Values != nil {
		out["values"] = s.Values
	}

	if s.Discriminator != nil {
		out["discriminator"] = s.Discriminator
	}

	return json.Marshal(out)
}

type Schema struct {
	ID                 *string            `json:"id"`
	Ref                *string            `json:"ref"`
	Definitions        map[string]*Schema `json:"definitions"`
	Type               *string            `json:"type"`
	Elements           *Schema            `json:"elements"`
	Properties         map[string]*Schema `json:"properties"`
	OptionalProperties map[string]*Schema `json:"optionalProperties"`
	Values             *Schema            `json:"values"`
	Discriminator      *Discriminator     `json:"discriminator"`
	AnyOf              []*Schema          `json:"anyOf"`

	baseURI url.URL
	refURI  url.URL
}

type Discriminator struct {
	PropertyName string             `json:"propertyName"`
	Mapping      map[string]*Schema `json:"mapping"`
}
