package types

import (
	"encoding/json"
)

// Section defines the type of section
type Section string

// FileReaderContent defines the content of a .do file
type FileReaderContent string

// RawSectionContent defines the raw content of a section
type RawSectionContent string

// NormalizedSectionContent defines the normalized content of a section
type NormalizedSectionContent string

// SectionExpressions defines the expressions of a section
type SectionExpressions []string

// Let defines the variables section
type Let struct {
	Variables map[string]interface{} `json:"variables"`
}

// Do defines the request section
type Do struct {
	Method  String  `json:"method"`
	URL     String  `json:"url"`
	Params  Map     `json:"params"`
	Query   Map     `json:"query"`
	Headers Map     `json:"headers"`
	Body    *String `json:"body"`
}

// DoFile is the representation of file.do
type DoFile struct {
	Let Let `json:"let"`
	Do  Do  `json:"do"`
}

// Response defines the response of a request
type Response struct {
	StatusCode int                    `json:"status_code"`
	Body       string                 `json:"body"`
	Headers    map[string]interface{} `json:"headers"`
}

// CommandLineOutput defines the output of the command line
type CommandLineOutput struct {
	DoFile   DoFile    `json:"do_file"`
	Response *Response `json:"response"`
	Error    *string   `json:"error"`
}

// MarshalIndent returns the JSON representation of CommandLineOutput
func (c CommandLineOutput) MarshalIndent() string {
	value, _ := json.MarshalIndent(c, "", "   ")
	return string(value)
}
