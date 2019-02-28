package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	packr "github.com/gobuffalo/packr/v2"
)

type timerType struct {
	start int64
	stop  int64
}

// For benchmark timing tests
var timer timerType

func RestartTimer() {
	timer.start = time.Now().UnixNano()
}

func LogTime(t *testing.T, whatWasMeasured string) {
	timer.stop = time.Now().UnixNano()
	deltaMs := (timer.stop - timer.start) / int64(time.Millisecond)
	t.Logf("%s %d ms", whatWasMeasured, deltaMs)
}

func TestGetFiles(t *testing.T) {
	box := packr.New("templates", "./templates")
	media := createMedia(box, "testmedia", ".", true, true)
	files, err := media.getFiles("")
	assertExpectNoErr(t, "", err)
	assertTrue(t, "No files found", len(files) > 5)
}

func TestGetFilesInvalid(t *testing.T) {
	box := packr.New("templates", "./templates")
	media := createMedia(box, "testmedia", ".", true, true)
	files, err := media.getFiles("invalidfolder")
	assertExpectErr(t, "invalid path shall give errors", err)
	assertTrue(t, "Should not find any files", len(files) == 0)
}

func TestGetFilesHacker(t *testing.T) {
	box := packr.New("templates", "./templates")
	media := createMedia(box, "testmedia", ".", true, true)
	files, err := media.getFiles("../..")
	assertExpectErr(t, "hacker path shall give errors", err)
	assertTrue(t, "Should not find any files", len(files) == 0)
}

func TestIsRotationNeeded(t *testing.T) {
	box := packr.New("templates", "./templates")
	media := createMedia(box, "testmedia", ".", true, true)

	rotationNeeded := media.isRotationNeeded("exif_rotate/180deg.jpg")
	assertTrue(t, "Rotation should be needed", rotationNeeded)

	rotationNeeded = media.isRotationNeeded("exif_rotate/mirror.jpg")
	assertTrue(t, "Rotation should be needed", rotationNeeded)

	rotationNeeded = media.isRotationNeeded("exif_rotate/mirror_rotate_90deg_cw.jpg")
	assertTrue(t, "Rotation should be needed", rotationNeeded)

	rotationNeeded = media.isRotationNeeded("exif_rotate/mirror_rotate_270deg.jpg")
	assertTrue(t, "Rotation should be needed", rotationNeeded)

	rotationNeeded = media.isRotationNeeded("exif_rotate/mirror_vertical.jpg")
	assertTrue(t, "Rotation should be needed", rotationNeeded)

	rotationNeeded = media.isRotationNeeded("exif_rotate/rotate_270deg_cw.jpg")
	assertTrue(t, "Rotation should be needed", rotationNeeded)

	rotationNeeded = media.isRotationNeeded("exif_rotate/rotate_90deg_cw.jpg")
	assertTrue(t, "Rotation should be needed", rotationNeeded)

	rotationNeeded = media.isRotationNeeded("exif_rotate/normal.jpg")
	assertFalse(t, "Rotation should not be needed", rotationNeeded)

	rotationNeeded = media.isRotationNeeded("exif_rotate/no_exif.jpg")
	assertFalse(t, "Rotation should not be needed", rotationNeeded)

	rotationNeeded = media.isRotationNeeded("non_existing.jpg")
	assertFalse(t, "Rotation should not be needed", rotationNeeded)

	rotationNeeded = media.isRotationNeeded("png.png")
	assertFalse(t, "Rotation should not be needed", rotationNeeded)

	rotationNeeded = media.isRotationNeeded("../../../hackerpath/secret.jpg")
	assertFalse(t, "Rotation should not be needed", rotationNeeded)

	// Turn of rotation
	media.autoRotate = false

	rotationNeeded = media.isRotationNeeded("exif_rotate/mirror_rotate_90deg_cw.jpg")
	assertFalse(t, "Rotation should not be needed when turned off", rotationNeeded)
}

func TestRotateAndWrite(t *testing.T) {
	outFileName := "tmpout/TestRotateAndWrite/jpeg_rotated_fixed.jpg"
	os.MkdirAll("tmpout/TestRotateAndWrite", os.ModePerm) // If already exist no problem
	os.Remove(outFileName)
	box := packr.New("templates", "./templates")
	media := createMedia(box, "testmedia", ".", true, true)
	outFile, err := os.Create(outFileName)
	assertExpectNoErr(t, "unable to create out", err)
	defer outFile.Close()
	RestartTimer()
	err = media.rotateAndWrite(outFile, "jpeg_rotated.jpg")
	LogTime(t, "rotate JPG")
	assertExpectNoErr(t, "unable to rotate out", err)
	t.Logf("Manually check that %s has been rotated correctly", outFileName)
}

func tEXIFThumbnail(t *testing.T, media *Media, filename string) {
	t.Helper()
	inFileName := "exif_rotate/" + filename
	outFileName := "tmpout/TestWriteEXIFThumbnail/thumb_" + filename
	os.Remove(outFileName)
	outFile, err := os.Create(outFileName)
	assertExpectNoErr(t, "unable to create out", err)
	defer outFile.Close()
	RestartTimer()
	err = media.writeEXIFThumbnail(outFile, inFileName)
	LogTime(t, inFileName+" thumbnail time")
	assertExpectNoErr(t, "unable to extract thumbnail", err)
	assertFileExist(t, "", outFileName)
	t.Logf("Manually check that %s thumbnail is ok", outFileName)
}

func TestWriteEXIFThumbnail(t *testing.T) {
	os.MkdirAll("tmpout/TestWriteEXIFThumbnail", os.ModePerm) // If already exist no problem
	box := packr.New("templates", "./templates")
	media := createMedia(box, "testmedia", ".", true, true)

	tEXIFThumbnail(t, media, "normal.jpg")
	tEXIFThumbnail(t, media, "180deg.jpg")
	tEXIFThumbnail(t, media, "mirror.jpg")
	tEXIFThumbnail(t, media, "mirror_rotate_90deg_cw.jpg")
	tEXIFThumbnail(t, media, "mirror_rotate_270deg.jpg")
	tEXIFThumbnail(t, media, "mirror_vertical.jpg")
	tEXIFThumbnail(t, media, "rotate_270deg_cw.jpg")
	tEXIFThumbnail(t, media, "rotate_90deg_cw.jpg")

	// Test some invalid
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	err := media.writeEXIFThumbnail(writer, "../../../hackerpath/secret.jpg")
	assertExpectErr(t, "hacker attack shall not be allowed", err)

	err = media.writeEXIFThumbnail(writer, "no_exif.jpg")
	assertExpectErr(t, "No EXIF shall not have thumbnail", err)
}

func TestFullPath(t *testing.T) {
	// Root path
	box := packr.New("templates", "./templates")
	media := createMedia(box, ".", ".", true, true)
	p, err := media.getFullMediaPath("afile.jpg")
	assertExpectNoErr(t, "unable to get valid full path", err)
	assertEqualsStr(t, "invalid path", "afile.jpg", p)

	_, err = media.getFullMediaPath("../../secret_file")
	assertExpectErr(t, "hackers shall not be allowed", err)

	// Relative path
	media = createMedia(box, "arelative/path", ".", true, true)
	p, err = media.getFullMediaPath("afile.jpg")
	assertExpectNoErr(t, "unable to get valid full path", err)
	assertEqualsStr(t, "invalid path", "arelative/path/afile.jpg", p)

	_, err = media.getFullMediaPath("../../secret_file")
	assertExpectErr(t, "hackers shall not be allowed", err)

	// Absolute path
	media = createMedia(box, "/root/absolute/path", ".", true, true)
	p, err = media.getFullMediaPath("afile.jpg")
	assertExpectNoErr(t, "unable to get valid full path", err)
	assertEqualsStr(t, "invalid path", "/root/absolute/path/afile.jpg", p)

	_, err = media.getFullMediaPath("../../secret_file")
	assertExpectErr(t, "hackers shall not be allowed", err)
}

func TestThumbnailPath(t *testing.T) {
	box := packr.New("templates", "./templates")
	media := createMedia(box, "/c/mediapath", "/d/thumbpath", true, true)

	thumbPath, err := media.thumbnailPath("myimage.jpg")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", "/d/thumbpath/_myimage.jpg", thumbPath)

	thumbPath, err = media.thumbnailPath("subdrive/myimage.jpg")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", "/d/thumbpath/subdrive/_myimage.jpg", thumbPath)

	thumbPath, err = media.thumbnailPath("subdrive/myimage.png")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", "/d/thumbpath/subdrive/_myimage.jpg", thumbPath)

	_, err = media.thumbnailPath("subdrive/myimage")
	assertExpectErr(t, "", err)

	_, err = media.thumbnailPath("subdrive/../../hacker")
	assertExpectErr(t, "", err)
}

func tGenerateImageThumbnail(t *testing.T, media *Media, inFileName, outFileName string) {
	t.Helper()
	os.Remove(outFileName)
	RestartTimer()
	err := media.generateImageThumbnail(inFileName, outFileName)
	LogTime(t, inFileName+"thumbnail generation: ")
	assertExpectNoErr(t, "", err)
	assertFileExist(t, "", outFileName)
	t.Logf("Manually check that %s thumbnail is ok", outFileName)
}

func TestGenerateImageThumbnail(t *testing.T) {
	os.MkdirAll("tmpout/TestGenerateImageThumbnail", os.ModePerm) // If already exist no problem

	box := packr.New("templates", "./templates")
	media := createMedia(box, "", "", true, true)

	tGenerateImageThumbnail(t, media, "testmedia/jpeg.jpg", "tmpout/TestGenerateImageThumbnail/jpeg_thumbnail.jpg")
	tGenerateImageThumbnail(t, media, "testmedia/jpeg_rotated.jpg", "tmpout/TestGenerateImageThumbnail/jpeg_rotated_thumbnail.jpg")
	tGenerateImageThumbnail(t, media, "testmedia/png.png", "tmpout/TestGenerateImageThumbnail/png_thumbnail.jpg")
	tGenerateImageThumbnail(t, media, "testmedia/gif.gif", "tmpout/TestGenerateImageThumbnail/gif_thumbnail.jpg")
	tGenerateImageThumbnail(t, media, "testmedia/tiff.tiff", "tmpout/TestGenerateImageThumbnail/tiff_thumbnail.jpg")
	tGenerateImageThumbnail(t, media, "testmedia/exif_rotate/no_exif.jpg", "tmpout/TestGenerateImageThumbnail/exif_rotate/no_exif.jpg")

	// Test some invalid
	err := media.generateImageThumbnail("nonexisting.png", "dont_matter.png")
	assertExpectErr(t, "", err)

	err = media.generateImageThumbnail("testmedia/invalid.jpg", "dont_matter.jpg")
	assertExpectErr(t, "", err)
}

func tWriteThumbnail(t *testing.T, media *Media, inFileName, outFileName string, failExpected bool) {
	t.Helper()
	os.Remove(outFileName)
	outFile, err := os.Create(outFileName)
	assertExpectNoErr(t, "unable to create out", err)
	defer outFile.Close()
	err = media.writeThumbnail(outFile, inFileName)
	if failExpected {
		assertExpectErr(t, "should fail", err)
	} else {
		assertExpectNoErr(t, "unable to write thumbnail", err)
		t.Logf("Manually check that %s thumbnail is ok", outFileName)
	}
}

func TestWriteThumbnail(t *testing.T) {
	os.MkdirAll("tmpcache/TestWriteThumbnail", os.ModePerm) // If already exist no problem
	os.RemoveAll("tmpcache/TestWriteThumbnail/*")
	os.MkdirAll("tmpout/TestWriteThumbnail", os.ModePerm) // If already exist no problem
	os.RemoveAll("tmpout/TestWriteThumbnail/*")

	box := packr.New("templates", "./templates")
	media := createMedia(box, "testmedia", "tmpcache/TestWriteThumbnail", true, true)

	// JPEG with embedded EXIF
	tWriteThumbnail(t, media, "jpeg.jpg", "tmpout/TestWriteThumbnail/jpeg.jpg", false)

	// JPEG without embedded EXIF
	tWriteThumbnail(t, media, "exif_rotate/no_exif.jpg", "tmpout/TestWriteThumbnail/jpeg_no_exif.jpg", false)

	// Non JPEG - no exif
	tWriteThumbnail(t, media, "png.png", "tmpout/TestWriteThumbnail/png.jpg", false)

	// Video - only if video is supported
	if media.videoThumbnailSupport() {
		tWriteThumbnail(t, media, "video.mp4", "tmpout/TestWriteThumbnail/video.jpg", false)
	}

	// Non existing file
	tWriteThumbnail(t, media, "dont_exist.jpg", "tmpout/TestWriteThumbnail/dont_exist.jpg", true)

	// Invalid file
	tWriteThumbnail(t, media, "invalid.jpg", "tmpout/TestWriteThumbnail/invalid.jpg", true)

	// Disable thumb cache
	media = createMedia(box, "testmedia", "tmpcache/TestWriteThumbnail", false, true)

	// JPEG with embedded EXIF
	tWriteThumbnail(t, media, "jpeg.jpg", "tmpout/TestWriteThumbnail/jpeg.jpg", false)

	// Non JPEG - no exif
	tWriteThumbnail(t, media, "png.png", "tmpout/TestWriteThumbnail/png.jpg", true)
}

func TestVideoThumbnailSupport(t *testing.T) {
	// Since we cannot guarantee that ffmpeg is available on the test
	// host we will replace the ffmpeg command temporary
	origCmd := ffmpegCmd
	defer func() {
		ffmpegCmd = origCmd
	}()

	box := packr.New("templates", "./templates")
	media := createMedia(box, "", "", true, true)

	t.Logf("ffmpeg supported: %v", media.videoThumbnailSupport())

	ffmpegCmd = "thiscommanddontexit"
	assertFalse(t, ffmpegCmd, media.videoThumbnailSupport())

	ffmpegCmd = "cmd"
	shallBeTrueOnWindows := media.videoThumbnailSupport()

	ffmpegCmd = "echo"
	shallBeTrueOnNonWindows := media.videoThumbnailSupport()

	assertTrue(t, "Shall be true on at least one platform", shallBeTrueOnWindows || shallBeTrueOnNonWindows)
}

func tGenerateVideoThumbnail(t *testing.T, media *Media, inFileName, outFileName string) {
	t.Helper()
	os.Remove(outFileName)
	RestartTimer()
	err := media.generateVideoThumbnail(inFileName, outFileName)
	LogTime(t, inFileName+"thumbnail generation: ")
	assertExpectNoErr(t, "", err)
	assertFileExist(t, "", outFileName)
	t.Logf("Manually check that %s thumbnail is ok", outFileName)
}

func TestGenerateVideoThumbnail(t *testing.T) {
	box := packr.New("templates", "./templates")
	media := createMedia(box, "", "", true, true)
	if !media.videoThumbnailSupport() {
		t.Skip("ffmpeg not installed skipping test")
		return
	}
	tmp := "tmpout/TestGenerateVideoThumbnail"
	os.MkdirAll(tmp, os.ModePerm) // If already exist no problem
	tmpSpace := "tmpout/TestGenerateVideoThumbnail/with space in path"
	os.MkdirAll(tmpSpace, os.ModePerm) // If already exist no problem

	tGenerateVideoThumbnail(t, media, "testmedia/video.mp4", tmp+"/video_thumbnail.jpg")
	tGenerateVideoThumbnail(t, media, "testmedia/video.mp4", tmpSpace+"/video_thumbnail.jpg")

	// Test some invalid
	err := media.generateVideoThumbnail("nonexisting.mp4", tmp+"dont_matter.jpg")
	assertExpectErr(t, "", err)
	err = media.generateVideoThumbnail("png.png", tmp+"dont_matter_png.jpg")
	assertExpectErr(t, "", err)
}

func TestGenerateThumbnails(t *testing.T) {
	cache := "tmpcache/TestGenerateThumbnails"
	os.RemoveAll(cache)
	os.MkdirAll(cache, os.ModePerm)

	box := packr.New("templates", "./templates")
	media := createMedia(box, "testmedia", cache, true, true)
	media.generateThumbnails("")

	// Check that thumbnails where generated
	assertFileExist(t, "", filepath.Join(cache, "_screenshot_browser.jpg"))
}
