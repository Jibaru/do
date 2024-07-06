package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Parse reads the file at the given path and returns a map of environment variables.
func Parse(filepath string) (map[string]string, error) {
	envMap := make(map[string]string)

	file, err := os.Open(filepath)
	if err != nil {
		return nil, NewCanNotReadFileError(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if isCommentOrEmpty(line) {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, err
		}

		envMap[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, NewCanNotReadFileError(err.Error())
	}

	return envMap, nil
}

// isCommentOrEmpty returns true if the line is a comment or empty.
func isCommentOrEmpty(line string) bool {
	return len(line) == 0 || strings.HasPrefix(line, "#")
}

// parseLine returns the key and value of the line.
func parseLine(line string) (string, string, error) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", NewCanNotReadFileError(fmt.Sprintf("invalid line: %s", line))
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	value = removeQuotes(value)

	return key, value, nil
}

// removeQuotes removes the quotes from the value.
func removeQuotes(value string) string {
	if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
		(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
		value = strings.Trim(value, "\"'")
	}
	return value
}

// ParseAndSet reads the file at the given path and sets the environment variables.
func ParseAndSet(filepath string) error {
	vars, err := Parse(filepath)
	if err != nil {
		return err
	}

	for key, value := range vars {
		err = os.Setenv(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}
