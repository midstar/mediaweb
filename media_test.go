package main

import (
	"bufio"
	"bytes"
	"os"
	"testing"
	"time"
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

func TestIsRotationNeeded(t *testing.T) {
	media := createMedia("testmedia", ".", true, true)

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
	media := createMedia("testmedia", ".", true, true)
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
	t.Logf("Manually check that %s thumbnail is ok", outFileName)
}

func TestWriteEXIFThumbnail(t *testing.T) {
	os.MkdirAll("tmpout/TestWriteEXIFThumbnail", os.ModePerm) // If already exist no problem
	media := createMedia("testmedia", ".", true, true)

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
	media := createMedia(".", ".", true, true)
	p, err := media.getFullMediaPath("afile.jpg")
	assertExpectNoErr(t, "unable to get valid full path", err)
	assertEqualsStr(t, "invalid path", "afile.jpg", p)

	p, err = media.getFullMediaPath("../../secret_file")
	assertExpectErr(t, "hackers shall not be allowed", err)

	// Relative path
	media = createMedia("arelative/path", ".", true, true)
	p, err = media.getFullMediaPath("afile.jpg")
	assertExpectNoErr(t, "unable to get valid full path", err)
	assertEqualsStr(t, "invalid path", "arelative/path/afile.jpg", p)

	p, err = media.getFullMediaPath("../../secret_file")
	assertExpectErr(t, "hackers shall not be allowed", err)

	// Absolute path
	media = createMedia("/root/absolute/path", ".", true, true)
	p, err = media.getFullMediaPath("afile.jpg")
	assertExpectNoErr(t, "unable to get valid full path", err)
	assertEqualsStr(t, "invalid path", "/root/absolute/path/afile.jpg", p)

	p, err = media.getFullMediaPath("../../secret_file")
	assertExpectErr(t, "hackers shall not be allowed", err)
}

func TestThumbnailPath(t *testing.T) {
	media := createMedia("/c/mediapath", "/d/thumbpath", true, true)

	thumbPath, err := media.thumbnailPath("myimage.jpg")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", "/d/thumbpath/_myimage.jpg", thumbPath)

	thumbPath, err = media.thumbnailPath("subdrive/myimage.jpg")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", "/d/thumbpath/subdrive/_myimage.jpg", thumbPath)

	thumbPath, err = media.thumbnailPath("subdrive/myimage.png")
	assertExpectNoErr(t, "", err)
	assertEqualsStr(t, "", "/d/thumbpath/subdrive/_myimage.jpg", thumbPath)

	thumbPath, err = media.thumbnailPath("subdrive/myimage")
	assertExpectErr(t, "", err)

	thumbPath, err = media.thumbnailPath("subdrive/../../hacker")
	assertExpectErr(t, "", err)
}

func tGenerateImageThumbnail(t *testing.T, media *Media, inFileName, outFileName string) {
	os.Remove(outFileName)
	RestartTimer()
	err := media.generateImageThumbnail(inFileName, outFileName)
	LogTime(t, inFileName+"thumbnail generation: ")
	assertExpectNoErr(t, "", err)
	t.Logf("Manually check that %s thumbnail is ok", outFileName)
}

func TestGenerateImageThumbnail(t *testing.T) {
	os.MkdirAll("tmpout/TestGenerateImageThumbnail", os.ModePerm) // If already exist no problem

	media := createMedia("", "", true, true)

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

	media := createMedia("testmedia", "tmpcache/TestWriteThumbnail", true, true)

	// JPEG with embedded EXIF
	tWriteThumbnail(t, media, "jpeg.jpg", "tmpout/TestWriteThumbnail/jpeg.jpg", false)

	// JPEG without embedded EXIF
	tWriteThumbnail(t, media, "exif_rotate/no_exif.jpg", "tmpout/TestWriteThumbnail/jpeg_no_exif.jpg", false)

	// Non JPEG - no exif
	tWriteThumbnail(t, media, "png.png", "tmpout/TestWriteThumbnail/png.jpg", false)

	// Non existing file
	tWriteThumbnail(t, media, "dont_exist.jpg", "tmpout/TestWriteThumbnail/dont_exist.jpg", true)

	// Invalid file
	tWriteThumbnail(t, media, "invalid.jpg", "tmpout/TestWriteThumbnail/invalid.jpg", true)

	// Disable thumb cache
	media = createMedia("testmedia", "tmpcache/TestWriteThumbnail", false, true)

	// JPEG with embedded EXIF
	tWriteThumbnail(t, media, "jpeg.jpg", "tmpout/TestWriteThumbnail/jpeg.jpg", false)

	// Non JPEG - no exif
	tWriteThumbnail(t, media, "png.png", "tmpout/TestWriteThumbnail/png.jpg", true)
}
