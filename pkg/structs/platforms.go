package structs

type Platforms map[string]*Versions

func (p Platforms) AddPlatform(platform string) *Versions {
	// check if plattform already exists
	if v, exists := p[platform]; exists {
		return v
	}
	// create a new entry and return it
	newv := &Versions{}
	p[platform] = newv
	return newv
}

func (p Platforms) GetPlatform(platform string) *Versions {
	// return the plattform entry
	if val, exists := p[platform]; exists {
		return val
	}
	return nil
}
