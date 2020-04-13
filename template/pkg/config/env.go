package config

import "github.com/kelseyhightower/envconfig"

// The Env struct holds all needed environmental variables
// that configure the service app.
type Env struct {
	Port string `default:"8080" envconfig:"port"`
}

// ProcessEnv processes the environment, filling an
// Env struct's fields with the values found.
func ProcessEnv() (e Env, err error) {
	err = envconfig.Process("svc", &e)
	return e, err
}
