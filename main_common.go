package main

import (
	packr "github.com/gobuffalo/packr/v2"
	"github.com/midstar/llog"
)

func mainCommon() *WebAPI {
	s := loadSettings(findConfFile())
	llog.SetLevel(s.logLevel)
	if s.logFile != "" {
		llog.Info("Logging will continue in file %s", s.logFile)
		llog.SetFile(s.logFile, 1024) // 1 MB logs
	}
	llog.Info("Version: %s", applicationVersion)
	llog.Info("Build time: %s", applicationBuildTime)
	llog.Info("Git hash: %s", applicationGitHash)
	box := packr.New("templates", "./templates")
	media := createMedia(box, s.mediaPath, s.thumbPath,
		s.enableThumbCache, s.genThumbsOnStartup, s.autoRotate)
	webAPI := CreateWebAPI(s.port, "templates", media, box, s.userName, s.password)
	return webAPI
}
