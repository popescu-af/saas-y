package generator

import (
	"fmt"
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
	CodeFormatter(path string) (st SymbolTable, err error)
	GenerateProject(name, path string) (err error)
}

// Do generates code and infrastructure declaration for the given Spec
// and dumps it in the specified output directory.
func Do(g Abstract, spec *model.Spec, outdir string) (err error) {
	st = make(SymbolTable)

	for _, svc := range spec.Services {
		err = service(g, svc, outdir)
		if err != nil {
			return
		}
	}

	fmt.Println(st)
	return
}

func service(g Abstract, svc model.Service, outdir string) (err error) {
	basePath := path.Join(outdir, "services", svc.Name)

	dirs, err := serviceDirs(g, basePath)
	if err != nil {
		return
	}

	if err = g.GenerateProject(svc.Name, basePath); err != nil {
		return
	}

	err = structs(g, svc.Structs, dirs[6])
	if err != nil {
		return
	}

	components := map[string]string{
		"api_example":    dirs[4],
		"api_definition": dirs[3],
		"main":           dirs[0],
		"env":            dirs[2],
		"http_router":    dirs[5],
		"http_wrapper":   dirs[5],
	}

	for component, outdir := range components {
		err = serviceComponent(g, svc, component, outdir)
		if err != nil {
			return
		}
	}

	// TODO:
	// - pkg/service/http_wrapper.go
	// - unit tests
	// - deploy/

	return
}

func serviceDirs(g Abstract, basePath string) (dirs []string, err error) {
	dirs = []string{
		path.Join(basePath, g.CommandPath()),
		path.Join(basePath, "deploy"),
		path.Join(basePath, g.PackagePath(), "config"),
		path.Join(basePath, g.PackagePath(), "logic"),
		path.Join(basePath, g.PackagePath(), "logic", "example"),
		path.Join(basePath, g.PackagePath(), "service"),
		path.Join(basePath, g.PackagePath(), "structs"),
	}

	for _, dir := range dirs {
		if err = os.MkdirAll(dir, 0770); err != nil {
			return
		}
	}
	return
}

func serviceComponent(g Abstract, svc model.Service, componentName, outdir string) (err error) {
	filler := templateFiller(g.GetTemplate(componentName), g.CodeFormatter)
	fPath := path.Join(outdir, componentName+g.FileExtension())
	err = filler(svc, fPath)
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

func templateFiller(templ string, codeFormatter func(string) (SymbolTable, error)) templateFillerFunction {
	loadedTempl := template.Must(template.New("templ").
		Funcs(template.FuncMap{
			"capitalize": func(s string) string { return strings.ToUpper(s[:1]) + s[1:] },
			"toLower":    strings.ToLower,
			"toUpper":    strings.ToUpper,
			"symbolize":  symbolize,
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

		newSymTable, err := codeFormatter(resultPath)
		if err == nil {
			for k, v := range newSymTable {
				st[k] = v
			}
		}
		return
	}
}

// SymbolTable is a map translating original symbol names to symbol names
// conforming to the naming style of the generated language.
type SymbolTable map[string]string

var st SymbolTable

func symbolize(originalName string) string {
	if translatedName, ok := st[originalName]; ok {
		return translatedName
	}
	return originalName
}
