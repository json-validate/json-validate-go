package jsonvalidate

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaJSONRoundtrip(t *testing.T) {
	// you can't take a pointer to a string literal, so here we are
	var strEmpty = ""
	var strA = "a"

	testCases := []struct {
		in  string
		out SchemaStruct
	}{
		{
			`{}`,
			SchemaStruct{
				Extra: map[string]interface{}{},
			},
		},
		{
			`{"id":""}`,
			SchemaStruct{
				ID:    &strEmpty,
				Extra: map[string]interface{}{},
			},
		},
		{
			`{"id":"a"}`,
			SchemaStruct{
				ID:    &strA,
				Extra: map[string]interface{}{},
			},
		},
		{
			`{"ref":"a"}`,
			SchemaStruct{
				Ref:   &strA,
				Extra: map[string]interface{}{},
			},
		},
		{
			`{"definitions":{"a":{"id":"a"},"b":{"ref":"a"}}}`,
			SchemaStruct{
				Definitions: &map[string]SchemaStruct{
					"a": SchemaStruct{
						ID:    &strA,
						Extra: map[string]interface{}{},
					},
					"b": SchemaStruct{
						Ref:   &strA,
						Extra: map[string]interface{}{},
					},
				},
				Extra: map[string]interface{}{},
			},
		},
		{
			`{"type":"a"}`,
			SchemaStruct{
				Type:  &strA,
				Extra: map[string]interface{}{},
			},
		},
		{
			`{"elements":{"id":""}}`,
			SchemaStruct{
				Elements: &SchemaStruct{
					ID:    &strEmpty,
					Extra: map[string]interface{}{},
				},
				Extra: map[string]interface{}{},
			},
		},
		{
			`{"properties":{"a":{"id":"a"},"b":{"ref":"a"}}}`,
			SchemaStruct{
				Properties: &map[string]SchemaStruct{
					"a": SchemaStruct{
						ID:    &strA,
						Extra: map[string]interface{}{},
					},
					"b": SchemaStruct{
						Ref:   &strA,
						Extra: map[string]interface{}{},
					},
				},
				Extra: map[string]interface{}{},
			},
		},
		{
			`{"optionalProperties":{"a":{"id":"a"},"b":{"ref":"a"}}}`,
			SchemaStruct{
				OptionalProperties: &map[string]SchemaStruct{
					"a": SchemaStruct{
						ID:    &strA,
						Extra: map[string]interface{}{},
					},
					"b": SchemaStruct{
						Ref:   &strA,
						Extra: map[string]interface{}{},
					},
				},
				Extra: map[string]interface{}{},
			},
		},
		{
			`{"values":{"id":""}}`,
			SchemaStruct{
				Values: &SchemaStruct{
					ID:    &strEmpty,
					Extra: map[string]interface{}{},
				},
				Extra: map[string]interface{}{},
			},
		},
		{
			`{"discriminator":{"propertyName":"a","mapping":{"a":{"id":"a"},"b":{"ref":"a"}}}}`,
			SchemaStruct{
				Discriminator: &SchemaStructDiscriminator{
					PropertyName: "a",
					Mapping: map[string]SchemaStruct{
						"a": SchemaStruct{
							ID:    &strA,
							Extra: map[string]interface{}{},
						},
						"b": SchemaStruct{
							Ref:   &strA,
							Extra: map[string]interface{}{},
						},
					},
				},
				Extra: map[string]interface{}{},
			},
		},
		{
			`{"a":[1],"b":{},"c":null,"definitions":{},"discriminator":{"propertyName":"","mapping":{}},"elements":{},"id":"","optionalProperties":{},"properties":{},"ref":"","type":"","values":{}}`,
			SchemaStruct{
				ID:                 &strEmpty,
				Ref:                &strEmpty,
				Definitions:        &map[string]SchemaStruct{},
				Type:               &strEmpty,
				Elements:           &SchemaStruct{Extra: map[string]interface{}{}},
				Properties:         &map[string]SchemaStruct{},
				OptionalProperties: &map[string]SchemaStruct{},
				Values:             &SchemaStruct{Extra: map[string]interface{}{}},
				Discriminator:      &SchemaStructDiscriminator{PropertyName: "", Mapping: map[string]SchemaStruct{}},
				Extra: map[string]interface{}{
					"a": []interface{}{1.0},
					"b": map[string]interface{}{},
					"c": nil,
				},
			},
		},
	}

	for i, tt := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var out SchemaStruct
			err := json.Unmarshal([]byte(tt.in), &out)
			assert.NoError(t, err)
			assert.Equal(t, tt.out, out)

			in, err := json.Marshal(out)
			assert.NoError(t, err)
			assert.Equal(t, []byte(tt.in), in)
		})
	}
}
