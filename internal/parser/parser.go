package parser

import (
	"errors"
	"fmt"

	"github.com/jibaru/do/internal/parser/caller"
	"github.com/jibaru/do/internal/parser/cleaner"
	"github.com/jibaru/do/internal/parser/extractor"
	"github.com/jibaru/do/internal/parser/replacer"
	"github.com/jibaru/do/internal/parser/resolver"
	"github.com/jibaru/do/internal/reader"
	"github.com/jibaru/do/internal/types"
)

type Parser interface {
	ParseFromFilename(filename string) (*types.DoFile, error)
}

type parser struct {
	doFileReader      reader.FileReader
	commentCleaner    cleaner.Cleaner
	sectionExtractor  extractor.Extractor
	variablesReplacer replacer.DoReplacer
	funcCaller        caller.Caller
	letResolver       resolver.LetResolver
}

func New(
	doFileReader reader.FileReader,
	commentCleaner cleaner.Cleaner,
	sectionExtractor extractor.Extractor,
	variablesReplacer replacer.DoReplacer,
	funcCaller caller.Caller,
	letResolver resolver.LetResolver,
) Parser {
	return &parser{
		doFileReader,
		commentCleaner,
		sectionExtractor,
		variablesReplacer,
		funcCaller,
		letResolver,
	}
}

func (p *parser) ParseFromFilename(filename string) (*types.DoFile, error) {
	content, err := p.doFileReader.Read(filename)
	if err != nil {
		return nil, err
	}

	cleanedContent, err := p.commentCleaner.Clean(content)
	if err != nil {
		return nil, err
	}

	letSentences, err := p.sectionExtractor.Extract(types.LetSection, cleanedContent)
	if err != nil {
		if !errors.Is(err, extractor.ErrSectionExtractorNoBlock) {
			return nil, err
		}
	}

	doSentences, err := p.sectionExtractor.Extract(types.DoSection, cleanedContent)
	if err != nil {
		return nil, err
	}

	if doSentences == nil {
		return nil, NewDoSectionEmptyError()
	}

	if !doSentences.Has(types.DoMethod) {
		return nil, NewMethodRequiredError()
	}

	if !doSentences.Has(types.DoURL) {
		return nil, NewURLRequiredError()
	}

	letSentences, err = p.letResolver.Resolve(letSentences)
	if err != nil {
		return nil, err
	}

	var letVariables map[string]interface{}
	if letSentences != nil {
		letVariables = letSentences.ToMap()
	}

	doVariables := doSentences.ToMap()
	err = p.variablesReplacer.Replace(doVariables, letVariables)
	if err != nil {
		return nil, err
	}

	err = p.funcCaller.Call(doVariables)
	if err != nil {
		return nil, err
	}

	if _, ok := doVariables[types.DoMethod].(types.String); !ok {
		return nil, NewTypeNotExpectedError(
			types.DoMethod,
			fmt.Sprintf("%T", types.String("")),
			fmt.Sprintf("%T", doVariables[types.DoMethod]),
		)
	}

	if _, ok := doVariables[types.DoURL].(types.String); !ok {
		return nil, NewTypeNotExpectedError(
			types.DoURL,
			fmt.Sprintf("%T", types.String("")),
			fmt.Sprintf("%T", doVariables[types.DoURL]),
		)
	}

	if doVariables[types.DoParams] != nil {
		mp, ok := doVariables[types.DoParams].(types.Map)
		if !ok || (mp != nil && !mp.HasBasicTypesValues()) {
			return nil, NewTypeNotExpectedError(
				types.DoParams,
				fmt.Sprintf("types.Map[string]basic types"),
				fmt.Sprintf("%T", doVariables[types.DoParams]),
			)
		}
	}

	if doVariables[types.DoQuery] != nil {
		mp, ok := doVariables[types.DoQuery].(types.Map)
		if !ok || (mp != nil && !mp.HasBasicTypesValues()) {
			return nil, NewTypeNotExpectedError(
				types.DoQuery,
				fmt.Sprintf("types.Map[string]basic types"),
				fmt.Sprintf("%T", doVariables[types.DoQuery]),
			)
		}
	}

	if doVariables[types.DoHeaders] != nil {
		mp, ok := doVariables[types.DoHeaders].(types.Map)
		if !ok || (mp != nil && !mp.HasBasicTypesValues()) {
			return nil, NewTypeNotExpectedError(
				types.DoHeaders,
				fmt.Sprintf("types.Map[string]string"),
				fmt.Sprintf("%T", doVariables[types.DoHeaders]),
			)
		}
	}

	doFile := &types.DoFile{
		Let: types.Let{
			Variables: letVariables,
		},
		Do: types.Do{
			Method: doVariables[types.DoMethod].(types.String),
			URL:    doVariables[types.DoURL].(types.String),
		},
	}

	if mp, ok := doVariables[types.DoParams]; ok {
		doFile.Do.Params = mp.(types.Map)
	}

	if mp, ok := doVariables[types.DoQuery]; ok {
		doFile.Do.Query = mp.(types.Map)
	}

	if mp, ok := doVariables[types.DoHeaders]; ok {
		doFile.Do.Headers = mp.(types.Map)
	}

	if mp, ok := doVariables[types.DoBody]; ok {
		switch mp.(type) {
		case types.String:
			doFile.Do.Body = mp.(types.String)
		case types.Map:
			doFile.Do.Body = mp.(types.Map)
		default:
			return nil, NewTypeNotExpectedError(
				types.DoBody,
				fmt.Sprintf("types.String or types.Map[string]interface{}"),
				fmt.Sprintf("%T", doVariables[types.DoBody]),
			)
		}
	}

	return doFile, nil
}
