package structs

type Versions map[string]*VersionEntry

func (v Versions) AddVersion(version string) *VersionEntry {
	// check if version already exists
	if ve, exists := v[version]; exists {
		return ve
	}
	// create a new entry and return it
	newv := NewVersionEntry()
	v[version] = newv
	return newv
}

func (v Versions) GetVersion(version string) *VersionEntry {
	// return the plattform entry
	if val, exists := v[version]; exists {
		return val
	}
	return nil
}
