package monitors

import (
	core "github.com/gerty-monit/core"
	jsc "github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"log"
	"net/http"
)

type JsonSchemaMonitor struct {
	*core.HttpMonitor
	schema string
}

// ensure we always implement Monitor
var _ core.Monitor = (*JsonSchemaMonitor)(nil)

func checkSchema(schemaFile string) core.SuccessChecker {
	return func(resp *http.Response) bool {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error reading response body, %v", err)
			return false
		}

		schema := jsc.NewReferenceLoader(schemaFile)
		json := jsc.NewStringLoader(string(body))
		result, err := jsc.Validate(schema, json)

		if err != nil {
			log.Printf("error validating schema, %v", err)
			return false
		}

		if result.Valid() {
			return true
		} else {
			log.Printf("%s schema errors:", schemaFile)
			for _, err := range result.Errors() {
				log.Printf("\t %v: \t %s (on field %s)", err.Value(), err.Description(), err.Field())
			}
			return false
		}
	}
}

func NewJsonSchemaMonitorWithOptions(title, description, url, schema string,
	opts *core.HttpMonitorOptions) *JsonSchemaMonitor {
	opts.Successful = checkSchema(schema)

	delegate := core.NewHttpMonitorWithOptions(title, description, url, opts)
	return &JsonSchemaMonitor{delegate, schema}
}

func NewJsonSchemaMonitor(title, description, url, schema string) *JsonSchemaMonitor {
	opts := &core.HttpMonitorOptions{Successful: checkSchema(schema)}
	return NewJsonSchemaMonitorWithOptions(title, description, url, schema, opts)
}
