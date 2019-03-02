package jsonvalidate

import (
	e "errors"
	"net/url"
	"strconv"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewRegistry(t *testing.T) {
	badURI := "::"
	emptyURI := ""

	testCases := []struct {
		in  []SchemaStruct
		out Registry
		err error
	}{
		{
			[]SchemaStruct{},
			Registry{Schemas: []Schema{}},
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
					Discriminator: &SchemaStructDiscriminator{
						PropertyName: "::",
						Mapping: map[string]SchemaStruct{
							"a": SchemaStruct{
								ID: &emptyURI,
							},
						},
					},
				},
			},
			Registry{},
			ErrBadSubSchema,
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
