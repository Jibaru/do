package parser

import (
	"errors"

	"github.com/jibaru/do/internal/parser/extractor"
	"github.com/jibaru/do/internal/parser/replacer"
	"github.com/jibaru/do/internal/reader"
	"github.com/jibaru/do/internal/types"
)

var (
	ErrDoSectionEmpty = errors.New("do section is empty")
	ErrMethodRequired = errors.New("method is required")
	ErrURLRequired    = errors.New("url is required")
)

type Parser interface {
	FromFilename(filename string) (*types.DoFile, error)
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

func (p *parser) FromFilename(filename string) (*types.DoFile, error) {
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
		return nil, ErrDoSectionEmpty
	}

	if doVariables[types.DoMethod] == nil {
		return nil, ErrMethodRequired
	}

	if doVariables[types.DoURL] == nil {
		return nil, ErrURLRequired
	}

	p.variablesReplacer.Replace(doVariables, letVariables)

	doFile := &types.DoFile{
		Let: types.Let{
			Variables: letVariables,
		},
		Do: types.Do{
			Method: doVariables[types.DoMethod].(string),
			URL:    doVariables[types.DoURL].(string),
		},
	}

	if mp, ok := doVariables[types.DoParams]; ok {
		doFile.Do.Params = mp.(map[string]interface{})
	}

	if mp, ok := doVariables[types.DoQuery]; ok {
		doFile.Do.Query = mp.(map[string]interface{})
	}

	if mp, ok := doVariables[types.DoHeaders]; ok {
		doFile.Do.Headers = mp.(map[string]interface{})
	}

	if mp, ok := doVariables[types.DoBody]; ok {
		doFile.Do.Body = mp.(string)
	}

	return doFile, nil
}
