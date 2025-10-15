package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Special constants indicating should keep current version
const KEEP = -2

// / Special constants indicating should increase current version by 1
const INCREASE = -1

func ParseVersionNumber(s string) (int, error) {
	switch strings.ToUpper(s) {
	case "KEEP":
		return KEEP, nil
	case "INC":
		return INCREASE, nil
	default:
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("invalid version number: %s", s)
		}
		if n < 0 {
			return 0, fmt.Errorf("version number must be non-negative: %s", s)
		}
		return n, nil
	}
}

// see https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string and  https://regex101.com/r/Ly7O1x/3/
var versionRe = regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

type Version struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
	Build      string
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func Parse(s string) (*Version, error) {
	matches := versionRe.FindStringSubmatch(s)
	if matches == nil {
		return nil, fmt.Errorf("invalid version string: %s", s)
	}

	v := Version{
		Major:      parseInt(matches[1]),
		Minor:      parseInt(matches[2]),
		Patch:      parseInt(matches[3]),
		Prerelease: matches[4],
		Build:      matches[5],
	}
	return &v, nil
}

func (v *Version) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch))
	if v.Prerelease != "" {
		b.WriteString("-")
		b.WriteString(v.Prerelease)
	}
	if v.Build != "" {
		b.WriteString("+")
		b.WriteString(v.Build)
	}
	return b.String()
}

func (v *Version) BumpMajor(ver int) {
	if ver == KEEP {
		return
	}
	if ver == INCREASE {
		v.Major++
	} else if ver >= 0 {
		v.Major = ver
	}
	v.Minor = 0
	v.Patch = 0
}

func (v *Version) BumpMinor(ver int) {
	if ver == KEEP {
		return
	}
	if ver == INCREASE {
		v.Minor++
	} else if ver >= 0 {
		v.Minor = ver
	}
	v.Patch = 0
}

func (v *Version) BumpPatch(ver int) {
	if ver == KEEP {
		return
	}
	if ver == INCREASE {
		v.Patch++
	} else if ver >= 0 {
		v.Patch = ver
	}
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
	if c := strings.Compare(v.Prerelease, v2.Prerelease); c > 0 {
		return true
	} else if c < 0 {
		return false
	}
	if c := strings.Compare(v.Build, v2.Build); c > 0 {
		return true
	} else if c < 0 {
		return false
	}
	return false
}

func (v *Version) LessThan(v2 *Version) bool {
	return v2.GreaterThan(v)
}

func (v *Version) Equal(v2 *Version) bool {
	return v.Major == v2.Major && v.Minor == v2.Minor && v.Patch == v2.Patch && v.Prerelease == v2.Prerelease && v.Build == v2.Build
}
