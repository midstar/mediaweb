package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/midstar/gocfg"
	"github.com/midstar/llog"
)

type settings struct {
	port                int        // Network port
	ip                  string     // Network IP ("" means any)
	mediaPath           string     // Top level path for media files
	cachePath           string     // Top level path for cache (thumbs and preview)
	enableThumbCache    bool       // Generate thumbnails
	genThumbsOnStartup  bool       // Generate all thumbnails on startup
	genThumbsOnAdd      bool       // Generate thumbnails when file added (start watcher)
	autoRotate          bool       // Rotate JPEG files when needed
	enablePreview       bool       // Generate preview files
	previewMaxSide      int        // Max height/width of preview file
	genPreviewOnStartup bool       // Generate all preview on startup
	genPreviewOnAdd     bool       // Generate preview when file added (start watcher)
	logLevel            llog.Level // Logging level
	logFile             string     // Log file ("" means stderr)
	userName            string     // User name ("" means no authentication)
	password            string     // Password
	tlsCertFile         string     // TLS certification file
	tlsKeyFile          string     // TLS key file
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

	// Load IP (OPTIONAL)
	// Default: ""
	ip := config.GetString("ip", "")
	result.ip = ip

	// Load mediaPath (MANDATORY)
	if !config.HasKey("mediapath") {
		llog.Panic("Mandatory property 'mediapath' is not defined in %s", fileName)
	}
	mediaPath := config.GetString("mediapath", "")
	result.mediaPath = mediaPath

	// Load cachePath (OPTIONAL)
	// Default: OS temp directory
	if config.HasKey("cachepath") {
		cachePath := config.GetString("cachepath", "")
		result.cachePath = cachePath
	} else {
		// For backwards compatibility with old versions
		if config.HasKey("thumbpath") {
			cachePath := config.GetString("thumbpath", "")
			result.cachePath = cachePath
		} else {
			// Use default temporary directory + mediaweb
			tempDir := os.TempDir()
			result.cachePath = filepath.Join(tempDir, "mediaweb")
		}
	}

	// Check that mediapath and cachepath are not the same
	if pathEquals(result.mediaPath, result.cachePath) {
		llog.Panic("cachepath and mediapath have the same value '%s'", result.mediaPath)
	}

	// Load enableThumbCache (OPTIONAL)
	// Default: true
	enableThumbCache, err := config.GetBool("enablethumbcache", true)
	if err != nil {
		llog.Warn("%s", err)
	}
	result.enableThumbCache = enableThumbCache

	// Load genthumbsonstartup (OPTIONAL)
	// Default: false
	genThumbsOnStartup, err := config.GetBool("genthumbsonstartup", false)
	if err != nil {
		llog.Warn("%s", err)
	}
	result.genThumbsOnStartup = genThumbsOnStartup

	// Load genthumbsonadd (OPTIONAL)
	// Default: true
	genThumbsOnAdd, err := config.GetBool("genthumbsonadd", true)
	if err != nil {
		llog.Warn("%s", err)
	}
	result.genThumbsOnAdd = genThumbsOnAdd

	// Load autoRotate (OPTIONAL)
	// Default: true
	autoRotate, err := config.GetBool("autorotate", true)
	if err != nil {
		llog.Warn("%s", err)
	}
	result.autoRotate = autoRotate

	// Load enablePreview (OPTIONAL)
	// Default: false
	enablePreview, err := config.GetBool("enablepreview", false)
	if err != nil {
		llog.Warn("%s", err)
	}
	result.enablePreview = enablePreview

	// Load previewMaxSide (OPTIONAL)
	// Default: 1280 (pixels)
	previewMaxSide, err := config.GetInt("previewmaxside", 1280)
	if err != nil {
		llog.Warn("%s", err)
	}
	result.previewMaxSide = previewMaxSide

	// Load genpreviewonstartup (OPTIONAL)
	// Default: false
	genPreviewOnStartup, err := config.GetBool("genpreviewonstartup", false)
	if err != nil {
		llog.Warn("%s", err)
	}
	result.genPreviewOnStartup = genPreviewOnStartup

	// Load genpreviewonadd (OPTIONAL)
	// Default: true
	genPreviewOnAdd, err := config.GetBool("genpreviewonadd", true)
	if err != nil {
		llog.Warn("%s", err)
	}
	result.genPreviewOnAdd = genPreviewOnAdd

	// Load logFile (OPTIONAL)
	// Default: "" (log to stderr)
	logFile := config.GetString("logfile", "")
	result.logFile = logFile

	// Load logLevel (OPTIONAL)
	// Default: info
	logLevel := config.GetString("loglevel", "info")
	result.logLevel = toLogLvl(logLevel)

	// Load username (OPTIONAL)
	// Default: "" (no authentication)
	userName := config.GetString("username", "")
	result.userName = userName

	// Load password (OPTIONAL)
	// Default: ""
	password := config.GetString("password", "")
	result.password = password

	// Load tlsCertFile (OPTIONAL)
	// Default: ""
	tlsCertFile := config.GetString("tlscertfile", "")
	result.tlsCertFile = tlsCertFile

	// Load tlsKeyFile (OPTIONAL)
	// Default: ""
	tlsKeyFile := config.GetString("tlskeyfile", "")
	result.tlsKeyFile = tlsKeyFile

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

func pathEquals(path1, path2 string) bool {
	diffPath, err := filepath.Rel(path1, path2)
	if err == nil && (diffPath == "" || diffPath == ".") {
		return true
	}
	return false
}
