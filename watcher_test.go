package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/fsnotify/fsnotify"
)

func copyFile(t *testing.T, sourceFile, destinationFile string) {
	t.Helper()
	input, err := ioutil.ReadFile(sourceFile)
	assertExpectNoErr(t, "", err)

	err = ioutil.WriteFile(destinationFile, input, 0644)
	assertExpectNoErr(t, "", err)
}

// assertFileCreated checks if a file is created within 50 seconds or fails.
func assertFileCreated(t *testing.T, message string, name string) {
	t.Helper()
	for i := 0; i < 500; i++ {
		if _, err := os.Stat(name); err == nil {
			return // File found
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("File %s not created. %s", name, message)
}

// assertFileRemoved check if a file is removed within 50 seconds or fails.
func assertFileRemoved(t *testing.T, message string, name string) {
	t.Helper()
	for i := 0; i < 500; i++ {
		if _, err := os.Stat(name); err != nil {
			return // File removed
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("File %s was never removed. %s", name, message)
}

func TestStartWatcher(t *testing.T) {
	mediaPath := "tmpout/TestStartWatcher"
	os.RemoveAll(mediaPath)
	os.MkdirAll(mediaPath, os.ModePerm)

	cache := "tmpcache/TestStartWatcher"
	os.RemoveAll(cache)
	os.MkdirAll(cache, os.ModePerm)

	box := rice.MustFindBox("templates")
	media := createMedia(box, mediaPath, cache, true, false, true, true)

	time.Sleep(100 * time.Millisecond) // Wait for watcher to start

	// Add a new file
	copyFile(t, "templates/icon_image.png", mediaPath+"/icon_image.png")

	// Verify that thumbnail was created
	assertFileCreated(t, "", cache+"/_icon_image.jpg")

	// Remove file
	os.Remove(mediaPath + "/icon_image.png")

	// Verify that thumbnail was removed
	assertFileRemoved(t, "", cache+"/_icon_image.jpg")

	// Add many files
	copyFile(t, "templates/icon_image.png", mediaPath+"/icon_image.png")
	copyFile(t, "testmedia/exif_rotate/no_exif.jpg", mediaPath+"/no_exif.jpg")
	copyFile(t, "testmedia/video.mp4", mediaPath+"/video.mp4")
	copyFile(t, "testmedia/gif.gif", mediaPath+"/gif.gif")
	copyFile(t, "testmedia/tiff.tiff", mediaPath+"/tiff.tiff")

	// Verify that thumbnails where created
	assertFileCreated(t, "", cache+"/_icon_image.jpg")
	assertFileCreated(t, "", cache+"/_no_exif.jpg")
	assertFileCreated(t, "", cache+"/_video.jpg")
	assertFileCreated(t, "", cache+"/_gif.jpg")
	assertFileCreated(t, "", cache+"/_tiff.jpg")

	media.stopWatcher()
}
func TestWatchFolder(t *testing.T) {
	box := rice.MustFindBox("templates")
	// Don't start the watcher, so that we can test its internal
	// functionality
	media := createMedia(box, "testmedia", ".", true, false, false, true)

	watcher, err := fsnotify.NewWatcher()

	// Test some valid
	err = media.watchFolder(watcher, "testmedia")
	assertExpectNoErr(t, "", err)
	err = media.watchFolder(watcher, "templates")
	assertExpectNoErr(t, "", err)

	// Test some invalid
	err = media.watchFolder(watcher, "dontexist")
	assertExpectErr(t, "", err)
	err = media.watchFolder(watcher, "testmedia/dontexit")
	assertExpectErr(t, "", err)
	err = media.watchFolder(watcher, "testmedia/jpeg.jpg")
	assertExpectErr(t, "", err)
}
