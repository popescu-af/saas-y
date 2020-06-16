package templates

// MainExample is the template for the main functionality in go code.
const MainExample = `package main

import (
	"fmt"
	"net/http"

	"github.com/popescu-af/saas-y/pkg/log"

	"{{.RepositoryURL}}/internal/config"
	"{{.RepositoryURL}}/internal/logic"
	"{{.RepositoryURL}}/internal/service"

	{{range $i, $d := .DependencyInfos -}}
	client{{$i}} "{{$d.RepositoryURL}}/pkg/client"
	{{end}}
)

func main() {
	defer log.Sync()

	log.Info("{{.Name}} started")

	env, err := config.ProcessEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

	api := logic.NewAPI(
		{{range $i, $d := .Dependencies}}
			{{- with $name := $d | cleanName | capitalize -}}
				client{{$i}}.New{{$name}}Client(env.{{$name}}Addr),
			{{- end -}}
		{{end}}
	)
	httpWrapper := service.NewHTTPWrapper(api)
	router := service.NewRouter(httpWrapper.Paths())

	log.Fatal(fmt.Sprintf("error serving - %v", http.ListenAndServe(fmt.Sprintf(":%s", env.Port), router)))
}`
