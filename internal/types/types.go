package types

import (
	"encoding/json"
)

// Section defines the type of section
type Section string

// FileReaderContent defines the content of a .do file
type FileReaderContent string

// CleanedContent defines the cleaned content of a .do file
type CleanedContent string

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
	Method  String      `json:"method"`
	URL     String      `json:"url"`
	Params  Map         `json:"params"`
	Query   Map         `json:"query"`
	Headers Map         `json:"headers"`
	Body    interface{} `json:"body"`
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

// Sentence defines a key-value pair
type Sentence struct {
	Key   string
	Value interface{}
}

// Sentences defines a ordered list of sentences
type Sentences struct {
	All         []Sentence
	KeysWithIdx map[string]int
}

// NewSentences creates a new Sentences
func NewSentences() *Sentences {
	return &Sentences{
		All:         make([]Sentence, 0),
		KeysWithIdx: make(map[string]int),
	}
}

// NewSentencesFromSlice creates a new Sentences from a slice of Sentence
func NewSentencesFromSlice(sentences []Sentence) *Sentences {
	sentencesMap := make(map[string]int)
	for i, sentence := range sentences {
		sentencesMap[sentence.Key] = i
	}

	return &Sentences{
		All:         sentences,
		KeysWithIdx: sentencesMap,
	}
}

// Entries returns all sentences
func (s *Sentences) Entries() []Sentence {
	return s.All
}

// Get returns the value of a key
func (s *Sentences) Get(key string) (interface{}, bool) {
	idx, ok := s.KeysWithIdx[key]
	if !ok {
		return nil, false
	}

	return s.All[idx].Value, true
}

// Has returns true if the key exists
func (s *Sentences) Has(key string) bool {
	_, ok := s.KeysWithIdx[key]
	return ok
}

// Set sets the value of a key
func (s *Sentences) Set(key string, value interface{}) {
	idx, ok := s.KeysWithIdx[key]
	if ok {
		s.All[idx].Value = value
	} else {
		s.All = append(s.All, Sentence{Key: key, Value: value})
		s.KeysWithIdx[key] = len(s.All) - 1
	}
}

// ToMap converts Sentences to a Map
func (s *Sentences) ToMap() Map {
	m := make(map[string]interface{})
	for _, sentence := range s.All {
		m[sentence.Key] = sentence.Value
	}

	return m
}
