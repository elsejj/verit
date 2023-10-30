package version

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
	Build string
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func Parse(s string) *Version {
	a := strings.Split(s, ".")
	v := Version{}
	for i, s := range a {
		switch i {
		case 0:
			v.Major = parseInt(s)
		case 1:
			v.Minor = parseInt(s)
		case 2:
			v.Patch = parseInt(s)
		case 3:
			v.Build = s
		}
	}
	return &v
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *Version) FullVersion() string {
	return fmt.Sprintf("%d.%d.%d.%s", v.Major, v.Minor, v.Patch, v.Build)
}

func (v *Version) BumpMajor() {
	v.Major++
	v.Minor = 0
	v.Patch = 0
}

func (v *Version) BumpMinor() {
	v.Minor++
	v.Patch = 0
}

func (v *Version) BumpPatch() {
	v.Patch++
}

func (v *Version) GreaterThan(v2 *Version) bool {
	if v.Major > v2.Major {
		return true
	} else if v.Major < v2.Major {
		return false
	}

	if v.Minor > v2.Minor {
		return true
	} else if v.Minor < v2.Minor {
		return false
	}

	if v.Patch > v2.Patch {
		return true
	} else if v.Patch < v2.Patch {
		return false
	}
	return false
}

func (v *Version) LessThan(v2 *Version) bool {
	return !v.GreaterThan(v2)
}

func (v *Version) Equal(v2 *Version) bool {
	return v.Major == v2.Major && v.Minor == v2.Minor && v.Patch == v2.Patch
}
