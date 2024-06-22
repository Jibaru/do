package types

var (
	ReservedKeywordNames = map[string]struct{}{
		"do":       {},
		"let":      {},
		"true":     {},
		"false":    {},
		"null":     {},
		"if":       {},
		"else":     {},
		"for":      {},
		"range":    {},
		"break":    {},
		"continue": {},
		"return":   {},
		"func":     {},
		"import":   {},
		"from":     {},
		"as":       {},
		"const":    {},
		"var":      {},
		"string":   {},
		"int":      {},
		"float":    {},
		"bool":     {},
		"map":      {},
		"list":     {},
	}
)

func IsReservedKeyword(name string) bool {
	_, ok := ReservedKeywordNames[name]
	return ok
}
