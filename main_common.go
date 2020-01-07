package main

import (
	"github.com/GeertJohan/go.rice"
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
	box := rice.MustFindBox("templates")
	media := createMedia(box, s.mediaPath, s.cachePath,
		s.enableThumbCache, s.genThumbsOnStartup,
		s.genThumbsOnAdd, s.autoRotate, s.enablePreview, s.previewMaxSide,
	    s.genPreviewOnStartup, s.genPreviewOnAdd)
	webAPI := CreateWebAPI(s.port, "templates", media, box, s.userName, s.password)
	return webAPI
}
