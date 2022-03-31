package dirscan

import (
	"os"
	"syscall"
)

// ExtFilter function
func ExtFilter(ext ...string) (f FilterFunc) {
	return func(fi *FileInfo) bool {
		for _, e := range ext {
			if fi.CheckExt(e) {
				return true
			}
		}
		return false
	}
}

// getWorkDir function
func getWorkDir(path string) (workDir string, err error) {
	if err = os.Chdir(path); err != nil {
		return
	}

	if workDir, err = os.Getwd(); err != nil {
		return
	}

	return
}

// getFileDescriptor function
func getFileDescriptor(path string) (fd int, err error) {

	if fd, err = syscall.Open(path, dsOpenFlags, dsOpenMask); err != nil {
		return
	}

	if fd, err = syscall.Openat(fd, path, dsOpenMask, dsOpenFlags); err != nil {
		return
	}

	return
}
