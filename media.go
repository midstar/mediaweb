package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

var imgExtensions = [...]string{".png", ".jpg", ".jpeg", ".tif", ".gif"}
var vidExtensions = [...]string{".avi", ".mov", ".vid", ".mkv", ".mp4"}

// Media represents the media including its base path
type Media struct {
	mediaPath  string // Top level path for media files
	thumbPath  string // Top level path for thumbnails
	autoRotate bool   // Rotate JPEG files when needed
}

// File represents a folder or any other file
type File struct {
	Type string // folder, image or video
	Name string
	Path string // Including Name. Always using / (even on Windows)
}

func createMedia(mediaPath string, thumbPath string, autoRotate bool) *Media {
	return &Media{mediaPath: mediaPath,
		thumbPath:  thumbPath,
		autoRotate: autoRotate}
}

// getFullMediaPath returns the full path of the provided path, i.e:
// media path + relative path. Returns error on security hacks,
// i.e. when someone tries to access ../../../ for example to
// get files that are not within configured media folder.
func (m *Media) getFullMediaPath(relativePath string) (string, error) {
	fullPath := filepath.Join(m.mediaPath, relativePath)
	return fullPath, nil
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

func (m *Media) extractEXIF(relativeFilePath string) *exif {
	return nil
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
	extension := filepath.Ext(relativeFilePath)
	if strings.EqualFold(extension, ".jpg") == false &&
		strings.EqualFold(extension, ".jpeg") == false {
		return false // Only JPEG can be rotaded
	}
	fullPath, err := m.getFullMediaPath(relativeFilePath)
	if err != nil {
		log.Printf("Unable to get full media path for %s\n", relativeFilePath)
		return false
	}
	efile, err := os.Open(fullPath)
	if err != nil {
		log.Printf("Could not open file for EXIF decoder: %s\n", fullPath)
		return false
	}
	defer efile.Close()
	ex, err := exif.Decode(efile)
	if err != nil {
		return false // No EXIF info exist
	}
	orientTag, _ := ex.Get(exif.Orientation)
	if orientTag == nil {
		return false // No Orientation
	}
	orientInt, err := strconv.Atoi(orientTag.String())
	if err != nil {
		return false // Orientation is not set correctly, assume no rotation
	}
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

// thumbnailPath returns the thumbnail file path. Thumbnails are always
// stored in JPEG format (.jpg extension) and starts with '_'.
// Returns error if the media path is invalid.
func (m *Media) thumbnailPath(w io.Writer, relativeMediaPath string) (string, error) {
	return "", nil
}

// writeThumbnail writes thumbnail for media to w. If no thumbnail exist
// and error will be returned.
// It will first check if it is a JPEG with an embedded thumbnail. If not
// it will check if a thumbnail is stored in the thumbnail folder.
func (m *Media) writeThumbnail(w io.Writer, relativeFilePath string) error {

	/*	fullMediaPath, err := m.getFullMediaPath(relativeFilePath)
		if err != nil {
			return err
		}
		fmt.Printf("%s", fullMediaPath)
		efile, err := os.Open("20161201_084009.jpg")
		if err != nil {
			log.Printf("Could not open file for EXIF decoder: %s\n", fullPath)
			return false
		}
		defer efile.Close()
		ex, err := exif.Decode(efile)
		if err != nil {
			return false // No EXIF info exist
		}
		orientTag, _ := ex.Get(exif.Orientation)
		if orientTag == nil {
			return false // No Orientation
		}	*/
	return nil
}
