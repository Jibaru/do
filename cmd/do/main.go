package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/jibaru/do/internal/env"
	"github.com/jibaru/do/internal/parser"
	"github.com/jibaru/do/internal/parser/analyzer"
	"github.com/jibaru/do/internal/parser/caller"
	"github.com/jibaru/do/internal/parser/cleaner"
	"github.com/jibaru/do/internal/parser/extractor"
	"github.com/jibaru/do/internal/parser/normalizer"
	"github.com/jibaru/do/internal/parser/partitioner"
	"github.com/jibaru/do/internal/parser/replacer"
	"github.com/jibaru/do/internal/parser/resolver"
	"github.com/jibaru/do/internal/parser/taker"
	"github.com/jibaru/do/internal/reader"
	"github.com/jibaru/do/internal/request"
	"github.com/jibaru/do/internal/types"
	"github.com/jibaru/do/internal/utils"
)

const Version = "v0.0.0"

type params struct {
	versionFlag bool
	envPath     string
	filename    string
}

func main() {
	output := types.CommandLineOutput{}
	p, err := readParams()
	if err != nil {
		output.Error = utils.Ptr(err.Error())
		fmt.Println(output.MarshalIndent())
		return
	}

	if p.versionFlag {
		fmt.Println(Version)
		return
	}

	if p.envPath != "" {
		err = env.ParseAndSet(p.envPath)
		if err != nil {
			output.Error = utils.Ptr(err.Error())
			fmt.Println(output.MarshalIndent())
			return
		}
	}

	doFileReader := reader.NewFileReader()
	commentCleaner := cleaner.New()
	sectionTaker := taker.New()
	sectionNormalizer := normalizer.New()
	sectionPartitioner := partitioner.New()
	expressionAnalyzer := analyzer.New()
	sectionExtractor := extractor.New(sectionTaker, sectionNormalizer, sectionPartitioner, expressionAnalyzer)
	variablesReplacer := replacer.New()
	funcCaller := caller.New()
	letResolver := resolver.NewLetResolver()
	theParser := parser.New(doFileReader, commentCleaner, sectionExtractor, variablesReplacer, funcCaller, letResolver)
	client := request.NewHttpClient(&http.Client{})

	doFile, err := theParser.ParseFromFilename(p.filename)
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

func readParams() (params, error) {
	var p params

	flag.BoolVar(&p.versionFlag, "version", false, "Version of the tool")
	flag.BoolVar(&p.versionFlag, "v", false, "Version of the tool")

	flag.StringVar(&p.filename, "file", "", "Path to the do file (required)")
	flag.StringVar(&p.filename, "f", "", "Path to the do file (required)")

	flag.StringVar(&p.envPath, "env", "", "Path to the env file (optional)")
	flag.StringVar(&p.envPath, "e", "", "Path to the env file (optional)")

	flag.Parse()

	return p, nil
}
