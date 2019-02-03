package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/midstar/llog"
)

func TestSettingsDefault(t *testing.T) {
	contents :=
		`
port = 9834
mediapath = Y:\pictures`
	fullPath := createConfigFile(t, "TestSettingsDefault.cfg", contents)
	s := loadSettings(fullPath)

	// Mandatory values
	assertEqualsInt(t, "port", 9834, s.port)
	assertEqualsStr(t, "mediaPath", "Y:\\pictures", s.mediaPath)

	// All default on optional
	assertEqualsStr(t, "thumbPath", filepath.Join(os.TempDir(), "mediaweb"), s.thumbPath)
	assertEqualsBool(t, "enableThumbCache", true, s.enableThumbCache)
	assertEqualsBool(t, "autoRotate", true, s.autoRotate)
	assertEqualsInt(t, "logLevel", int(llog.LvlInfo), int(s.logLevel))
	assertEqualsStr(t, "logFile", "", s.logFile)

}

func TestSettings(t *testing.T) {
	contents :=
		`
port = 80
mediapath = /media/usb/pictures
thumbpath = /tmp/thumb
enablethumbcache = off
autorotate = false
loglevel = debug
logfile = /tmp/log/mediaweb.log
`
	fullPath := createConfigFile(t, "TestSettings.cfg", contents)
	s := loadSettings(fullPath)

	// Mandatory values
	assertEqualsInt(t, "port", 80, s.port)
	assertEqualsStr(t, "mediaPath", "/media/usb/pictures", s.mediaPath)

	// Check set values on optional
	assertEqualsStr(t, "thumbPath", "/tmp/thumb", s.thumbPath)
	assertEqualsBool(t, "enableThumbCache", false, s.enableThumbCache)
	assertEqualsBool(t, "autoRotate", false, s.autoRotate)
	assertEqualsInt(t, "logLevel", int(llog.LvlDebug), int(s.logLevel))
	assertEqualsStr(t, "logFile", "/tmp/log/mediaweb.log", s.logFile)

}

func TestSettingsInvalidOptional(t *testing.T) {
	contents :=
		`
port = 80
mediapath = /media/usb/pictures
thumbpath = /tmp/thumb
enablethumbcache = 33
autorotate = invalid
loglevel = debug
logfile = /tmp/log/mediaweb.log
`
	fullPath := createConfigFile(t, "TestSettings.cfg", contents)
	s := loadSettings(fullPath)

	// Mandatory values
	assertEqualsInt(t, "port", 80, s.port)
	assertEqualsStr(t, "mediaPath", "/media/usb/pictures", s.mediaPath)

	// Check set values on optional
	assertEqualsStr(t, "thumbPath", "/tmp/thumb", s.thumbPath)
	assertEqualsInt(t, "logLevel", int(llog.LvlDebug), int(s.logLevel))
	assertEqualsStr(t, "logFile", "/tmp/log/mediaweb.log", s.logFile)

	// Should be default on invalid values
	assertEqualsBool(t, "enableThumbCache", true, s.enableThumbCache)
	assertEqualsBool(t, "autoRotate", true, s.autoRotate)

}

func expectPanic(t *testing.T) {
	// Panic handler (panic is expected)
	recover()
	t.Log("No worry. Panic is expected in the test!!")
}

func TestSettingsNotExisting(t *testing.T) {
	defer expectPanic(t)
	loadSettings("dontexist.cfg")
	t.Fatal("Non existing file. Panic expected")
}

func TestSettingsMissingPort(t *testing.T) {
	contents :=
		`
mediapath = Y:\pictures`
	fullPath := createConfigFile(t, "TestSettingsMissingPort.cfg", contents)
	defer expectPanic(t)
	loadSettings(fullPath)
	t.Fatal("Panic expected")
}

func TestSettingsInvalidPort(t *testing.T) {
	contents :=
		`port=nonint
mediapath = Y:\pictures`
	fullPath := createConfigFile(t, "TestSettingsInvalidPort.cfg", contents)
	defer expectPanic(t)
	loadSettings(fullPath)
	t.Fatal("Panic expected")
}

func TestSettingsMissingMediaPath(t *testing.T) {
	contents :=
		`port=80`
	fullPath := createConfigFile(t, "TestSettingsMissingMediaPath.cfg", contents)
	defer expectPanic(t)
	loadSettings(fullPath)
	t.Fatal("Panic expected")
}

func TestToLogLvl(t *testing.T) {
	checkLvl(t, llog.LvlTrace, "trace")
	checkLvl(t, llog.LvlDebug, "debug")
	checkLvl(t, llog.LvlInfo, "info")
	checkLvl(t, llog.LvlWarn, "warn")
	checkLvl(t, llog.LvlError, "error")
	checkLvl(t, llog.LvlPanic, "panic")

	// Invalid shall be info
	checkLvl(t, llog.LvlInfo, "")
	checkLvl(t, llog.LvlInfo, "invalid")

}

func checkLvl(t *testing.T, expected llog.Level, strLevel string) {
	level := toLogLvl(strLevel)
	if level != expected {
		t.Fatalf("%s should be level %d but was %d", strLevel, int(expected), int(level))
	}

}

// createConfigFile creates a configuration file. Returns the full path to it.
func createConfigFile(t *testing.T, name, contents string) string {
	os.MkdirAll("tmpout", os.ModeDir)
	fullName := "tmpout/" + name
	os.Remove(fullName) // Remove old if it exist
	err := ioutil.WriteFile(fullName, []byte(contents), 0644)
	if err != nil {
		t.Fatalf("Unable to create configuration file. Reason: %s", err)
	}
	return fullName
}
