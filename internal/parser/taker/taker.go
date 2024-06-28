package taker

import (
	"strings"

	"github.com/jibaru/do/internal/types"
)

type Taker interface {
	Take(section types.Section, text types.CleanedContent) (types.RawSectionContent, error)
}

type taker struct{}

func New() Taker {
	return &taker{}
}

func (t *taker) Take(section types.Section, text types.CleanedContent) (types.RawSectionContent, error) {
	inBlock := false
	textInNoBlocks := strings.Builder{}

	for _, ch := range text {
		toWrite := ch
		if inBlock {
			toWrite = ' '
		}

		if ch == '{' {
			inBlock = true
			toWrite = ' '
		}

		if ch == '}' {
			inBlock = false
			toWrite = ' '
		}

		textInNoBlocks.WriteRune(toWrite)
	}

	textWithOnlyBlocks := textInNoBlocks.String()
	startIndex := strings.Index(textWithOnlyBlocks, string(section))
	if startIndex == -1 {
		return "", NewNoBlockError()
	}

	content := strings.Builder{}
	inBlock = false
	inString := false
	openBracesCount := 0
	foundOpeningBrace := false

	for i := startIndex + len(string(section)); i < len(text); i++ {
		if text[i] == '{' && !inString {
			if !foundOpeningBrace {
				foundOpeningBrace = true
				inBlock = true
			}
			openBracesCount++
			if openBracesCount > 1 {
				content.WriteByte(text[i])
			}
			continue
		}

		if text[i] == '}' && !inString {
			openBracesCount--
			if openBracesCount == 0 {
				inBlock = false
				break
			}
			content.WriteByte(text[i])
			continue
		}

		if text[i] == '"' {
			inString = !inString
		}

		if inBlock {
			content.WriteByte(text[i])
		}
	}

	if !foundOpeningBrace {
		return "", NewMissingOpeningBraceError()
	}

	if openBracesCount != 0 {
		return "", NewMissingClosingBraceError()
	}

	sectionContent := content.String()

	if strings.TrimSpace(sectionContent) == "" {
		return "", NewNoBlockError()
	}

	return types.RawSectionContent(sectionContent), nil
}
