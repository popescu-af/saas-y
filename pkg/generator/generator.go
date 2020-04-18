package generator

import (
	"html/template"
	"os"
	"path"

	"github.com/popescu-af/saas-y/pkg/model"
)

// Abstract is the interface for saas-y code & infrastructure generators.
type Abstract interface {
	FileExtension() string

	CommandPath() string
	PackagePath() string
	StructsTemplate() string
	MainTemplate() string

	CodeFormatter(path string) (err error)
}

// Do generates code and infrastructure declaration for the given Spec
// and dumps it in the specified output directory.
func Do(g Abstract, spec *model.Spec, outdir string) (err error) {
	err = structs(g, spec.Structs, outdir)
	if err != nil {
		return
	}

	for _, svc := range spec.Services {
		err = service(g, svc, outdir)
		if err != nil {
			return
		}
	}

	return
}

func structs(g Abstract, structs []model.Struct, outdir string) (err error) {
	dir := path.Join(outdir, g.PackagePath(), "structs")
	if err = os.MkdirAll(dir, 0770); err != nil {
		return
	}

	filler := templateFiller(g.StructsTemplate, g.CodeFormatter)
	for _, s := range structs {
		fPath := path.Join(dir, s.Name+g.FileExtension())
		err = filler(s, fPath)
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
	}

	for _, dir := range dirs {
		if err = os.MkdirAll(dir, 0770); err != nil {
			return
		}
	}

	filler := templateFiller(g.MainTemplate, g.CodeFormatter)
	fPath := path.Join(dirs[0], "main"+g.FileExtension())
	err = filler(svc, fPath)
	if err != nil {
		return
	}

	// TODO: continue

	return
}

type templateFillerFunction func(interface{}, string) error

func templateFiller(templateGetter func() string, codeFormatter func(string) error) templateFillerFunction {
	templ := template.Must(template.New("templ").Parse(templateGetter()))

	return func(s interface{}, resultPath string) (err error) {
		var f *os.File
		f, err = os.Create(resultPath)
		if err != nil {
			return
		}
		templ.Execute(f, s)
		f.Close()
		err = codeFormatter(resultPath)
		return
	}
}
