package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jibaru/do/internal/types"
)

func readLetSection(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	content := string(data)

	letRegex := regexp.MustCompile(`(?s)let\s*\{(.*?)\}`)

	letMatch := letRegex.FindStringSubmatch(content)
	if len(letMatch) > 1 {
		letContent := letMatch[1]
		return strings.TrimSpace(letContent), nil
	}

	return "", errors.New("let section not found")
}

func parseLetSection(letContent string) (map[string]interface{}, error) {
	if strings.TrimSpace(letContent) == "" {
		return nil, nil
	}

	letVariables := make(map[string]interface{})

	keyValueRegex := regexp.MustCompile(`(\w+)\s*=\s*(".+?"|\d+|\d+\.\d+)`)

	matches := keyValueRegex.FindAllStringSubmatch(letContent, -1)
	for _, match := range matches {
		key := match[1]
		value := match[2]

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			letVariables[key] = strings.Trim(value, "\"")
		} else if i, err := strconv.Atoi(value); err == nil {
			letVariables[key] = i
		} else if f, err := strconv.ParseFloat(value, 64); err == nil {
			letVariables[key] = f
		}
	}

	return letVariables, nil
}

func readDoSection(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	content := string(data)

	doRegex := regexp.MustCompile(`do\s*\{([^{}]*(\{[^{}]*\}[^{}]*)*)\}`)

	doMatch := doRegex.FindStringSubmatch(content)
	if len(doMatch) > 1 {
		return strings.TrimSpace(doMatch[1]), nil
	}

	return "", errors.New("do section not found")
}

func parseDoSection(doContent string) (map[string]interface{}, error) {
	doVariables := make(map[string]interface{})

	keyValueRegex := regexp.MustCompile(`(\w+)\s*=\s*(\{.*?\}|".*?"|\w+)`)

	matches := keyValueRegex.FindAllStringSubmatch(doContent, -1)
	for _, match := range matches {
		key := match[1]
		value := match[2]

		// if json
		if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(value), &obj); err != nil {
				return nil, err
			}
			doVariables[key] = obj
		} else {
			// if string
			value = strings.Trim(value, "\"")
			doVariables[key] = value
		}
	}

	return doVariables, nil
}

func replaceVariables(doVariables map[string]interface{}, letVariables map[string]interface{}) {
	if letVariables == nil {
		return
	}

	for key, value := range doVariables {
		switch val := value.(type) {
		case string:
			doVariables[key] = replaceStringVariables(val, letVariables)
		case map[string]interface{}:
			replaceVariables(val, letVariables)
		}
	}
}

func replaceStringVariables(value string, letVariables map[string]interface{}) string {
	for key, val := range letVariables {
		value = strings.ReplaceAll(value, fmt.Sprintf("$%s", key), fmt.Sprintf("%v", val))
	}
	return value
}

func Filename(filename string) (*types.DoFile, error) {
	letContent, err := readLetSection(filename)
	if err != nil {
		return nil, err
	}

	letVariables, err := parseLetSection(letContent)
	if err != nil {
		return nil, err
	}

	doContent, err := readDoSection(filename)
	if err != nil {
		return nil, err
	}

	doVariables, err := parseDoSection(doContent)
	if err != nil {
		return nil, err
	}

	if doVariables == nil {
		return nil, errors.New("do section is empty")
	}

	if doVariables["method"] == nil {
		return nil, errors.New("method is required")
	}

	if doVariables["url"] == nil {
		return nil, errors.New("url is required")
	}

	replaceVariables(doVariables, letVariables)

	doFile := &types.DoFile{
		Let: types.Let{
			Variables: letVariables,
		},
		Do: types.Do{
			Method: doVariables["method"].(string),
			URL:    doVariables["url"].(string),
		},
	}

	if mp, ok := doVariables["params"]; ok {
		doFile.Do.Params = mp.(map[string]interface{})
	}

	if mp, ok := doVariables["query"]; ok {
		doFile.Do.Query = mp.(map[string]interface{})
	}

	if mp, ok := doVariables["headers"]; ok {
		doFile.Do.Headers = mp.(map[string]interface{})
	}

	if mp, ok := doVariables["body"]; ok {
		doFile.Do.Body = mp.(map[string]interface{})
	}

	return doFile, nil
}
