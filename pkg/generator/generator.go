package generator

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	common_templ "github.com/popescu-af/saas-y/pkg/generator/common/templates"
	"github.com/popescu-af/saas-y/pkg/generator/common/templates/k8s"
	"github.com/popescu-af/saas-y/pkg/model"
)

// Abstract is the interface for saas-y code & infrastructure generators.
type Abstract interface {
	FileExtension() string
	CommandPath() string
	InternalPath() string
	PackagePath() string
	GetTemplate(name string) string
	CodeFormatter(path string) (st SymbolTable, err error)
	GenerateProject(name, path string) (err error)
}

// Init initializes the generator.
func Init() {
	st = make(SymbolTable)
}

// Do generates code and infrastructure declaration for the given Spec
// and dumps it in the specified output directory.
func Do(g Abstract, spec *model.Spec, outdir string) (err error) {
	deployDir := path.Join(outdir, "deploy")
	if err = os.MkdirAll(deployDir, 0770); err != nil {
		return
	}

	entities := []struct {
		template string
		outpath  string
	}{
		{k8s.Ingress, path.Join(deployDir, "ingress.yaml")},
		{k8s.ClusterIssuer, path.Join(deployDir, "letsencrypt-issuer.yaml")},
		{k8s.Certificate, path.Join(deployDir, "letsencrypt-certificate.yaml")},
		{k8s.Registry, path.Join(deployDir, "docker-registry.yaml")},
	}

	for _, e := range entities {
		err = CommonEntity(spec, e.template, e.outpath)
		if err != nil {
			return
		}
	}

	for _, svc := range spec.Services {
		err = Service(g, svc, outdir)
		if err != nil {
			return
		}
	}

	fmt.Println("Done.")
	return
}

// Service generates all files for a service entity.
func Service(g Abstract, svc model.Service, outdir string) (err error) {
	basePath := path.Join(outdir, "services", svc.Name)

	dirs, err := serviceDirs(g, basePath)
	if err != nil {
		return
	}

	if err = g.GenerateProject(svc.Name, basePath); err != nil {
		return
	}

	entities := []struct {
		template string
		outpath  string
	}{
		{g.GetTemplate("dockerfile"), path.Join(dirs[0], "Dockerfile.example")},
		{common_templ.Makefile, path.Join(dirs[0], "Makefile.example")},
		{common_templ.Readme, path.Join(dirs[0], "README-example.md")},
		{k8s.DeplSvc, path.Join(dirs[2], svc.Name+".yaml")},
	}

	for _, e := range entities {
		err = CommonEntity(svc, e.template, e.outpath)
		if err != nil {
			return
		}
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
		{"api_definition", dirs[6]},
		{"errors_example", dirs[4]},
		{"main", dirs[1]},
		{"env", dirs[3]},
		{"http_router", dirs[5]},
		{"http_wrapper", dirs[5]},
		{"http_error_handler_example", dirs[5]},
	}

	for _, c := range components {
		err = serviceComponent(g, svc, c.template, c.outdir)
		if err != nil {
			return
		}
	}

	// TODO:
	// - line number when validation error occurs (save originating line number when parsing JSON)
	//
	// - README.md on how to use saas-y
	//   - test the usage of readme
	// - other files for github (contributors, license, etc.)
	// - open source and post on several channels (+donate?)
	//
	// New features:
	// - move code to internal, put client code into pkg
	// - support [] of data (POD or struct) as structure member attribute
	// - support null return from API
	// - generate client code snippets
	//   - add env variable for connectivity to the dependencies
	//   - generate wrapper over HTTP client code to be easily accessible by logic package
	// - linkage between saas-y generated services
	// - code/example for talking to well-known services/tools (redis, etc.)
	// - CORS
	// - authentication
	// - unit tests for the generated service (everything excluding the pure logic)
	//
	// Ideas:
	// - generate from docker-compose file (yaml file)?

	return
}

func serviceDirs(g Abstract, basePath string) (dirs []string, err error) {
	dirs = []string{
		basePath,
		path.Join(basePath, g.CommandPath()),
		path.Join(basePath, "deploy"),
		path.Join(basePath, g.InternalPath(), "config"),
		path.Join(basePath, g.InternalPath(), "logic"),
		path.Join(basePath, g.InternalPath(), "service"),
		path.Join(basePath, g.PackagePath(), "exports"),
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

// CommonEntity generates an entity that is common to all languages.
func CommonEntity(obj interface{}, templ string, resultPath string) (err error) {
	loadedTempl := template.Must(template.New("templ").
		Funcs(template.FuncMap{
			"toUpper": strings.ToUpper,
			"yamlify": func(s string) string {
				for _, r := range []string{".", "_"} {
					s = strings.ReplaceAll(s, r, "-")
				}
				return s
			},
		}).
		Parse(templ))
	err = applyObjectToTemplateAndSaveToFile(obj, loadedTempl, resultPath)
	return
}

type templateFillerFunction func(interface{}, string) error

func templateFiller(templ string, codeFormatter func(string) (SymbolTable, error)) templateFillerFunction {
	paramStack := ""

	loadedTempl := template.Must(template.New("templ").
		Funcs(template.FuncMap{
			"decapitalize": func(s string) string { return strings.ToLower(s[:1]) + s[1:] },
			"capitalize":   func(s string) string { return strings.ToUpper(s[:1]) + s[1:] },
			"toLower":      strings.ToLower,
			"toUpper":      strings.ToUpper,
			"symbolize":    symbolize,
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
				for _, p := range ss {
					if !strings.Contains(p, "}") {
						continue
					}

					tokens := strings.Split(p, ":")
					if len(tokens) != 2 {
						log.Fatalf("invalid path parameter spec: %s (should be in the form name:type)", p)
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
			"pushParam": func(p string) string {
				paramStack += p + ", "
				return p
			},
			"printParamStack": func() string {
				var temp string
				paramStack, temp = temp, paramStack
				l := len(temp)
				if l > 0 {
					return temp[:l-2]
				}
				return ""
			},
			"typeName": func(t string) string {
				switch t {
				case "int":
					return "int64"
				case "uint":
					return "uint64"
				case "float":
					return "float64"
				case "string":
					return "string"
				}

				return ""
			},
			"cleanPath": func(p string) string {
				re := regexp.MustCompile(`:(int|uint|float|string)`)
				return re.ReplaceAllString(p, "")
			},
		}).
		Parse(templ))

	return func(s interface{}, resultPath string) (err error) {
		err = applyObjectToTemplateAndSaveToFile(s, loadedTempl, resultPath)
		if err != nil {
			return
		}

		newSymTable, err := codeFormatter(resultPath)
		if err == nil {
			for k, v := range newSymTable {
				st[k] = v
			}
		}
		return
	}
}

func applyObjectToTemplateAndSaveToFile(obj interface{}, loadedTempl *template.Template, resultPath string) (err error) {
	var f *os.File
	f, err = os.Create(resultPath)
	if err != nil {
		return
	}
	loadedTempl.Execute(f, obj)
	f.Close()
	return
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
