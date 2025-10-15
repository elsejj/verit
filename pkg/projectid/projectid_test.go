package projectid

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/elsejj/verit/pkg/version"
)

func TestWhichDetectsMixProjects(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected ProjectID
	}{
		{
			name: "python only",
			files: map[string]string{
				"pyproject.toml": `
[tool.poetry]
name = "demo"
version = "1.2.3"
`,
			},
			expected: Python,
		},
		{
			name: "node only",
			files: map[string]string{
				"package.json": `{"name":"demo","version":"1.2.3"}`,
			},
			expected: Node,
		},
		{
			name: "flutter only",
			files: map[string]string{
				"pubspec.yaml": `
name: demo
version: 1.2.3+4
`,
			},
			expected: Flutter,
		},
		{
			name: "python and node",
			files: map[string]string{
				"pyproject.toml": `
[tool.poetry]
name = "demo"
version = "1.2.3"
`,
				"package.json": `{"name":"demo","version":"1.2.3"}`,
			},
			expected: Mix,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for name, content := range tt.files {
				writeFile(t, dir, name, content)
			}

			if got := Which(dir); got != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestMixProjectVersionOperations(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "pyproject.toml", `
[tool.poetry]
name = "demo"
version = "1.2.3"
`)
	writeFile(t, dir, "package.json", `{"name":"demo","version":"1.2.3"}`)

	project := Mix.Project(dir)
	if project == nil {
		t.Fatalf("expected Mix project, got nil")
	}

	v, err := project.GetVersion()
	if err != nil {
		t.Fatalf("get version: %v", err)
	}
	if v.String() != "1.2.3" {
		t.Fatalf("expected version 1.2.3, got %s", v)
	}

	newVersion, err := version.Parse("1.2.4")
	if err != nil {
		t.Fatalf("parse version: %v", err)
	}

	if err := project.SetVersion(newVersion); err != nil {
		t.Fatalf("set version: %v", err)
	}

	assertFileContains(t, filepath.Join(dir, "pyproject.toml"), `version = "1.2.4"`)
	assertFileContains(t, filepath.Join(dir, "package.json"), `"version":"1.2.4"`)
}

func TestMixProjectMismatchedVersions(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "pyproject.toml", `
[tool.poetry]
name = "demo"
version = "1.2.3"
`)
	writeFile(t, dir, "package.json", `{"name":"demo","version":"1.2.4"}`)

	project := Mix.Project(dir)
	if project == nil {
		t.Fatalf("expected Mix project, got nil")
	}

	if _, err := project.GetVersion(); err == nil {
		t.Fatalf("expected error for mismatched versions")
	}
}

func TestFlutterProjectVersionOperations(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "pubspec.yaml", `
name: demo
version: 1.2.3+4
`)

	project := Flutter.Project(dir)
	if project == nil {
		t.Fatalf("expected Flutter project, got nil")
	}

	v, err := project.GetVersion()
	if err != nil {
		t.Fatalf("get version: %v", err)
	}
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 || v.Build != "4" {
		t.Fatalf("unexpected version: %+v", v)
	}

	newVersion, err := version.Parse("1.2.4")
	if err != nil {
		t.Fatalf("parse version: %v", err)
	}
	newVersion.Build = "5"

	if err := project.SetVersion(newVersion); err != nil {
		t.Fatalf("set version: %v", err)
	}

	assertFileContains(t, filepath.Join(dir, "pubspec.yaml"), `version: 1.2.4+5`)
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}

func assertFileContains(t *testing.T, path, needle string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if !strings.Contains(string(data), needle) {
		t.Fatalf("expected %s to contain %q, got %s", path, needle, string(data))
	}
}
