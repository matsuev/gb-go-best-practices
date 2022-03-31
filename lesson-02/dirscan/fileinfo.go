package dirscan

import (
	"os"
	"syscall"
	"unsafe"
)

// FileInfo struct
type FileInfo struct {
	fType uint8
	fName []byte
	fPath string
}

const (
	dotChar      = byte(46)
	zeroChar     = int8(0)
	findStep     = 4
	findStepBack = findStep - 1
)

// Create function
func NewFileInfo(dirent *syscall.Dirent, path string) (*FileInfo, error) {
	fName := parseName(dirent.Name)
	fType := dirent.Type

	if fType == syscall.DT_UNKNOWN {
		fStat, err := os.Stat(path + "/" + string(fName))
		if err != nil {
			return nil, err
		}
		if fStat.IsDir() {
			fType = syscall.DT_DIR
		}
	}

	fi := &FileInfo{fType, fName, path}

	return fi, nil
}

// String function
func (fi *FileInfo) String() string {
	return string(fi.fName)
}

// Name function
func (fi *FileInfo) Name() string {
	return fi.String()
}

func (fi *FileInfo) Path() string {
	return fi.fPath
}

// IsDir function
func (fi *FileInfo) IsDir() bool {
	if fi.fType == syscall.DT_DIR {
		nameLen := len(fi.fName)
		if (nameLen == 1 && fi.fName[0] == dotChar) || (nameLen == 2 && fi.fName[0] == dotChar && fi.fName[1] == dotChar) {
			return false
		}
		return true
	}
	return false
}

// IsFile function
func (fi *FileInfo) IsFile() bool {
	if fi.fType == syscall.DT_REG {
		return true
	}
	return false
}

// CheckExt function
func (fi *FileInfo) CheckExt(ext string) bool {
	extBytes := []byte(ext)
	extLen := len(extBytes)

	testBuf := fi.fName[len(fi.fName)-extLen:]
	for i := 0; i < extLen; i++ {
		if testBuf[i] != extBytes[i] {
			return false
		}
	}

	return true
}

// parseName function
func parseName(buf DirName) (result []byte) {
	name := (*[dsNameSize]byte)(unsafe.Pointer(&buf[0]))
	result = name[0:nameLen(buf)]
	return
}

// nameLen function
func nameLen(buf DirName) int {
	var i int
	for i = 0; i < dsNameSize; i += findStep {
		if buf[i] == zeroChar {
			break
		}
	}
	for i -= findStepBack; i < dsNameSize; i++ {
		if buf[i] == zeroChar {
			return i
		}
	}

	return dsNameSize
}
