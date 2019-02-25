package jsonvalidate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpec(t *testing.T) {
	type testCase struct {
		Name      string   `json:"name"`
		Registry  []Schema `json:"registry"`
		Schema    Schema   `json:"schema"`
		Instances []struct {
			Instance interface{}       `json:"instance"`
			Errors   []ValidationError `json:"errors"`
		} `json:"instances"`
	}

	err := filepath.Walk("spec/tests", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		var testCases []testCase
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&testCases)
		if err != nil {
			return err
		}

		for _, tt := range testCases {
			t.Run(path+"/"+tt.Name, func(t *testing.T) {
				for i, instance := range tt.Instances {
					t.Run(strconv.Itoa(i), func(t *testing.T) {
						schemas := []Schema{}
						schemas = append(schemas, tt.Registry...)
						schemas = append(schemas, tt.Schema)

						validator, err := NewValidator(schemas)
						assert.NoError(t, err)

						result, err := validator.Validate(instance.Instance)
						assert.NoError(t, err)
						assert.Equal(t, instance.Errors, result.Errors)
					})
				}
			})
		}

		return nil
	})

	assert.NoError(t, err)
}
