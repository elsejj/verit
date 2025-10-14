package utils

import (
	"fmt"
	"os"
	"regexp"
)

// Grep searches fileName for the first match of the provided regular expression
// and returns the matched content (full match followed by any capture groups).
func Grep(fileName string, re *regexp.Regexp) (string, error) {
	info, err := os.Stat(fileName)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("%s should be a file", fileName)
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	found := re.FindSubmatch(data)
	if len(found) == 0 {
		return "", fmt.Errorf("pattern %s not found in %s", re.String(), fileName)
	}
	if len(found) == 1 {
		// only one match, no capture group
		return string(found[0]), nil
	}
	// one capture group
	if len(found) == 2 {
		return string(found[1]), nil
	}
	return "", fmt.Errorf("multiple results found, which may be unexpected. ensure your regex has at most one capture group")
}
