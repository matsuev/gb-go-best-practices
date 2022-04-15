package dirscan

import "path/filepath"

// FilterFunc type.
type FilterFunc func(FileListItem) bool

// ExtFilter function.
func ExtFilter(ext ...string) (f FilterFunc) {
	return func(fi FileListItem) bool {
		for _, e := range ext {
			if filepath.Ext(fi.Name()) == e {
				return true
			}
		}

		return false
	}
}
