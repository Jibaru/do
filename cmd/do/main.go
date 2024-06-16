package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jibaru/do/internal/parser"
	"github.com/jibaru/do/internal/parser/analyzer"
	"github.com/jibaru/do/internal/parser/extractor"
	"github.com/jibaru/do/internal/parser/normalizer"
	"github.com/jibaru/do/internal/parser/partitioner"
	"github.com/jibaru/do/internal/parser/replacer"
	"github.com/jibaru/do/internal/parser/taker"
	"github.com/jibaru/do/internal/reader"
	"github.com/jibaru/do/internal/request"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("use: do <filename.do>")
		return
	}

	filename := os.Args[1]

	doFileReader := reader.NewFileReader()
	sectionTaker := taker.New()
	sectionNormalizer := normalizer.New()
	sectionPartitioner := partitioner.New()
	expressionAnalyzer := analyzer.New()
	sectionExtractor := extractor.New(sectionTaker, sectionNormalizer, sectionPartitioner, expressionAnalyzer)
	variablesReplacer := replacer.New()
	theParser := parser.New(doFileReader, sectionExtractor, variablesReplacer)
	client := request.NewHttpClient(&http.Client{})

	doFile, err := theParser.FromFilename(filename)
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
