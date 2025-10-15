package changelog

import (
	"os"
	"path/filepath"
	"strings"
)

// EnsureUpdated checks whether the changelog in workDir contains a section for the
// provided version. If CHANGELOG.md is missing, the function reports success.
func EnsureUpdated(workDir, version string) bool {
	version = strings.TrimSpace(version)
	if version == "" {
		return true
	}

	changelogPath := filepath.Join(workDir, "CHANGELOG.md")

	data, err := os.ReadFile(changelogPath)
	if err != nil {
		if os.IsNotExist(err) {
			return true
		}
		return false
	}

	target := newVersionTarget(version)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if isVersionLine(line, target) {
			return true
		}
	}

	return false
}

type versionTarget struct {
	exact         string
	withPrefix    string
	withoutPrefix string
}

func newVersionTarget(version string) versionTarget {
	version = strings.TrimSpace(version)
	if version == "" {
		return versionTarget{}
	}

	without := version
	if len(without) > 0 && (without[0] == 'v' || without[0] == 'V') {
		without = without[1:]
	}

	with := without
	if len(with) > 0 {
		with = "v" + without
	}

	return versionTarget{
		exact:         version,
		withPrefix:    with,
		withoutPrefix: without,
	}
}

func isVersionLine(line string, target versionTarget) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return false
	}
	if !strings.HasPrefix(trimmed, "#") {
		return false
	}

	header := strings.TrimLeft(trimmed, "#")
	header = strings.TrimSpace(header)
	if header == "" {
		return false
	}

	tokens := strings.Fields(header)
	for _, token := range tokens {
		clean := normalizeToken(token)
		if matchesTarget(clean, target) {
			return true
		}
	}

	return false
}

func normalizeToken(token string) string {
	clean := strings.TrimSpace(token)
	clean = strings.Trim(clean, "[](){}<>`*_")
	for _, suffix := range []string{":", ",", ";", ".", "!", "?"} {
		clean = strings.TrimSuffix(clean, suffix)
	}
	return clean
}

func matchesTarget(token string, target versionTarget) bool {
	if token == "" {
		return false
	}

	if token == target.exact || token == target.withPrefix || token == target.withoutPrefix {
		return true
	}

	trimmed := trimVersionPrefix(token)
	return trimmed == target.withoutPrefix
}

func trimVersionPrefix(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 'v' || s[0] == 'V' {
		return s[1:]
	}
	return s
}
