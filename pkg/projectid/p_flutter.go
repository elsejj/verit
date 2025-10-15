package projectid

import (
	"fmt"
	"path"
	"regexp"
	"strings"

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
	raw, err := utils.Grep(p.versionFile(), flutterVersionRE)
	if err != nil {
		return nil, fmt.Errorf("version not found")
	}

	base := raw
	build := ""
	if idx := strings.Index(base, "+"); idx >= 0 {
		base = base[:idx]
		build = raw[idx+1:]
	}

	v, err := version.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("parse version: %w", err)
	}
	if build != "" {
		v.Build = build
	}
	return v, nil
}

func (p *FlutterProject) SetVersion(v *version.Version) error {
	value := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Build != "" {
		value = fmt.Sprintf("%s+%s", value, v.Build)
	}

	return utils.Sed(p.versionFile(), flutterVersionRE, value)
}

var _ Project = &FlutterProject{}
