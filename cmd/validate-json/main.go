package main

import (
	"fmt"
	"log"
	"os"

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
	}

	app.CustomAppHelpTemplate = cli.AppHelpTemplate + exampleMessage

	app.Action = func(c *cli.Context) error {
		fmt.Println("hello, world")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
