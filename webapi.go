package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/GeertJohan/go.rice"
	"github.com/midstar/llog"
)

// WebAPI represents the REST API server.
type WebAPI struct {
	server       *http.Server
	templatePath string // Path to the templates
	media        *Media
	box          *rice.Box
	userName     string // User name ("" means no authentication)
	password     string // Password
}

// CreateWebAPI creates a new Web API instance
func CreateWebAPI(port int, templatePath string, media *Media, box *rice.Box, userName, password string) *WebAPI {
	portStr := fmt.Sprintf(":%d", port)
	server := &http.Server{Addr: portStr}
	webAPI := &WebAPI{
		server:       server,
		templatePath: templatePath,
		media:        media,
		box:          box,
		userName:     userName,
		password:     password}
	http.Handle("/", webAPI)
	return webAPI
}

// Start starts the HTTP server. Stop it using the Stop function. Non-blocking.
// Returns a channel that is written to when the HTTP server has stopped.
func (wa *WebAPI) Start() chan bool {
	done := make(chan bool)

	go func() {
		llog.Info("Starting Web API on port %s\n", wa.server.Addr)
		if err := wa.server.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			llog.Info("WebAPI: ListenAndServe() shutdown reason: %s", err)
		}
		// TODO fix this wa.media.stopWatcher() // Stop the folder watcher (if it is running)
		done <- true // Signal that http server has stopped
	}()
	return done
}

// Stop stops the HTTP server.
func (wa *WebAPI) Stop() {
	wa.server.Shutdown(context.Background())
}

// ServeHTTP handles incoming HTTP requests
func (wa *WebAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Handle authentication
	if wa.userName != "" {
		// Authentication required
		user, pass, _ := r.BasicAuth()
		if wa.userName != user || wa.password != pass {
			llog.Info("Invalid user login attempt. user: %s, password: %s", user, pass)
			w.Header().Set("WWW-Authenticate", "Basic realm=\"MediaWEB requires username and password\"")
			http.Error(w, "Unauthorized. Invalid username or password.", http.StatusUnauthorized)
			return
		}
	}

	// Handle request
	var head string
	originalURL := r.URL.Path
	llog.Trace("Got request: %s", r.URL.Path)
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "shutdown" && r.Method == "POST" {
		wa.Stop()
	} else if head == "folder" && r.Method == "GET" {
		wa.serveHTTPFolder(w, r)
	} else if head == "media" && r.Method == "GET" {
		wa.serveHTTPMedia(w, r)
	} else if head == "thumb" && r.Method == "GET" {
		wa.serveHTTPThumbnail(w, r)
	} else if head == "isThumbGenInProgress" && r.Method == "GET" {
		toJSON(w, wa.media.isThumbGenInProgress())
	} else if r.Method == "GET" {
		r.URL.Path = originalURL
		wa.serveHTTPStatic(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "This is not a valid path: %s or method %s!", r.URL.Path, r.Method)
	}
}

func (wa *WebAPI) serveHTTPStatic(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Path
	if len(r.URL.Path) > 0 {
		fileName = r.URL.Path[1:] // Remove '/'
	}
	if fileName == "" {
		// Default is index page
		fileName = "index.html"
	}
	bytes, err := wa.box.Bytes(fileName)
	if err != nil || len(bytes) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Unable to find: %s!", fileName)
	} else {
		if filepath.Ext(fileName) == ".html" {
			w.Header().Set("Content-Type", "text/html")
		} else {
			w.Header().Set("Content-Type", "image/png")
		}
		w.Write(bytes)
	}
}

// serveHTTPFolder generates JSON will files in folder
func (wa *WebAPI) serveHTTPFolder(w http.ResponseWriter, r *http.Request) {
	folder := ""
	if len(r.URL.Path) > 0 {
		folder = r.URL.Path[1:] // Remove '/'
	}
	files, err := wa.media.getFiles(folder)
	if err != nil {
		http.Error(w, "Get files: "+err.Error(), http.StatusNotFound)
		return
	}
	toJSON(w, files)
}

// serveHTTPMedia opens the media
func (wa *WebAPI) serveHTTPMedia(w http.ResponseWriter, r *http.Request) {
	relativePath := r.URL.Path
	// Only accept media files of security reasons
	if wa.media.getFileType(relativePath) == "" {
		http.Error(w, "Not a valid media file: "+relativePath, http.StatusNotFound)
		return
	}
	if wa.media.isRotationNeeded(relativePath) {
		// This is a JPEG file which requires rotation.
		w.Header().Set("Content-Type", "image/jpeg")
		err := wa.media.rotateAndWrite(w, relativePath)
		if err != nil {
			http.Error(w, "Rotate file: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// This is any other media file
		fullPath, err := wa.media.getFullMediaPath(relativePath)
		if err != nil {
			http.Error(w, "Get files: "+err.Error(), http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, fullPath)
	}
}

// serveHTTPThumbnail opens the media thumbnail or the default thumbnail
// if no thumbnail exist.
func (wa *WebAPI) serveHTTPThumbnail(w http.ResponseWriter, r *http.Request) {
	relativePath := r.URL.Path
	err := wa.media.writeThumbnail(w, relativePath)
	if err == nil {
		w.Header().Set("Content-Type", "image/jpeg")
	} else {
		// No thumbnail. Use the default
		w.Header().Set("Content-Type", "image/png")
		fileType := wa.media.getFileType(relativePath)
		if fileType == "image" {
			iconImage, _ := wa.box.Bytes("icon_image.png")
			w.Write(iconImage)
			//http.ServeFile(w, r, wa.templatePath+"/icon_image.png")
		} else if fileType == "video" {
			iconVideo, _ := wa.box.Bytes("icon_video.png")
			w.Write(iconVideo)
			//http.ServeFile(w, r, wa.templatePath+"/icon_video.png")
		} else {
			// Folder
			iconFolder, _ := wa.box.Bytes("icon_folder.png")
			w.Write(iconFolder)
			//http.ServeFile(w, r, wa.templatePath+"/icon_folder.png")
		}
	}
}

// toJSON converts the v object to JSON and writes result to the response
func toJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// shiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
