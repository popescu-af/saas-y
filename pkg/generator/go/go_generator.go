package gengo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/popescu-af/saas-y/pkg/generator"
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

// GetTemplate returns the named template for go code.
func (g *Generator) GetTemplate(name string) string {
	switch name {
	case "api_definition":
		return templates.APIDefinition
	case "api_example":
		return templates.APIExample
	case "env":
		return templates.Env
	case "http_router":
		return templates.HTTPRouter
	case "http_wrapper":
		return templates.HTTPWrapper
	case "main":
		return templates.Main
	case "struct":
		return templates.Struct
	}

	fmt.Println("Returning empty template!")
	return ""
}

// CodeFormatter returns the code formatter for go code.
func (g *Generator) CodeFormatter(path string) (st generator.SymbolTable, err error) {
	fmt.Println("Formatting", path)

	var out, errout bytes.Buffer

	cmd := exec.Command("golint", path)
	cmd.Stdout = &out
	cmd.Stderr = &errout

	err = cmd.Run()
	if err != nil {
		fmt.Println(errout.String())
		return
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	st = make(generator.SymbolTable)
	contents := string(b)
	reLint := regexp.MustCompile(`(var|method|struct field|type) ([a-zA-Z0-9_]+) should be ([a-zA-Z0-9]+)`)

	for _, matchLint := range reLint.FindAllStringSubmatch(out.String(), -1) {
		old := matchLint[2]
		new := matchLint[3]

		st[old] = new
		reDefinition := regexp.MustCompile(old + `([^"'])`)

		for _, matchDefinition := range reDefinition.FindAllStringSubmatch(contents, -1) {
			contents = strings.ReplaceAll(contents, old+matchDefinition[1], new+matchDefinition[1])
		}
	}

	out.Reset()
	errout.Reset()

	cmd = exec.Command("gofmt")
	cmd.Stdin = bytes.NewBufferString(contents)
	cmd.Stdout = &out
	cmd.Stderr = &errout

	err = cmd.Run()
	if err != nil {
		fmt.Printf("%s", errout.String())
		return
	}

	err = ioutil.WriteFile(path, out.Bytes(), 0660)
	return
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
	var out, errout bytes.Buffer

	cmd := exec.Command("go", "env")
	cmd.Stdout = &out
	cmd.Stderr = &errout
	if err = cmd.Run(); err != nil {
		fmt.Printf("%s", errout.String())
		return
	}

	cmd = exec.Command("go", "mod", "init", name)
	cmd.Env = []string{"GO111MODULE=on"}
	cmd.Dir = path
	cmd.Stdout = &out
	cmd.Stderr = &errout

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
	errout.Reset()

	if err = cmd.Run(); err != nil {
		fmt.Printf("%s", errout.String())
		err = nil
	}
	return
}
