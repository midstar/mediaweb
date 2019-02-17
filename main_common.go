package main

import (
	packr "github.com/gobuffalo/packr/v2"
	"github.com/midstar/llog"
)

func mainCommon() *WebAPI {
	s := loadSettings(findConfFile())
	llog.SetLevel(s.logLevel)
	if s.logFile != "" {
		llog.SetFile(s.logFile, 1024) // 1 MB logs
	}
	llog.Info("Version: %s", applicationVersion)
	llog.Info("Build time: %s", applicationBuildTime)
	llog.Info("Git hash: %s", applicationGitHash)
	media := createMedia(s.mediaPath, s.thumbPath, s.enableThumbCache, s.autoRotate)
	box := packr.New("templates", "./templates")
	webAPI := CreateWebAPI(s.port, "templates", media, box, s.userName, s.password)
	return webAPI
}
