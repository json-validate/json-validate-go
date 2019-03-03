package jsonvalidate

import (
	"net/url"
	"strconv"

	"github.com/json-validate/json-pointer-go"
)

type vm struct {
	maxErrors      int
	maxDepth       int
	registry       Registry
	instanceTokens []string
	schemas        []schemaStack
	errors         []ValidationError
}

type schemaStack struct {
	uri    *url.URL
	tokens []string
}

func (vm *vm) eval(schema *Schema, instance interface{}) error {
	if schema.RefSchema != nil {
		tokens := []string{}
		if schema.Ref.Fragment != "" {
			tokens = []string{"definitions", schema.Ref.Fragment}
		}

		if err := vm.pushSchema(schema.RefSchema.Base, tokens); err != nil {
			return err
		}

		if err := vm.eval(schema.RefSchema, instance); err != nil {
			return err
		}

		vm.popSchema()
	}

	switch schema.Kind {
	case SchemaKindEmpty:
		return nil
	case SchemaKindType:
		switch schema.Type {
		case SchemaTypeNull:
			if instance != nil {
				vm.pushSchemaToken("type")
				if err := vm.reportError(); err != nil {
					return err
				}
				vm.popSchemaToken()
			}
		case SchemaTypeBoolean:
			if _, ok := instance.(bool); !ok {
				vm.pushSchemaToken("type")
				if err := vm.reportError(); err != nil {
					return err
				}
				vm.popSchemaToken()
			}
		case SchemaTypeNumber:
			if _, ok := instance.(float64); !ok {
				vm.pushSchemaToken("type")
				if err := vm.reportError(); err != nil {
					return err
				}
				vm.popSchemaToken()
			}
		case SchemaTypeString:
			if _, ok := instance.(string); !ok {
				vm.pushSchemaToken("type")
				if err := vm.reportError(); err != nil {
					return err
				}
				vm.popSchemaToken()
			}
		}
	case SchemaKindElements:
		vm.pushSchemaToken("elements")

		if elems, ok := instance.([]interface{}); ok {
			for i, elem := range elems {
				vm.pushInstanceToken(strconv.Itoa(i))
				if err := vm.eval(schema.Elements, elem); err != nil {
					return err
				}
				vm.popInstanceToken()
			}
		} else {
			if err := vm.reportError(); err != nil {
				return err
			}
		}

		vm.popSchemaToken()
	case SchemaKindProperties:
		if object, ok := instance.(map[string]interface{}); ok {
			// First, required properties.
			vm.pushSchemaToken("properties")
			for property, subSchema := range schema.Properties {
				vm.pushSchemaToken(property)

				if value, ok := object[property]; ok {
					vm.pushInstanceToken(property)
					if err := vm.eval(subSchema, value); err != nil {
						return err
					}
					vm.popInstanceToken()
				} else {
					if err := vm.reportError(); err != nil {
						return err
					}
				}

				vm.popSchemaToken()
			}
			vm.popSchemaToken()

			// Then, optional properties.
			vm.pushSchemaToken("optionalProperties")
			for property, subSchema := range schema.OptionalProperties {
				vm.pushSchemaToken(property)

				if value, ok := object[property]; ok {
					vm.pushInstanceToken(property)
					if err := vm.eval(subSchema, value); err != nil {
						return err
					}
					vm.popInstanceToken()
				}

				vm.popSchemaToken()
			}
			vm.popSchemaToken()

		} else {
			// Which errors we're gonna produce has to do with which keywords appeared
			// in the schema.

			if schema.Properties != nil {
				vm.pushSchemaToken("properties")
				if err := vm.reportError(); err != nil {
					return err
				}
				vm.popSchemaToken()
			}

			if schema.OptionalProperties != nil {
				vm.pushSchemaToken("optionalProperties")
				if err := vm.reportError(); err != nil {
					return err
				}
				vm.popSchemaToken()
			}
		}
	case SchemaKindValues:
		vm.pushSchemaToken("values")

		if object, ok := instance.(map[string]interface{}); ok {
			for key, value := range object {
				vm.pushInstanceToken(key)
				if err := vm.eval(schema.Values, value); err != nil {
					return err
				}
				vm.popInstanceToken()
			}
		} else {
			if err := vm.reportError(); err != nil {
				return err
			}
		}

		vm.popSchemaToken()
	case SchemaKindDiscriminator:
		vm.pushSchemaToken("discriminator")

		if object, ok := instance.(map[string]interface{}); ok {
			if prop, ok := object[schema.DiscriminatorPropertyName]; ok {
				if propStr, ok := prop.(string); ok {
					if subSchema, ok := schema.DiscriminatorMapping[propStr]; ok {
						vm.pushSchemaToken("mapping")
						vm.pushSchemaToken(propStr)
						if err := vm.eval(subSchema, instance); err != nil {
							return err
						}
						vm.popSchemaToken()
						vm.popSchemaToken()
					} else {
						vm.pushSchemaToken("mapping")
						vm.pushInstanceToken(schema.DiscriminatorPropertyName)
						if err := vm.reportError(); err != nil {
							return err
						}
						vm.popInstanceToken()
						vm.popSchemaToken()
					}
				} else {
					vm.pushSchemaToken("propertyName")
					vm.pushInstanceToken(schema.DiscriminatorPropertyName)
					if err := vm.reportError(); err != nil {
						return err
					}
					vm.popInstanceToken()
					vm.popSchemaToken()
				}
			} else {
				vm.pushSchemaToken("propertyName")
				if err := vm.reportError(); err != nil {
					return err
				}
				vm.popSchemaToken()
			}

		} else {
			if err := vm.reportError(); err != nil {
				return err
			}
		}

		vm.popSchemaToken()
	}

	return nil
}

func (vm *vm) reportError() error {
	schemaStack := vm.schemas[len(vm.schemas)-1]
	instancePath := make([]string, len(vm.instanceTokens))
	schemaPath := make([]string, len(schemaStack.tokens))

	copy(instancePath, vm.instanceTokens)
	copy(schemaPath, schemaStack.tokens)

	vm.errors = append(vm.errors, ValidationError{
		InstancePath: jsonpointer.Ptr{Tokens: instancePath},
		SchemaPath:   jsonpointer.Ptr{Tokens: schemaPath},
		SchemaURI:    *schemaStack.uri,
	})

	if len(vm.errors) == vm.maxErrors {
		return errMaxErrors
	}

	return nil
}

func (vm *vm) pushInstanceToken(t string) {
	vm.instanceTokens = append(vm.instanceTokens, t)
}

func (vm *vm) popInstanceToken() {
	vm.instanceTokens = vm.instanceTokens[:len(vm.instanceTokens)-1]
}

func (vm *vm) pushSchema(uri *url.URL, tokens []string) error {
	if len(vm.schemas) == vm.maxDepth {
		return ErrMaxDepth
	}

	vm.schemas = append(vm.schemas, schemaStack{
		uri:    uri,
		tokens: tokens,
	})

	return nil
}

func (vm *vm) popSchema() {
	vm.schemas = vm.schemas[:len(vm.schemas)-1]
}

func (vm *vm) pushSchemaToken(t string) {
	stack := &vm.schemas[len(vm.schemas)-1]
	stack.tokens = append(stack.tokens, t)
}

func (vm *vm) popSchemaToken() {
	stack := &vm.schemas[len(vm.schemas)-1]
	stack.tokens = stack.tokens[:len(stack.tokens)-1]
}
