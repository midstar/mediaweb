package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/midstar/llog"
)

// stopWatcher stops the media watcher go-routine if it is running.
// It is perfectly ok to call this function even if the watcher is
// not running.
func (m Media) stopWatcher() {
	m.stopWatcherChan <- true
}

// startWatcher identifies all folders within the mediaPath including
// subfolders and starts the folder watcher go routine.
func (m *Media) startWatcher() {
	llog.Info("Starting media watcher")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		llog.Error("Unable to watch for new media since: %s", err)
		return
	}

	go m.mediaWatcher(watcher)

	m.watchFolder(watcher, m.mediaPath) // TODO put back
}

// watchFolder with watch the provided folder including its
// sub folders (i.e. recursively).
// The error return value is just for test purposes.
func (m *Media) watchFolder(watcher *fsnotify.Watcher, path string) error {
	err := watcher.Add(path)
	if err != nil {
		llog.Error("Watch folder %s error: %s", path, err)
	}
	// Go through its subfolders and watch these
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		errMsg := fmt.Sprintf("Reading folder %s error: %s", path, err)
		llog.Error(errMsg)
		return fmt.Errorf(errMsg)
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			m.watchFolder(watcher, filepath.Join(path, fileInfo.Name()))
		}
	}
	return nil
}

// mediaWatcher contains the loop that watches the file events.
// Call stopWatcher to exit.
//
// Limitation 1:
// The Write event is fired of several reasons and there might
// be multiple Write events for one operation. Therefore we
// ignore this event. Basically this means that if the media
// is modified we won't detect it.
//
// Limitation 2:
// If a subfolder is created we probably won't detect the first
// file(s) created in that folder, since the create file events
// of the first files will be generated after we have been able
// to watch the event.
//
// Limitation 3:
// If a directory is removed we will not remove the thumbnails
// generated in that directory
func (m *Media) mediaWatcher(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if ok {
				llog.Info("Watcher event: %s", event) // TODO change to llog.Debug
				path, err := m.getRelativeMediaPath(event.Name)
				if err == nil {
					if m.isImage(event.Name) || m.isVideo(event.Name) {
						if event.Op&fsnotify.Rename == fsnotify.Rename ||
							event.Op&fsnotify.Remove == fsnotify.Remove {
							// Remove thumbnail if it exist
							thumbPath, err := m.thumbnailPath(path)
							if err == nil {
								llog.Info("Removing thumbnail if it exist: %s", thumbPath)
								os.Remove(thumbPath)
							}
						} else if event.Op&fsnotify.Create == fsnotify.Create {
							// Create thumbnail
							time.Sleep(500 * time.Millisecond)
							m.generateThumbnail(path)
						}
					} else if event.Op&fsnotify.Create == fsnotify.Create {
						// Check if it was a diretory that was created
						if _, err := ioutil.ReadDir(event.Name); err == nil {
							llog.Info("Watching new folder %s", event.Name)
							m.watchFolder(watcher, event.Name)
						}
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if ok {
				llog.Error("Watcher error: %s", err)
			}
		case <-m.stopWatcherChan:
			llog.Info("Shutting down media watcher")
			watcher.Close()
			return
		}
	}
}
