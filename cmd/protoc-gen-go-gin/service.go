package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"html/template"
	"os"
	"path/filepath"
)

const tplMainService = `package {{.SnakeCaseType}}

import (
	"github/invokerw/gintos/log"

	"{{.ImportPath}}"
)

type {{.ServiceType}}Service struct {
	log *log.Helper
}

var _ {{.ImportName}}.I{{.ServiceType}}Server = (*{{.ServiceType}}Service)(nil)

func New{{.ServiceType}}Service(logger log.Logger) *{{.ServiceType}}Service {
	return &{{.ServiceType}}Service{
		log: log.NewHelper(logger),
	}
}

`

const tplMethodService = `package {{.PackageName}}

import (
	"github.com/gin-gonic/gin"
	"{{.RequestImportPath}}"
	{{if ne .RequestImportPath .ReplyImportPath}}
	"{{.ReplyImportPath}}"
	{{- end }}
)

func (s *{{.ServiceType}}Service) {{.Name}}(*gin.Context, *{{.RequestPackageName}}.{{.Request}}) (*{{.ReplyPackageName}}.{{.Reply}}, error) {
	// Implement your logic here
	
	return nil, nil
}

`

func serviceGenerate(gen *protogen.Plugin, savePath string, allServiceDesc []*serviceDesc) error {

	type SInfo struct {
		serviceDesc
		ImportName string // Import name for the service
	}

	type MInfo struct {
		methodDesc
		PackageName string
		ServiceType string
	}

	do := func(sd *serviceDesc) error {
		path := filepath.Join(savePath, sd.SnakeCaseType)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
		// Create the service go file if it does not exist
		mainFilePath := filepath.Join(path, sd.SnakeCaseType+".go")
		if _, err := os.Stat(mainFilePath); err != nil {
			mainFile, err := os.Create(mainFilePath)
			if err != nil {
				return fmt.Errorf("failed to create service file %s: %w", sd.SnakeCaseType, err)
			}
			defer mainFile.Close()
			tmpl, err := template.New("service_main").Parse(tplMainService)
			if err != nil {
				return fmt.Errorf("failed to parse template for service %s: %w", sd.ServiceName, err)
			}
			if err := tmpl.Execute(mainFile, SInfo{
				serviceDesc: *sd,
				ImportName:  getPackageName(sd.ImportPath),
			}); err != nil {
				return fmt.Errorf("failed to execute template for service %s: %w", sd.ServiceName, err)
			}
		}

		for _, m := range sd.Methods {
			methodFilePath := filepath.Join(path, m.SnakeCaseName+".go")
			if _, err := os.Stat(methodFilePath); err != nil {
				methodFile, err := os.Create(methodFilePath)
				if err != nil {
					return fmt.Errorf("failed to create method file %s: %w", m.Name, err)
				}
				defer methodFile.Close()
				tmpl, err := template.New("method").Parse(tplMethodService)
				if err != nil {
					return fmt.Errorf("failed to parse template for method %s: %w", m.Name, err)
				}
				info := &MInfo{
					methodDesc:  *m,
					PackageName: sd.SnakeCaseType,
					ServiceType: sd.ServiceType,
				}
				info.Request = stringsLastIndex(info.Request, ".")
				info.Reply = stringsLastIndex(info.Reply, ".")

				if err := tmpl.Execute(methodFile, info); err != nil {
					return fmt.Errorf("failed to execute template for method %s: %w", m.Name, err)
				}
			}
		}

		return nil
	}
	for _, sd := range allServiceDesc {
		if err := do(sd); err != nil {
			return err
		}
	}
	return nil
}
