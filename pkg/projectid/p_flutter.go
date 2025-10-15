package projectid

import (
	"fmt"
	"path"
	"regexp"

	"github.com/elsejj/verit/internal/utils"
	"github.com/elsejj/verit/pkg/version"
)

type FlutterProject struct {
	workdir string
}

func (p *FlutterProject) versionFile() string {
	return path.Join(p.workdir, "pubspec.yaml")
}

func isFlutter(workdir string) bool {
	return utils.FileExists(path.Join(workdir, "pubspec.yaml"))
}

func (p *FlutterProject) IsMe(workdir string) bool {
	return isFlutter(workdir)
}

func (p *FlutterProject) ID() ProjectID {
	return Flutter
}

func (p *FlutterProject) WorkDir() string {
	return p.workdir
}

var flutterVersionRE = regexp.MustCompile(`(?m)^\s*version\s*:\s*["']?([^\s"']+)["']?`)

func (p *FlutterProject) GetVersion() (*version.Version, error) {
	v, err := utils.Grep(p.versionFile(), flutterVersionRE)
	if err != nil {
		return nil, fmt.Errorf("version not found")
	}

	return version.Parse(v)
}

func (p *FlutterProject) SetVersion(v *version.Version) error {
	return utils.Sed(p.versionFile(), flutterVersionRE, v.String())
}

var _ Project = &FlutterProject{}
