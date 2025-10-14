package main

import (
	"flag"
	"fmt"

	"github.com/elsejj/verit/pkg/projectid"
	"github.com/elsejj/verit/pkg/version"

	_ "embed"
)

var flagVersion string
var flagVerbose bool
var flagBuild string
var flagWorkDir string
var flagHelp bool
var flagShowVersion bool
var flagBumpMajor = version.KEEP
var flagBumpMinor = version.KEEP
var flagBumpPatch = version.KEEP

//go:embed version.txt
var ver string

func initFlags() {
	flag.StringVar(&flagVersion, "v", "", "version to set (shorthand)")
	flag.StringVar(&flagVersion, "version", "", "version to set")

	flag.BoolVar(&flagVerbose, "vv", false, "verbose output(shorthand)")
	flag.BoolVar(&flagVerbose, "verbose", false, "verbose output")

	flag.IntVar(&flagBumpMajor, "M", version.INCREASE, "bump major version, default -1 for increase current major by 1 (shorthand)")
	flag.IntVar(&flagBumpMajor, "major", version.INCREASE, "bump major version")

	flag.IntVar(&flagBumpMinor, "m", version.INCREASE, "bump minor version, default -1 for increase current minor by 1 (shorthand)")
	flag.IntVar(&flagBumpMinor, "minor", version.INCREASE, "bump minor version")

	flag.IntVar(&flagBumpPatch, "p", version.INCREASE, "bump patch version, default -1 for increase current patch by 1 (shorthand)")
	flag.IntVar(&flagBumpPatch, "patch", version.INCREASE, "bump patch version")

	flag.StringVar(&flagBuild, "b", "", "set build version (shorthand)")
	flag.StringVar(&flagBuild, "build", "", "set build version")

	flag.StringVar(&flagWorkDir, "d", "", "work directory (shorthand)")
	flag.StringVar(&flagWorkDir, "dir", "", "work directory")

	flag.BoolVar(&flagHelp, "h", false, "show help (shorthand)")

	flag.BoolVar(&flagShowVersion, "V", false, "show version (shorthand)")

	flag.Parse()
}

func main() {

	initFlags()

	if flagShowVersion {
		fmt.Println(ver)
		return
	}

	if flagHelp {
		showHelp()
		return
	}

	workdir := projectid.Pwd()

	if len(flagWorkDir) > 0 {
		workdir = flagWorkDir
	}

	id := projectid.Which(workdir)

	p := id.Project(workdir)

	if len(flagVersion) > 0 {
		v := version.Parse(flagVersion)
		setVersion(p, v)
	} else {
		changed := bumpVersion(p)
		if len(flagBuild) > 0 {
			changed = true
			setBuildVersion(p)
		}
		if changed {
			fmt.Println("version changed")
		}
		showVersion(p)
	}
}

func showHelp() {
	fmt.Printf("verit (%s) - manage project version", ver)
	fmt.Println("verit [options]")
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

	v.BumpMajor(flagBumpMajor)
	v.BumpMinor(flagBumpMinor)
	v.BumpPatch(flagBumpPatch)

	setVersion(p, v)

	return flagBumpMajor != version.KEEP || flagBumpMinor != version.KEEP || flagBumpPatch != version.KEEP
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
