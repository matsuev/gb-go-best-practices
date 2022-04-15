package dirscan

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

const (
	Recursive = -1
	SingleDir = 0
)

// ScanResult interface.
type ScanResult interface {
	Append(...FileListItem)
	Result() []FileListItem
}

// DirScanner struct.
type DirScanner struct {
	path    string
	depth   int
	delay   time.Duration
	wg      *sync.WaitGroup
	sigChan chan os.Signal
	sigWg   *sync.WaitGroup
	list    ScanResult
	logger  *zap.Logger
}

// Create function.
func Create(path string, depth int, delay time.Duration, logger *zap.Logger) (*DirScanner, error) {
	dirPath, err := getFullPath(path)
	if err != nil {
		logger.Error("Error in getFullPath", zap.String("path", path), zap.Error(err))

		return nil, err
	}

	logger.Debug("Create new scanner:", zap.String("path", dirPath))

	if depth < 0 {
		depth = -1
	}

	ds := &DirScanner{
		path:    dirPath,
		depth:   depth,
		delay:   delay,
		wg:      new(sync.WaitGroup),
		sigChan: make(chan os.Signal, 1),
		sigWg:   new(sync.WaitGroup),
		list:    new(FileList),
		logger:  logger,
	}

	return ds, nil
}

// Find function.
func (ds *DirScanner) Find(ctx context.Context, filter FilterFunc) error {
	dirEntry, err := os.ReadDir(ds.path)
	if err != nil {
		return err
	}

	ds.sigWg.Add(1)

	go ds.processSignals()

	ds.logger.Debug("Start scan dir:", zap.String("path", ds.path))

	for _, de := range dirEntry {
		select {
		case <-ctx.Done():
			ds.logger.Info("Terminated by user", zap.String("path", ds.path))

			return nil
		default:
			if err = ds.processDirEntry(ctx, de, filter); err != nil {
				return err
			}
		}
	}

	ds.wg.Wait()
	close(ds.sigChan)
	ds.sigWg.Wait()

	return err
}

// Result function.
func (ds *DirScanner) Result() []FileListItem {
	return ds.list.Result()
}

// getFullPath function.
func getFullPath(path string) (string, error) {
	if err := os.Chdir(path); err != nil {
		return "", err
	}

	result, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return result, nil
}

// processDirEntry function.
func (ds *DirScanner) processDirEntry(ctx context.Context, de os.DirEntry, filter FilterFunc) (err error) {
	time.Sleep(ds.delay)

	switch {
	case de.Type().IsRegular():
		ds.processFile(de, filter)
	case de.Type().IsDir() && ds.depth != 0:
		ds.wg.Add(1)

		go func() {
			defer ds.wg.Done()
			err = ds.processDirectory(ctx, de, filter)
		}()
	}

	return
}

// processFile function.
func (ds *DirScanner) processFile(de os.DirEntry, filter FilterFunc) {
	fe := NewFileEntry(de.Name(), ds.path)
	ds.logger.Debug("Process file", zap.String("filepath", fe.FullPath()))

	if filter != nil && filter(fe) {
		ds.logger.Debug("Append file to results", zap.String("filepath", fe.FullPath()))
		ds.list.Append(fe)
	}
}

// processDirectory function.
func (ds *DirScanner) processDirectory(ctx context.Context, de os.DirEntry, filter FilterFunc) error {
	dirPath := filepath.Join(ds.path, de.Name())

	dir, err := Create(dirPath, ds.depth-1, ds.delay, ds.logger)
	if err != nil {
		ds.logger.Error("Can't read directory:", zap.String("path", dirPath), zap.Error(err))

		return err
	}

	ds.logger.Debug("Process directory", zap.String("path", dirPath))

	if err = dir.Find(ctx, filter); err != nil {
		ds.logger.Error("File search error: ", zap.Error(err))
	} else {
		ds.list.Append(dir.Result()...)
	}

	return err
}

// processSignal function.
func (ds *DirScanner) processSignals() {
	ds.logger.Debug("Start signal processing", zap.String("path", ds.path))
	signal.Notify(ds.sigChan, syscall.SIGUSR1, syscall.SIGUSR2)

	for {
		sig := <-ds.sigChan
		ds.logger.Debug("Received signal:", zap.Any("sig", sig), zap.String("path", ds.path))

		switch sig {
		case syscall.SIGUSR1:
			ds.logger.Info("Current search path:", zap.String("path", ds.path))
		case syscall.SIGUSR2:
			ds.logger.Info("Increase depth level:", zap.Int("depth", ds.depth))
			ds.depth += 2
		default:
			signal.Stop(ds.sigChan)
			ds.logger.Debug("Stop signal processing", zap.String("path", ds.path))
			ds.sigWg.Done()

			return
		}
	}
}
