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
REM packr2 has a bug in Windows where absolute paths are generated instead
REM of relative paths. We use --legacy to fix this
REM set PACKRCMD=packr2 --legacy
REM echo %PACKRCMD%
REM %PACKRCMD%
set INSTALLCMD=packr2 build -ldflags="-X 'main.applicationBuildTime=%DATE% %TIME%' -X main.applicationVersion=%VERSION% -X main.applicationGitHash=%GITHASH%" github.com/midstar/mediaweb
echo %INSTALLCMD%
%INSTALLCMD%
