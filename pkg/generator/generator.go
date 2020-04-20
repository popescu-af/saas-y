package generator

import (
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/popescu-af/saas-y/pkg/model"
)

// Abstract is the interface for saas-y code & infrastructure generators.
type Abstract interface {
	FileExtension() string
	CommandPath() string
	PackagePath() string
	GetTemplate(name string) string
	CodeFormatter(path string) (err error)
	GenerateProject(name, path string) (err error)
}

// Do generates code and infrastructure declaration for the given Spec
// and dumps it in the specified output directory.
func Do(g Abstract, spec *model.Spec, outdir string) (err error) {
	for _, svc := range spec.Services {
		err = service(g, svc, outdir)
		if err != nil {
			return
		}
	}

	return
}

func service(g Abstract, svc model.Service, outdir string) (err error) {
	base := path.Join(outdir, "services", svc.Name)

	dirs := []string{
		path.Join(base, g.CommandPath()),
		path.Join(base, "deploy"),
		path.Join(base, g.PackagePath(), "config"),
		path.Join(base, g.PackagePath(), "logic", "example"),
		path.Join(base, g.PackagePath(), "service"),
		path.Join(base, g.PackagePath(), "structs"),
	}

	if err = g.GenerateProject(svc.Name, base); err != nil {
		return
	}

	for _, dir := range dirs {
		if err = os.MkdirAll(dir, 0770); err != nil {
			return
		}
	}

	filler := templateFiller(g.GetTemplate("main"), g.CodeFormatter)
	fPath := path.Join(dirs[0], "main"+g.FileExtension())
	err = filler(svc, fPath)
	if err != nil {
		return
	}

	err = envVars(g, svc, dirs[2])
	if err != nil {
		return
	}

	err = structs(g, svc.Structs, dirs[5])
	if err != nil {
		return
	}

	// TODO:
	// - pkg/config/env.go
	// - pkg/logic/api_implementation.go
	// - pkg/service/api_definition.go
	// - pkg/service/http_router.go
	// - pkg/service/http_wrapper.go
	// - deploy/

	return
}

func envVars(g Abstract, svc model.Service, outdir string) (err error) {
	filler := templateFiller(g.GetTemplate("config"), g.CodeFormatter)
	fPath := path.Join(outdir, "config"+g.FileExtension())
	err = filler(svc, fPath)
	if err != nil {
		return
	}

	return
}

func structs(g Abstract, structs []model.Struct, outdir string) (err error) {
	filler := templateFiller(g.GetTemplate("struct"), g.CodeFormatter)
	for _, s := range structs {
		fPath := path.Join(outdir, s.Name+g.FileExtension())
		err = filler(s, fPath)
		if err != nil {
			return
		}
	}

	return
}

// TODO:
// - ingress
// - external services

type templateFillerFunction func(interface{}, string) error

func templateFiller(templ string, codeFormatter func(string) error) templateFillerFunction {
	loadedTempl := template.Must(template.New("templ").
		Funcs(template.FuncMap{
			"tolower": strings.ToLower,
			"toupper": strings.ToUpper,
		}).
		Parse(templ))

	return func(s interface{}, resultPath string) (err error) {
		var f *os.File
		f, err = os.Create(resultPath)
		if err != nil {
			return
		}
		loadedTempl.Execute(f, s)
		f.Close()
		err = codeFormatter(resultPath)
		return
	}
}
