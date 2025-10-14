package projectid

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/elsejj/verit/pkg/version"
)

type NodeProject struct {
	workdir string
}

func (p *NodeProject) versionFile() string {
	return path.Join(p.workdir, "package.json")
}

func isNode(workdir string) bool {
	return fileExists(path.Join(workdir, "package.json"))
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

func (p *NodeProject) GetVersion() (*version.Version, error) {
	versionFile := p.versionFile()
	fp, err := os.Open(versionFile)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	versionRE := regexp.MustCompile(`^\s*"version"\s*:\s*"(.+)"`)

	for scanner.Scan() {
		if m := versionRE.FindSubmatch(scanner.Bytes()); m != nil {
			return version.Parse(string(m[1])), nil
		}
	}
	return nil, fmt.Errorf("version not found")
}

func (p *NodeProject) SetVersion(v *version.Version) error {
	versionFile := p.versionFile()
	fp, err := os.Open(versionFile)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(fp)

	versionRE := regexp.MustCompile(`^\s*"version"\s*:\s*"(.+)"`)

	buf := make([]byte, 0, 1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		if m := versionRE.FindSubmatch(line); m != nil {
			newVersion := bytes.Replace(line, m[1], []byte(v.String()), 1)
			buf = append(buf, newVersion...)
			buf = append(buf, '\n')
		} else {
			buf = append(buf, line...)
			buf = append(buf, '\n')
		}
	}
	fp.Close()

	return os.WriteFile(versionFile, buf, 0644)
}

var _ Project = &NodeProject{}
