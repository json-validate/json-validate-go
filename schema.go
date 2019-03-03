package jsonvalidate

import (
	"encoding/json"
	"net/url"
)

// SchemaStruct is a JSON-friendly representation of a JSON Validate schema.
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

// Schema is an abstract representation of a JSON Validate schema.
//
// Whereas SchemaStruct is meant for marshaling/unmarshaling data, Schema is
// meant for higher-level processing of schemas. If you intend to manipulate
// schemas in order to perform code or UI generation, you should use Schema.
type Schema struct {
	// The base URI of the schema; this is the ID of the schema's root.
	Base *url.URL

	// Whether this schema is a root schema. Only root schemas may have nonzero ID
	// or Definitions.
	IsRoot bool

	// The ID of the Schema. ID is meaningful iff IsRoot is true. The Fragment
	// part is guaranteed to be empty.
	ID *url.URL

	// The definitions for this Schema. Definitions is meaningful iff IsRoot is
	// true.
	Definitions map[string]*Schema

	// Indicates which keywords may be set on this schema.
	Kind SchemaKind

	// Meaningful iff Kind is SchemaKindRef.
	Ref       *url.URL // the parsed URI that was referred to
	RefSchema *Schema  // the schema the ref resolved to

	// Meaningful iff Kind is SchemaKindType
	Type SchemaType

	// Meaningful iff Kind is SchemaKindElements
	Elements *Schema

	// Meaningful iff Kind is SchemaKindProperties
	Properties         map[string]*Schema // required properties
	OptionalProperties map[string]*Schema // optional properties

	// Meaningful iff Kind is SchemaKindValues
	Values *Schema

	// Meaningful iff Kind is SchemaKindDiscriminator
	DiscriminatorPropertyName string             // property to switch on
	DiscriminatorMapping      map[string]*Schema // mapping from value to schema

	// Extra stores data that's in a schema, but isn't part of the formal spec.
	//
	// If a schema contains non-formalized data like `title` or `description`,
	// those data will be present in this field.
	Extra map[string]interface{}
}

// SchemaKind is an enum of possible types of schemas.
//
// Most JSON Validate keywords are mutually exclusive. This enum serves to
// indicate which keywords amy possibly appear.
type SchemaKind int

const (
	// SchemaKindEmpty indicates a schema without any keywords.
	SchemaKindEmpty SchemaKind = iota

	// SchemaKindRef indicates a schema with the "ref" keyword.
	SchemaKindRef

	// SchemaKindType indicates a schema with the "type" keyword.
	SchemaKindType

	// SchemaKindElements indicates a schema with the "elements" keyword.
	SchemaKindElements

	// SchemaKindProperties indicates a schema with the "properties" or
	// "optionalProperties" keywords, or both.
	SchemaKindProperties

	// SchemaKindValues indicates a schema with the "values" keyword.
	SchemaKindValues

	// SchemaKindDiscriminator indidcates a schema with the "discriminator"
	// keyword.
	SchemaKindDiscriminator
)

// SchemaType indicates possible values of the "type" keyword.
type SchemaType int

const (
	// SchemaTypeNull indicates the value "null".
	SchemaTypeNull SchemaType = iota

	// SchemaTypeBoolean indicates the value "boolean".
	SchemaTypeBoolean

	// SchemaTypeNumber indicates the value "number".
	SchemaTypeNumber

	// SchemaTypeString indicates the type "string".
	SchemaTypeString
)
