package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/elsejj/verit/pkg/projectid"
)

// EnsureClean verifies the working tree is clean and optionally that all commits are pushed.
func EnsureClean(dir string, checkPushed bool) error {
	if err := ensureCommitted(dir); err != nil {
		return err
	}

	if checkPushed {
		if err := ensurePushed(dir); err != nil {
			return err
		}
	}

	return nil
}

func ensureCommitted(dir string) error {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("git status failed: %s", strings.TrimSpace(string(ee.Stderr)))
		}
		return fmt.Errorf("git status failed: %w", err)
	}

	if len(bytes.TrimSpace(out)) > 0 {
		return fmt.Errorf("git working tree has uncommitted changes")
	}

	return nil
}

func ensurePushed(dir string) error {
	branchCmd := exec.Command("git", "status", "--porcelain=v2", "--branch")
	branchCmd.Dir = dir
	branchOut, err := branchCmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("git status --branch failed: %s", strings.TrimSpace(string(ee.Stderr)))
		}
		return fmt.Errorf("git status --branch failed: %w", err)
	}

	var (
		upstreamFound bool
		upstreamName  string
	)

	for _, line := range strings.Split(string(branchOut), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# branch.upstream") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				upstreamName = fields[2]
			}
			upstreamFound = true
			continue
		}

		if !upstreamFound {
			continue
		}

		if strings.HasPrefix(line, "# branch.ab") {
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}
			aheadCount, err := strconv.Atoi(strings.TrimPrefix(fields[2], "+"))
			if err != nil {
				continue
			}
			if aheadCount > 0 {
				if upstreamName == "" {
					upstreamName = "@{u}"
				}
				return fmt.Errorf("git branch is ahead of %s by %d commit(s); push your changes first", upstreamName, aheadCount)
			}
			break
		}
	}

	return nil
}

// CreateTag writes a git tag for the project's current version and optionally pushes it.
func CreateTag(p projectid.Project, push bool) (string, error) {
	v, err := p.GetVersion()
	if err != nil {
		return "", fmt.Errorf("resolve version before tagging failed: %w", err)
	}

	if err := EnsureClean(p.WorkDir(), push); err != nil {
		return "", err
	}

	tagName := "v" + v.String()
	if err := Run(p.WorkDir(), "tag", tagName); err != nil {
		return "", err
	}

	if push {
		if err := Run(p.WorkDir(), "push", "--force", "origin", tagName); err != nil {
			return "", err
		}
	}

	return tagName, nil
}

// Run executes a git command within dir and surfaces stderr on failure.
func Run(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg != "" {
			return fmt.Errorf("git %s failed: %s", strings.Join(args, " "), msg)
		}
		return fmt.Errorf("git %s failed: %w", strings.Join(args, " "), err)
	}
	return nil
}
