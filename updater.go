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
	generateThumbnails(relativePath string, recursive bool) *ThumbnailStatistics
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
	exit                  bool                 // Flag to indicate that updater shall exit
	media                 mediaInterface       // To run the actual thumbnail update
}

func createUpdater(media mediaInterface) *Updater {
	return &Updater{
		directories: make(map[string]time.Time),
		mutex:       sync.Mutex{},
		minTimeSinceChangeSec: 5, // Five seconds
		exit:  false,
		media: media}
}

func (u *Updater) startUpdater() {
	go u.updaterThread()
}

func (u *Updater) stopUpdater() {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.exit = true
}

func (u *Updater) exitIsSet() bool {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return u.exit
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
		time.Sleep(1 * time.Second)
		if u.exitIsSet() {
			break
		}
		path, ok := u.nextDirectoryToUpdate()
		if ok {
			u.media.generateThumbnails(path, false)
		}
	}
	llog.Info("Shutting down updater")
}
