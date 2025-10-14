package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/elsejj/verit/pkg/projectid"
	"github.com/elsejj/verit/pkg/version"

	_ "embed"
)

var flagBuild string
var flagWorkDir string
var flagSetVersion string
var flagAppVersion bool
var flagBumpMajor string
var flagBumpMinor string
var flagBumpPatch string
var flagHelp bool
var flagVerbose bool
var flagGitTag bool
var flagGitTagPush bool

//go:embed version.txt
var ver string

func initFlags() {

	flag.BoolVar(&flagVerbose, "verbose", false, "verbose output")

	flag.StringVarP(&flagBumpMajor, "major", "M", "KEEP", "bump major version, no argument value to increase current major by 1")

	flag.StringVarP(&flagBumpMinor, "minor", "m", "KEEP", "bump minor version, no argument value to increase current minor by 1")

	flag.StringVarP(&flagBumpPatch, "patch", "p", "KEEP", "bump patch version, no argument value to increase current patch by 1")

	flag.StringVarP(&flagBuild, "build", "b", "", "set build version")

	flag.StringVarP(&flagWorkDir, "work-dir", "w", "", "work directory of the project, default to current directory")

	flag.BoolVarP(&flagHelp, "help", "h", false, "show help (shorthand)")

	flag.StringVarP(&flagSetVersion, "version", "v", "", "version to set, like 1.2.3")
	flag.BoolVarP(&flagAppVersion, "app-version", "V", false, "show app version")
	flag.BoolVarP(&flagGitTag, "tag", "t", false, "create git tag using current version")
	flag.BoolVarP(&flagGitTagPush, "tag-push", "T", false, "create git tag and push it with --force")

	flag.Lookup("major").NoOptDefVal = "INC"
	flag.Lookup("minor").NoOptDefVal = "INC"
	flag.Lookup("patch").NoOptDefVal = "INC"

	flag.Parse()

	if flagGitTagPush {
		flagGitTag = true
	}
}

func main() {

	initFlags()

	if flagHelp {
		showHelp()
		return
	}

	if flagAppVersion {
		fmt.Printf("v%s\n", ver)
		return
	}

	workdir := projectid.Pwd()

	if len(flagWorkDir) > 0 {
		workdir = flagWorkDir
	}

	id := projectid.Which(workdir)

	p := id.Project(workdir)

	if len(flagSetVersion) > 0 {
		v, err := version.Parse(flagSetVersion)
		if err != nil {
			fmt.Println("invalid version:", err)
			return
		}
		setVersion(p, v)
	} else {
		changed := bumpVersion(p)
		if len(flagBuild) > 0 {
			changed = true
			setBuildVersion(p)
		}
		if changed {
			if flagVerbose {
				fmt.Println("version changed")
			}
		}
	}
	if flagGitTag {
		if err := createGitTag(p); err != nil {
			fmt.Println(err)
			return
		}
	}

	showVersion(p)
}

func showHelp() {
	fmt.Println("verit - manage project version")
	fmt.Println("version:", ver)
	fmt.Println("usage: verit [options]")
	fmt.Println("options:")
	flag.PrintDefaults()
}

func bumpVersion(p projectid.Project) bool {
	v, err := p.GetVersion()
	if err != nil {
		if flagVerbose {
			fmt.Println("get version failed", err, "use default version '0.0.0'")
		}
		v = &version.Version{}
	}

	major, err := version.ParseVersionNumber(flagBumpMajor)
	if err != nil {
		fmt.Println("invalid major version:", err)
		return false
	}

	minor, err := version.ParseVersionNumber(flagBumpMinor)
	if err != nil {
		fmt.Println("invalid minor version:", err)
		return false
	}

	patch, err := version.ParseVersionNumber(flagBumpPatch)
	if err != nil {
		fmt.Println("invalid patch version:", err)
		return false
	}

	if major == version.KEEP && minor == version.KEEP && patch == version.KEEP {
		if flagVerbose {
			fmt.Println("no version change")
		}
		return false
	}

	v.BumpMajor(major)
	v.BumpMinor(minor)
	v.BumpPatch(patch)

	setVersion(p, v)

	return true
}

func setBuildVersion(p projectid.Project) {
	v, err := p.GetVersion()
	if err != nil {
		if flagVerbose {
			fmt.Println("get version failed", err, "use default version '0.0.0'")
		}
		v = &version.Version{}
	}
	v.Build = flagBuild
	setVersion(p, v)
}

func setVersion(p projectid.Project, v *version.Version) {
	err := p.SetVersion(v)
	if err != nil {
		fmt.Println(err)
		return
	}
	if flagVerbose {
		fmt.Printf("'%s' project in '%s' set to version '%s'\n", p.ID(), p.WorkDir(), v)
	}
}

func createGitTag(p projectid.Project) error {
	v, err := p.GetVersion()
	if err != nil {
		return fmt.Errorf("resolve version before tagging failed: %w", err)
	}

	if err := ensureCleanGit(p.WorkDir()); err != nil {
		return err
	}

	tagName := "v" + v.String()
	if err := runGit(p.WorkDir(), "tag", tagName); err != nil {
		return err
	}

	if flagGitTagPush {
		if err := runGit(p.WorkDir(), "push", "--force", "origin", tagName); err != nil {
			return err
		}
	}

	if flagVerbose {
		if flagGitTagPush {
			fmt.Printf("created and pushed tag '%s'\n", tagName)
		} else {
			fmt.Printf("created tag '%s'\n", tagName)
		}
	}
	return nil
}

func ensureCleanGit(dir string) error {
	err := gitCheckModified(dir)
	if err != nil {
		return err
	}

	if flagGitTagPush {
		err = gitCheckPushed(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

func gitCheckPushed(dir string) error {
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

func gitCheckModified(dir string) error {
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

func runGit(dir string, args ...string) error {
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

func showVersion(p projectid.Project) {
	v, err := p.GetVersion()
	if err != nil {
		fmt.Println(err)
		return
	}
	if flagVerbose {
		fmt.Printf("'%s' project in '%s' version is '%s'\n", p.ID(), p.WorkDir(), v)
	} else {
		fmt.Println(v)
	}
}
