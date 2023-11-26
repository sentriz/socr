package directories

// alias -> path
type Directories map[string]string

func (d Directories) PathByAlias(alias string) (string, bool) {
	alias, ok := d[alias]
	return alias, ok
}

func (d Directories) AliasByPath(path string) (string, bool) {
	for k, v := range d {
		if v == path {
			return k, true
		}
	}
	return "", false
}
