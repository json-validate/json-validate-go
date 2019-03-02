package jsonvalidate

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

// Registry is a collection of schemas which may refer to each other.
type Registry struct {
	Schemas []Schema
}

// NewRegistry constructs a new registry from a set of schemas.
//
// It is guaranteed that schemas within the returned registry shall point to,
// and be pointed to by, only one another. Therefore, if the caller does not
// create new pointers into the registry's schemas, then discarding the registry
// will allow all of the contained schemas to be garbage collected.
func NewRegistry(schemaStructs []SchemaStruct) (Registry, error) {
	schemas := []Schema{}
	for i, schema := range schemaStructs {
		s, err := parseSchemaStruct(true, schema)
		if err != nil {
			return Registry{}, errors.Wrapf(err, "error parsing schema at index %d", i)
		}

		schemas = append(schemas, s)
	}

	return Registry{Schemas: schemas}, nil
}

func parseSchemaStruct(root bool, s SchemaStruct) (Schema, error) {
	out := Schema{}

	if s.ID != nil {
		if !root {
			return Schema{}, ErrBadSubSchema
		}

		id, err := url.Parse(*s.ID)
		if err != nil {
			return Schema{}, err
		}

		out.ID = id
	}

	if s.Definitions != nil {
		if !root {
			return Schema{}, ErrBadSubSchema
		}

		out.Definitions = make(map[string]*Schema, len(*s.Definitions))
		for k, v := range *s.Definitions {
			schema, err := parseSchemaStruct(false, v)
			if err != nil {
				return Schema{}, errors.Wrapf(err, "error parsing definition %s", k)
			}

			out.Definitions[k] = &schema
		}
	}

	out.Kind = SchemaKindEmpty

	if s.Ref != nil {
		if out.Kind == SchemaKindEmpty {
			out.Kind = SchemaKindRef
		} else {
			return Schema{}, ErrBadSchemaKind
		}

		ref, err := url.Parse(*s.Ref)
		if err != nil {
			return Schema{}, err
		}

		out.Ref = ref
	}

	if s.Type != nil {
		if out.Kind == SchemaKindEmpty {
			out.Kind = SchemaKindType
		} else {
			return Schema{}, ErrBadSchemaKind
		}

		switch *s.Type {
		case "null":
			out.Type = SchemaTypeNull
		case "boolean":
			out.Type = SchemaTypeBoolean
		case "number":
			out.Type = SchemaTypeNumber
		case "string":
			out.Type = SchemaTypeString
		default:
			return Schema{}, fmt.Errorf("invalid type: %s", *s.Type)
		}
	}

	if s.Elements != nil {
		if out.Kind == SchemaKindEmpty {
			out.Kind = SchemaKindElements
		} else {
			return Schema{}, ErrBadSchemaKind
		}

		schema, err := parseSchemaStruct(false, *s.Elements)
		if err != nil {
			return Schema{}, errors.Wrap(err, "error parsing elements")
		}

		out.Elements = &schema
	}

	if s.Properties != nil {
		if out.Kind == SchemaKindEmpty {
			out.Kind = SchemaKindProperties
		} else {
			return Schema{}, ErrBadSchemaKind
		}

		out.Properties = make(map[string]*Schema, len(*s.Properties))
		for k, v := range *s.Properties {
			schema, err := parseSchemaStruct(false, v)
			if err != nil {
				return Schema{}, errors.Wrapf(err, "error parsing properties %s", k)
			}

			out.Properties[k] = &schema
		}
	}

	if s.OptionalProperties != nil {
		if out.Kind == SchemaKindEmpty || out.Kind == SchemaKindProperties {
			out.Kind = SchemaKindProperties
		} else {
			return Schema{}, ErrBadSchemaKind
		}

		out.OptionalProperties = make(map[string]*Schema, len(*s.OptionalProperties))
		for k, v := range *s.OptionalProperties {
			schema, err := parseSchemaStruct(false, v)
			if err != nil {
				return Schema{}, errors.Wrapf(err, "error parsing optionalProperties %s", k)
			}

			out.OptionalProperties[k] = &schema
		}
	}

	if s.Values != nil {
		if out.Kind == SchemaKindEmpty {
			out.Kind = SchemaKindValues
		} else {
			return Schema{}, ErrBadSchemaKind
		}

		schema, err := parseSchemaStruct(false, *s.Values)
		if err != nil {
			return Schema{}, errors.Wrap(err, "error parsing values")
		}

		out.Values = &schema
	}

	if s.Discriminator != nil {
		if out.Kind == SchemaKindEmpty {
			out.Kind = SchemaKindDiscriminator
		} else {
			return Schema{}, ErrBadSchemaKind
		}

		out.DiscriminatorPropertyName = s.Discriminator.PropertyName
		out.DiscriminatorMapping = make(map[string]*Schema, len(s.Discriminator.Mapping))
		for k, v := range s.Discriminator.Mapping {
			schema, err := parseSchemaStruct(false, v)
			if err != nil {
				return Schema{}, errors.Wrapf(err, "error parsing mapping %s", k)
			}

			out.DiscriminatorMapping[k] = &schema
		}
	}

	return Schema{}, nil
}
