@echo off
REM Test of service.bat - both installation and uninstallation
cd %GOPATH%\src\github.com\midstar\mediaweb\
call scripts\service.bat install || exit \b 1
echo Waiting 2 seconds
timeout 2 > NUL
echo Testing connection
FOR /F "tokens=*" %%g IN ('curl -s -o /dev/null -w "%%{http_code}" http://localhost:9834') do (SET HTTP_STATUS=%%g)
if [%HTTP_STATUS%] NEQ [200] (
	echo Test Failed! Unable to connect to MediaWEB.
	echo.
	echo Expected status code 200, but got %HTTP_STATUS%
	exit /b 1
)
echo MediaWEB connected
call scripts\service.bat uninstall || exit \b 1
echo Test passed!