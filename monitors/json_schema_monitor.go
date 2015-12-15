package monitors

import (
	core "github.com/gerty-monit/core/monitors"
	jsc "github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"log"
	"net/http"
)

type JsonSchemaMonitor struct {
	delegate *core.HttpMonitor
	schema   string
}

func checkSchema(schemaFile string) core.SuccessChecker {
	return func(resp *http.Response) bool {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("error reading response body, %v", err)
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
			log.Printf("schema errors:")
			for _, err := range result.Errors() {
				log.Printf("\t %s: \t %s", err.Field(), err.Description())
			}
			return false
		}
	}
}

func (monitor *JsonSchemaMonitor) Check() int {
	return monitor.delegate.Check()
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
