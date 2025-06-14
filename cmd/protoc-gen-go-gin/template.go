package main

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

//go:embed httpTemplate.tpl
var httpTemplate string

type serviceDesc struct {
	ServiceType   string // Greeter
	SnakeCaseType string // greeter
	ServiceName   string // helloworld.Greeter
	Metadata      string // api/helloworld/helloworld.proto
	Methods       []*methodDesc
	MethodSets    map[string]*methodDesc

	PackageName string
	ImportPath  string
}

type methodDesc struct {
	// method
	Name               string
	SnakeCaseName      string // say_hello
	OriginalName       string // The parsed original name
	Num                int
	Request            string
	RequestImportPath  string
	RequestPackageName string
	Reply              string
	ReplyImportPath    string
	ReplyPackageName   string
	Comment            string
	// http_rule
	Path         string
	Method       string
	HasVars      bool
	HasBody      bool
	Body         string
	ResponseBody string
}

func (s *serviceDesc) execute() string {
	s.MethodSets = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSets[m.Name] = m
	}
	buf := new(bytes.Buffer)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(httpTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return strings.Trim(buf.String(), "\r\n")
}
