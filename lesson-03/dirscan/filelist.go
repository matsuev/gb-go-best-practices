package dirscan

import (
	"sync"
)

// FileListItem interface.
type FileListItem interface {
	Name() string
	Path() string
}

// FileList struct.
type FileList struct {
	sync.Mutex
	list []FileListItem
}

// NewFileList function.
func NewFileList() *FileList {
	fl := new(FileList)
	fl.Clear()

	return fl
}

// Clear function.
func (fl *FileList) Clear() {
	fl.Lock()
	fl.list = make([]FileListItem, 0)
	fl.Unlock()
}

// Append function.
func (fl *FileList) Append(fi ...FileListItem) {
	fl.Lock()
	fl.list = append(fl.list, fi...)
	fl.Unlock()
}

// Result function.
func (fl *FileList) Result() (result []FileListItem) {
	fl.Lock()
	result = append(result, fl.list...)
	fl.Unlock()

	return
}
