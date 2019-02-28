# json-validate

This package provides a Golang implementation of JSON Validate. In particular,
it offers:

1. A validator you can use to check that data is valid against a schema. It has
   full support for the standard JSON Validate error format.
1. Structs you can you to parse out JSON Validate schemas.
1. Functions to validate that a schema is "valid" beyond just having the right
   property names and types.
1. Utility functions to iterate over JSON Schemas, so you can implement your own
   code generation or UI generation on top of JSON Validate.

## Usage

For detailed usage, see the reference docs:

https://godoc.org/github.com/json-validate/json-validate-go

### Parsing JSON Validate schemas

If you're storing your schemas as JSON or YAML, you can use the `SchemaStruct`
type to deserialize them:

```go
func main() {
  data := `
    {
      "
    }
  `
}
```
