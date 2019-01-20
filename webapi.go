package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

// WebAPI represents the REST API server.
type WebAPI struct {
	server       *http.Server
	templatePath string // Path to the templates
	media        *Media
}

// CreateWebAPI creates a new Web API instance
func CreateWebAPI(port int, templatePath string, media *Media) *WebAPI {
	portStr := fmt.Sprintf(":%d", port)
	server := &http.Server{Addr: portStr}
	webAPI := &WebAPI{
		server:       server,
		templatePath: templatePath,
		media:        media}
	http.Handle("/", webAPI)
	return webAPI
}

// Start starts the HTTP server. Stop it using the Stop function. Non-blocking.
// Returns a channel that is written to when the HTTP server has stopped.
func (wa *WebAPI) Start() chan bool {
	done := make(chan bool)

	go func() {
		log.Printf("Starting Web API on port %s\n", wa.server.Addr)
		if err := wa.server.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("WebAPI: ListenAndServe() shutdown reason: %s", err)
		}
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
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "shutdown" && r.Method == "POST" {
		wa.Stop()
	} else if head == "" && r.Method == "GET" {
		wa.serveHTTPIndex(w, r)
	} else if head == "folder" && r.Method == "GET" {
		wa.serveHTTPFolder(w, r)
	} else if head == "media" && r.Method == "GET" {
		wa.serveHTTPMedia(w, r)
	} else if head == "static" && r.Method == "GET" {
		wa.serveHTTPStatic(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "This is not a valid path: %s or method %s!", r.URL.Path, r.Method)
	}
}

func (wa *WebAPI) serveHTTPIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(wa.templatePath, "index.html"))
}

func (wa *WebAPI) serveHTTPStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(wa.templatePath, r.URL.Path))
}

// serveHTTPFolder generates JSON will files in folder
func (wa *WebAPI) serveHTTPFolder(w http.ResponseWriter, r *http.Request) {
	// TBD secure that this is an allowed folder - SECURITY RISK
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
	// TBD secure that this is an allowed file - SECURITY RISK
	http.ServeFile(w, r, filepath.Join(wa.media.basePath, r.URL.Path))
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
