package dirscan

import "path/filepath"

// FileEntry struct.
type FileEntry struct {
	name string
	path string
}

// NewFileEntry function.
func NewFileEntry(name string, path string) *FileEntry {
	fe := &FileEntry{
		name: name,
		path: path,
	}

	return fe
}

// Name function.
func (fe *FileEntry) Name() string {
	return fe.name
}

// Path function.
func (fe *FileEntry) Path() string {
	return fe.path
}

// FullPath function.
func (fe *FileEntry) FullPath() string {
	return filepath.Join(fe.path, fe.name)
}
