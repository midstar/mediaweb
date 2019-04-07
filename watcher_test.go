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
	// TBD

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
