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
	outFileName := "testmedia/tmpout/jpeg_rotated_fixed.jpg"
	os.Mkdir("testmedia/tmpout", os.ModeDir) // If already exist no problem
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
	outFileName := "testmedia/tmpout/thumb_" + filename
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
	os.Mkdir("testmedia/tmpout", os.ModeDir) // If already exist no problem
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
	os.Mkdir("testmedia/tmpout", os.ModeDir)                            // If already exist no problem
	os.Mkdir("testmedia/tmpout/TestGenerateImageThumbnail", os.ModeDir) // If already exist no problem

	media := createMedia("", "", true, true)

	tGenerateImageThumbnail(t, media, "testmedia/jpeg.jpg", "testmedia/tmpout/TestGenerateImageThumbnail/jpeg_thumbnail.jpg")
	tGenerateImageThumbnail(t, media, "testmedia/jpeg_rotated.jpg", "testmedia/tmpout/TestGenerateImageThumbnail/jpeg_rotated_thumbnail.jpg")
	tGenerateImageThumbnail(t, media, "testmedia/png.png", "testmedia/tmpout/TestGenerateImageThumbnail/png_thumbnail.jpg")
	tGenerateImageThumbnail(t, media, "testmedia/gif.gif", "testmedia/tmpout/TestGenerateImageThumbnail/gif_thumbnail.jpg")
	tGenerateImageThumbnail(t, media, "testmedia/tiff.tiff", "testmedia/tmpout/TestGenerateImageThumbnail/tiff_thumbnail.jpg")

	// Test some invalid
	err := media.generateImageThumbnail("nonexisting.png", "dont_matter.png")
	assertExpectErr(t, "", err)

	err = media.generateImageThumbnail("testmedia/jpeg.jpg", "/invalid/path/out.jpg")
	assertExpectErr(t, "", err)

	err = media.generateImageThumbnail("testmedia/invalid.jpg", "/invalid/path/out.jpg")
	assertExpectErr(t, "", err)
}
