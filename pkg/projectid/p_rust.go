package projectid

import (
	"fmt"
	"path"
	"regexp"

	"github.com/elsejj/verit/internal/utils"
	"github.com/elsejj/verit/pkg/version"
)

type RustProject struct {
	workdir string
}

func (p *RustProject) versionFile() string {
	return path.Join(p.workdir, "Cargo.toml")
}

func isRust(workdir string) bool {
	return utils.FileExists(path.Join(workdir, "Cargo.toml"))
}

func (p *RustProject) IsMe(workdir string) bool {
	return isRust(workdir)
}

func (p *RustProject) ID() ProjectID {
	return Rust
}

func (p *RustProject) WorkDir() string {
	return p.workdir
}

var RustVersionRE = regexp.MustCompile(`\s*version\s*=\s*"(.+)"`)

func (p *RustProject) GetVersion() (*version.Version, error) {
	v, err := utils.Grep(p.versionFile(), RustVersionRE)
	if err != nil {
		return nil, fmt.Errorf("version not found")
	}

	return version.Parse(v)
}

func (p *RustProject) SetVersion(v *version.Version) error {
	return utils.Sed(p.versionFile(), RustVersionRE, v.String())
}

var _ Project = &RustProject{}
