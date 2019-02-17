@echo off

REM MediaWEB windows service (un)installation batch script

if [%1] EQU [install] (

  if not exist mediaweb.exe (
    echo ERROR! mediaweb.exe needs to be in current directory
    exit /b 1
  )

  if not exist mediaweb.conf (
    echo ERROR! mediaweb.conf needs to be in current directory
    exit /b 1
  )
  
  echo ------------------------------------------
  echo Installing MediaWEB windows service
  echo ------------------------------------------
  sc create mediaweb binpath="\"%cd%\mediaweb.exe\" \"%cd%\"" start= auto DisplayName= "MediaWEB" || (
      echo ERROR! Unable to create mediaweb service
      exit /b 1
  )
  sc description mediaweb "MediaWEB Service" || (
      echo ERROR! Unable to create description for mediaweb service
      echo.
      echo Make sure you are running cmd.exe as an administrator
      exit /b 1
  )
  sc start mediaweb || (
      echo ERROR! Unable to start mediaweb service
      exit /b 1
  )
  sc query mediaweb || (
      echo ERROR! Unable to query mediaweb service
      exit /b 1
  )

  echo MediaWEB service successfully installed!
  exit /b 0

) else if [%1] EQU [uninstall] (

  echo ------------------------------------------
  echo Uninstalling MediaWEB windows service
  echo ------------------------------------------

  sc stop mediaweb
  sc delete mediaweb || (
      exit /b 1
  )

  echo Uninstallation complete!
  exit /b 0

) else if [%1] EQU [] (
  call :print_usage
) else (
  echo ERROR! Unknown command '%1'
  call :print_usage
)

exit /b 0

:print_usage
  echo.
  echo Usage:
  echo.
  echo.NOTE! Start cmd.exe with administrator privileges.
  echo.
  echo.Update mediaweb.conf before installation.
  echo.
  echo For MediaWEB service installation:
  echo.
  echo   service.bat install
  echo.
  echo.
  echo For MediaWEB service uninstallation:
  echo.
  echo   sudo sh service.sh uninstall
  echo.
  exit /b 1
