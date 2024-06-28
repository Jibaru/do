package cleaner

import (
	"strings"

	"github.com/jibaru/do/internal/types"
)

type Cleaner interface {
	Clean(rawContent types.FileReaderContent) (types.CleanedContent, error)
}

type cleaner struct{}

func New() Cleaner {
	return &cleaner{}
}

// Clean removes all the comments from the content
func (d *cleaner) Clean(rawContent types.FileReaderContent) (types.CleanedContent, error) {
	// Comments starts with // and ends with \n, and comments inside strings are not removed
	// Consider string with " and `
	// When a comment is found, the rest of the line should be omitted
	content := string(rawContent)

	var result strings.Builder
	inQuotes := false
	inBackticks := false
	inComment := false

	for i, ch := range content {
		if ch == '"' && !inBackticks {
			inQuotes = !inQuotes
		}
		if ch == '`' && !inQuotes {
			inBackticks = !inBackticks
		}
		if !inQuotes && !inBackticks && !inComment && ch == '/' && i+1 < len(content) && content[i+1] == '/' {
			inComment = true
		}
		if inComment && ch == '\n' {
			inComment = false
		}
		if !inComment {
			result.WriteByte(content[i])
		}
	}

	return types.CleanedContent(result.String()), nil
}
