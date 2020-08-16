package templates

// Main is the template for the main functionality in go code.
const Main = `package main

import (
	"fmt"
	"net/http"

	"github.com/popescu-af/saas-y/pkg/log"

	"{{.RepositoryURL}}/internal/config"
	"{{.RepositoryURL}}/internal/logic"
	"{{.RepositoryURL}}/internal/service"

	{{range $d := .DependencyInfos -}}
	{{$d.Name | cleanName | toLower}} "{{$d.RepositoryURL}}/pkg/client"
	{{end}}
)

func main() {
	defer log.Sync()

	log.Info("{{.Name}} started")

	env, err := config.ProcessEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

	impl := logic.NewImpl(
		env,
		{{range $d := .DependencyInfos}}
			{{- with $name := $d.Name | cleanName | capitalize -}}
				{{$d.Name | cleanName | toLower}}.New{{$name}}Client(env.{{$name}}Addr),
			{{- end}}
		{{end}}
	)
	httpWrapper := service.NewHTTPWrapper(impl)
	router := service.NewRouter(httpWrapper.Paths())

	log.Fatal(fmt.Sprintf("error serving - %v", http.ListenAndServe(fmt.Sprintf(":%s", env.Port), router)))
}`
