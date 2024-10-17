package plgmgr

type PluginVersion struct {
	Major int
	Minor int
	Patch int
}

func (v PluginVersion) Less(another PluginVersion) bool {
	if v.Major < v.Major {
		return true
	}
	if v.Major > v.Major {
		return false
	}
	if v.Minor < v.Minor {
		return true
	}
	if v.Minor > v.Minor {
		return false
	}
	if v.Patch < v.Patch {
		return true
	}
	if v.Patch > v.Patch {
		return false
	}
	return false
}

func ParseVersion(ver []byte) *PluginVersion {
	if len(ver) < 3 {
		// invalid version cfg
		return nil
	}
	// major, err := strconv.ParseInt(string(ver[0]), 10, 32)
	// if err != nil {
	// 	return nil
	// }
	// minor, err := strconv.ParseInt(string(ver[1]), 10, 32)
	// if err != nil {
	// 	return nil
	// }
	// patch, err := strconv.ParseInt(string(ver[2]), 10, 32)
	// if err != nil {
	// 	return nil
	// }
	// return &PluginVersion{
	// 	Major: int(major),
	// 	Minor: int(minor),
	// 	Patch: int(patch),
	// }
	return &PluginVersion{
		Major: 1,
		Minor: 0,
		Patch: 0,
	}
}
