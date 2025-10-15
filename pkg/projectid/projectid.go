package projectid

import (
	"os"

	"github.com/elsejj/verit/pkg/version"
)

var projectDetectionOrder = []ProjectID{
	Python,
	Go,
	Node,
	Flutter,
	Rust,
}

var projectCheckers = map[ProjectID]func(string) bool{
	Node:    isNode,
	Python:  isPython,
	Go:      isGo,
	Flutter: isFlutter,
	Rust:    isRust,
}

// Project represents a generic project with versioning capabilities
type Project interface {
	// Test if the project is of the specified type
	IsMe(workdir string) bool
	// Get the project ID
	ID() ProjectID
	// Get the working directory of the project
	WorkDir() string
	// Get the current version of the project
	GetVersion() (*version.Version, error)
	// Set a new version for the project
	SetVersion(v *version.Version) error
}

// Pwd returns the current working directory
func Pwd() string {
	pwd, _ := os.Getwd()
	return pwd
}

// Which returns the type of project in current directory
func Which(workdir string) ProjectID {
	var matches []ProjectID
	for _, id := range projectDetectionOrder {
		checker, ok := projectCheckers[id]
		if !ok {
			continue
		}
		if checker(workdir) {
			matches = append(matches, id)
		}
	}

	if len(matches) > 1 {
		return Multiple
	}
	if len(matches) == 1 {
		return matches[0]
	}
	return 0
}
