package main

import (
	"encoding/json"
	"fmt"
	"github.com/jibaru/do/internal/reader"
	"log"
	"net/http"
	"os"

	"github.com/jibaru/do/internal/parser"
	"github.com/jibaru/do/internal/request"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("use: do <filename.do>")
		return
	}

	filename := os.Args[1]

	doFileReader := reader.NewFileReader()
	doSectionNormalizer := parser.NewNormalizer()
	sectionExtractor := parser.NewSectionExtractor(doSectionNormalizer)
	variablesReplacer := parser.NewVariablesReplacer()
	psr := parser.New(doFileReader, sectionExtractor, variablesReplacer)
	client := request.NewHttpClient(&http.Client{})

	doFile, err := psr.FromFilename(filename)
	if err != nil {
		log.Printf("error parsing: %v\n", err)
		return
	}

	doFileAsJson, _ := json.MarshalIndent(doFile, "", "   ")
	fmt.Printf("request: %s\n", string(doFileAsJson))

	response, err := client.Do(*doFile)
	if err != nil {
		fmt.Printf("error in request: %v\n", err)
		return
	}

	responseAsJson, _ := json.MarshalIndent(response, "", "   ")
	fmt.Printf("response: %s\n", string(responseAsJson))
}
