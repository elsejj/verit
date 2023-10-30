package main

import (
	"flag"
	"fmt"

	pv "github.com/elsejj/verit/internal/version"
	"github.com/elsejj/verit/pkg/projectid"
	"github.com/elsejj/verit/pkg/version"
)

var flagVersion string
var flagVerbose bool
var flagBump int = 0
var flagBuild string
var flagWorkDir string
var flagHelp bool
var flagShowVersion bool

func initFlags() {
	flag.StringVar(&flagVersion, "v", "", "version to set (shorthand)")
	flag.StringVar(&flagVersion, "version", "", "version to set")

	flag.BoolVar(&flagVerbose, "vv", false, "verbose output(shorthand)")
	flag.BoolVar(&flagVerbose, "verbose", false, "verbose output")

	flag.IntVar(&flagBump, "b", 0, "bump version (shorthand)")
	flag.IntVar(&flagBump, "bump", 0, "bump version")

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
		fmt.Println(pv.Version)
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
		changed := false
		if flagBump > 0 {
			bumpVersion(p)
			changed = true
		} else if len(flagBuild) > 0 {
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
	fmt.Printf("verit (%s) - manage project version", pv.Version)
	fmt.Println("verit [options]")
	fmt.Println("options:")
	flag.PrintDefaults()
}

func bumpVersion(p projectid.Project) {
	v, err := p.GetVersion()
	if err != nil {
		if flagVerbose {
			fmt.Println("get version failed", err, "use default version '0.0.0'")
		}
		v = &version.Version{}
	}

	switch flagBump {
	case 1:
		v.BumpMajor()
	case 2:
		v.BumpMinor()
	case 3:
		v.BumpPatch()
	default:
		fmt.Println("invalid bump version type: ", flagBump)
		return
	}
	setVersion(p, v)
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
