package projectid

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/elsejj/verit/internal/utils"
	"github.com/elsejj/verit/pkg/version"
)

/*
GoProject represents a Go project with versioning capabilities
Because Go projects do not have a standard version file, we flow the rules below:
  - lookup the project by checking the existence of "version.txt"
  - this file can be embedded to a go variable use `go:embed` directive
  - the file content should be like `x.y.z`
*/
type GoProject struct {
	workdir           string
	_versionFile      string
	_versionFileFound bool
}

func isGo(workdir string) bool {
	return utils.FileExists(path.Join(workdir, "go.mod"))
}

func (p *GoProject) versionFile() string {
	if p._versionFile != "" {
		return p._versionFile
	}
	if !p._versionFileFound {
		p._versionFile, _ = utils.FindFileDown(p.workdir, "version.txt")
		p._versionFileFound = true
	}
	return p._versionFile
}

func (p *GoProject) IsMe(workdir string) bool {
	return isGo(workdir)
}

func (p *GoProject) ID() ProjectID {
	return Go
}

func (p *GoProject) WorkDir() string {
	return p.workdir
}

func (p *GoProject) GetVersion() (*version.Version, error) {
	versionFile := p.versionFile()
	data, err := os.ReadFile(versionFile)
	if err != nil {
		return nil, fmt.Errorf("version.txt not found")
	}
	v, err := version.Parse(string(bytes.TrimSpace(data)))
	if err != nil {
		return nil, fmt.Errorf("parse version from %s failed: %w", versionFile, err)
	}
	return v, nil
}

func (p *GoProject) SetVersion(v *version.Version) error {
	versionFile := p.versionFile()
	if versionFile == "" {
		return fmt.Errorf("version.txt not found, please create one")
	}
	fp, err := os.Create(versionFile)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.WriteString(v.String())
	return err
}

var _ Project = &GoProject{}
