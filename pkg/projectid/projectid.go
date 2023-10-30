package projectid

import (
	"os"
	"path"

	"github.com/elsejj/verit/pkg/version"
)

type Project interface {
	IsMe(workdir string) bool
	ID() ProjectID
	WorkDir() string
	GetVersion() (*version.Version, error)
	SetVersion(v *version.Version) error
}

// Pwd returns the current working directory
func Pwd() string {
	pwd, _ := os.Getwd()
	return pwd
}

// Which returns the type of project in current directory
func Which(workdir string) ProjectID {
	checkers := map[ProjectID]func(string) bool{
		Node:   isNode,
		Python: isPython,
		Go:     isGo,
	}

	for id, checker := range checkers {
		if checker(workdir) {
			return id
		}
	}
	return 0
}

func fileExists(paths ...string) bool {
	fname := path.Join(paths...)
	_, err := os.Stat(fname)
	return err == nil
}

func isNode(workdir string) bool {
	return fileExists(workdir, "package.json")
}

func isPython(workdir string) bool {
	return fileExists(workdir, "pyproject.toml")
}

func isGo(workdir string) bool {
	return fileExists(workdir, "go.mod")
}
