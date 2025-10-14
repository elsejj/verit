package utils

import (
	"os"
	"path"
)

// FileExists checks if a file or directory exists at the given path.
func FileExists(paths ...string) bool {
	fName := path.Join(paths...)
	_, err := os.Stat(fName)
	return err == nil
}

// FindFileDown searches for a file named fileName starting from startDir and it's subdirectories.
// Returns the full path to the file and true if found, otherwise returns an empty string and false.
func FindFileDown(startDir, fileName string) (string, bool) {
	info, err := os.Stat(startDir)
	if err != nil {
		return "", false
	}

	if !info.IsDir() {
		if path.Base(startDir) == fileName {
			return startDir, true
		}
		return "", false
	}

	stack := []string{startDir}

	for len(stack) > 0 {
		dir := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			fullPath := path.Join(dir, entry.Name())

			if entry.Name() == fileName {
				return fullPath, true
			}

			if entry.IsDir() {
				stack = append(stack, fullPath)
			}
		}
	}

	return "", false
}
