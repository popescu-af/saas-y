package templates

// Main is the template for the main functionality in go code.
const Main = `package main

import (
	"fmt"
	"net/http"

	"github.com/popescu-af/saas-y/pkg/logutil"

	"{{.Name}}/internal/config"
	"{{.Name}}/internal/logic"
	"{{.Name}}/internal/service"
)

func main() {
	defer logutil.Sync()

	logutil.Info("{{.Name}} started")

	env, err := config.ProcessEnv()
	if err != nil {
		logutil.Fatal(err.Error())
	}

	api := logic.NewAPI()
	httpWrapper := service.NewHTTPWrapper(api)
	router := service.NewRouter(httpWrapper.Paths())

	logutil.Fatal(fmt.Sprintf("error serving - %v", http.ListenAndServe(fmt.Sprintf(":%s", env.Port), router)))
}`
