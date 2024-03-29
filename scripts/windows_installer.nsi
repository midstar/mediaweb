; MediaWEB installer creator NSIS script
;
; Prerequisities:
;  - GOPATH environment variable needs to be set correctly
;  - mediaweb needs to be built (go build github.com/midstar/mediaweb)
;    and the exe must be in the mediaweb folder
;
; Usage:
;  - makensis -DVERSION=<version> windows_installer.nsi
;     (<version> should be in the format 1.1.1.1)
;
; The installer will be put in GOPATH\src\github.com\midstar\mediaweb folder
;    
;
;-------------------------------------------------

;--------------------------------
;External dependencies / libraries

; Use the NSIS Modern UI 2
!include MUI2.nsh
!include x64.nsh

;--------------------------------
;Common definitions

!define APPLICATION_NAME "MediaWEB"
!define APPLICATION_FOLDER "MediaWEB"
!define APPLICATION_SOURCE "$%GOPATH%\src\github.com\midstar\mediaweb"
!define APPLICATION_BINARY "$%GOPATH%\src\github.com\midstar\mediaweb"

; The application version. Override with /DVERSION flag
!ifndef VERSION
!define VERSION "0.0.0.0"
!endif

; The name of the installer
Name "${APPLICATION_NAME} ${VERSION}"

; The file to write
OutFile "${APPLICATION_SOURCE}\mediaweb_windows_x64_setup.exe"

; The default installation directory
InstallDir $PROGRAMFILES64\${APPLICATION_FOLDER}

; Registry key to check for directory (so if you install again, it will 
; overwrite the old one automatically)
InstallDirRegKey HKLM "Software\${APPLICATION_FOLDER}" "Install_Dir"

; Request application privileges
RequestExecutionLevel admin

;--------------------------------
;Interface Settings

!define MUI_ABORTWARNING
!define MUI_ICON "${APPLICATION_SOURCE}\testmedia\logo.ico"

;--------------------------------
;Pages

!insertmacro MUI_PAGE_LICENSE "${APPLICATION_SOURCE}\LICENSE.txt"
!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_DIRECTORY
Page custom selectMediaPathPage selectMediaPathPageLeave
!insertmacro MUI_PAGE_INSTFILES
!define MUI_FINISHPAGE_RUN
!define MUI_FINISHPAGE_RUN_FUNCTION "LaunchLink"
!define MUI_FINISHPAGE_RUN_TEXT "Launch MediaWEB User Interface"
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

;--------------------------------
;Languages
 
!insertmacro MUI_LANGUAGE "English"

;--------------------------------
;Version Information

VIProductVersion "${VERSION}"
VIAddVersionKey /LANG=${LANG_ENGLISH} "ProductName" "${APPLICATION_NAME}"
VIAddVersionKey /LANG=${LANG_ENGLISH} "Comments" "Share photos and videos on Internet"
VIAddVersionKey /LANG=${LANG_ENGLISH} "CompanyName" "Joel Midstjarna"
VIAddVersionKey /LANG=${LANG_ENGLISH} "LegalTrademarks" "-"
VIAddVersionKey /LANG=${LANG_ENGLISH} "LegalCopyright" "Copyright Joel Midstjarna"
VIAddVersionKey /LANG=${LANG_ENGLISH} "FileDescription" "${APPLICATION_NAME} Setup"
VIAddVersionKey /LANG=${LANG_ENGLISH} "FileVersion" "${VERSION}"
VIAddVersionKey /LANG=${LANG_ENGLISH} "ProductVersion" "${VERSION}"

;-----------------------------------------------------------------------------
; Init function - executed before the installation starts
Function .onInit
  
  ;---------------------------------------------------------------------------
  ; Check if already installed 
 
  ReadRegStr $R0 HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPLICATION_FOLDER}"  "UninstallString"
  StrCmp $R0 "" noPreviousInstaller
 
  MessageBox MB_OKCANCEL|MB_ICONEXCLAMATION "${APPLICATION_NAME} is already installed. $\n$\n\
    Click `OK` to remove the previous version or `Cancel` to cancel this upgrade."  IDOK uninst
  Abort
 
  ;Run the uninstaller
  uninst:
     ClearErrors
      Exec $R0
   
  noPreviousInstaller:

FunctionEnd

;======================================================================================================================
; Custom dialog Select Media Path

Var Dialog
Var TextMediaPath

Function selectMediaPathPage
  !insertmacro MUI_HEADER_TEXT "Media Directory" "Provide photo and video directory to share."

  #Create Dialog and quit if error
  nsDialogs::Create 1018
  Pop $Dialog
  ${If} $Dialog == error
          Abort
  ${EndIf}       

  ${NSD_CreateGroupBox} 5% 16u 90% 34u "Media directory"
  Pop $0

    ReadEnvStr $0 USERPROFILE
    ${NSD_CreateDirRequest} 15% 30u 49% 12u "$0\Pictures"
    Pop $TextMediaPath

    ${NSD_CreateBrowseButton} 65% 30u 20% 12u "Browse..."
    Pop $0
    ${NSD_OnClick} $0 OnDirBrowse

  ${NSD_CreateLabel} 5% 70u 100% 12u "Note! You can always change this later by updating mediaweb.conf"
  Pop $0

  nsDialogs::Show
FunctionEnd

Function OnDirBrowse
  ${NSD_GetText} $TextMediaPath $0
  nsDialogs::SelectFolderDialog "Select Media Directory" "$0" 
  Pop $0
  ${If} $0 != error
      ${NSD_SetText} $TextMediaPath "$0"
  ${EndIf}
FunctionEnd

Function selectMediaPathPageLeave
    ${NSD_GetText} $TextMediaPath $0
FunctionEnd

;======================================================================================================================
; Application install section
Section "${APPLICATION_NAME}" SectionMain

  SectionIn RO
  
  ; Set output path to the installation directory.
  SetOutPath $INSTDIR
  
  ; Create configuration file
  FileOpen $4 "$INSTDIR\mediaweb.conf" w

  FileWrite $4 "# Server network port.$\r$\n"
  FileWrite $4 "# This parameter is MANDATORY$\r$\n"
  FileWrite $4 "port = 9834$\r$\n"
  FileWrite $4 "$\r$\n"
  FileWrite $4 "# Network inteface to listen to.$\r$\n"
  FileWrite $4 "# If parameter is not set, the server will listen to$\r$\n"
  FileWrite $4 "# all interfaces.$\r$\n"
  FileWrite $4 "#ip = 127.0.0.1$\r$\n"
  FileWrite $4 "$\r$\n"
  FileWrite $4 "# Media path, i.e. where is your media located$\r$\n"
  FileWrite $4 "# This parameter is MANADTORY$\r$\n"
  FileWrite $4 "mediapath = $0$\r$\n"
  FileWrite $4 "$\r$\n"
  FileWrite $4 "# Cache path is by default your operating systems$\r$\n"
  FileWrite $4 "# temp folder + mediaweb. Cache path is where$\r$\n"
  FileWrite $4 "# thumbnails and preview images are stored.$\r$\n"
  FileWrite $4 "#$\r$\n"
  FileWrite $4 "# Is not allowed to be the same as mediapath.$\r$\n"
  FileWrite $4 "#$\r$\n"
  FileWrite $4 "#cachepath =$\r$\n"
  FileWrite $4 "$\r$\n"
  FileWrite $4 "# Thumbnail cache is on by default$\r$\n"
  FileWrite $4 "#enablethumbcache = off$\r$\n"
  FileWrite $4 "$\r$\n"
  FileWrite $4 "# Generate thumbs on startup is by default off. Uncomment$\r$\n"
  FileWrite $4 "# below to generate thumbs every time Media WEB startup.$\r$\n"
  FileWrite $4 "#genthumbsonstartup = on$\r$\n"
  FileWrite $4 "$\r$\n"
  FileWrite $4 "# Watch media path for updates is by default on.$\r$\n"
  FileWrite $4 "# Uncomment below to don't generate new thumbs for files$\r$\n"
  FileWrite $4 "# that are added in the media path$\r$\n"
  FileWrite $4 "#genthumbsonadd = off$\r$\n"
  FileWrite $4 "$\r$\n"
  FileWrite $4 "# Auto rotate of JPEG is by default on. Uncomment below$\r$\n"
  FileWrite $4 "# to disable auto rotate of JPEG.$\r$\n"
  FileWrite $4 "#autorotate = off$\r$\n"
  FileWrite $4 "$\r$\n"
  FileWrite $4 "# Resize images before providing them to the client. The$\r$\n"
  FileWrite $4 "# resized images are cached in the same location as the$\r$\n"
  FileWrite $4 "# thumbnails.$\r$\n"
  FileWrite $4 "#$\r$\n"
  FileWrite $4 "# The advantages are:$\r$\n"
  FileWrite $4 "#  1. Lower network bandwith required$\r$\n"
  FileWrite $4 "#  2. Smoother navigation at the client. Particular if$\r$\n"
  FileWrite $4 "#     browsing the images using a mobile client.$\r$\n"
  FileWrite $4 "#$\r$\n"
  FileWrite $4 "# Disadvantages are:$\r$\n"
  FileWrite $4 "#  1. Slower response time to view the image first time$\r$\n"
  FileWrite $4 "#     since resizing image might take several seconds.$\r$\n"
  FileWrite $4 "#  2. Increased cache storage required.$\r$\n" 
  FileWrite $4 "#$\r$\n"
  FileWrite $4 "# Previews are default off.$\r$\n"
  FileWrite $4 "#enablepreview = on$\r$\n"
  FileWrite $4 "$\r$\n"
  FileWrite $4 "# Max size of preview images in pixels. The image will$\r$\n"
  FileWrite $4 "# be resized so that width and hight is not larger than$\r$\n"
  FileWrite $4 "# this value.$\r$\n"
  FileWrite $4 "#previewmaxside = 1280$\r$\n"
  FileWrite $4 "$\r$\n"
  FileWrite $4 "# Generate preview images on startup is by default off. Uncomment$\r$\n"
  FileWrite $4 "# below to generate preview every time Media WEB startup.$\r$\n"
  FileWrite $4 "#$\r$\n"
  FileWrite $4 "# Warning! A lot of cache space might be required$\r$\n"
  FileWrite $4 "#genpreviewonstartup = on$\r$\n"
  FileWrite $4 "$\r$\n"
  FileWrite $4 "# Watch media path for updates is by default on.$\r$\n"
  FileWrite $4 "# Uncomment below to don't generate new image previews for$\r$\n"
  FileWrite $4 "# files that are added in the media path$\r$\n"
  FileWrite $4 "#genpreviewonadd = off$\r$\n"
  FileWrite $4 "$\r$\n"
  FileWrite $4 "# Remove unnecessary files from cache is by default off.$\r$\n"
  FileWrite $4 "# Uncomment below to remove cache files for media files$\r$\n"
  FileWrite $4 "# that has been removed.$\r$\n"
  FileWrite $4 "#enablecachecleanup = on$\r$\n"
  FileWrite $4 "$\r$\n"
  FileWrite $4 "logfile = mediaweb.log$\r$\n"
  FileWrite $4 "$\r$\n"
  FileWrite $4 "# Log level is 'info' by default$\r$\n"
  FileWrite $4 "#loglevel = trace$\r$\n"
  FileWrite $4 "$\r$\n"
  FileWrite $4 "# User name and password for authentication. Leave commented$\r$\n"
  FileWrite $4 "# for no authentication$\r$\n"
  FileWrite $4 "#username = myusername$\r$\n"
  FileWrite $4 "#password = mypassword$\r$\n"
  FileWrite $4 "$\r$\n"
  FileWrite $4 "# TLS (HTTPS) certification file and key file. Leave commented\r$\n"
  FileWrite $4 "# for no encryption (HTTP). If both parameters are set TlS\r$\n"
  FileWrite $4 "# will be enabled.\r$\n"
  FileWrite $4 "#tlscertfile = public.crt\r$\n"
  FileWrite $4 "#tlskeyfile = private.key\r$\n"
  FileWrite $4 "$\r$\n"

  FileClose $4

  ; Copy mediaweb binary
  File "${APPLICATION_BINARY}\mediaweb.exe"
	
	; Copy mediaweb URL
  File "${APPLICATION_SOURCE}\MediaWEB.url"
	
	; Copy mediaweb icon
  File "${APPLICATION_SOURCE}\testmedia\logo.ico"
	
	
  
  ; Write the installation path into the registry
  WriteRegStr HKLM SOFTWARE\${APPLICATION_FOLDER} "Install_Dir" "$INSTDIR"
  
  ; Write the version into the registry
  WriteRegStr HKLM SOFTWARE\${APPLICATION_FOLDER} "Version" "${VERSION}"
  

  ; Write the uninstall keys for Windows
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPLICATION_FOLDER}" "DisplayName" "${APPLICATION_NAME} ${VERSION}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPLICATION_FOLDER}" "Publisher" "Joel Midstjarna"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPLICATION_FOLDER}" "UninstallString" '"$INSTDIR\uninstall.exe"'
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPLICATION_FOLDER}" "DisplayIcon" "$\"$INSTDIR\logo.ico$\""
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPLICATION_FOLDER}" "NoModify" 1
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPLICATION_FOLDER}" "NoRepair" 1
  WriteUninstaller "uninstall.exe"
  
	;---------------------------------------------------------------
  ; Install and start the windows service
	ClearErrors
  ExecWait "sc create mediaweb binpath= $\"\$\"$INSTDIR\mediaweb.exe\$\" \$\"$INSTDIR\$\"$\" start= auto DisplayName= $\"MediaWEB$\""
  IfErrors 0 createOk
  MessageBox MB_OK|MB_ICONSTOP "Unable to install mediaweb as a service."
  Goto endService
	
	createOk:
  ExecWait "sc description mediaweb $\"MediaWEB Service$\""
  IfErrors 0 descriptionOk
  MessageBox MB_OK|MB_ICONSTOP "Unable to add description to mediaweb service."
  Goto endService
	
	descriptionOk:
  ExecWait "sc start mediaweb"
  IfErrors 0 endService
  MessageBox MB_OK|MB_ICONSTOP "Unable to start mediaweb service."
  Goto endService
	
	endService:
	
SectionEnd


;======================================================================================================================
; Start menu shortcuts install section (can be disabled by the user)
Section "Start Menu Shortcuts" SectionStartMenu

  CreateDirectory "$SMPROGRAMS\${APPLICATION_FOLDER}"
  CreateShortcut "$SMPROGRAMS\${APPLICATION_FOLDER}\Uninstall.lnk" "$INSTDIR\uninstall.exe" "" "$INSTDIR\logo.ico" 0
  CreateShortcut "$SMPROGRAMS\${APPLICATION_FOLDER}\MediaWEB.lnk" "$INSTDIR\MediaWEB.url" "" "$INSTDIR\logo.ico" 0
	
SectionEnd


;======================================================================================================================
; Description of the sections
!insertmacro MUI_FUNCTION_DESCRIPTION_BEGIN
	!insertmacro MUI_DESCRIPTION_TEXT ${SectionMain} "Install and start ${APPLICATION_NAME}."
	!insertmacro MUI_DESCRIPTION_TEXT ${SectionStartMenu} "Create Shortcuts on Start Menu."
!insertmacro MUI_FUNCTION_DESCRIPTION_END


;======================================================================================================================
; Uninstaller section
Section "Uninstall"

  ; --------------------------------------------------------------------------  
  ; Uninstall and stop mediaweb service
  execWait "sc stop mediaweb"
	execWait "sc delete mediaweb"
 
  ; Remove registry keys
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPLICATION_FOLDER}"
  DeleteRegKey HKLM SOFTWARE\${APPLICATION_FOLDER}

  ; Remove shortcuts, if any
  Delete "$SMPROGRAMS\${APPLICATION_FOLDER}\*.*"
  RMDir "$SMPROGRAMS\${APPLICATION_FOLDER}"
	
	; Remove installation directory
	RMDir /r $INSTDIR\*

SectionEnd

;======================================================================================================================
; Helper functions

Function LaunchLink
  ExecShell "" "$INSTDIR\MediaWEB.url"
FunctionEnd
