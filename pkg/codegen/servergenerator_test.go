package codegen

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"text/template"
	"unicode"
)

func TestBaseGenerator_Generate(t *testing.T) {

	operations := []OperationDefinition{
		{
			OperationId: "GetCatStatus",
			Responses: []ResponseDefinition{{
				StatusCode:  "200",
				Description: "Success",
				Ref:         "",
			}},
			Summary: "Get cat status",
			Method:  "GET",
			Path:    "/cat",
			Spec:    &openapi3.Operation{Tags: []string{"cat"}},
		},
	}

	type args struct {
		conf Configuration
	}
	tests := []struct {
		name      string
		generator ServerGenerator
		args      args
		want      string
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name:      "StdHttpGenerator",
			generator: NewStdHttp(),
			args: args{
				conf: Configuration{
					Generate: GenerateOptions{StdHTTPServer: true},
				},
			}, want: expectedStdHttp, wantErr: assert.NoError,
		},
		{
			name:      "StdHttpGenerator Grouped Tags",
			generator: NewStdHttp(),
			args: args{
				conf: Configuration{
					Generate:      GenerateOptions{StdHTTPServer: true},
					OutputOptions: OutputOptions{GroupByTag: true},
				},
			}, want: expectedStdHttpGrouped, wantErr: assert.NoError,
		},
		{
			name:      "GinGenerator Grouped Tags",
			generator: NewGinGenerator(),
			args: args{
				conf: Configuration{
					Generate:      GenerateOptions{StdHTTPServer: true},
					OutputOptions: OutputOptions{GroupByTag: true},
				},
			}, want: expectedGinGrouped, wantErr: assert.NoError,
		},
		{
			name:      "GinGenerator Grouped Tags check interface",
			generator: NewGinGenerator(),
			args: args{
				conf: Configuration{
					Generate:      GenerateOptions{StdHTTPServer: true},
					OutputOptions: OutputOptions{GroupByTag: true},
				},
			}, want: expectedGinGroupedInterface, wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TemplateFunctions["opts"] = func() Configuration { return tt.args.conf }
			tem := template.New("oapi-codegen").Funcs(TemplateFunctions)
			err := LoadTemplates(templates, tem)
			if err != nil {
				t.Errorf("error parsing oapi-codegen templates: %e", err)
			}
			got, err := tt.generator.Generate(tem, operations)
			if !tt.wantErr(t, err, fmt.Sprintf("Generate(%v, %v)", tt.generator, tt.args.conf)) {
				return
			}
			assert.Contains(t, removeWhitespace(got), removeWhitespace(tt.want), "Generate(%v, %v)", tt.generator, tt.args.conf)
		})
	}
}

func removeWhitespace(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, ch := range s {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

const expectedStdHttp = `
type ServerInterface interface {
	// Get cat status
	// (GET /cat)
	GetCatStatus(w http.ResponseWriter, r *http.Request)
}
`
const expectedStdHttpGrouped = `
type ServerInterface interface {
	CatAPI
}
`

const expectedGinGrouped = `
type ServerInterface interface {
	CatAPI
}
`
const expectedGinGroupedInterface = `
// CatAPI handlers for tag Cat
type CatAPI interface {
	// Get cat status
	// (GET /cat)
	GetCatStatus(c *gin.Context)
}`
