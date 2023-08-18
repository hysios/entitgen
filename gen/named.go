package gen

import "strconv"

type Named struct {
	names    map[string]bool
	sections map[string]string
}

// init
func (named *Named) init() {
	if named.names == nil {
		named.names = make(map[string]bool)
		named.sections = make(map[string]string)
	}
}

// SuggestName returns a name that is not already in use.
func (m *Named) SuggestName(section, name string) string {
	m.init()
	if regname, ok := m.sections[section]; ok {
		return regname
	}

	if m.names[name] {
		return m.SuggestName(section, name+strconv.Itoa(len(m.names)))
	}

	m.names[name] = true
	m.sections[section] = name
	return name
}
