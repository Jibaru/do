package parser

import (
	"errors"
	"fmt"

	"github.com/jibaru/do/internal/parser/extractor"
	"github.com/jibaru/do/internal/parser/replacer"
	"github.com/jibaru/do/internal/reader"
	"github.com/jibaru/do/internal/types"
)

type Parser interface {
	ParseFromFilename(filename string) (*types.DoFile, error)
}

type parser struct {
	doFileReader      reader.FileReader
	sectionExtractor  extractor.Extractor
	variablesReplacer replacer.Replacer
}

func New(
	doFileReader reader.FileReader,
	sectionExtractor extractor.Extractor,
	variablesReplacer replacer.Replacer,
) Parser {
	return &parser{
		doFileReader,
		sectionExtractor,
		variablesReplacer,
	}
}

func (p *parser) ParseFromFilename(filename string) (*types.DoFile, error) {
	content, err := p.doFileReader.Read(filename)
	if err != nil {
		return nil, err
	}

	letVariables, err := p.sectionExtractor.Extract(types.LetSection, content)
	if err != nil {
		if errors.Is(err, extractor.ErrSectionExtractorNoBlock) {
			letVariables = nil
		} else {
			return nil, err
		}
	}

	doVariables, err := p.sectionExtractor.Extract(types.DoSection, content)
	if err != nil {
		return nil, err
	}

	if doVariables == nil {
		return nil, NewDoSectionEmptyError()
	}

	if doVariables[types.DoMethod] == nil {
		return nil, NewMethodRequiredError()
	}

	if doVariables[types.DoURL] == nil {
		return nil, NewURLRequiredError()
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

	err = p.variablesReplacer.Replace(doVariables, letVariables)
	if err != nil {
		return nil, err
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
		body := mp.(types.String)
		if body != "" {
			doFile.Do.Body = &body
		}
	}

	return doFile, nil
}
