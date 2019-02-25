package jsonvalidate

import "net/url"

type Schema struct {
	Base               *string            `json:"base"`
	Ref                *string            `json:"ref"`
	Definitions        map[string]*Schema `json:"definitions"`
	Type               *string            `json:"type"`
	Elements           *Schema            `json:"elements"`
	Properties         map[string]*Schema `json:"properties"`
	OptionalProperties map[string]*Schema `json:"optionalProperties"`
	Values             *Schema            `json:"values"`
	Disciminator       *Disciminator      `json:"discriminator"`
	AnyOf              []*Schema          `json:"anyOf"`

	baseURI url.URL
	refURI  url.URL
}

type Disciminator struct {
	PropertyName string             `json:"propertyName"`
	Mappings     map[string]*Schema `json:"mappings"`
}
