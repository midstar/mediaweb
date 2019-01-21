package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

var imgExtensions = [...]string{".png", ".jpg", ".jpeg", ".tiff", ".gif"}
var vidExtensions = [...]string{".avi", ".mov", ".vid"}

// Media represents the media including its base path
type Media struct {
	basePath string
}

// File represents a folder or any other file
type File struct {
	Type string // folder, image or video
	Name string
	Path string // Including Name
}

func createMedia(basePath string) *Media {
	return &Media{basePath: basePath}
}

// getFiles returns a slice of File's sorted on file name
func (m *Media) getFiles(relativePath string) ([]File, error) {
	//var files []File
	files := make([]File, 0, 500)
	fullPath := filepath.Join(m.basePath, relativePath)
	fileInfo, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return files, err
	}

	for _, fileInfo := range fileInfo {
		fileType := ""
		if fileInfo.IsDir() {
			fileType = "folder"
		} else {
			extension := filepath.Ext(fileInfo.Name())

			// Check if this is an image
			for _, imgExtension := range imgExtensions {
				if strings.EqualFold(extension, imgExtension) {
					fileType = "image"
					break
				}
			}

			// Check if this is a video
			if fileType == "" {
				for _, vidExtension := range vidExtensions {
					if strings.EqualFold(extension, vidExtension) {
						fileType = "video"
						break
					}
				}
			}
		}
		// Only add directories, videos and images
		if fileType != "" {
			file := File{
				Type: fileType,
				Name: fileInfo.Name(),
				Path: filepath.Join(relativePath, fileInfo.Name())}
			files = append(files, file)
		}
	}
	return files, nil
}
