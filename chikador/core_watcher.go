package chikador

import (
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

// add scans through the directory for more directories than adds those directories before scanning
// for more subdirectories.
func add(path string, watcher *fsnotify.Watcher) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		fp := filepath.Join(path, file.Name())
		if file.IsDir() {
			if err := add(fp, watcher); err != nil {
				return err
			}
			continue
		}
	}
	return watcher.Add(path)
}
