package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/elsejj/verit/internal/changelog"
	"github.com/elsejj/verit/internal/git"
	"github.com/elsejj/verit/pkg/projectid"
	"github.com/elsejj/verit/pkg/version"

	_ "embed"
)

var flagWorkDir string
var flagSetVersion string
var flagAppVersion bool
var flagBumpMajor string
var flagBumpMinor string
var flagBumpPatch string
var flagSetBuild string
var flagSetPrerelease string
var flagHelp bool
var flagVerbose bool
var flagGitTag bool
var flagGitTagPush bool
var flagRecursive bool

//go:embed version.txt
var ver string

func initFlags() {

	flag.BoolVar(&flagVerbose, "verbose", false, "verbose output")

	flag.StringVarP(&flagBumpMajor, "major", "M", "KEEP", "bump major version, no argument value to increase current major by 1")

	flag.StringVarP(&flagBumpMinor, "minor", "m", "KEEP", "bump minor version, no argument value to increase current minor by 1")

	flag.StringVarP(&flagBumpPatch, "patch", "p", "KEEP", "bump patch version, no argument value to increase current patch by 1")

	flag.StringVarP(&flagSetBuild, "build", "b", "", "set build version")

	flag.StringVarP(&flagSetPrerelease, "prerelease", "r", "", "set prerelease version")

	flag.StringVarP(&flagWorkDir, "work-dir", "w", "", "work directory of the project, default to current directory")
	flag.BoolVarP(&flagRecursive, "recursive", "R", false, "recursively bump version")

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

	processDir(workdir, flagRecursive)

}

func showHelp() {
	fmt.Println("verit - manage project version")
	fmt.Println("version:", ver)
	fmt.Println("usage: verit [options]")
	fmt.Println("options:")
	flag.PrintDefaults()
}

// well known dependencies dir
var wellknownDependencies = []string{
	"node_modules",
	".venv",
}

func isWellknownDependency(workdir string) bool {
	for _, d := range wellknownDependencies {
		if strings.Contains(workdir, d) {
			return true
		}
	}
	return false
}

func processDir(workdir string, recursive bool) {
	id := projectid.Which(workdir)

	p := id.Project(workdir)

	if p == nil {
		if flagVerbose {
			fmt.Println("unsupported project in", workdir)
		}
		return
	}

	if len(flagSetVersion) > 0 {
		v, err := version.Parse(flagSetVersion)
		if err != nil {
			fmt.Println("invalid version:", err)
			return
		}
		setVersion(p, v)
	} else {
		changed := bumpVersion(p)
		if changed {
			if flagVerbose {
				fmt.Println("version changed")
			}
		}
	}
	if recursive {
		filepath.WalkDir(workdir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() || isWellknownDependency(path) || path == workdir {
				return nil
			}
			processDir(path, false)
			return nil
		})
	}

	// only do git operations in root project
	if flagGitTag {
		v, err := p.GetVersion()
		if err != nil {
			fmt.Println(err)
			return
		}
		if !changelog.EnsureUpdated(workdir, v.String()) {
			fmt.Println("changelog not updated for version", v)
			return
		}
		tagName, err := git.CreateTag(p, flagGitTagPush)
		if err != nil {
			fmt.Println(err)
			return
		}
		if flagVerbose {
			if flagGitTagPush {
				fmt.Printf("created and pushed tag '%s'\n", tagName)
			} else {
				fmt.Printf("created tag '%s'\n", tagName)
			}
		}
	}
}

func bumpVersion(p projectid.Project) bool {
	v, err := p.GetVersion()
	if err != nil {
		if flagVerbose {
			fmt.Println("get version failed", err)
		}
		return false
	}

	changed := false

	major, err := version.ParseVersionNumber(flagBumpMajor)
	if err != nil {
		fmt.Println("invalid major version:", err)
		return false
	}
	changed = changed || (major != version.KEEP)

	minor, err := version.ParseVersionNumber(flagBumpMinor)
	if err != nil {
		fmt.Println("invalid minor version:", err)
		return false
	}
	changed = changed || (minor != version.KEEP)

	patch, err := version.ParseVersionNumber(flagBumpPatch)
	if err != nil {
		fmt.Println("invalid patch version:", err)
		return false
	}
	changed = changed || (patch != version.KEEP)

	if len(flagSetPrerelease) > 0 {
		changed = true
	}

	if len(flagSetBuild) > 0 {
		changed = true
	}

	if !changed {
		if flagVerbose {
			fmt.Println("no version change")
		}
		return false
	}

	oldVersion := v.String()
	v.BumpMajor(major)
	v.BumpMinor(minor)
	v.BumpPatch(patch)

	v.Prerelease = flagSetPrerelease
	v.Build = flagSetBuild

	setVersion(p, v)
	newVersion := v.String()

	fmt.Println(p.ID().String(), "project in", p.WorkDir(), "bumped from", oldVersion, "to", newVersion)

	return true
}

func setVersion(p projectid.Project, v *version.Version) {
	err := p.SetVersion(v)
	if err != nil {
		if flagVerbose {
			fmt.Println(err)

		}
		return
	}
	if flagVerbose {
		fmt.Printf("'%s' project in '%s' set to version '%s'\n", p.ID(), p.WorkDir(), v)
	}
}
