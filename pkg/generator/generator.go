package generator

import (
	"fmt"
	"html/template"
	"log"
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

	fmt.Println("Done.")
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

	components := []struct {
		template string
		outdir   string
	}{
		{"api_example", dirs[4]},
		{"api_definition", dirs[3]},
		{"main", dirs[0]},
		{"env", dirs[2]},
		{"http_router", dirs[5]},
		{"http_wrapper", dirs[5]},
	}

	for _, c := range components {
		err = serviceComponent(g, svc, c.template, c.outdir)
		if err != nil {
			return
		}
	}

	// TODO:
	// - header params handling
	// - query params handling
	// - method signature creation from body, parameters
	// - different handling for options method
	// - params passing to inner API method
	// - unit tests
	// - validate method type, combination of params
	// - deploy/
	// New features:
	// - code/example for talking to well-known services/tools (redis, etc.)
	// - linkage between saas-y generated services
	// - k8s yaml files for all good to have stuff:
	//   - ingress
	//   - cert-manager
	//   - docker-register
	// - unit tests for the generated service (everything excluding the pure logic)

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
			"pathHasParameters": func(s string) string {
				ss := strings.Split(s, "/")
				if strings.Contains(ss[len(ss)-1], "}") {
					return "yes"
				}
				return "" // empty value means false in {{if $x}} template conditional
			},
			"pathParameters": func(s string) (result []string) {
				paramMap := make(map[string]string)

				ss := strings.Split(s, "/")
				for i := len(ss) - 1; i >= 0; i-- {
					if !strings.Contains(ss[i], "}") {
						break
					}

					tokens := strings.Split(ss[i], ":")
					if len(tokens) != 2 {
						log.Fatalf("invalid path parameter spec: %s (should be in the form name:type)", ss[i])
					}

					pName := tokens[0][1:]
					if _, ok := paramMap[pName]; ok {
						log.Fatalf("path paramerter already defined: %s", pName)
					}

					pType := tokens[1][:len(tokens[1])-1]
					for _, t := range []string{"int", "uint", "float", "string"} {
						if t == pType {
							result = append(result, pName)
							result = append(result, pType)
							paramMap[pName] = pType
						}
					}
					if _, ok := paramMap[pName]; !ok {
						log.Fatalf("invalid type for parameter '%s': '%s'", pName, pType)
					}
				}
				return result
			},
			"indicesParameters": func(parameters []string) []int {
				var indices []int
				for i := 0; i < len(parameters); i += 2 {
					indices = append(indices, i)
				}
				return indices
			},
			"inc": func(i int) int {
				return i + 1
			},
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
