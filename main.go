package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func test() {
	startTime := time.Now().UnixNano()
	efile, err := os.Open("IMG_3023.jpg")
	if err != nil {
		log.Printf("Could not open file for EXIF decoder\n")
	}
	defer efile.Close()
	ex, err := exif.Decode(efile)
	if err != nil {
		log.Printf("No EXIF exists: %s\n", err)
	}
	thumbBytes, err := ex.JpegThumbnail()
	if err != nil {
		log.Printf("No thumbnail exists: %s\n", err)
	}
	err = ioutil.WriteFile("thumbnail.jpg", thumbBytes, 0644)
	if err != nil {
		log.Printf("Unable to write thumbnail %s\n", err)
	}
	stopTime := time.Now().UnixNano()
	delta := (stopTime - startTime) / int64(time.Millisecond)
	log.Printf("Delta: %d\n", delta)
}

func main() {
	test()
	media := createMedia(".", ".", false) //createMedia("Y:", ".")
	webAPI := CreateWebAPI(9834, "templates", media)
	httpServerDone := webAPI.Start()
	<-httpServerDone // Block until http server is done
}
