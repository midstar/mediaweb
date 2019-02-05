package main

import (
	packr "github.com/gobuffalo/packr/v2"
	"github.com/midstar/llog"
)

func main() {
	s := loadSettings(findConfFile())
	llog.SetLevel(s.logLevel)
	if s.logFile != "" {
		llog.SetFile(s.logFile, 1024) // 1 MB logs
	}
	media := createMedia(s.mediaPath, s.thumbPath, s.enableThumbCache, s.autoRotate)
	box := packr.New("templates", "./templates")
	webAPI := CreateWebAPI(s.port, "templates", media, box)
	httpServerDone := webAPI.Start()
	<-httpServerDone // Block until http server is done
}
