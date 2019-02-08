@echo off
REM Builds MediaWEB using the correct version, build time and git hash
if [%1] EQU [] (
  echo No version provided. Using unofficial.
  set VERSION=unofficial
) else (
  set VERSION=%1
)
echo version: %VERSION%
for /f %%i in ('git rev-parse HEAD') do set GITHASH=%%i
echo git hash: %GITHASH%
echo building / installing
cd %GOPATH%\src\github.com\midstar\mediaweb
set PACKRCMD=packr2
echo %PACKRCMD%
%PACKRCMD%
type main-packr.go
set INSTALLCMD=go build -ldflags="-X 'main.applicationBuildTime=%DATE% %TIME%' -X main.applicationVersion=%VERSION% -X main.applicationGitHash=%GITHASH%" github.com/midstar/mediaweb
echo %INSTALLCMD%
%INSTALLCMD%
