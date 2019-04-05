package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/midstar/llog"
)

// stopWatcher stops the folder watcher go-routine if it is running.
// It is perfectly ok to call this function even if the watcher is
// not running.
func (m Media) stopWatcher() {
	m.stopWatcherChan <- true
}

// startWatcher identifies all folders within the mediaPath including
// subfolders and starts the folder watcher go routine.
func (m *Media) startWatcher() {
	llog.Info("Starting folder watcher")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		llog.Panic("%s", err)
	}

	go m.folderWatcher(watcher)

	err = watcher.Add("tmpout")
	if err != nil {
		llog.Panic("%s", err)
	}
}

func (m *Media) folderWatcher(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			llog.Info("event: %s", event)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			llog.Info("error: %s", err)
		case <-m.stopWatcherChan:
			llog.Info("Folder watcher stopped intentionally")
			return
		}
	}
}
