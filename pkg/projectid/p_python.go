package projectid

import (
	"fmt"
	"path"
	"regexp"

	"github.com/elsejj/verit/internal/utils"
	"github.com/elsejj/verit/pkg/version"
)

type PythonProject struct {
	workdir string
}

func (p *PythonProject) versionFile() string {
	return path.Join(p.workdir, "pyproject.toml")
}

func isPython(workdir string) bool {
	return utils.FileExists(path.Join(workdir, "pyproject.toml"))
}

func (p *PythonProject) IsMe(workdir string) bool {
	return isPython(workdir)
}

func (p *PythonProject) ID() ProjectID {
	return Python
}

func (p *PythonProject) WorkDir() string {
	return p.workdir
}

var pythonVersionRE = regexp.MustCompile(`\s*version\s*=\s*"(.+)"`)

func (p *PythonProject) GetVersion() (*version.Version, error) {
	v, err := utils.Grep(p.versionFile(), pythonVersionRE)
	if err != nil {
		return nil, fmt.Errorf("version not found")
	}

	return version.Parse(v)
}

func (p *PythonProject) SetVersion(v *version.Version) error {
	return utils.Sed(p.versionFile(), pythonVersionRE, v.String())
}

var _ Project = &PythonProject{}
