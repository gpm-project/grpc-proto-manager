package repo

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var semanticMatcher = regexp.MustCompile("^v[0-9]+.[0-9]+.[0-9]+$")

// Version object to capture the uploaded version on the repo.
type Version struct {
	// Major version.
	Major int
	// Minor version.
	Minor int
	// Patch version.
	Patch int
}

// EmptyVersion returns an empty version.
func EmptyVersion() *Version {
	return &Version{}
}

// FromTag builds a Version from a given tag on the repo.
func FromTag(repoTag string) (*Version, error) {
	repoTag = strings.TrimRight(repoTag, "\n")

	// Check if the version matches the regex
	if !semanticMatcher.MatchString(repoTag) {
		return nil, fmt.Errorf("version %s does not match %s", repoTag, semanticMatcher.String())
	}

	noVersion := strings.ReplaceAll(repoTag, "v", "")
	split := strings.Split(noVersion, ".")

	major, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, err
	}
	minor, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, err
	}
	patch, err := strconv.Atoi(split[2])
	if err != nil {
		return nil, err
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

// IncrementMinor the minor version.
func (v *Version) IncrementMinor() {
	v.Minor++
}

// String representation of this version.
func (v *Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}
