package types

// Let defines the variables section
type Let struct {
	Variables map[string]interface{} `json:"variables"`
}

// Do defines the request section
type Do struct {
	Method  string                 `json:"method"`
	URL     string                 `json:"url"`
	Params  map[string]interface{} `json:"params"`
	Query   map[string]interface{} `json:"query"`
	Headers map[string]interface{} `json:"headers"`
	Body    string                 `json:"body"`
}

// DoFile is the representation of file.do
type DoFile struct {
	Let Let `json:"let"`
	Do  Do  `json:"do"`
}
