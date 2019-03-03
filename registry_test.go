package jsonvalidate

import (
	e "errors"
	"net/url"
	"strconv"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewRegistryRef(t *testing.T) {
	// Test the happy case for creating a new registry. This involves creating
	// cyclical data structures, so it's difficult to express this as a table
	// test.
	registry1 := Schema{
		IsRoot: true,
		ID:     &url.URL{Scheme: "http", Host: "example.com", Path: "/foo"},
		Definitions: map[string]*Schema{
			"a": &Schema{
				Kind: SchemaKindRef,
				Ref:  &url.URL{},
			},
			"b": &Schema{
				Kind: SchemaKindRef,
				Ref:  &url.URL{Scheme: "http", Host: "example.com", Path: "/foo"},
			},
			"c": &Schema{
				Kind: SchemaKindRef,
				Ref:  &url.URL{Fragment: "a"},
			},
			"d": &Schema{
				Kind: SchemaKindRef,
				Ref:  &url.URL{Scheme: "http", Host: "example.com", Path: "/foo", Fragment: "a"},
			},
			"e": &Schema{
				Kind: SchemaKindRef,
				Ref:  &url.URL{Scheme: "http", Host: "example.com", Path: "/bar"},
			},
		},
		Kind: SchemaKindEmpty,
	}

	registry2 := Schema{
		IsRoot: true,
		ID:     &url.URL{Scheme: "http", Host: "example.com", Path: "/bar"},
		Kind:   SchemaKindEmpty,
	}

	registry1.Definitions["a"].RefSchema = &registry1
	registry1.Definitions["b"].RefSchema = &registry1
	registry1.Definitions["c"].RefSchema = registry1.Definitions["a"]
	registry1.Definitions["d"].RefSchema = registry1.Definitions["a"]
	registry1.Definitions["e"].RefSchema = &registry2

	expected := Registry{
		Schemas: map[url.URL]*Schema{
			url.URL{Scheme: "http", Host: "example.com", Path: "/foo"}: &registry1,
			url.URL{Scheme: "http", Host: "example.com", Path: "/bar"}: &registry2,
		},
	}

	id1 := "http://example.com/foo"
	id2 := "http://example.com/bar"
	refA := ""
	refB := "http://example.com/foo"
	refC := "#a"
	refD := "http://example.com/foo#a"
	refE := "http://example.com/bar"

	input := []SchemaStruct{
		SchemaStruct{
			ID: &id1,
			Definitions: &map[string]SchemaStruct{
				"a": SchemaStruct{Ref: &refA},
				"b": SchemaStruct{Ref: &refB},
				"c": SchemaStruct{Ref: &refC},
				"d": SchemaStruct{Ref: &refD},
				"e": SchemaStruct{Ref: &refE},
			},
		},
		SchemaStruct{
			ID: &id2,
		},
	}

	actual, err := NewRegistry(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestNewRegistry(t *testing.T) {
	badURI := "::"
	emptyURI := ""
	typeNull := "null"
	undefinedURI1 := "http://example.com/foo"
	undefinedURI2 := "http://example.com/bar"

	testCases := []struct {
		in  []SchemaStruct
		out Registry
		err error
	}{
		{
			[]SchemaStruct{},
			Registry{Schemas: map[url.URL]*Schema{}},
			nil,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Type: &typeNull,
				},
			},
			Registry{
				Schemas: map[url.URL]*Schema{
					url.URL{}: &Schema{
						IsRoot: true,
						ID:     &url.URL{},
						Kind:   SchemaKindType,
						Type:   SchemaTypeNull,
					},
				},
			},
			nil,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Elements: &SchemaStruct{},
				},
			},
			Registry{
				Schemas: map[url.URL]*Schema{
					url.URL{}: &Schema{
						IsRoot: true,
						ID:     &url.URL{},
						Kind:   SchemaKindElements,
						Elements: &Schema{
							Kind: SchemaKindEmpty,
						},
					},
				},
			},
			nil,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Properties: &map[string]SchemaStruct{
						"a": SchemaStruct{},
					},
					OptionalProperties: &map[string]SchemaStruct{
						"a": SchemaStruct{},
					},
				},
			},
			Registry{
				Schemas: map[url.URL]*Schema{
					url.URL{}: &Schema{
						IsRoot: true,
						ID:     &url.URL{},
						Kind:   SchemaKindProperties,
						Properties: map[string]*Schema{
							"a": &Schema{Kind: SchemaKindEmpty},
						},
						OptionalProperties: map[string]*Schema{
							"a": &Schema{Kind: SchemaKindEmpty},
						},
					},
				},
			},
			nil,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Values: &SchemaStruct{},
				},
			},
			Registry{
				Schemas: map[url.URL]*Schema{
					url.URL{}: &Schema{
						IsRoot: true,
						ID:     &url.URL{},
						Kind:   SchemaKindValues,
						Values: &Schema{
							Kind: SchemaKindEmpty,
						},
					},
				},
			},
			nil,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Discriminator: &SchemaStructDiscriminator{
						PropertyName: "::",
						Mapping: map[string]SchemaStruct{
							"a": SchemaStruct{},
						},
					},
				},
			},
			Registry{
				Schemas: map[url.URL]*Schema{
					url.URL{}: &Schema{
						IsRoot: true,
						ID:     &url.URL{},
						Kind:   SchemaKindDiscriminator,
						DiscriminatorPropertyName: "::",
						DiscriminatorMapping: map[string]*Schema{
							"a": &Schema{Kind: SchemaKindEmpty},
						},
					},
				},
			},
			nil,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Extra: map[string]interface{}{
						"a": []interface{}{1.0},
						"b": map[string]interface{}{},
						"c": nil,
					},
				},
			},
			Registry{
				Schemas: map[url.URL]*Schema{
					url.URL{}: &Schema{
						IsRoot: true,
						ID:     &url.URL{},
						Kind:   SchemaKindEmpty,
						Extra: map[string]interface{}{
							"a": []interface{}{1.0},
							"b": map[string]interface{}{},
							"c": nil,
						},
					},
				},
			},
			nil,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					ID: &badURI,
				},
			},
			Registry{},
			&url.Error{
				Op:  "parse",
				URL: "::",
				Err: e.New("missing protocol scheme"),
			},
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Ref: &badURI,
				},
			},
			Registry{},
			&url.Error{
				Op:  "parse",
				URL: "::",
				Err: e.New("missing protocol scheme"),
			},
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Definitions: &map[string]SchemaStruct{
						"a": SchemaStruct{
							Ref: &badURI,
						},
					},
				},
			},
			Registry{},
			&url.Error{
				Op:  "parse",
				URL: "::",
				Err: e.New("missing protocol scheme"),
			},
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Type: &badURI,
				},
			},
			Registry{},
			e.New("invalid type: ::"),
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Elements: &SchemaStruct{
						Ref: &badURI,
					},
				},
			},
			Registry{},
			&url.Error{
				Op:  "parse",
				URL: "::",
				Err: e.New("missing protocol scheme"),
			},
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Properties: &map[string]SchemaStruct{
						"a": SchemaStruct{
							Ref: &badURI,
						},
					},
				},
			},
			Registry{},
			&url.Error{
				Op:  "parse",
				URL: "::",
				Err: e.New("missing protocol scheme"),
			},
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					OptionalProperties: &map[string]SchemaStruct{
						"a": SchemaStruct{
							Ref: &badURI,
						},
					},
				},
			},
			Registry{},
			&url.Error{
				Op:  "parse",
				URL: "::",
				Err: e.New("missing protocol scheme"),
			},
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Values: &SchemaStruct{
						Ref: &badURI,
					},
				},
			},
			Registry{},
			&url.Error{
				Op:  "parse",
				URL: "::",
				Err: e.New("missing protocol scheme"),
			},
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Discriminator: &SchemaStructDiscriminator{
						PropertyName: "::",
						Mapping: map[string]SchemaStruct{
							"a": SchemaStruct{
								Ref: &badURI,
							},
						},
					},
				},
			},
			Registry{},
			&url.Error{
				Op:  "parse",
				URL: "::",
				Err: e.New("missing protocol scheme"),
			},
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Definitions: &map[string]SchemaStruct{
						"a": SchemaStruct{
							ID: &emptyURI,
						},
					},
				},
			},
			Registry{},
			ErrBadSubSchema,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Elements: &SchemaStruct{
						ID: &emptyURI,
					},
				},
			},
			Registry{},
			ErrBadSubSchema,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Properties: &map[string]SchemaStruct{
						"a": SchemaStruct{
							ID: &emptyURI,
						},
					},
				},
			},
			Registry{},
			ErrBadSubSchema,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					OptionalProperties: &map[string]SchemaStruct{
						"a": SchemaStruct{
							ID: &emptyURI,
						},
					},
				},
			},
			Registry{},
			ErrBadSubSchema,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Values: &SchemaStruct{
						ID: &emptyURI,
					},
				},
			},
			Registry{},
			ErrBadSubSchema,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Definitions: &map[string]SchemaStruct{
						"a": SchemaStruct{
							ID: &emptyURI,
						},
					},
				},
			},
			Registry{},
			ErrBadSubSchema,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Elements: &SchemaStruct{
						Definitions: &map[string]SchemaStruct{},
					},
				},
			},
			Registry{},
			ErrBadSubSchema,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Properties: &map[string]SchemaStruct{
						"a": SchemaStruct{
							Definitions: &map[string]SchemaStruct{},
						},
					},
				},
			},
			Registry{},
			ErrBadSubSchema,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					OptionalProperties: &map[string]SchemaStruct{
						"a": SchemaStruct{
							Definitions: &map[string]SchemaStruct{},
						},
					},
				},
			},
			Registry{},
			ErrBadSubSchema,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Values: &SchemaStruct{
						Definitions: &map[string]SchemaStruct{},
					},
				},
			},
			Registry{},
			ErrBadSubSchema,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Discriminator: &SchemaStructDiscriminator{
						PropertyName: "::",
						Mapping: map[string]SchemaStruct{
							"a": SchemaStruct{
								Definitions: &map[string]SchemaStruct{},
							},
						},
					},
				},
			},
			Registry{},
			ErrBadSubSchema,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Ref:  &emptyURI,
					Type: &typeNull,
				},
			},
			Registry{},
			ErrBadSchemaKind,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Type:     &typeNull,
					Elements: &SchemaStruct{},
				},
			},
			Registry{},
			ErrBadSchemaKind,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Elements:           &SchemaStruct{},
					Properties:         &map[string]SchemaStruct{},
					OptionalProperties: &map[string]SchemaStruct{},
				},
			},
			Registry{},
			ErrBadSchemaKind,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Properties:         &map[string]SchemaStruct{},
					OptionalProperties: &map[string]SchemaStruct{},
					Values:             &SchemaStruct{},
				},
			},
			Registry{},
			ErrBadSchemaKind,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Values:        &SchemaStruct{},
					Discriminator: &SchemaStructDiscriminator{},
				},
			},
			Registry{},
			ErrBadSchemaKind,
		},
		{
			[]SchemaStruct{
				SchemaStruct{
					Ref: &undefinedURI1,
					Definitions: &map[string]SchemaStruct{
						"a": SchemaStruct{
							Ref: &undefinedURI2,
						},
					},
				},
			},
			Registry{},
			ErrMissingSchemas{
				URIs: []url.URL{
					url.URL{Scheme: "http", Host: "example.com", Path: "/foo"},
					url.URL{Scheme: "http", Host: "example.com", Path: "/bar"},
				},
			},
		},
	}

	for i, tt := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			out, err := NewRegistry(tt.in)
			assert.Equal(t, tt.out, out)
			assert.Equal(t, tt.err, errors.Cause(err))
		})
	}
}
