package filesystem

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type Multistat struct {
	followSymlinks bool
}

func NewMultistat(followSymlinks bool) *Multistat {
	return &Multistat{
		followSymlinks: followSymlinks,
	}
}

func (r *Multistat) Run(ctx context.Context, path string) (<-chan *FileInfo, error) {
	info, err := stat(path, r.followSymlinks)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, errors.New("Path must be a directory")
	}
	wg := new(sync.WaitGroup)
	collector := make(chan *FileInfo)
	wg.Add(1)
	go r.scrape(ctx, wg, collector, path)
	go func() {
		wg.Wait()
		close(collector)
	}()
	return collector, nil
}

func (r *Multistat) scrape(ctx context.Context, wg *sync.WaitGroup, collector chan<- *FileInfo, path string) {
	defer wg.Done()
	entries, err := os.ReadDir(path)
	if err != nil {
		collector <- &FileInfo{err: err}
		return
	}
	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		if r.canceled(ctx) {
			return
		}
		if entry.IsDir() {
			wg.Add(1)
			go r.scrape(ctx, wg, collector, entryPath)
		} else {
			info, _ := stat(entryPath, r.followSymlinks)
			collector <- info
		}
	}
}

func (r *Multistat) canceled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
	}
	return false
}
