package codegen

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type ServerGenerator interface {
	Generate(t *template.Template, ops interface{}) (string, error)
}

type BaseGenerator struct {
	templates []string
	conf      Configuration
}

func (g *BaseGenerator) Generate(t *template.Template, ops interface{}) (string, error) {
	var generatedTemplates []string
	for _, tmpl := range g.templates {
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)

		if err := t.ExecuteTemplate(w, tmpl, ops); err != nil {
			return "", fmt.Errorf("error generating %s: %s", tmpl, err)
		}
		if err := w.Flush(); err != nil {
			return "", fmt.Errorf("error flushing output buffer for %s: %s", tmpl, err)
		}
		generatedTemplates = append(generatedTemplates, buf.String())
	}

	return strings.Join(generatedTemplates, "\n"), nil
}

// NoOpGenerator used when no generator is selected
type NoOpGenerator struct{}

func (g NoOpGenerator) Generate(*template.Template, interface{}) (string, error) {
	return "", nil
}

// StdHttp used to define a default set of templates and its order
type StdHttp struct {
	BaseGenerator
}

func NewStdHttp() *StdHttp {
	return &StdHttp{BaseGenerator{templates: []string{"stdhttp/std-http-interface.tmpl", "stdhttp/std-http-middleware.tmpl", "stdhttp/std-http-handler.tmpl"}}}
}

type GinGenerator struct {
	BaseGenerator
}

func NewGinGenerator() *GinGenerator {
	return &GinGenerator{BaseGenerator{templates: []string{"gin/gin-interface.tmpl", "gin/gin-wrappers.tmpl", "gin/gin-register.tmpl"}}}
}

type IrisGenerator struct {
	BaseGenerator
}

func NewIrisGenerator() *IrisGenerator {
	return &IrisGenerator{BaseGenerator{templates: []string{"iris/iris-interface.tmpl", "iris/iris-middleware.tmpl", "iris/iris-handler.tmpl"}}}
}

type EchoGenerator struct {
	BaseGenerator
}

func NewEchoGenerator() *EchoGenerator {
	return &EchoGenerator{BaseGenerator{templates: []string{"echo/echo-interface.tmpl", "echo/echo-wrappers.tmpl", "echo/echo-register.tmpl"}}}
}

type ChiGenerator struct {
	BaseGenerator
}

func NewChiGenerator() *ChiGenerator {
	return &ChiGenerator{BaseGenerator{templates: []string{"chi/chi-interface.tmpl", "chi/chi-middleware.tmpl", "chi/chi-handler.tmpl"}}}
}

type FiberGenerator struct {
	BaseGenerator
}

func NewFiberGenerator() *FiberGenerator {
	return &FiberGenerator{BaseGenerator{templates: []string{"fiber/fiber-interface.tmpl", "fiber/fiber-middleware.tmpl", "fiber/fiber-handler.tmpl"}}}
}

type GorillaGenerator struct {
	BaseGenerator
}

func NewGorillaGenerator() *GorillaGenerator {
	return &GorillaGenerator{BaseGenerator{templates: []string{"gorilla/gorilla-interface.tmpl", "gorilla/gorilla-middleware.tmpl", "gorilla/gorilla-register.tmpl"}}}
}
