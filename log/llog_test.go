// Unit tests for the llog package. The unit tests are semi
// manual, i.e. you have to manually check the output and
// have test logging turned on.
package llog

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

type levels struct {
	trace bool
	debug bool
	info  bool
	warn  bool
	err   bool
	panic bool
}

func logAndGetLevelsLogged() (levels, string) {
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	log.SetOutput(buffer)
	logAllLevelsExceptPanic()
	result := buffer.String()
	var l levels
	l.trace = strings.Contains(result, "TRACE")
	l.debug = strings.Contains(result, "DEBUG")
	l.info = strings.Contains(result, "INFO")
	l.warn = strings.Contains(result, "WARN")
	l.err = strings.Contains(result, "ERROR")
	return l, result
}

func logAllLevelsExceptPanic() {
	Trace("this is a trace - param %d and %s", 1, "param 2")
	Debug("this is a debug - param %d and %s", 1, "param 2")
	Info("this is an info - param %d and %s", 1, "param 2")
	Warn("this is a warn - param %d and %s", 1, "param 2")
	Error("this is an error - param %d and %s", 1, "param 2")
}

func TestLevelTrace(t *testing.T) {
	SetLevel(LvlTrace)
	l, result := logAndGetLevelsLogged()
	if !l.trace || !l.debug || !l.info || !l.warn || !l.err {
		t.Fatalf("All levels should have been logged! Result:\n%s", result)
	}
}
func TestLevelDebug(t *testing.T) {
	SetLevel(LvlDebug)
	l, result := logAndGetLevelsLogged()
	if l.trace || !l.debug || !l.info || !l.warn || !l.err {
		t.Fatalf("Trace shall not be in log! Result:\n%s", result)
	}
}
func TestLevelInfo(t *testing.T) {
	SetLevel(LvlInfo)
	l, result := logAndGetLevelsLogged()
	if l.trace || l.debug || !l.info || !l.warn || !l.err {
		t.Fatalf("Trace and debug shall not be in log! Result:\n%s", result)
	}
}
func TestLevelWarn(t *testing.T) {
	SetLevel(LvlWarn)
	l, result := logAndGetLevelsLogged()
	if l.trace || l.debug || l.info || !l.warn || !l.err {
		t.Fatalf("Trace, debug and info shall not be in log! Result:\n%s", result)
	}
}
func TestLevelError(t *testing.T) {
	SetLevel(LvlError)
	l, result := logAndGetLevelsLogged()
	if l.trace || l.debug || l.info || l.warn || !l.err {
		t.Fatalf("Trace, debug and info shall not be in log! Result:\n%s", result)
	}
}
func TestLevelPanic(t *testing.T) {
	SetLevel(LvlPanic)
	l, result := logAndGetLevelsLogged()
	// Nothing should have been logged
	if l.trace || l.debug || l.info || l.warn || l.err {
		t.Fatalf("All levels should have been logged! Result:\n%s", result)
	}
}

func TestPanic(t *testing.T) {
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	log.SetOutput(buffer)
	defer func() {
		// Panic handler (panic is expected)
		if r := recover(); r != nil {
			result := buffer.String()
			if !strings.Contains(result, "PANIC") {
				t.Fatalf("Log contains no PANIC entry")
			}
		}
	}()
	Panic("Panic here")
	t.Fatalf("The code did not panic")
}

func TestSetFile(t *testing.T) {
	SetLevel(LvlInfo)
	err := SetFile("tmplog.txt", 100)
	if err != nil {
		t.Fatalf("Unable to log to file. Reason: %s", err)
	}
	Info("Hello")

	// Cleanup
	log.SetOutput(os.Stderr)
	globFile.Sync()
	globFile.Close()
	globFile = nil
	err = os.Remove("tmplog.txt")
	if err != nil {
		t.Fatalf("Unable to cleanup. Reason: %s", err)
	}
}

func TestSetInvalidFile(t *testing.T) {
	err := SetFile("thispathdontexit/log.txt", 100)
	if err == nil {
		t.Fatalf("Invalid file shall give an error")
	}
}

func fileExist(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func fileSizeKB(fileName string) int {
	file, err := os.Open(fileName)
	if err != nil {
		panic("Unable to open file")
	}
	defer file.Close()
	info, _ := file.Stat()
	return int(info.Size() / 1024)
}

func fileModTime(fileName string) time.Time {
	file, err := os.Open(fileName)
	if err != nil {
		panic("Unable to open file")
	}
	defer file.Close()
	info, _ := file.Stat()
	return info.ModTime()
}

func TestLogWrap(t *testing.T) {
	logFileName := "tmplog.txt"
	backupFileName := logFileName + ".1"
	// Remove old files if exist
	os.Remove(logFileName)
	os.Remove(backupFileName)
	// Check that no backup file exist
	if fileExist(backupFileName) {
		t.Fatalf("Unable to remove backup file")
	}
	SetLevel(LvlInfo)
	err := SetFile(logFileName, 7) // Max size 7 KB
	if err != nil {
		t.Fatalf("Unable to log to file. Reason: %s", err)
	}

	nbrOfWraps := 0
	lastModTime := time.Now()
	for i := 0; i < 700; i++ {
		Info("entry %d", i)
		if fileSizeKB(logFileName) > 8 {
			t.Fatalf("Log file is exceeding 7 + 1 KB")
		}
		if fileExist(backupFileName) && lastModTime != fileModTime(backupFileName) {
			lastModTime = fileModTime(backupFileName)
			nbrOfWraps++
			fileSize := fileSizeKB(backupFileName)
			if fileSize < 7 || fileSize > 8 {
				t.Fatalf("Wrong backup file size (expected 7-8 KB): %d", fileSize)
			}
		}
	}
	t.Logf("During test it was %d number of wraps", nbrOfWraps)
	if nbrOfWraps < 3 {
		t.Fatalf("To few wraps: %d", nbrOfWraps)
	}
	// Cleanup
	log.SetOutput(os.Stderr)
	globFile.Sync()
	globFile.Close()
	globFile = nil
	os.Remove(logFileName)
	os.Remove(backupFileName)
}
