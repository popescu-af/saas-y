package templates

// Env is the template for the environment config in go code.
const Env = `package config

import "github.com/kelseyhightower/envconfig"

// Env holds all environmental variables for the service app.
type Env struct {
	Port string ` + "`" + `default:"{{.Port}}" envconfig:"PORT"` + "`" + `
	{{range .Environment -}}
	{{.Name | toLower | capitalize}} {{.Type}} ` + "`" + `default:"{{.Value}}" envconfig:"{{.Name | toUpper}}"` + "`" + `
	{{- end}}
	{{range $d := .Dependencies -}}
	{{$d | replaceHyphens | toLower | capitalize}}Addr string ` + "`" + `default:"" envconfig:"{{$d | replaceHyphens | toUpper}}_ADDR"` + "`" + `
	{{- end}}
}

// ProcessEnv processes the environment, filling an
// Env struct's fields with the found values.
func ProcessEnv() (e Env, err error) {
	err = envconfig.Process("app", &e)
	return e, err
}`
