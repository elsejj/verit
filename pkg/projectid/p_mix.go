package projectid

import (
	"fmt"

	"github.com/elsejj/verit/pkg/version"
)

type MixProject struct {
	workdir  string
	projects []Project
}

func (p *MixProject) scanProjects() []Project {
	if p.projects != nil {
		return p.projects
	}

	for _, id := range projectDetectionOrder {
		if id == Mix {
			continue
		}
		checker, ok := projectCheckers[id]
		if !ok || !checker(p.workdir) {
			continue
		}
		sub := id.Project(p.workdir)
		if sub == nil {
			continue
		}
		p.projects = append(p.projects, sub)
	}
	return p.projects
}

func (p *MixProject) IsMe(workdir string) bool {
	if workdir != p.workdir {
		return false
	}
	return len(p.scanProjects()) > 1
}

func (p *MixProject) ID() ProjectID {
	return Mix
}

func (p *MixProject) WorkDir() string {
	return p.workdir
}

func (p *MixProject) GetVersion() (*version.Version, error) {
	if len(p.projects) == 0 {
		return nil, fmt.Errorf("no supported projects detected in %s", p.workdir)
	}

	var current *version.Version
	var currentID ProjectID

	for _, sub := range p.projects {
		v, err := sub.GetVersion()
		if err != nil {
			return nil, fmt.Errorf("%s project: %w", sub.ID(), err)
		}
		if current == nil {
			current = v
			currentID = sub.ID()
			continue
		}
		if current.Build != v.Build || !current.Equal(v) {
			return nil, fmt.Errorf("version mismatch between %s and %s projects", currentID, sub.ID())
		}
	}

	return current, nil
}

func (p *MixProject) SetVersion(v *version.Version) error {
	if len(p.projects) == 0 {
		return fmt.Errorf("no supported projects detected in %s", p.workdir)
	}

	for _, sub := range p.projects {
		if err := sub.SetVersion(v); err != nil {
			return fmt.Errorf("%s project: %w", sub.ID(), err)
		}
	}

	return nil
}

var _ Project = &MixProject{}
