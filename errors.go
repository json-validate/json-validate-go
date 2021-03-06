package jsonvalidate

import (
	"errors"
	"fmt"
	"net/url"
)

var ErrBadSubSchema = errors.New("invalid sub-schema")
var ErrBadSchemaKind = errors.New("invalid keyword combination")
var ErrMaxDepth = errors.New("max recursion depth reached during validation")
var errMaxErrors = errors.New("max errors reached")

type ErrMissingSchemas struct {
	URIs []url.URL
}

func (e ErrMissingSchemas) Error() string {
	return fmt.Sprintf("missing schemas: %v", e.URIs)
}
