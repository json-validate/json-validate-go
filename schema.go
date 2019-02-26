package jsonvalidate

import "net/url"

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
