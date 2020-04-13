package main

import (
	"fmt"
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/popescu-af/saas-y/template/pkg/config"
	"github.com/popescu-af/saas-y/template/pkg/logic/example"
	"github.com/popescu-af/saas-y/template/pkg/service"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("server started")

	env, err := config.ProcessEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

	api := example.NewAPI(logger)
	httpWrapper := service.NewHTTPWrapper(api)
	router := service.NewRouter(httpWrapper.Paths(), logger)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", env.Port), router))
}
