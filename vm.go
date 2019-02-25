package jsonvalidate

import (
	"net/url"

	"github.com/json-validate/json-pointer-go"
)

type vm struct {
	maxErrors      int
	maxDepth       int
	registry       map[url.URL]*Schema
	instanceTokens []string
	schemas        []schemaStack
	errors         []ValidationError
}

type schemaStack struct {
	uri    url.URL
	tokens []string
}

func (vm *vm) eval(schema *Schema, instance interface{}) error {
	switch instance.(type) {
	case nil:
		if schema.Type != nil && *schema.Type != "null" {
			vm.pushSchemaToken("type")
			if err := vm.reportError(); err != nil {
				return err
			}
			vm.popSchemaToken()
		}
	case bool:
		if schema.Type != nil && *schema.Type != "boolean" {
			vm.pushSchemaToken("type")
			if err := vm.reportError(); err != nil {
				return err
			}
			vm.popSchemaToken()
		}
	case float64:
		if schema.Type != nil && *schema.Type != "number" {
			vm.pushSchemaToken("type")
			if err := vm.reportError(); err != nil {
				return err
			}
			vm.popSchemaToken()
		}
	case string:
		if schema.Type != nil && *schema.Type != "string" {
			vm.pushSchemaToken("type")
			if err := vm.reportError(); err != nil {
				return err
			}
			vm.popSchemaToken()
		}
	case []interface{}:
		if schema.Type != nil {
			vm.pushSchemaToken("type")
			if err := vm.reportError(); err != nil {
				return err
			}
			vm.popSchemaToken()
		}
	case map[string]interface{}:
		if schema.Type != nil {
			vm.pushSchemaToken("type")
			if err := vm.reportError(); err != nil {
				return err
			}
			vm.popSchemaToken()
		}
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
		SchemaURI:    schemaStack.uri,
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

func (vm *vm) pushSchema(uri url.URL, tokens []string) error {
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
