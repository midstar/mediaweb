@echo off
REM Builds plm using the correct version, build time and git hash
if [%1] EQU [] (
  echo No version provided. Using unofficial.
  echo VERSION=unofficial
) else (
  set VERSION=%1
)
echo version: %1
for /f %%i in ('git rev-parse HEAD') do set GITHASH=%%i
echo git hash: %GITHASH%
echo building / installing
pushd %GOPATH%\src\github.com\midstar\mediaweb
set PACKRCMD=packr2
echo %PACKRCMD%
%PACKRCMD%
popd
set INSTALLCMD=go install -ldflags="-X 'main.applicationBuildTime=%DATE% %TIME%' -X main.applicationVersion=%1 -X main.applicationGitHash=%GITHASH%" github.com/midstar/mediaweb
echo %INSTALLCMD%
%INSTALLCMD%
