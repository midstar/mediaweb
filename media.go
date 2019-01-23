package main

import (
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

var imgExtensions = [...]string{".png", ".jpg", ".jpeg", ".tiff", ".gif"}
var vidExtensions = [...]string{".avi", ".mov", ".vid", ".mkv", ".mp4"}

// Media represents the media including its base path
type Media struct {
	basePath string
}

// File represents a folder or any other file
type File struct {
	Type string // folder, image or video
	Name string
	Path string // Including Name. Always using / (even on Windows)
}

func createMedia(basePath string) *Media {
	return &Media{basePath: basePath}
}

// getFullPath returns the full path of the provided path, i.e:
// media path + relative path. Returns error on security hacks,
// i.e. when someone tries to access ../../../ for example to
// get files that are not within configured media folder.
func (m *Media) getFullPath(relativePath string) (string, error) {
	fullPath := filepath.Join(m.basePath, relativePath)
	return fullPath, nil
}

// getFiles returns a slice of File's sorted on file name
func (m *Media) getFiles(relativePath string) ([]File, error) {
	//var files []File
	files := make([]File, 0, 500)
	fullPath, err := m.getFullPath(relativePath)
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

// rotateAndWriteJPEG opens and rotates a JPG/JPEG file according to
// EXIF rotation information. Then it writes the rotated image
// to the io.Writer.
func (m *Media) rotateAndWriteJPEG(w io.Writer, relativeFilePath string) error {
	fullPath, err := m.getFullPath(relativeFilePath)
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
