package projectid

import (
	"fmt"
	"path"
	"regexp"

	"github.com/elsejj/verit/internal/utils"
	"github.com/elsejj/verit/pkg/version"
)

type NodeProject struct {
	workdir string
}

func (p *NodeProject) versionFile() string {
	return path.Join(p.workdir, "package.json")
}

func isNode(workdir string) bool {
	return utils.FileExists(path.Join(workdir, "package.json"))
}

func (p *NodeProject) IsMe(workdir string) bool {
	return isNode(workdir)
}

func (p *NodeProject) ID() ProjectID {
	return Node
}

func (p *NodeProject) WorkDir() string {
	return p.workdir
}

var nodeVersionRE = regexp.MustCompile(`^\s*"version"\s*:\s*"(.+)"`)

func (p *NodeProject) GetVersion() (*version.Version, error) {
	v, err := utils.Grep(p.versionFile(), nodeVersionRE)
	if err != nil {
		return nil, fmt.Errorf("version not found")
	}

	return version.Parse(v)

}

func (p *NodeProject) SetVersion(v *version.Version) error {
	return utils.Sed(p.versionFile(), nodeVersionRE, v.String())
}

var _ Project = &NodeProject{}
