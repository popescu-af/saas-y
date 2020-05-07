package config

import "github.com/kelseyhightower/envconfig"

// Env holds all environmental variables for the service app.
type Env struct {
	Port       string `default:"80" envconfig:"PORT"`
	EnvVarName int64  `default:"42" envconfig:"ENV_VAR_NAME"`
}

// ProcessEnv processes the environment, filling an
// Env struct's fields with the found values.
func ProcessEnv() (e Env, err error) {
	err = envconfig.Process("app", &e)
	return e, err
}
