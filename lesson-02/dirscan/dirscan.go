package dirscan

import (
	"context"
	"log"
	"syscall"
	"unsafe"
)

const (
	dsOpenFlags = syscall.O_RDONLY | syscall.O_NONBLOCK
	dsOpenMask  = 0
	dsBufSize   = 4096
	dsNameSize  = 256
)

// DirScanner struct
type DirScanner struct{}

// DirName type
type DirName [dsNameSize]int8

// FilterFunc type
type FilterFunc func(*FileInfo) bool

// Create function
func Create() *DirScanner {
	return new(DirScanner)
}

// FindFiles function
func (ds *DirScanner) FindFiles(ctx context.Context, path string, depth int, filter FilterFunc) (result []FileInfo, err error) {
	if path, err = getWorkDir(path); err != nil {
		return
	}

	fd, err := getFileDescriptor(path)
	if err != nil {
		return
	}

	scanBuffer := make([]byte, dsBufSize)
	var n int

	for {
		if n, err = syscall.ReadDirent(fd, scanBuffer); err != nil || n <= 0 {
			return
		}

		if result, err = ds.parseDirent(ctx, scanBuffer[0:n], path, depth, filter); err != nil {
			return
		}
	}
}

// parseDirent function
func (ds *DirScanner) parseDirent(ctx context.Context, buf []byte, path string, depth int, filter FilterFunc) (result []FileInfo, err error) {
	for len(buf) > 0 {
		select {
		case <-ctx.Done():
			return
		default:
			dirent := (*syscall.Dirent)(unsafe.Pointer(&buf[0]))

			if fi, ferr := NewFileInfo(dirent, path); ferr != nil {
				log.Println(err)
			} else {
				switch {
				case fi.IsFile():
					if filter == nil || filter(fi) {
						result = append(result, *fi)
					}
				case fi.IsDir():
					if depth != 0 {
						if children, cerr := ds.FindFiles(ctx, path+"/"+fi.Name(), depth-1, filter); cerr != nil {
							log.Println(err)
						} else {
							result = append(result, children...)
						}
					}
				}
			}
			buf = buf[dirent.Reclen:]
		}
	}
	return
}
