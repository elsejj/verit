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

type PythonProject struct {
	workdir string
}

func (p *PythonProject) IsMe(workdir string) bool {
	return isNode(workdir)
}

func (p *PythonProject) ID() ProjectID {
	return Python
}

func (p *PythonProject) WorkDir() string {
	return p.workdir
}

func (p *PythonProject) GetVersion() (*version.Version, error) {
	versionFile := path.Join(p.workdir, "pyproject.toml")
	fp, err := os.Open(versionFile)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	scaner := bufio.NewScanner(fp)

	versionRE := regexp.MustCompile(`^\s*version\s*=\s*"(.+)"`)

	for scaner.Scan() {
		if m := versionRE.FindSubmatch(scaner.Bytes()); m != nil {
			return version.Parse(string(m[1])), nil
		}
	}
	return nil, fmt.Errorf("version not found")
}

func (p *PythonProject) SetVersion(v *version.Version) error {
	versionFile := path.Join(p.workdir, "pyproject.toml")
	fp, err := os.Open(versionFile)
	if err != nil {
		return err
	}

	scaner := bufio.NewScanner(fp)

	versionRE := regexp.MustCompile(`^\s*version\s*=\s*"(.+)"`)

	buf := make([]byte, 0, 1024)
	for scaner.Scan() {
		line := scaner.Bytes()
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

var _ Project = &PythonProject{}
