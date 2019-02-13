package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/midstar/gocfg"
	"github.com/midstar/llog"
)

type settings struct {
	port             int        // Network port
	mediaPath        string     // Top level path for media files
	thumbPath        string     // Top level path for thumbnails
	enableThumbCache bool       // Generate thumbnails
	autoRotate       bool       // Rotate JPEG files when needed
	logLevel         llog.Level // Logging level
	logFile          string     // Log file ("" means stderr)
}

// defaultConfPath holds configuration file paths in priority order
var defaultConfPaths = []string{"mediaweb.conf", "/etc/mediaweb.conf", "/etc/mediaweb/mediaweb.conf"}

// For unit test purposes we do it like this (to be able to change confPaths)
var confPaths = defaultConfPaths

// findConfFile finds the location of the configuration file depending on confPaths
// panics if no configuration file was found
func findConfFile() string {
	result := ""
	for _, path := range confPaths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			result = path
			break
		}
	}
	if result == "" {
		llog.Panic("No configuration file found. Looked in %s", strings.Join(confPaths, ", "))
	}
	return result
}

// loadSettings loads settings from a .conf file. Panics if configuration file
// don't exist or if any of the mandatory settings don't exist.
func loadSettings(fileName string) settings {
	result := settings{}
	llog.Info("Loading configuration: %s", fileName)
	config, err := gocfg.LoadConfiguration(fileName)
	if err != nil {
		llog.Panic("%s", err)
	}

	// Load port (MANDATORY)
	if !config.HasKey("port") {
		llog.Panic("Mandatory property 'port' is not defined in %s", fileName)
	}
	port, err := config.GetInt("port", 0)
	if err != nil {
		llog.Panic("%s", err)
	}
	result.port = port

	// Load mediaPath (MANDATORY)
	if !config.HasKey("mediapath") {
		llog.Panic("Mandatory property 'mediapath' is not defined in %s", fileName)
	}
	mediaPath := config.GetString("mediapath", "")
	result.mediaPath = mediaPath

	// Load thumbPath (OPTIONAL)
	// Default: OS temp directory
	if config.HasKey("thumbpath") {
		thumbPath := config.GetString("thumbpath", "")
		result.thumbPath = thumbPath
	} else {
		// Use default temporary directory + mediaweb
		tempDir := os.TempDir()
		result.thumbPath = filepath.Join(tempDir, "mediaweb")
	}

	// Load enableThumbCache (OPTIONAL)
	// Default: true
	enableThumbCache, err := config.GetBool("enablethumbcache", true)
	if err != nil {
		llog.Warn("%s", err)
	}
	result.enableThumbCache = enableThumbCache

	// Load autoRotate (OPTIONAL)
	// Default: true
	autoRotate, err := config.GetBool("autorotate", true)
	if err != nil {
		llog.Warn("%s", err)
	}
	result.autoRotate = autoRotate

	// Load logFile (OPTIONAL)
	// Default: "" (log to stderr)
	logFile := config.GetString("logfile", "")
	result.logFile = logFile

	// Load logLevel (OPTIONAL)
	// Default: info
	logLevel := config.GetString("loglevel", "info")
	result.logLevel = toLogLvl(logLevel)

	return result
}

func toLogLvl(level string) llog.Level {
	var logLevel llog.Level
	switch level {
	case "trace":
		logLevel = llog.LvlTrace
	case "debug":
		logLevel = llog.LvlDebug
	case "info":
		logLevel = llog.LvlInfo
	case "warn":
		logLevel = llog.LvlWarn
	case "error":
		logLevel = llog.LvlError
	case "panic":
		logLevel = llog.LvlPanic
	default:
		llog.Warn("Invalid loglevel '%s'. Using info level.", level)
		logLevel = llog.LvlInfo
	}

	return logLevel
}
