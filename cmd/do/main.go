package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jibaru/do/internal/parser"
	"github.com/jibaru/do/internal/parser/analyzer"
	"github.com/jibaru/do/internal/parser/caller"
	"github.com/jibaru/do/internal/parser/cleaner"
	"github.com/jibaru/do/internal/parser/extractor"
	"github.com/jibaru/do/internal/parser/normalizer"
	"github.com/jibaru/do/internal/parser/partitioner"
	"github.com/jibaru/do/internal/parser/replacer"
	"github.com/jibaru/do/internal/parser/taker"
	"github.com/jibaru/do/internal/reader"
	"github.com/jibaru/do/internal/request"
	"github.com/jibaru/do/internal/types"
	"github.com/jibaru/do/internal/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("use: do <filename.do>")
		return
	}

	filename := os.Args[1]

	doFileReader := reader.NewFileReader()
	commentCleaner := cleaner.New()
	sectionTaker := taker.New()
	sectionNormalizer := normalizer.New()
	sectionPartitioner := partitioner.New()
	expressionAnalyzer := analyzer.New()
	sectionExtractor := extractor.New(sectionTaker, sectionNormalizer, sectionPartitioner, expressionAnalyzer)
	variablesReplacer := replacer.New()
	funcCaller := caller.New()
	theParser := parser.New(doFileReader, commentCleaner, sectionExtractor, variablesReplacer, funcCaller)
	client := request.NewHttpClient(&http.Client{})

	output := types.CommandLineOutput{}

	doFile, err := theParser.ParseFromFilename(filename)
	if err != nil {
		output.Error = utils.Ptr(err.Error())
		fmt.Println(output.MarshalIndent())
		return
	}

	output.DoFile = *doFile

	response, err := client.Do(*doFile)
	if err != nil {
		output.Error = utils.Ptr(err.Error())
		fmt.Println(output.MarshalIndent())
		return
	}

	output.Response = response
	fmt.Println(output.MarshalIndent())
}
