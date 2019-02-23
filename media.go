package main

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/midstar/llog"
	"github.com/rwcarlsen/goexif/exif"
)

var imgExtensions = [...]string{".png", ".jpg", ".jpeg", ".tif", ".tiff", ".gif"}
var vidExtensions = [...]string{".avi", ".mov", ".vid", ".mkv", ".mp4"}

// Media represents the media including its base path
type Media struct {
	mediaPath        string // Top level path for media files
	thumbPath        string // Top level path for thumbnails
	enableThumbCache bool   // Generate thumbnails
	autoRotate       bool   // Rotate JPEG files when needed
}

// File represents a folder or any other file
type File struct {
	Type string // folder, image or video
	Name string
	Path string // Including Name. Always using / (even on Windows)
}

// createMedia creates a new media. If thumb cache is enabled the path is
// created when needed.
func createMedia(mediaPath string, thumbPath string, enableThumbCache bool, autoRotate bool) *Media {
	llog.Info("Media path: %s", mediaPath)
	if enableThumbCache {
		directory := filepath.Dir(thumbPath)
		err := os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			llog.Error("Unable to create thumbnail cache path %s. Reason: %s", thumbPath, err)
			llog.Info("Thumbnail cache will be disabled")
			enableThumbCache = false
		} else {
			llog.Info("Thumbnail cache path: %s", thumbPath)
		}
	} else {
		llog.Info("Thumbnail cache disabled")
	}
	llog.Info("JPEG auto rotate: %t", autoRotate)
	return &Media{mediaPath: filepath.ToSlash(filepath.Clean(mediaPath)),
		thumbPath:        filepath.ToSlash(filepath.Clean(thumbPath)),
		enableThumbCache: enableThumbCache,
		autoRotate:       autoRotate}
}

// getFullPath returns the full path from an absolute base
// path and a relative path. Returns error on security hacks,
// i.e. when someone tries to access ../../../ for example to
// get files that are not within configured base path.
//
// Always returning front slashes / as path separator
func (m *Media) getFullPath(basePath string, relativePath string) (string, error) {
	fullPath := filepath.ToSlash(filepath.Join(basePath, relativePath))
	diffPath, err := filepath.Rel(basePath, fullPath)
	diffPath = filepath.ToSlash(diffPath)
	if err != nil || strings.HasPrefix(diffPath, "../") {
		return m.mediaPath, fmt.Errorf("Hacker attack. Someone tries to access: %s", fullPath)
	}
	return fullPath, nil
}

// getFullMediaPath returns the full path of the provided path, i.e:
// media path + relative path.
func (m *Media) getFullMediaPath(relativePath string) (string, error) {
	return m.getFullPath(m.mediaPath, relativePath)
}

// getFullThumbPath returns the full path of the provided path, i.e:
// thumb path + relative path.
func (m *Media) getFullThumbPath(relativePath string) (string, error) {
	return m.getFullPath(m.thumbPath, relativePath)
}

// getFiles returns a slice of File's sorted on file name
func (m *Media) getFiles(relativePath string) ([]File, error) {
	//var files []File
	files := make([]File, 0, 500)
	fullPath, err := m.getFullMediaPath(relativePath)
	if err != nil {
		return files, err
	}
	fileInfo, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return files, err
	}

	for _, fileInfo := range fileInfo {
		fileType := ""
		if fileInfo.IsDir() {
			fileType = "folder"
		} else {
			fileType = m.getFileType(fileInfo.Name())
		}
		// Only add directories, videos and images
		if fileType != "" {
			// Use path with / slash
			pathOriginal := filepath.Join(relativePath, fileInfo.Name())
			pathNew := filepath.ToSlash(pathOriginal)

			file := File{
				Type: fileType,
				Name: fileInfo.Name(),
				Path: pathNew}
			files = append(files, file)
		} else {
			llog.Debug("getFiles - omitting: %s", fileInfo.Name())
		}
	}
	return files, nil
}

// getFileType returns "video" for video files and "image" for image files.
// For all other files (including folders) "" is returned.
// relativeFileName can also include an absolute or relative path.
func (m *Media) getFileType(relativeFileName string) string {
	extension := filepath.Ext(relativeFileName)

	// Check if this is an image
	for _, imgExtension := range imgExtensions {
		if strings.EqualFold(extension, imgExtension) {
			return "image"
		}
	}

	// Check if this is a video
	for _, vidExtension := range vidExtensions {
		if strings.EqualFold(extension, vidExtension) {
			return "video"
		}
	}

	return "" // Not a video or an image
}

func (m *Media) extractEXIF(fullFilePath string) *exif.Exif {
	if !m.isJPEG(fullFilePath) {
		return nil // Only JPEG has EXIF
	}
	efile, err := os.Open(fullFilePath)
	if err != nil {
		llog.Warn("Could not open file for EXIF decoding. File: %s reason: %s\n", fullFilePath, err)
		return nil
	}
	defer efile.Close()
	ex, err := exif.Decode(efile)
	if err != nil {
		llog.Debug("No EXIF. file %s reason: %s\n", fullFilePath, err)
		return nil
	}
	return ex
}

func (m *Media) isJPEG(pathAndFile string) bool {
	extension := filepath.Ext(pathAndFile)
	if strings.EqualFold(extension, ".jpg") == false &&
		strings.EqualFold(extension, ".jpeg") == false {
		return false
	}
	return true
}

// isRotationNeeded returns true if the file needs to be rotated.
// It finds this out by reading the EXIF rotation information
// in the file.
// If Media.autoRotate is false this function will always return
// false.
func (m *Media) isRotationNeeded(relativeFilePath string) bool {
	if m.autoRotate == false {
		return false
	}
	fullPath, err := m.getFullMediaPath(relativeFilePath)
	if err != nil {
		llog.Info("Unable to get full media path for %s\n", relativeFilePath)
		return false
	}
	ex := m.extractEXIF(fullPath)
	if ex == nil {
		return false // No EXIF info exist
	}
	orientTag, _ := ex.Get(exif.Orientation)
	if orientTag == nil {
		return false // No Orientation
	}
	orientInt, _ := orientTag.Int(0)
	if orientInt > 1 && orientInt < 9 {
		return true // Rotation is needed
	}
	return false
}

// rotateAndWrite opens and rotates a JPG/JPEG file according to
// EXIF rotation information. Then it writes the rotated image
// to the io.Writer. NOTE! This process requires Decoding and
// encoding of the image which takes a LOT of time (2-3 sec).
// Check if image needs rotation with isRotationNeeded first.
func (m *Media) rotateAndWrite(w io.Writer, relativeFilePath string) error {
	fullPath, err := m.getFullMediaPath(relativeFilePath)
	if err != nil {
		return err
	}

	img, err := imaging.Open(fullPath, imaging.AutoOrientation(true))
	if err != nil {
		return err
	}
	err = imaging.Encode(w, img, imaging.JPEG)
	if err != nil {
		return err
	}
	return nil
}

// writeEXIFThumbnail extracts the EXIF thumbnail from a JPEG file
// and rotates it when needed (based on the EXIF orientation tag).
// Returns err if no thumbnail exist.
func (m *Media) writeEXIFThumbnail(w io.Writer, relativeFilePath string) error {
	fullPath, err := m.getFullMediaPath(relativeFilePath)
	if err != nil {
		return err
	}
	ex := m.extractEXIF(fullPath)
	if ex == nil {
		return fmt.Errorf("No EXIF info for %s", fullPath)
	}
	thumbBytes, err := ex.JpegThumbnail()
	if err != nil {
		return fmt.Errorf("No EXIF thumbnail for %s", fullPath)
	}
	orientTag, _ := ex.Get(exif.Orientation)
	if orientTag == nil {
		// No Orientation assume no rotation needed
		w.Write(thumbBytes)
		return nil
	}
	orientInt, _ := orientTag.Int(0)
	if orientInt > 1 && orientInt < 9 {
		// Rotation is needed
		img, err := imaging.Decode(bytes.NewReader(thumbBytes))
		if err != nil {
			llog.Warn("Unable to decode EXIF thumbnail for %s", fullPath)
			w.Write(thumbBytes)
			return nil
		}
		var outImg *image.NRGBA
		switch orientInt {
		case 2:
			outImg = imaging.FlipV(img)
		case 3:
			outImg = imaging.Rotate180(img)
		case 4:
			outImg = imaging.Rotate180(imaging.FlipV(img))
		case 5:
			outImg = imaging.Rotate270(imaging.FlipV(img))
		case 6:
			outImg = imaging.Rotate270(img)
		case 7:
			outImg = imaging.Rotate90(imaging.FlipV(img))
		case 8:
			outImg = imaging.Rotate90(img)
		}
		imaging.Encode(w, outImg, imaging.JPEG)
	} else {
		// No rotation is needed
		w.Write(thumbBytes)
	}
	return nil
}

// thumbnailPath returns the absolute thumbnail file path from a
// media path. Thumbnails are always stored in JPEG format (.jpg
// extension) and starts with '_'.
// Returns error if the media path is invalid.
func (m *Media) thumbnailPath(relativeMediaPath string) (string, error) {
	path, file := filepath.Split(relativeMediaPath)
	if !m.isJPEG(file) {
		// Replace extension with .jpg
		ext := filepath.Ext(file)
		if ext == "" {
			return "", fmt.Errorf("File has no extension: %s", file)
		}
		file = strings.Replace(file, ext, ".jpg", -1)
	}
	file = "_" + file
	relativeThumbnailPath := filepath.Join(path, file)
	return m.getFullThumbPath(relativeThumbnailPath)
}

// generateImageThumbnail generates a thumbnail from any of the supported
// images. Will create necessary subdirectories in the thumbpath.
func (m *Media) generateImageThumbnail(fullMediaPath, fullThumbPath string) error {
	img, err := imaging.Open(fullMediaPath, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("Unable to open image %s. Reason: %s", fullMediaPath, err)
	}
	thumbImg := imaging.Thumbnail(img, 256, 256, imaging.Lanczos)

	// Create subdirectories if needed
	directory := filepath.Dir(fullThumbPath)
	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Unable to create directories in %s for creating thumbnail. Reason %s", fullThumbPath, err)
	}

	// Write thumbnail to file
	outFile, err := os.Create(fullThumbPath)
	if err != nil {
		return fmt.Errorf("Unable to open %s for creating thumbnail. Reason %s", fullThumbPath, err)
	}
	defer outFile.Close()
	err = imaging.Encode(outFile, thumbImg, imaging.JPEG)

	return err
}

// writeThumbnail writes thumbnail for media to w.
//
// It has following sequence/priority:
//  1. Write embedded EXIF thumbnail if it exist (only JPEG)
//  2. Write a cached thumbnail file exist in thumbPath
//  3. Generate a thumbnail to cache and write
//  4. If all above fails return error
func (m *Media) writeThumbnail(w io.Writer, relativeFilePath string) error {
	err := m.writeEXIFThumbnail(w, relativeFilePath)
	if err != nil && m.enableThumbCache {
		err = nil

		// No EXIF, check thumb cache
		thumbFileName, err := m.thumbnailPath(relativeFilePath)
		if err != nil {
			return err
		}
		thumbFile, err := os.Open(thumbFileName)
		if err != nil {
			// No thumb exist. Create it
			llog.Info("Creating new thumbnail for %s", relativeFilePath)
			fullMediaPath, err := m.getFullMediaPath(relativeFilePath)
			if err != nil {
				return err
			}
			err = m.generateImageThumbnail(fullMediaPath, thumbFileName)
			if err != nil {
				return err
			}
			thumbFile, err = os.Open(thumbFileName)
			if err != nil {
				return err
			}
		}
		defer thumbFile.Close()
		_, err = io.Copy(w, thumbFile)
	}
	return err
}

// For testing purposes
var ffmpegCmd = "ffmpeg"

// videoThumbnailSupport returns true if ffmpeg is installed, and thus
// video thumbnails is supported
func videoThumbnailSupport() bool {
	_, err := exec.LookPath(ffmpegCmd)
	return err == nil
}

// generateVideoThumbnail generates a thumbnail from any of the supported
// videos. Will create necessary subdirectories in the thumbpath.
//
// Utilizes the external ffmpeg software.
func (m *Media) generateVideoThumbnail(fullMediaPath, fullThumbPath string) error {
	if !videoThumbnailSupport() {
		return fmt.Errorf("video thumbnails not supported. ffmpeg not installed")
	}
	ffmpegArgs := []string{
		"-i",
		fullMediaPath,
		"-ss",
		"00:00:05", // 5 seconds into movie
		"-vframes",
		"1",
		fullThumbPath}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	//cmd := exec.Command(ffmpegCmd, ffmpegArg)
	cmd := exec.Command(ffmpegCmd, ffmpegArgs...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s %s\nError: %s\nStdout: %s\nStderr: %s",
			ffmpegCmd, strings.Join(ffmpegArgs, " "), err, stdout.String(), stderr.String())
	}
	return err
}
