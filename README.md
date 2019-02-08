![GitHub Logo](/templates/logo.png)

# MediaWEB - Access your photos and videos over Internet

[![Go Report Card](https://goreportcard.com/badge/github.com/midstar/mediaweb)](https://goreportcard.com/report/github.com/midstar/mediaweb)
[![AppVeyor](https://ci.appveyor.com/api/projects/status/github/midstar/mediaweb?svg=true)](https://ci.appveyor.com/api/projects/status/github/midstar/mediaweb)
[![Coverage Status](https://coveralls.io/repos/github/midstar/mediaweb/badge.svg?branch=master)](https://coveralls.io/github/midstar/mediaweb?branch=master)

MediaWeb is a small self-contained web server software to enable you to access your photos and videos over the Internet in your WEB browser.

The main design goal of MediaWEB is that no additional dependencies shall be needed, such as a server (Apache, NGINX, IIS etc.) or database (MySQL, sqlite etc.). The only files required to run MediaWEB are:

* The mediaweb executable
* A configuration file, mediaweb.conf

No additional stuff, such as dockers and similar is required. 

MediaWEB is well suited to run on small platforms such as Raspberry Pi, Banana Pi, ROCK64 and similar. It is still very fast and can be used with advantage on PC:s running Windows, Linux or Mac OS.

## Features

* Simple WEB GUI for viewing your images and videos
* Thumbnail support, primary by reading of EXIF thumbnail if it exist, otherwise thumbnails will be created and stored in a thumbnail cache
* Automatic rotation JPEG images when needed (based on EXIF information)

## Install and configure

Download the binaries for your platform [here on GitHub](https://github.com/midstar/mediaweb/releases).

If the binaries don't exist for your platform, see "Build from source" below.

Update the mediapath setting in mediaweb.conf. You might also want to change the port setting.

Just start the mediaweb executable. It will look for the mediaweb.conf in the same folder.

## Build from source

To build from source on any platform you need to:

* Install Golang 
* Set the GOPATH environment variable

On Windows execute (from cmd.exe):

    cd %GOPATH%\src\github.com\midstar\mediaweb\scripts
    install_deps.bat
    build.bat

On Linux/Mac execute (from a shell):

    cd $GOPATH/src/github.com/midstar/mediaweb/scripts
    sh install_deps.sh
    sh build.sh

The mediaweb executable will be in the GOPATH/bin directory.


## Future improvements

* Create thumbnails for videos (probably using ffmpeg)
* Add support for TLS/SSL
* Add Windows installer (install MediaWEB as a service on Windows)

## Author and license

This application is written by Joel Midstjärna and is licensed under the MIT License.