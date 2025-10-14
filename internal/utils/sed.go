package utils

import (
	"fmt"
	"os"
	"regexp"
)

// Sed performs a search-and-replace operation on the specified file using the provided regular expression and replacement string.
func Sed(fileName string, re *regexp.Regexp, replace string) error {

	stat, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return fmt.Errorf("%s should be a file", fileName)
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	index := re.FindSubmatchIndex(data)
	if len(index) == 0 {
		return fmt.Errorf("pattern %s not found in %s", re.String(), fileName)
	}

	parts := [][]byte{}
	if len(index) == 2 {
		// only one match, no capture group
		parts = append(parts, data[:index[0]])
		parts = append(parts, []byte(replace))
		parts = append(parts, data[index[1]:])
	}
	if len(index) == 4 {
		// one capture group
		parts = append(parts,
			data[:index[2]],
			[]byte(replace),
			data[index[3]:],
		)
	}

	fp, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, stat.Mode())
	if err != nil {
		return err
	}
	defer fp.Close()

	for _, part := range parts {
		_, err = fp.Write(part)
		if err != nil {
			return err
		}
	}
	return nil
}
