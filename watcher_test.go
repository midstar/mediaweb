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

// assertFileCreated checks if a file is created within 10 seconds or fails.
func assertFileCreated(t *testing.T, message string, name string) {
	t.Helper()
	for i := 0; i < 100; i++ {
		if _, err := os.Stat(name); err == nil {
			return // File found
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("File %s not created. %s", name, message)
}

// assertFileRemoved check if a file is removed within 10 seconds or fails.
func assertFileRemoved(t *testing.T, message string, name string) {
	t.Helper()
	for i := 0; i < 100; i++ {
		if _, err := os.Stat(name); err != nil {
			return // File removed
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("File %s was never removed. %s", name, message)
}

func TestWatcherImages(t *testing.T) {
	mediaPath := "tmpout/TestWatcherImages"
	os.RemoveAll(mediaPath)
	os.MkdirAll(mediaPath, os.ModePerm)

	cache := "tmpcache/TestWatcherImages"
	os.RemoveAll(cache)
	os.MkdirAll(cache, os.ModePerm)

	box := rice.MustFindBox("templates")
	media := createMedia(box, mediaPath, cache, true, false, true, true)
	defer media.stopWatcher()

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
	copyFile(t, "testmedia/gif.gif", mediaPath+"/gif.gif")
	copyFile(t, "testmedia/tiff.tiff", mediaPath+"/tiff.tiff")

	// Verify that thumbnails where created
	assertFileCreated(t, "", cache+"/_icon_image.jpg")
	assertFileCreated(t, "", cache+"/_no_exif.jpg")
	assertFileCreated(t, "", cache+"/_gif.jpg")
	assertFileCreated(t, "", cache+"/_tiff.jpg")

}

func TestWatcherSubfolder(t *testing.T) {
	mediaPath := "tmpout/TestWatcherSubfolder"
	os.RemoveAll(mediaPath)
	os.MkdirAll(mediaPath, os.ModePerm)

	cache := "tmpcache/TestWatcherSubfolder"
	os.RemoveAll(cache)
	os.MkdirAll(cache, os.ModePerm)

	box := rice.MustFindBox("templates")
	media := createMedia(box, mediaPath, cache, true, false, true, true)
	defer media.stopWatcher()

	time.Sleep(100 * time.Millisecond) // Wait for watcher to start

	// Add a subdirectory with files
	os.MkdirAll(mediaPath+"/subdir", os.ModePerm)
	time.Sleep(500 * time.Millisecond) // Wait for subfolder to be watched
	copyFile(t, "templates/icon_image.png", mediaPath+"/subdir/icon_image.png")
	copyFile(t, "testmedia/exif_rotate/no_exif.jpg", mediaPath+"/subdir/no_exif.jpg")
	copyFile(t, "testmedia/gif.gif", mediaPath+"/subdir/gif.gif")
	copyFile(t, "testmedia/tiff.tiff", mediaPath+"/subdir/tiff.tiff")

	// Verify that thumbnails where created for subdirectory
	assertFileCreated(t, "", cache+"/subdir/_icon_image.jpg")
	assertFileCreated(t, "", cache+"/subdir/_no_exif.jpg")
	assertFileCreated(t, "", cache+"/subdir/_gif.jpg")
	assertFileCreated(t, "", cache+"/subdir/_tiff.jpg")

	// Add a subdirectory of the subdiretory
	os.MkdirAll(mediaPath+"/subdir/submore", os.ModePerm)
	time.Sleep(500 * time.Millisecond) // Wait for subfolder to be watched
	copyFile(t, "testmedia/exif_rotate/no_exif.jpg", mediaPath+"/subdir/submore/no_exif.jpg")
	assertFileCreated(t, "", cache+"/subdir/submore/_no_exif.jpg")
}

func TestWatcherVideo(t *testing.T) {
	mediaPath := "tmpout/TestWatcherVideo"
	os.RemoveAll(mediaPath)
	os.MkdirAll(mediaPath, os.ModePerm)

	cache := "tmpcache/TestWatcherVideo"
	os.RemoveAll(cache)
	os.MkdirAll(cache, os.ModePerm)

	box := rice.MustFindBox("templates")
	media := createMedia(box, mediaPath, cache, true, false, true, true)
	defer media.stopWatcher()

	if !media.videoThumbnailSupport() {
		t.Skip("ffmpeg not installed skipping test")
		return
	}

	time.Sleep(100 * time.Millisecond) // Wait for watcher to start

	// Add a new video file
	copyFile(t, "testmedia/video.mp4", mediaPath+"/video.mp4")

	// Verify that thumbnail was created
	assertFileCreated(t, "", cache+"/_video.jpg")
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
