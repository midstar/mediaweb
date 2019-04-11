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
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/disintegration/imaging"
	"github.com/midstar/llog"
	"github.com/rwcarlsen/goexif/exif"
)

var imgExtensions = [...]string{".png", ".jpg", ".jpeg", ".tif", ".tiff", ".gif"}
var vidExtensions = [...]string{".avi", ".mov", ".vid", ".mkv", ".mp4"}

// Media represents the media including its base path
type Media struct {
	mediaPath          string    // Top level path for media files
	thumbPath          string    // Top level path for thumbnails
	enableThumbCache   bool      // Generate thumbnails
	autoRotate         bool      // Rotate JPEG files when needed
	box                *rice.Box // For icons
	thumbGenInProgress bool      // True if thumbnail generation in progress
	stopWatcherChan    chan bool // Set to true to stop the Watcher go-routine
}

// File represents a folder or any other file
type File struct {
	Type string // folder, image or video
	Name string
	Path string // Including Name. Always using / (even on Windows)
}

// createMedia creates a new media. If thumb cache is enabled the path is
// created when needed.
func createMedia(box *rice.Box, mediaPath string, thumbPath string, enableThumbCache, genThumbsOnStartup, startWatcher, autoRotate bool) *Media {
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
	media := &Media{mediaPath: filepath.ToSlash(filepath.Clean(mediaPath)),
		thumbPath:          filepath.ToSlash(filepath.Clean(thumbPath)),
		enableThumbCache:   enableThumbCache,
		autoRotate:         autoRotate,
		box:                box,
		thumbGenInProgress: false,
		stopWatcherChan:    make(chan bool)}
	llog.Info("Video thumbnails supported (ffmpeg installed): %v", media.videoThumbnailSupport())
	if enableThumbCache && genThumbsOnStartup {
		go media.generateAllThumbnails()
	}
	if enableThumbCache && startWatcher {
		go media.startWatcher()
	}
	return media
}

// getFullPath returns the full path from an absolute base
// path and a relative path. Returns error on security hacks,
// i.e. when someone tries to access ../../../ for example to
// get files that are not within configured base path.
//
// Always returning front slashes / as path separator
func (m *Media) getFullPath(basePath, relativePath string) (string, error) {
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

// getRelativePath returns the relative path from an absolute base
// path and a full path path. Returns error if the base path is
// not in the full path.
//
// Always returning front slashes / as path separator
func (m *Media) getRelativePath(basePath, fullPath string) (string, error) {
	relativePath, err := filepath.Rel(basePath, fullPath)
	if err == nil {
		relativePathSlash := filepath.ToSlash(relativePath)
		if strings.HasPrefix(relativePathSlash, "../") {
			return "", fmt.Errorf("%s is not a sub-path of %s", fullPath, basePath)
		}
		return relativePathSlash, nil
	}
	return "", err
}

// getRelativeMediaPath returns the relative media path of the provided path, i.e:
// full path - media path.
func (m *Media) getRelativeMediaPath(fullPath string) (string, error) {
	return m.getRelativePath(m.mediaPath, fullPath)
}

// getFiles returns a slice of File's sorted on file name
func (m *Media) getFiles(relativePath string) ([]File, error) {
	//var files []File
	files := make([]File, 0, 500)
	fullPath, err := m.getFullMediaPath(relativePath)
	if err != nil {
		return files, err
	}
	fileInfos, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return files, err
	}

	for _, fileInfo := range fileInfos {
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

	// Check if this is an image
	if m.isImage(relativeFileName) {
		return "image"
	}

	// Check if this is a video
	if m.isVideo(relativeFileName) {
		return "video"
	}

	return "" // Not a video nor an image
}

func (m *Media) isImage(pathAndFile string) bool {
	extension := filepath.Ext(pathAndFile)
	for _, imgExtension := range imgExtensions {
		if strings.EqualFold(extension, imgExtension) {
			return true
		}
	}
	return false
}

func (m *Media) isVideo(pathAndFile string) bool {
	extension := filepath.Ext(pathAndFile)
	for _, vidExtension := range vidExtensions {
		if strings.EqualFold(extension, vidExtension) {
			return true
		}
	}
	return false
}

func (m *Media) isJPEG(pathAndFile string) bool {
	extension := filepath.Ext(pathAndFile)
	if strings.EqualFold(extension, ".jpg") == false &&
		strings.EqualFold(extension, ".jpeg") == false {
		return false
	}
	return true
}

func (m *Media) extractEXIF(relativeFilePath string) *exif.Exif {
	fullFilePath, err := m.getFullMediaPath(relativeFilePath)
	if err != nil {
		llog.Info("Unable to get full media path for %s\n", relativeFilePath)
		return nil
	}
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

// isRotationNeeded returns true if the file needs to be rotated.
// It finds this out by reading the EXIF rotation information
// in the file.
// If Media.autoRotate is false this function will always return
// false.
func (m *Media) isRotationNeeded(relativeFilePath string) bool {
	if m.autoRotate == false {
		return false
	}
	ex := m.extractEXIF(relativeFilePath)
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
	ex := m.extractEXIF(relativeFilePath)
	if ex == nil {
		return fmt.Errorf("No EXIF info for %s", relativeFilePath)
	}
	thumbBytes, err := ex.JpegThumbnail()
	if err != nil {
		return fmt.Errorf("No EXIF thumbnail for %s", relativeFilePath)
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
			llog.Warn("Unable to decode EXIF thumbnail for %s", relativeFilePath)
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
	thumbImg := imaging.Thumbnail(img, 256, 256, imaging.Box)

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

// generateTumbnail generates a thumbnail for an image or video
// and returns the file name of the thumbnail. If a thumbnail already
// exist the file name will be returned.
func (m *Media) generateThumbnail(relativeFilePath string) (string, error) {
	thumbFileName, err := m.thumbnailPath(relativeFilePath)
	if err != nil {
		llog.Error("%s", err)
		return "", err
	}
	_, err = os.Stat(thumbFileName) // Check if file exist
	if err == nil {
		return thumbFileName, nil // Thumb already generated
	}

	// No thumb exist. Create it
	llog.Info("Creating new thumbnail for %s", relativeFilePath)
	startTime := time.Now().UnixNano()
	fullMediaPath, err := m.getFullMediaPath(relativeFilePath)
	if err != nil {
		llog.Error("%s", err)
		return thumbFileName, err
	}
	if m.isVideo(fullMediaPath) {
		err = m.generateVideoThumbnail(fullMediaPath, thumbFileName)
	} else {
		err = m.generateImageThumbnail(fullMediaPath, thumbFileName)
	}
	if err != nil {
		llog.Error("%s", err)
		return thumbFileName, err
	}
	deltaTime := (time.Now().UnixNano() - startTime) / int64(time.Millisecond)
	llog.Info("Thumbnail done for %s (conversion time: %d ms)",
		relativeFilePath, deltaTime)
	return thumbFileName, nil
}

// writeThumbnail writes thumbnail for media to w.
//
// It has following sequence/priority:
//  1. Write embedded EXIF thumbnail if it exist (only JPEG)
//  2. Write a cached thumbnail file exist in thumbPath
//  3. Generate a thumbnail to cache and write
//  4. If all above fails return error
func (m *Media) writeThumbnail(w io.Writer, relativeFilePath string) error {
	if !m.isImage(relativeFilePath) && !m.isVideo(relativeFilePath) {
		return fmt.Errorf("not a supported media type")
	}
	if m.writeEXIFThumbnail(w, relativeFilePath) == nil {
		return nil
	}
	if !m.enableThumbCache {
		return fmt.Errorf("Thumbnail cache disabled")
	}

	// No EXIF, check thumb cache (and generate if necessary)
	thumbFileName, err := m.generateThumbnail(relativeFilePath)
	if err != nil {
		return err // Logging handled in generateThumbnail
	}

	thumbFile, err := os.Open(thumbFileName)
	if err != nil {
		return err
	}
	defer thumbFile.Close()

	_, err = io.Copy(w, thumbFile)
	if err != nil {
		return err
	}

	return nil
}

// For testing purposes
var ffmpegCmd = "ffmpeg"

// videoThumbnailSupport returns true if ffmpeg is installed, and thus
// video thumbnails is supported
func (m *Media) videoThumbnailSupport() bool {
	_, err := exec.LookPath(ffmpegCmd)
	return err == nil
}

// generateVideoThumbnail generates a thumbnail from any of the supported
// videos. Will create necessary subdirectories in the thumbpath.
func (m *Media) generateVideoThumbnail(fullMediaPath, fullThumbPath string) error {
	// The temporary file for the screenshot
	screenShot := fullThumbPath + ".sh.jpg"

	// Extract the screenshot
	err := m.extractVideoScreenshot(fullMediaPath, screenShot)
	if err != nil {
		return err
	}
	defer os.Remove(screenShot) // Remove temporary file

	// Generate thumbnail from the screenshot
	img, err := imaging.Open(screenShot, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("Unable to open screenshot image %s. Reason: %s", screenShot, err)
	}
	thumbImg := imaging.Thumbnail(img, 256, 256, imaging.Box)

	// Add small video icon i upper right corner to indicate that this is
	// a video
	iconVideoImg, err := m.getVideoIcon()
	if err != nil {
		return err
	}
	thumbImg = imaging.Overlay(thumbImg, iconVideoImg, image.Pt(155, 11), 1.0)

	// Write thumbnail to file
	outFile, err := os.Create(fullThumbPath)
	if err != nil {
		return fmt.Errorf("Unable to open %s for creating thumbnail. Reason %s", fullThumbPath, err)
	}
	defer outFile.Close()
	err = imaging.Encode(outFile, thumbImg, imaging.JPEG)

	return err
}

// Cache to avoid regenerate icon each time (do it once)
var videoIcon image.Image

func (m *Media) getVideoIcon() (image.Image, error) {
	if videoIcon != nil {
		// To avoid re-generate
		return videoIcon, nil
	}
	var err error
	videoIconBytes, _ := m.box.Bytes("icon_video.png")
	videoIcon, err = imaging.Decode(bytes.NewReader(videoIconBytes))
	if err != nil {
		return nil, err
	}
	videoIcon = imaging.Resize(videoIcon, 90, 90, imaging.Box)
	return videoIcon, nil
}

// extractVideoScreenshot extracts a screenshot from a video using external
// ffmpeg software. Will create necessary directories in the outFilePath
func (m *Media) extractVideoScreenshot(inFilePath, outFilePath string) error {
	if !m.videoThumbnailSupport() {
		return fmt.Errorf("video thumbnails not supported. ffmpeg not installed")
	}

	// Create subdirectories if needed
	directory := filepath.Dir(outFilePath)
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Unable to create directories in %s for extracting screenshot. Reason %s", outFilePath, err)
	}

	// Define argments for ffmpeg
	ffmpegArgs := []string{
		"-i",
		inFilePath,
		"-ss",
		"00:00:05", // 5 seconds into movie
		"-vframes",
		"1",
		outFilePath}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	//cmd := exec.Command(ffmpegCmd, ffmpegArg)
	cmd := exec.Command(ffmpegCmd, ffmpegArgs...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%s %s\nError: %s\nStdout: %s\nStderr: %s",
			ffmpegCmd, strings.Join(ffmpegArgs, " "), err, stdout.String(), stderr.String())
	}
	return nil
}

// generateAllThumbnails goes through all files in the media path
// and generates thumbnails for these
func (m *Media) generateAllThumbnails() {
	m.thumbGenInProgress = true
	llog.Info("Generating all thumbnails")
	startTime := time.Now().UnixNano()
	stat := m.generateThumbnails("")
	deltaTime := (time.Now().UnixNano() - startTime) / int64(time.Second)
	minutes := int(deltaTime / 60)
	seconds := int(deltaTime) - minutes*60
	llog.Info(`Generating all thumbnails took %d minutes and %d seconds
  Number of folders: %d
  Number of images: %d
  Number of vidos: %d
  Number of images with embedded EXIF: %d
  Number of failed folders: %d
  Number of failed images: %d
  Number of failed videos: %d`, minutes, seconds, stat.NbrOfFolders, stat.NbrOfImages,
		stat.NbrOfVideos, stat.NbrOfExif, stat.NbrOfFailedFolders, stat.NbrOfFailedImages, stat.NbrOfFailedVideos)
	m.thumbGenInProgress = false
}

// ThumbnailStatistics statistics results from generateThumbnails
type ThumbnailStatistics struct {
	NbrOfFolders       int
	NbrOfImages        int
	NbrOfVideos        int
	NbrOfExif          int
	NbrOfFailedFolders int // I.e. unable to list contents of folder
	NbrOfFailedImages  int
	NbrOfFailedVideos  int
}

// generateThumbnails recursively goes through all files relativePath
// and its subdirectories and generates thumbnails for these. If
// relativePath is "" it means generate for all files.
func (m *Media) generateThumbnails(relativePath string) *ThumbnailStatistics {
	stat := ThumbnailStatistics{}
	files, err := m.getFiles(relativePath)
	if err != nil {
		stat.NbrOfFailedFolders = 1
		return &stat
	}
	for _, file := range files {
		if file.Type == "folder" {
			stat.NbrOfFolders++
			newStat := m.generateThumbnails(file.Path) // Recursive
			stat.NbrOfFolders += newStat.NbrOfFolders
			stat.NbrOfImages += newStat.NbrOfImages
			stat.NbrOfVideos += newStat.NbrOfVideos
			stat.NbrOfExif += newStat.NbrOfExif
			stat.NbrOfFailedFolders += newStat.NbrOfFailedFolders
			stat.NbrOfFailedImages += newStat.NbrOfFailedImages
			stat.NbrOfFailedVideos += newStat.NbrOfFailedVideos
		} else {
			if file.Type == "image" {
				stat.NbrOfImages++
			} else if file.Type == "video" {
				stat.NbrOfVideos++
			}
			// Check if file has EXIF thumbnail
			ex := m.extractEXIF(file.Path)
			if ex != nil {
				_, err := ex.JpegThumbnail()
				if err == nil {
					// Media has EXIF thumbnail
					stat.NbrOfExif++
					continue // Next file
				}
			}
			// Generate new thumbnail
			_, err = m.generateThumbnail(file.Path)
			if err != nil {
				if file.Type == "image" {
					stat.NbrOfFailedImages++
				} else if file.Type == "video" {
					stat.NbrOfFailedVideos++
				}
			}
		}
	}
	return &stat
}
