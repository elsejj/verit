package main

import (
	"fmt"

	flag "github.com/spf13/pflag"

	"github.com/elsejj/verit/internal/git"
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
