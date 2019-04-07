package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	rice "github.com/GeertJohan/go.rice"
)

var baseURL = "http://localhost:9834"

func respToString(response io.ReadCloser) string {
	defer response.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response)
	return buf.String()
}

func getHTML(t *testing.T, path string) string {
	t.Helper()
	resp, err := http.Get(fmt.Sprintf("%s/%s", baseURL, path))
	assertExpectNoErr(t, "", err)
	assertEqualsInt(t, "", int(http.StatusOK), int(resp.StatusCode))
	assertEqualsStr(t, "", "text/html", resp.Header.Get("content-type"))
	defer resp.Body.Close()
	return respToString(resp.Body)
}

func getHTMLAuthenticate(t *testing.T, path, user, pass string, expectFail bool) string {
	t.Helper()
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", baseURL, path), nil)
	assertExpectNoErr(t, "", err)
	req.SetBasicAuth(user, pass)
	resp, err := client.Do(req)
	assertExpectNoErr(t, "", err)
	if expectFail {
		assertEqualsInt(t, "", int(http.StatusUnauthorized), int(resp.StatusCode))
		return ""
	}
	assertEqualsInt(t, "", int(http.StatusOK), int(resp.StatusCode))
	assertEqualsStr(t, "", "text/html", resp.Header.Get("content-type"))
	defer resp.Body.Close()
	return respToString(resp.Body)
}

func getBinary(t *testing.T, path, contentType string) []byte {
	t.Helper()
	resp, err := http.Get(fmt.Sprintf("%s/%s", baseURL, path))
	assertExpectNoErr(t, "", err)
	assertEqualsInt(t, "", int(http.StatusOK), int(resp.StatusCode))
	assertEqualsStr(t, "", contentType, resp.Header.Get("content-type"))
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assertExpectNoErr(t, "", err)
	return body
}

func getObject(t *testing.T, path string, v interface{}) {
	t.Helper()
	resp, err := http.Get(fmt.Sprintf("%s/%s", baseURL, path))
	if err != nil {
		t.Fatalf("Unable to get path %s. Reason: %s", path, err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code for path %s: %d (%s)",
			path, resp.StatusCode, respToString(resp.Body))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to read body for %s. Reason: %s", path, err)
	}
	err = json.Unmarshal(body, &v)
	if err != nil {
		t.Fatalf("Unable decode path %s. Reason: %s", path, err)
	}
}

func startserver(t *testing.T) {
	t.Helper()
	go main()
	waitserver(t)
}

// waitserver waits for the server to be up and running
func waitserver(t *testing.T) {
	t.Helper()
	client := http.Client{Timeout: 100 * time.Millisecond}
	maxTries := 50
	for i := 0; i < maxTries; i++ {
		_, err := client.Get(fmt.Sprintf("%s", baseURL))
		if err == nil {
			// Up and running :-)
			return
		}
	}
	t.Fatalf("Server never started")
}

// shutdown shuts down server and clears the serveMux
func shutdown(t *testing.T) {
	// No answer expecetd on POST shutdown (short timeout)
	client := http.Client{Timeout: 1 * time.Second}
	client.Post(fmt.Sprintf("%s/shutdown", baseURL), "", nil)

	// Reset the serveMux
	http.DefaultServeMux = new(http.ServeMux)
}

func TestStatic(t *testing.T) {
	startserver(t)
	defer shutdown(t)

	// Get default (index)
	index := getHTML(t, "")
	if !strings.Contains(index, "<title>MediaWEB</title>") {
		t.Fatal("Index html title missing")
	}

	// Get index
	index = getHTML(t, "index.html")
	if !strings.Contains(index, "<title>MediaWEB</title>") {
		t.Fatal("Index html title missing")
	}

	// Get a png
	image := getBinary(t, "icon_folder.png", "image/png")
	assertTrue(t, "", len(image) > 100)

	// Get a non-existing png
	resp, err := http.Get(fmt.Sprintf("%s/invalid.html", baseURL))
	assertExpectNoErr(t, "", err)
	assertEqualsInt(t, "", int(http.StatusNotFound), int(resp.StatusCode))
}

func TestListFolders(t *testing.T) {
	startserver(t)
	defer shutdown(t)

	var files []File
	getObject(t, "folder", &files)
	assertTrue(t, "", len(files) > 5)

	// Test list folder that don't exist
	resp, err := http.Get(fmt.Sprintf("%s/folder/dont/exist", baseURL))
	assertExpectNoErr(t, "", err)
	assertEqualsInt(t, "", int(http.StatusNotFound), int(resp.StatusCode))
}

func TestGetMedia(t *testing.T) {
	startserver(t)
	defer shutdown(t)

	image := getBinary(t, "media/gif.gif", "image/gif")
	assertTrue(t, "", len(image) > 100)

	image = getBinary(t, "media/jpeg.jpg", "image/jpeg")
	assertTrue(t, "", len(image) > 100)

	image = getBinary(t, "media/jpeg_rotated.jpg", "image/jpeg")
	assertTrue(t, "", len(image) > 100)

	image = getBinary(t, "media/exif_rotate/no_exif.jpg", "image/jpeg")
	assertTrue(t, "", len(image) > 100)

	image = getBinary(t, "media/video.mp4", "video/mp4")
	assertTrue(t, "", len(image) > 100)

	resp, err := http.Get(fmt.Sprintf("%s/media/dont_exist.jpg", baseURL))
	assertExpectNoErr(t, "", err)
	assertEqualsInt(t, "", int(http.StatusNotFound), int(resp.StatusCode))

	resp, err = http.Get(fmt.Sprintf("%s/media/exif_rotate", baseURL))
	assertExpectNoErr(t, "", err)
	assertEqualsInt(t, "", int(http.StatusNotFound), int(resp.StatusCode))

	resp, err = http.Get(fmt.Sprintf("%s/media/../../hacker.png", baseURL))
	assertExpectNoErr(t, "", err)
	assertEqualsInt(t, "", int(http.StatusNotFound), int(resp.StatusCode))
}

func TestGetThumbnail(t *testing.T) {
	startserver(t)
	defer shutdown(t)

	image := getBinary(t, "thumb/gif.gif", "image/jpeg")
	assertTrue(t, "", len(image) > 100)

	image = getBinary(t, "thumb/jpeg.jpg", "image/jpeg")
	assertTrue(t, "", len(image) > 100)

	image = getBinary(t, "thumb/exif_rotate/no_exif.jpg", "image/jpeg")
	assertTrue(t, "", len(image) > 100)

	// Below will be png if ffmpeg is not installed and jpeg if ffmpeg is installed
	//image = getBinary(t, "thumb/video.mp4", "image/jpeg")
	//assertTrue(t, "", len(image) > 100)

	image = getBinary(t, "thumb/exif_rotate", "image/png")
	assertTrue(t, "", len(image) > 100)

	/*
		Non existing files will give a folder thumbnail by design
		resp, err := http.Get(fmt.Sprintf("%s/thumb/dont_exist.jpg", baseURL))
		assertExpectNoErr(t, "", err)
		assertEqualsInt(t, "", int(http.StatusNotFound), int(resp.StatusCode))
	*/
}

func TestGetThumbnailNoCache(t *testing.T) {
	box := rice.MustFindBox("templates")
	media := createMedia(box, "testmedia", "", false, false, false, true)
	webAPI := CreateWebAPI(9834, "templates", media, box, "", "")
	webAPI.Start()
	waitserver(t)
	defer shutdown(t)

	image := getBinary(t, "thumb/gif.gif", "image/png")
	assertTrue(t, "", len(image) > 100)

	// Has EXIF thumb, i.e. a jpeg is returned
	image = getBinary(t, "thumb/jpeg.jpg", "image/jpeg")
	assertTrue(t, "", len(image) > 100)

	image = getBinary(t, "thumb/exif_rotate/no_exif.jpg", "image/png")
	assertTrue(t, "", len(image) > 100)

	image = getBinary(t, "thumb/video.mp4", "image/png")
	assertTrue(t, "", len(image) > 100)

	image = getBinary(t, "thumb/exif_rotate", "image/png")
	assertTrue(t, "", len(image) > 100)

	/*
		Non existing files will give a folder thumbnail by design
		resp, err := http.Get(fmt.Sprintf("%s/thumb/dont_exist.jpg", baseURL))
		assertExpectNoErr(t, "", err)
		assertEqualsInt(t, "", int(http.StatusNotFound), int(resp.StatusCode))
	*/
}

func TestInvalidPath(t *testing.T) {
	startserver(t)
	defer shutdown(t)

	resp, err := http.Post(fmt.Sprintf("%s/invalid", baseURL), "", nil)
	assertExpectNoErr(t, "", err)
	assertEqualsInt(t, "", int(http.StatusNotFound), int(resp.StatusCode))
}

func TestAuthentication(t *testing.T) {
	box := rice.MustFindBox("templates")
	media := createMedia(box, "testmedia", "", true, false, false, true)
	webAPI := CreateWebAPI(9834, "templates", media, box, "myuser", "mypass")
	webAPI.Start()
	waitserver(t)
	defer shutdown(t)

	// Try to get without any authentication header
	resp, err := http.Get(baseURL)
	assertExpectNoErr(t, "", err)
	assertEqualsInt(t, "", int(http.StatusUnauthorized), int(resp.StatusCode))

	// Try to get with a valid user and password
	index := getHTMLAuthenticate(t, "index.html", "myuser", "mypass", false)
	if !strings.Contains(index, "<title>MediaWEB</title>") {
		t.Fatal("Index html title missing")
	}

	// Try to get with an invalid user but valid password
	getHTMLAuthenticate(t, "index.html", "invalid", "mypass", true)

	// Try to get with a valid user but invalid password
	getHTMLAuthenticate(t, "index.html", "myuser", "invalid", true)

	// Try to get with a valid user and password again
	index = getHTMLAuthenticate(t, "index.html", "myuser", "mypass", false)
	if !strings.Contains(index, "<title>MediaWEB</title>") {
		t.Fatal("Index html title missing")
	}

}
