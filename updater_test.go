package main

import (
	"testing"
	"time"
)

type mediaMock struct {
}

var lastPathGenerated string

func (m *mediaMock) generateThumbnails(relativePath string, recursive bool) *ThumbnailStatistics {
	lastPathGenerated = relativePath
	return nil
}

func TestUpdaterMarkAndTouch(t *testing.T) {
	u := createUpdater(&mediaMock{})
	t1 := time.Now()
	time.Sleep(2 * time.Millisecond)

	// Mark
	u.markDirectoryAsUpdated("dir1")
	tDir1, hasItem := u.directories["dir1"]
	assertTrue(t, "", hasItem)
	assertTrue(t, "", tDir1.After(t1))
	time.Sleep(2 * time.Millisecond)
	t2 := time.Now()
	assertTrue(t, "", tDir1.Before(t2))
	time.Sleep(2 * time.Millisecond)

	// Touch non-updated directory
	u.touchDirectory("dir2")
	_, hasItem = u.directories["dir2"]
	assertFalse(t, "", hasItem)

	// Touch existing directory
	u.touchDirectory("dir1")
	tDir1, hasItem = u.directories["dir1"]
	assertTrue(t, "", hasItem)
	assertTrue(t, "", tDir1.After(t2))

}

func TestNextDirectoryToUpdate(t *testing.T) {
	u := createUpdater(&mediaMock{})
	u.minTimeSinceChangeSec = 1 // 5 -> 1 sec to reduce test time

	// Add some directories
	u.markDirectoryAsUpdated("dir1")
	time.Sleep(2 * time.Millisecond)
	u.markDirectoryAsUpdated("dir2")
	time.Sleep(2 * time.Millisecond)
	u.markDirectoryAsUpdated("dir3")

	// No directory pass the 1 s limit yet
	path, ok := u.nextDirectoryToUpdate()
	assertFalse(t, "", ok)
	assertEqualsStr(t, "", "", path)
	assertEqualsInt(t, "", 3, len(u.directories))

	// Wait to pass the 1 s limit
	time.Sleep(1 * time.Second)
	path, ok = u.nextDirectoryToUpdate()
	assertTrue(t, "", ok)
	assertEqualsStr(t, "", "dir1", path)
	assertEqualsInt(t, "", 2, len(u.directories))
	_, hasItem := u.directories["dir1"]
	assertFalse(t, "", hasItem)
}

func TestUpdaterThread(t *testing.T) {
	u := createUpdater(&mediaMock{})
	u.minTimeSinceChangeSec = 1 // 5 -> 1 sec to reduce test time
	lastPathGenerated = ""

	// Start the thread
	u.startUpdater()

	// A a directory
	u.markDirectoryAsUpdated("dir1")
	time.Sleep(10 * time.Millisecond)

	// Nothing should be updated now
	assertEqualsStr(t, "", "", lastPathGenerated)

	// Wait for next update
	time.Sleep(2000 * time.Millisecond)
	assertEqualsStr(t, "", "dir1", lastPathGenerated)

	u.stopUpdaterAndWait()

}
