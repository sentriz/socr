package directories

// Directories maps an alias to its filesystem path.
type Directories map[string]string

func (d Directories) AliasByPath(path string) (string, bool) {
	for k, v := range d {
		if v == path {
			return k, true
		}
	}
	return "", false
}
