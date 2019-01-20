package main

import (
	"io/ioutil"
	"path"
)

// Media represents the media including its base path
type Media struct {
	basePath string
}

// File represents a directory or any other file
type File struct {
	IsDir bool
	Name  string
	Path  string // Including Name
}

func createMedia(basePath string) *Media {
	return &Media{basePath: basePath}
}

// getFiles returns a slice of File's sorted on file name
func (m *Media) getFiles(relativePath string) ([]File, error) {
	//var files []File
	files := make([]File, 0, 500)
	fullPath := path.Join(m.basePath, relativePath)
	fileInfo, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return files, err
	}

	for _, fileInfo := range fileInfo {
		file := File{
			IsDir: fileInfo.IsDir(),
			Name:  fileInfo.Name(),
			Path:  path.Join(relativePath, fileInfo.Name())}
		files = append(files, file)
	}
	return files, nil
}
