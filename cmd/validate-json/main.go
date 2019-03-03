package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/json-validate/json-pointer-go"
	jsonvalidate "github.com/json-validate/json-validate-go"
	"github.com/urfave/cli"
)

const exampleMessage = `
EXAMPLE:
		 Parse STDIN, and ensure that it is valid according to schema.json:

					validate-json schema.json

		 Same as above, but do not output any errors. Just as before, the exit code
		 will be nonzero if validation fails:

          validate-json -q schema.json

		 Parse STDIN, and ensure that it is valid according to schema.json -- but
		 schema.json relies on stuff defined in defs1.json and defs2.json:

					validate-json defs1.json defs2.json schema.json

		 If you want to validate STDIN against something other than the default
		 schema (the default schema is the one that doesn't use the "base" keyword),
		 then you can specify that URI using --schema-uri or -u:

					validate-json -u http://foo.com/bar defs1.json defs2.json schema.json

		 The order of the arguments in the two examples above does not matter.
`

type outputFormat int

const (
	outputFormatString outputFormat = iota + 1
	outputFormatJSON
)

func main() {
	app := cli.NewApp()
	app.Name = "validate-json"
	app.Usage = "Validate JSON data against a JSON Validate schema"
	app.ArgsUsage = "schemas..."

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "suppress error messages",
		},
		cli.StringFlag{
			Name:  "schema-uri, u",
			Usage: "the URI of the schema to validate against",
		},
		cli.StringFlag{
			Name:  "format, f",
			Usage: "how to format validation errors",
			Value: "string",
		},
	}

	app.CustomAppHelpTemplate = cli.AppHelpTemplate + exampleMessage

	app.Action = func(c *cli.Context) error {
		format := outputFormatString

		switch c.String("format") {
		case "string":
			format = outputFormatString
		case "json":
			format = outputFormatJSON
		default:
			return fmt.Errorf("unknown format: %s", c.String("format"))
		}

		return run(c.Args(), format)
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run(schemaPaths []string, format outputFormat) error {
	// parse each of the inputted paths into Schema structs
	schemas := make([]jsonvalidate.SchemaStruct, len(schemaPaths))
	for i, schemaPath := range schemaPaths {
		reader, err := os.Open(schemaPath)
		if err != nil {
			return err
		}

		decoder := json.NewDecoder(reader)
		err = decoder.Decode(&schemas[i])
		if err != nil {
			return err
		}
	}

	// construct a new validator from the given schemas
	registry, err := jsonvalidate.NewRegistry(schemas)
	if err != nil {
		return err
	}

	validator := jsonvalidate.Validator{Registry: registry}

	decoder := json.NewDecoder(os.Stdin)  // parses JSON from stdin
	encoder := json.NewEncoder(os.Stdout) // outputs JSON to stdout (for json output format)

	// i keeps track of which instance we're evaluating
	for i := 0; true; i++ {
		// read a JSON value out of stdin
		var instance interface{}
		err := decoder.Decode(&instance)
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		// validate the parsed JSON value
		result, err := validator.Validate(instance)
		if err != nil {
			return err
		}

		// output the errors
		for _, vErr := range result.Errors {
			switch format {
			case outputFormatString:
				fmt.Printf(
					"%d: error at: %#v (due to %#v) (schema id: %#v)\n",
					i, vErr.InstancePath.String(), vErr.SchemaPath.String(), vErr.SchemaURI.String(),
				)
			case outputFormatJSON:
				out := struct {
					Instance     int             `json:"instance"`
					InstancePath jsonpointer.Ptr `json:"instancePath"`
					SchemaPath   jsonpointer.Ptr `json:"schemaPath"`
					SchemaURI    string          `json:"schemaURI"`
				}{
					i,
					vErr.InstancePath,
					vErr.SchemaPath,
					vErr.SchemaURI.String(),
				}

				err := encoder.Encode(out)
				if err != nil {
					return err
				}
			}
		}
	}

	return err
}
