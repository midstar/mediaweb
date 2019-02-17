@echo off
REM Test of service.bat - both installation and uninstallation

pushd %GOPATH%\src\github.com\midstar\mediaweb\

REM Move files to a temporary folder to secure that packr 
REM (embedded resources) works as expected
if not exist "tmpout" mkdir "tmpout" || exit \b 1
mkdir tmpout\servicetest
copy mediaweb.exe tmpout\servicetest\mediaweb.exe  || exit \b 1
copy configs\mediaweb.conf tmpout\servicetest\mediaweb.conf || exit \b 1
cd tmpout\servicetest || exit \b 1

set SCRIPTPATH=%GOPATH%\src\github.com\midstar\mediaweb\scripts

call %SCRIPTPATH%\service.bat install || exit \b 1

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

call %SCRIPTPATH%\service.bat uninstall || exit \b 1

if not exist mediaweb.log (
	echo Fail: mediaweb.log was not created!
	exit \b 1
)

echo Test passed!

popd