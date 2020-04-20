package gengo

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/popescu-af/saas-y/pkg/generator/go/templates"
)

// Generator generates go code + infrastructure specification from a saas-y Spec.
type Generator struct {
}

// FileExtension returns the file extension for go code.
func (g *Generator) FileExtension() string {
	return ".go"
}

// CommandPath returns the command path for go code.
func (g *Generator) CommandPath() string {
	return path.Join("cmd")
}

// PackagePath returns the package path for go code.
func (g *Generator) PackagePath() string {
	return path.Join("pkg")
}

// StructsTemplate returns the structs template for go code.
func (g *Generator) StructsTemplate() string {
	return templates.Struct
}

// MainTemplate returns the main template for go code.
func (g *Generator) MainTemplate() string {
	return templates.Main
}

// CodeFormatter returns the code formatter for go code.
func (g *Generator) CodeFormatter(path string) (err error) {
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

// GenerateProject creates project-specific files.
func (g *Generator) GenerateProject(name, path string) (err error) {
	err = g.createGoModFile(name, path)

	// TODO:
	// Makefile
	// Dockerfile

	return
}

func (g *Generator) createGoModFile(name, path string) (err error) {
	var out bytes.Buffer

	cmd := exec.Command("go", "env")
	cmd.Stdout = &out
	if err = cmd.Run(); err != nil {
		return
	}

	cmd = exec.Command("go", "mod", "init", name)
	cmd.Env = []string{"GO111MODULE=on"}
	cmd.Dir = path
	cmd.Stderr = &out

	neededVars := []string{"GOCACHE", "GOPATH"}
	lines := strings.Split(out.String(), "\n")
	for _, l := range lines {
		tokens := strings.Split(l, "=")
		for _, v := range neededVars {
			if tokens[0] == v {
				var unquoted string
				if unquoted, err = strconv.Unquote(tokens[1]); err != nil {
					return
				}
				cmd.Env = append(cmd.Env, v+"="+unquoted)
				break
			}
		}
	}
	out.Reset()

	err = cmd.Run()
	return
}
