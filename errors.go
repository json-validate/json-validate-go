package jsonvalidate

import (
	"errors"
	"fmt"
	"net/url"
)

var ErrMaxDepth = errors.New("max recursion depth reached during validation")
var errMaxErrors = errors.New("max errors reached")

type ErrMissingSchemas struct {
	URIs []url.URL
}

func (e ErrMissingSchemas) Error() string {
	return fmt.Sprintf("missing schemas: %v", e.URIs)
}

type ErrNoSuchSchema struct {
	URI url.URL
}

func (e ErrNoSuchSchema) Error() string {
	return fmt.Sprintf("no such schema: %s", e.URI.String())
}
