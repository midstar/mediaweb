package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	rice "github.com/GeertJohan/go.rice"
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
	box := rice.MustFindBox("templates")
	media := createMedia(box, "testmedia", ".", true, false, false, true, false, 0)
	files, err := media.getFiles("")
	assertExpectNoErr(t, "", err)
	assertTrue(t, "No files found", len(files) > 5)
}

func TestGetFilesInvalid(t *testing.T) {
	box := rice.MustFindBox("templates")
	media := createMedia(box, "testmedia", ".", true, false, false, true, false, 0)
	files, err := media.getFiles("invalidfolder")
	assertExpectErr(t, "invalid path shall give errors", err)
	assertTrue(t, "Should not find any files", len(files) == 0)
}

func TestGetFilesHacker(t *testing.T) {
	box := rice.MustFindBox("templates")
	media := createMedia(box, "testmedia", ".", true, false, false, true, false, 0)
	files, err := media.getFiles("../..")
	assertExpectErr(t, "hacker path shall give errors", err)
	assertTrue(t, "Should not find any files", len(files) == 0)
}

func TestIsRotationNeeded(t *testing.T) {
	box := rice.MustFindBox("templates")
	media := createMedia(box, "testmedia", ".", true, false, false, true, false, 0)

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
	box := rice.MustFindBox("templates")
	media := createMedia(box, "testmedia", ".", true, false, false, true, false, 0)
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
	box := rice.MustFindBox("templates")
	media := createMedia(box, "testmedia", ".", true, false, false, true, false, 0)

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
	box := rice.MustFindBox("templates")
	media := createMedia(box, ".", ".", true, false, false, true, false, 0)
	p, err := media.getFullMediaPath("afile.jpg")
	assertExpectNoErr(t, "unable to get valid full path", err)
	assertEqualsStr(t, "invalid path", "afile.jpg", p)

	_, err = media.getFullMediaPath("../../secret_file")
	assertExpectErr(t, "hackers shall not be allowed", err)

	// Relative path
	media = createMedia(box, "arelative/path", ".", true, false, false, true, false, 0)
	p, err = media.getFullMediaPath("afile.jpg")
	assertExpectNoErr(t, "unable to get valid full path", err)
	assertEqualsStr(t, "invalid path", "arelative/path/afile.jpg", p)

	_, err = media.getFullMediaPath("../../secret_file")
	assertExpectErr(t, "hackers shall not be allowed", err)

	// Absolute path
	media = createMedia(box, "/root/absolute/path", ".", true, false, false, true, false, 0)
	p, err = media.getFullMediaPath("afile.jpg")
	assertExpectNoErr(t, "unable to get valid full path", err)
	assertEqualsStr(t, "invalid path", "/root/absolute/path/afile.jpg", p)

	_, err = media.getFullMediaPath("../../secret_file")
	assertExpectErr(t, "hackers shall not be allowed", err)
}

func TestRelativePath(t *testing.T) {
	// Root path
	box := rice.MustFindBox("templates")
	media := createMedia(box, "testmedia", ".", true, false, false, true, false, 0)

	result, err := media.getRelativePath("", "")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", ".", result)

	result, err = media.getRelativePath("", "directory")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", "directory", result)

	// Unix slashes
	result, err = media.getRelativePath("", "dir1/dir2/dir3")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", "dir1/dir2/dir3", result)

	result, err = media.getRelativePath("dir1", "dir1/dir2/dir3")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", "dir2/dir3", result)

	result, err = media.getRelativePath("dir1/", "dir1/dir2/dir3")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", "dir2/dir3", result)

	result, err = media.getRelativePath("dir1/dir2", "dir1/dir2/dir3")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", "dir3", result)

	result, err = media.getRelativePath("dir1/dir2/", "dir1/dir2/dir3")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", "dir3", result)

	result, err = media.getRelativePath("dir1/dir2/dir3", "dir1/dir2/dir3")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", ".", result)

	// Windows slashes - this will only work on windows
	if os.PathSeparator == '\\' {
		result, err = media.getRelativePath("", "dir1\\dir2\\dir3")
		assertExpectNoErr(t, "", err)
		assertEqualsStr(t, "", "dir1/dir2/dir3", result)

		result, err = media.getRelativePath("dir1\\", "dir1\\dir2\\dir3")
		assertExpectNoErr(t, "", err)
		assertEqualsStr(t, "", "dir2/dir3", result)
	}

	// Errors
	_, err = media.getRelativePath("another", "directory")
	assertExpectErr(t, "", err)

	_, err = media.getRelativePath("/a", "b")
	assertExpectErr(t, "", err)

	// getRelativeMediaPath
	result, err = media.getRelativeMediaPath("testmedia/dir1/dir2")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", "dir1/dir2", result)

	_, err = media.getRelativeMediaPath("another/dir1/dir2")
	assertExpectErr(t, "", err)
}

func TestThumbnailPath(t *testing.T) {
	box := rice.MustFindBox("templates")
	media := createMedia(box, "/c/mediapath", "/d/thumbpath", true, false, false, true, false, 0)

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
	LogTime(t, inFileName+" thumbnail generation: ")
	assertExpectNoErr(t, "", err)
	assertFileExist(t, "", outFileName)
	t.Logf("Manually check that %s thumbnail is ok", outFileName)
}

func TestGenerateImageThumbnail(t *testing.T) {
	os.MkdirAll("tmpout/TestGenerateImageThumbnail", os.ModePerm) // If already exist no problem

	box := rice.MustFindBox("templates")
	media := createMedia(box, "", "", true, false, false, true, false, 0)

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
	os.RemoveAll("tmpcache/TestWriteThumbnail")
	os.MkdirAll("tmpcache/TestWriteThumbnail", os.ModePerm)
	os.RemoveAll("tmpout/TestWriteThumbnail")
	os.MkdirAll("tmpout/TestWriteThumbnail", os.ModePerm)

	box := rice.MustFindBox("templates")
	media := createMedia(box, "testmedia", "tmpcache/TestWriteThumbnail", true, false, false, true, false, 0)

	// JPEG with embedded EXIF
	tWriteThumbnail(t, media, "jpeg.jpg", "tmpout/TestWriteThumbnail/jpeg.jpg", false)

	// JPEG without embedded EXIF
	tWriteThumbnail(t, media, "exif_rotate/no_exif.jpg", "tmpout/TestWriteThumbnail/jpeg_no_exif.jpg", false)

	// Non JPEG - no exif
	tWriteThumbnail(t, media, "png.png", "tmpout/TestWriteThumbnail/png.jpg", false)

	// Video - only if video is supported
	if media.videoThumbnailSupport() {
		tWriteThumbnail(t, media, "video.mp4", "tmpout/TestWriteThumbnail/video.jpg", false)

		// Test invalid
		tWriteThumbnail(t, media, "invalidvideo.mp4", "tmpout/TestWriteThumbnail/invalidvideo.jpg", true)
		// Check that error indication file is created
		assertFileExist(t, "", "tmpcache/TestWriteThumbnail/_invalidvideo.err.txt")
	}

	// Non existing file
	tWriteThumbnail(t, media, "dont_exist.jpg", "tmpout/TestWriteThumbnail/dont_exist.jpg", true)

	// Invalid file
	tWriteThumbnail(t, media, "invalid.jpg", "tmpout/TestWriteThumbnail/invalid.jpg", true)
	// Check that error indication file is created
	assertFileExist(t, "", "tmpcache/TestWriteThumbnail/_invalid.err.txt")
	// Generate again - just for coverage
	tWriteThumbnail(t, media, "invalid.jpg", "tmpout/TestWriteThumbnail/invalid.jpg", true)

	// Disable thumb cache
	media = createMedia(box, "testmedia", "tmpcache/TestWriteThumbnail", false, false, false, true, false, 0)

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

	box := rice.MustFindBox("templates")
	media := createMedia(box, "testmedia", ".", true, false, false, true, false, 0)

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
	box := rice.MustFindBox("templates")
	media := createMedia(box, "testmedia", ".", true, false, false, true, false, 0)
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
	err = media.generateVideoThumbnail("invalidvideo.mp4", tmp+"/invalidvideo.jpg")
	assertExpectErr(t, "", err)
}

func TestGenerateThumbnails(t *testing.T) {
	cache := "tmpcache/TestGenerateThumbnails"
	os.RemoveAll(cache)
	os.MkdirAll(cache, os.ModePerm)

	box := rice.MustFindBox("templates")
	media := createMedia(box, "testmedia", cache, true, false, false, true, false, 0)
	stat := media.generateThumbnails("", true)
	assertEqualsInt(t, "", 1, stat.NbrOfFolders)
	assertEqualsInt(t, "", 19, stat.NbrOfImages)
	assertEqualsInt(t, "", 2, stat.NbrOfVideos)
	assertEqualsInt(t, "", 10, stat.NbrOfExif)
	assertEqualsInt(t, "", 0, stat.NbrOfFailedFolders)
	assertEqualsInt(t, "", 1, stat.NbrOfFailedImages)
	if media.videoThumbnailSupport() {
		assertEqualsInt(t, "", 1, stat.NbrOfFailedVideos)
		assertFileExist(t, "", filepath.Join(cache, "_video.jpg"))
	} else {
		assertEqualsInt(t, "", 2, stat.NbrOfFailedVideos)
	}

	// Check that thumbnails where generated
	assertFileExist(t, "", filepath.Join(cache, "_png.jpg"))
	assertFileExist(t, "", filepath.Join(cache, "_gif.jpg"))
	assertFileExist(t, "", filepath.Join(cache, "_tiff.jpg"))
	assertFileExist(t, "", filepath.Join(cache, "exif_rotate", "_no_exif.jpg"))

	// Check that thumbnails where not generated for EXIF images
	assertFileNotExist(t, "", filepath.Join(cache, "_jpeg.jpg"))
	assertFileNotExist(t, "", filepath.Join(cache, "_jpeg_rotated.jpg"))
	assertFileNotExist(t, "", filepath.Join(cache, "exif_rotate", "_180deg.jpg"))
	assertFileNotExist(t, "", filepath.Join(cache, "exif_rotate", "_mirror.jpg"))

}

func TestGenerateAllThumbnails(t *testing.T) {
	cache := "tmpcache/TestGenerateAllThumbnails"
	os.RemoveAll(cache)
	os.MkdirAll(cache, os.ModePerm)

	box := rice.MustFindBox("templates")
	media := createMedia(box, "testmedia", cache, true, true, false, true, false, 0)

	for i := 0; i < 300; i++ {
		time.Sleep(100 * time.Millisecond)
		if !media.isThumbGenInProgress() {
			break
		}
	}
	assertFalse(t, "", media.isThumbGenInProgress())

	// Check that thumbnails where generated
	assertFileExist(t, "", filepath.Join(cache, "_png.jpg"))
	assertFileExist(t, "", filepath.Join(cache, "_gif.jpg"))
	assertFileExist(t, "", filepath.Join(cache, "_tiff.jpg"))
	assertFileExist(t, "", filepath.Join(cache, "exif_rotate", "_no_exif.jpg"))

	// Check that thumbnails where not generated for EXIF images
	assertFileNotExist(t, "", filepath.Join(cache, "_jpeg.jpg"))
	assertFileNotExist(t, "", filepath.Join(cache, "_jpeg_rotated.jpg"))
	assertFileNotExist(t, "", filepath.Join(cache, "exif_rotate", "_180deg.jpg"))
	assertFileNotExist(t, "", filepath.Join(cache, "exif_rotate", "_mirror.jpg"))

	if media.videoThumbnailSupport() {
		assertFileExist(t, "", filepath.Join(cache, "_video.jpg"))
	}
}

func TestGenerateNoThumbnails(t *testing.T) {
	cache := "tmpcache/TestGenerateNoThumbnails"
	os.RemoveAll(cache)
	os.MkdirAll(cache, os.ModePerm)

	box := rice.MustFindBox("templates")
	media := createMedia(box, "testmedia", cache, true, false, false, true, false, 0)

	assertFalse(t, "", media.isThumbGenInProgress())
	time.Sleep(100 * time.Millisecond)
	assertFalse(t, "", media.isThumbGenInProgress())

	// Check that no thumbnails where generated
	assertFileNotExist(t, "", filepath.Join(cache, "_png.jpg"))
	assertFileNotExist(t, "", filepath.Join(cache, "_gif.jpg"))
	assertFileNotExist(t, "", filepath.Join(cache, "_tiff.jpg"))
	assertFileNotExist(t, "", filepath.Join(cache, "_video.jpg"))
	assertFileNotExist(t, "", filepath.Join(cache, "exif_rotate", "_no_exif.jpg"))

}

func TestGetImageWidthAndHeight(t *testing.T) {
	box := rice.MustFindBox("templates")
	media := createMedia(box, "", "", true, false, false, true, false, 0)

	width, height, err := media.getImageWidthAndHeight("testmedia/jpeg.jpg")
	assertExpectNoErr(t, "", err)
	assertEqualsInt(t, "image width", 4128, width)
	assertEqualsInt(t, "image height", 2322, height)

	width, height, err = media.getImageWidthAndHeight("testmedia/gif.gif")
	assertExpectNoErr(t, "", err)
	assertEqualsInt(t, "image width", 3264, width)
	assertEqualsInt(t, "image height", 2448, height)

	width, height, err = media.getImageWidthAndHeight("testmedia/png.png")
	assertExpectNoErr(t, "", err)
	assertEqualsInt(t, "image width", 1632, width)
	assertEqualsInt(t, "image height", 1224, height)

	width, height, err = media.getImageWidthAndHeight("testmedia/tiff.tiff")
	assertExpectNoErr(t, "", err)
	assertEqualsInt(t, "image width", 979, width)
	assertEqualsInt(t, "image height", 734, height)

	// Test invalid
	_, _, err = media.getImageWidthAndHeight("testmedia/invalid.jpg")
	assertExpectErr(t, "", err)
}

func TestPreviewPath(t *testing.T) {
	box := rice.MustFindBox("templates")
	media := createMedia(box, "/c/mediapath", "/d/thumbpath", true, false, false, true, true, 1280)

	previewPath, err := media.previewPath("myimage.jpg")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", "/d/thumbpath/view_myimage.jpg", previewPath)

	previewPath, err = media.previewPath("subdrive/myimage.jpg")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", "/d/thumbpath/subdrive/view_myimage.jpg", previewPath)

	previewPath, err = media.previewPath("subdrive/myimage.png")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", "/d/thumbpath/subdrive/view_myimage.jpg", previewPath)

	_, err = media.previewPath("subdrive/myimage")
	assertExpectErr(t, "", err)

	_, err = media.previewPath("subdrive/../../hacker")
	assertExpectErr(t, "", err)
}

func tGenerateImagePreview(t *testing.T, media *Media, inFileName, outFileName string) {
	t.Helper()
	os.Remove(outFileName)
	RestartTimer()
	err := media.generateImagePreview(inFileName, outFileName)
	LogTime(t, inFileName+" preview generation: ")
	assertExpectNoErr(t, "", err)
	assertFileExist(t, "", outFileName)
	// Check dimensions
	width, height, err := media.getImageWidthAndHeight(outFileName)
	assertExpectNoErr(t, "reading dimensions", err)
	assertFalse(t, "preview width", width > media.previewMaxSide)
	assertFalse(t, "preview height", height > media.previewMaxSide)
}

func TestGenerateImagePreview(t *testing.T) {
	os.MkdirAll("tmpout/TestGenerateImagePreview", os.ModePerm) // If already exist no problem

	box := rice.MustFindBox("templates")
	media := createMedia(box, "", "", true, false, false, true, true, 1280)

	tGenerateImagePreview(t, media, "testmedia/jpeg.jpg", "tmpout/TestGenerateImagePreview/jpeg_preview.jpg")
	tGenerateImagePreview(t, media, "testmedia/jpeg_rotated.jpg", "tmpout/TestGenerateImagePreview/jpeg_rotated_preview.jpg")
	tGenerateImagePreview(t, media, "testmedia/png.png", "tmpout/TestGenerateImagePreview/png_preview.jpg")
	tGenerateImagePreview(t, media, "testmedia/gif.gif", "tmpout/TestGenerateImagePreview/gif_preview.jpg")
	tGenerateImagePreview(t, media, "testmedia/tiff.tiff", "tmpout/TestGenerateImagePreview/tiff_preview.jpg")
	tGenerateImagePreview(t, media, "testmedia/exif_rotate/no_exif.jpg", "tmpout/TestGenerateImagePreview/exif_rotate/no_exif_preview.jpg")

	// Test some invalid
	err := media.generateImagePreview("nonexisting.png", "dont_matter.png")
	assertExpectErr(t, "", err)

	err = media.generateImagePreview("testmedia/invalid.jpg", "dont_matter.jpg")
	assertExpectErr(t, "", err)
}

func tWritePreview(t *testing.T, media *Media, inFileName, outFileName string, failExpected bool) {
	t.Helper()
	os.Remove(outFileName)
	outFile, err := os.Create(outFileName)
	assertExpectNoErr(t, "unable to create out", err)
	defer outFile.Close()
	err = media.writePreview(outFile, inFileName)
	if failExpected {
		assertExpectErr(t, "should fail", err)
	} else {
		assertExpectNoErr(t, "unable to write preview", err)
		t.Logf("Manually check that %s preview is ok", outFileName)
	}
}

func TestWritePreview(t *testing.T) {
	os.RemoveAll("tmpcache/TestWritePreview")
	os.MkdirAll("tmpcache/TestWritePreview", os.ModePerm)
	os.RemoveAll("tmpout/TestWritePreview")
	os.MkdirAll("tmpout/TestWritePreview", os.ModePerm)

	box := rice.MustFindBox("templates")
	media := createMedia(box, "testmedia", "tmpcache/TestWritePreview", true, false, false, true, true, 970)

	// JPEG
	tWritePreview(t, media, "jpeg.jpg", "tmpout/TestWritePreview/jpeg.jpg", false)

	// Same file again, get cached result
	tWritePreview(t, media, "jpeg.jpg", "tmpout/TestWritePreview/jpeg.jpg", false)

	// PNG
	tWritePreview(t, media, "png.png", "tmpout/TestWritePreview/png.jpg", false)

	// TIFF
	tWritePreview(t, media, "tiff.tiff", "tmpout/TestWritePreview/tiff.tiff", false)

	// Video - should fail
	tWritePreview(t, media, "video.mp4", "tmpout/TestWritePreview/video.jpg", true)

	// Image smaller than previewMaxSide should fail
	tWritePreview(t, media, "screenshot_browser.jpg", "tmpout/TestWritePreview/screenshot_browser.jpg", true)

	// Non existing file
	tWritePreview(t, media, "dont_exist.jpg", "tmpout/TestWritePreview/dont_exist.jpg", true)

	// Invalid file
	tWritePreview(t, media, "invalid.jpg", "tmpout/TestWritePreview/invalid.jpg", true)

	// Invalid path
	tWritePreview(t, media, "../../secret.jpg", "tmpout/TestWritePreview/invalid.jpg", true)

	// Disable preview
	media = createMedia(box, "testmedia", "tmpcache/TestWritePreview", false, false, false, true, false, 0)

	// Should fail since preview is disabled now
	tWritePreview(t, media, "jpeg.jpg", "tmpout/TestWritePreview/jpeg.jpg", true)
}
