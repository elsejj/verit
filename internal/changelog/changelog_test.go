package changelog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureUpdated(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		version     string
		changelog   string
		createFile  bool
		wantUpdated bool
	}{
		{
			name:        "missing changelog",
			version:     "1.0.0",
			createFile:  false,
			wantUpdated: true,
		},
		{
			name:        "empty version",
			version:     "",
			createFile:  false,
			wantUpdated: true,
		},
		{
			name:    "version section with brackets",
			version: "0.2.0",
			changelog: `# Changelog

## [0.2.0] - 2025-10-14

### Added

- Something nice
`,
			createFile:  true,
			wantUpdated: true,
		},
		{
			name:    "version section with leading v",
			version: "1.3.0",
			changelog: `# Changelog

## v1.3.0
`,
			createFile:  true,
			wantUpdated: true,
		},
		{
			name:    "missing version section",
			version: "0.9.0",
			changelog: `# Changelog

## [Unreleased]

No release info.
`,
			createFile:  true,
			wantUpdated: false,
		},
		{
			name:    "version mentioned outside heading",
			version: "2.0.0",
			changelog: `# Changelog

## [Unreleased]

See version 2.0.0 for more.
`,
			createFile:  true,
			wantUpdated: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()

			if tc.createFile {
				path := filepath.Join(dir, "CHANGELOG.md")
				if err := os.WriteFile(path, []byte(tc.changelog), 0o644); err != nil {
					t.Fatalf("failed to write changelog: %v", err)
				}
			}

			if got := EnsureUpdated(dir, tc.version); got != tc.wantUpdated {
				t.Fatalf("EnsureUpdated(%q, %q) = %v, want %v", dir, tc.version, got, tc.wantUpdated)
			}
		})
	}
}
