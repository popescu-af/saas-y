package main

import (
	"flag"
	"log"
	"os"

	"github.com/popescu-af/saas-y/pkg/engine/golang"
)

func main() {
	inputFilePath, outputDirPath := parseArgs()

	stat, err := os.Stat(inputFilePath)
	if os.IsNotExist(err) {
		log.Fatalln("file does not exist - " + inputFilePath)
	}

	if stat.IsDir() {
		log.Fatalln("input is a directory - " + inputFilePath)
	}

	stat, err = os.Stat(outputDirPath)
	if err == nil && !stat.IsDir() {
		log.Fatalln("output is not a directory - " + outputDirPath)
	}

	err = golang.GenerateSourcesFromJSONSpec(inputFilePath, outputDirPath)
	if err != nil {
		log.Fatalf("saas-y error: %v", err)
	}
}

func parseArgs() (inputFilePath string, outputDirPath string) {
	flag.StringVar(&inputFilePath, "input", "./spec.json", "path to the saas-y specification JSON")
	flag.StringVar(&outputDirPath, "output", "./output-saas-y", "path to the output directory")
	flag.Parse()
	return
}
