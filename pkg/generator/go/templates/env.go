package templates

// Env is a template for the environment config go file.
const Env = `package config

import "github.com/kelseyhightower/envconfig"

// Env holds all environmental variables for the service app.
type Env struct {
	Port string ` + "`" + `default:"{{.Port}}" envconfig:"PORT"` + "`" + `
	{{range .Environment}}{{.Name | toLower}} {{.Type}} ` + "`" + `default:"{{.Value}}" envconfig:"{{.Name | toUpper}}"` + "`" + `
	{{end}}
}

// ProcessEnv processes the environment, filling an
// Env struct's fields with the found values.
func ProcessEnv() (e Env, err error) {
	err = envconfig.Process("app", &e)
	return e, err
}`
