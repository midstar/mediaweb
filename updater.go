package main

import (
	"sync"
	"time"

	"github.com/midstar/llog"
)

// mediaInterface is an interface representing the methods
// required from the Media type. The reason why we have an
// interface is to be able to mock it during testing.
type mediaInterface interface {
	generateCache(relativePath string, recursive bool) *PreCacheStatistics
	isPreCacheInProgress() bool
}

// directory represents one directory
type directory struct {
	path string // path to directory
}

// Updater represents the updater
type Updater struct {
	directories           map[string]time.Time // Key: Path, Value: Last update
	mutex                 sync.Mutex           // For thread safety
	minTimeSinceChangeSec int                  // Minimum time since change of directory before update
	stopUpdaterChan       chan bool            // Set to true to stop the updater go-routine
	done                  chan bool            // Set to true when updater go-routine has stopped
	media                 mediaInterface       // To run the actual thumbnail update
}

func createUpdater(media mediaInterface) *Updater {
	return &Updater{
		directories: make(map[string]time.Time),
		mutex:       sync.Mutex{},
		minTimeSinceChangeSec: 5, // Five seconds
		stopUpdaterChan:       make(chan bool),
		done:                  make(chan bool),
		media:                 media}
}

func (u *Updater) startUpdater() {
	llog.Info("Starting updater")
	go u.updaterThread()
}

func (u *Updater) stopUpdater() chan bool {
	u.stopUpdaterChan <- true
	return u.done
}

// stopUpdaterAndWait similar to stopUpdater but waits
// for updater go-routine to stop
func (u *Updater) stopUpdaterAndWait() {
	<-u.stopUpdater()
}

// markDirectoryAsUpdated adds the directory for update if
// it does not already exist. Otherwise it just updates the
// last updated time for the directory.
func (u *Updater) markDirectoryAsUpdated(path string) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.directories[path] = time.Now()
}

// touchDirectory is similar to markDirectoryAsUpdated but
// it is required that the directory exist, and if it does
// the last update time for the directory is updated.
func (u *Updater) touchDirectory(path string) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	_, directoryKnown := u.directories[path]
	if directoryKnown {
		u.directories[path] = time.Now()
	}
}

// nextDirectoryToUpdate returns the path of the directory
// that was updated first.
// The selected directory is removed from the directories
// map.
// If no directory fulfills the updating criteria that
// update was made= minTimeSinceChangeSec > second ago false
// with be returned.
func (u *Updater) nextDirectoryToUpdate() (string, bool) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	selectedPath := ""
	var selectedDuration time.Duration
	minTime := time.Duration(u.minTimeSinceChangeSec) * time.Second

	for path, t := range u.directories {
		duration := time.Since(t)
		if duration >= minTime && duration > selectedDuration {
			selectedPath = path
			selectedDuration = duration
		}
	}

	if selectedPath == "" {
		return "", false
	}
	delete(u.directories, selectedPath)
	return selectedPath, true
}

func (u *Updater) updaterThread() {
	for {
		select {
		case <-time.After(1 * time.Second):
			// If mediaweb is configured to update the thumbs on
			// startup we don't want to conflict with this
			if !u.media.isPreCacheInProgress() {
				path, ok := u.nextDirectoryToUpdate()
				if ok {
					llog.Info("Updating thumbs in %s", path)
					u.media.generateCache(path, false)
				}
			}
		case <-u.stopUpdaterChan:
			llog.Info("Shutting down updater")
			u.done <- true
			return
		}
	}
}
