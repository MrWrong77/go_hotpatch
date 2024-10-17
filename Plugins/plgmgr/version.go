package main

import (
	"strconv"
	"strings"
)

type PluginVersion struct {
	Major   int
	Minor   int
	Patch   int
	Version string
}

func (v PluginVersion) Less(another PluginVersion) bool {
	if v.Major < another.Major {
		return true
	}
	if v.Major > another.Major {
		return false
	}
	if v.Minor < another.Minor {
		return true
	}
	if v.Minor > another.Minor {
		return false
	}
	if v.Patch < another.Patch {
		return true
	}
	if v.Patch > another.Patch {
		return false
	}
	return false
}

func (v PluginVersion) String() string {
	return v.Version
}

func ParseVersion(version string) *PluginVersion {
	ver := strings.Split(version, ".")
	if len(ver) < 3 {
		// invalid version cfg
		return nil
	}
	major, err := strconv.ParseInt(string(ver[0]), 10, 32)
	if err != nil {
		return nil
	}
	minor, err := strconv.ParseInt(string(ver[1]), 10, 32)
	if err != nil {
		return nil
	}
	patch, err := strconv.ParseInt(string(ver[2]), 10, 32)
	if err != nil {
		return nil
	}
	return &PluginVersion{
		Major:   int(major),
		Minor:   int(minor),
		Patch:   int(patch),
		Version: version,
	}
}
