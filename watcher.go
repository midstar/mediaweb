package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/midstar/llog"
)

// Watcher represents the watcher type
type Watcher struct {
	media           *Media
	updater         *Updater
	stopWatcherChan chan bool // Set to true to stop the watcher go-routine
	done            chan bool // Set to true when watcher go-routine has stopped
}

func createWatcher(media *Media, thumbnails, preview bool) *Watcher {
	return &Watcher{
		media:           media,
		updater:         createUpdater(media, thumbnails, preview),
		stopWatcherChan: make(chan bool),
		done:            make(chan bool)}
}

// stopWatcher stops the media watcher go-routine if it is running.
// It is perfectly ok to call this function even if the watcher is
// not running.
// Also stops the updater go-routine.
// Returns the Watcher done channel and the Updater done channel
func (w *Watcher) stopWatcher() (chan bool, chan bool) {
	updaterDone := w.updater.stopUpdater()
	w.stopWatcherChan <- true
	return w.done, updaterDone
}

// stopWatcherAndWait similar to stopUpdater but waits
// for the watcher and updater go-routines to stop
func (w *Watcher) stopWatcherAndWait() {
	watcherDone, updaterDone := w.stopWatcher()
	<-watcherDone
	<-updaterDone
}

// startWatcher identifies all folders within the mediaPath including
// subfolders and starts the folder watcher go routine.
func (w *Watcher) startWatcher() {
	llog.Info("Starting media watcher")
	w.updater.startUpdater()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		llog.Error("Unable to watch for new media since: %s", err)
		return
	}

	go w.mediaWatcher(watcher)

	w.watchFolder(watcher, w.media.mediaPath)
}

// watchFolder with watch the provided folder including its
// sub folders (i.e. recursively).
// The error return value is just for test purposes.
func (w *Watcher) watchFolder(watcher *fsnotify.Watcher, path string) error {
	llog.Debug("Watching folder: %s", path)
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
		if fileInfo.IsDir() || fileInfo.Mode()&os.ModeSymlink != 0 {
			w.watchFolder(watcher, filepath.Join(path, fileInfo.Name()))
		}
	}
	return nil
}

// mediaWatcher contains the loop that watches the file events.
// Call stopWatcher to exit.
//
// Note that we ignore rename and delete events, i.e. there
// is no clean up.
func (w *Watcher) mediaWatcher(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if ok {
				llog.Debug("Watcher event: %s", event)
				path := event.Name
				// relativeMediaPath is always the last diretory, never a file
				// (because we call getDir)
				relativeMediaPath, err := w.media.getRelativeMediaPath(getDir(path))
				if err == nil {
					if event.Op&fsnotify.Create == fsnotify.Create {
						if isDir(path) {
							// This is an new diretory
							// Watch it
							w.watchFolder(watcher, path)
						}
						// Mark the directory as changed so that updater eventually
						// will create the thumbnails
						w.updater.markDirectoryAsUpdated(relativeMediaPath)
					} else if (event.Op&fsnotify.Remove == fsnotify.Remove) ||
						(event.Op&fsnotify.Rename == fsnotify.Rename) {
						// Files has been removed, renamed or moved
						// Mark the directory as changed so that updater eventually
						// will create the thumbnails
						w.updater.markDirectoryAsUpdated(relativeMediaPath)
					} else if event.Op&fsnotify.Write == fsnotify.Write {
						// Tell updater that there is operations performed in the
						// directory (i.e. wait for a while before generating the
						// thumbnails)
						w.updater.touchDirectory(relativeMediaPath)
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if ok {
				llog.Warn("Watcher error: %s", err)
			}
		case <-w.stopWatcherChan:
			llog.Info("Shutting down media watcher")
			watcher.Close()
			w.done <- true
			return
		}
	}
}

// isDir return true if the path is a directory
func isDir(path string) bool {
	_, err := ioutil.ReadDir(path)
	return err == nil
}

// getDir removes the file (if such exist) from a path.
// Always returns path with front slash separator
func getDir(path string) string {
	var result string
	if isDir(path) {
		result = path
	} else {
		result = filepath.Dir(path)
	}
	return filepath.ToSlash(result)
}
