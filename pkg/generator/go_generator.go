package generator

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/popescu-af/saas-y/pkg/model"
)

// Go generates go code + infrastructure specification from a saas-y Spec.
type Go struct {
}

// type Variable struct {
// 	Name  string `json:"name"`
// 	Type  string `json:"type"`
// 	Value string `json:"value"`
// }

// type Struct struct {
// 	Name   string     `json:"name"`
// 	Fields []Variable `json:"fields"`
// }

// Generate does the actual code + infrastructure specification generation.
func (g *Go) structs(structs []model.Struct, outdir string) (err error) {
	dir := path.Join(outdir, "pkg", "structs")
	if err = os.MkdirAll(dir, 0770); err != nil {
		return
	}

	templ := template.Must(template.New("struct").Parse(templateStruct))
	for _, s := range structs {
		var f *os.File
		fPath := path.Join(dir, s.Name+".go")
		f, err = os.Create(fPath)
		if err != nil {
			return
		}

		templ.Execute(f, s)
		f.Close()
		g.fmtAndLint(fPath)
	}

	return
}

func (g *Go) fmtAndLint(path string) (err error) {
	re := regexp.MustCompile(`(type|struct field) ([a-zA-Z0-9_]+) should be ([a-zA-Z0-9]+)`)
	lint := exec.Command("golint", path)

	var out bytes.Buffer
	lint.Stdout = &out

	err = lint.Run()
	if err != nil {
		return
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	contents := strings.ReplaceAll(string(b), "\n\n", "\n")

	for _, subst := range re.FindAllStringSubmatch(out.String(), -1) {
		old := subst[2] + " "
		new := subst[3] + " "
		new = strings.ToUpper(new[:1]) + new[1:]
		contents = strings.ReplaceAll(contents, old, new)
	}
	out.Reset()

	format := exec.Command("gofmt")
	format.Stdin = bytes.NewBufferString(contents)
	format.Stdout = &out

	err = format.Run()
	if err != nil {
		return
	}

	return ioutil.WriteFile(path, out.Bytes(), 0660)
}

const templateStruct = `package structs

// {{.Name}} - generated API structure
type {{.Name}} struct {
{{range .Fields}}
	{{.Name}} {{.Type}} ` + "`" + `json:"{{.Name}}"` + "`" +
	`{{end}}
}`
