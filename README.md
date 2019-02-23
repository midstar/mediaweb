![GitHub Logo](/templates/logo.png)

# MediaWEB - Access your photos and videos over Internet

[![Go Report Card](https://goreportcard.com/badge/github.com/midstar/mediaweb)](https://goreportcard.com/report/github.com/midstar/mediaweb)
[![AppVeyor](https://ci.appveyor.com/api/projects/status/github/midstar/mediaweb?svg=true)](https://ci.appveyor.com/api/projects/status/github/midstar/mediaweb)
[![Coverage Status](https://coveralls.io/repos/github/midstar/mediaweb/badge.svg?branch=master)](https://coveralls.io/github/midstar/mediaweb?branch=master)

MediaWeb is a small self-contained web server software to enable you to access your photos and videos over the Internet in your WEB browser.

The main design goal of MediaWEB is that no additional dependencies shall be needed, such as a script engine (Java, Python, Perl, Ruby etc.), or server (Apache, NGINX, IIS etc.) or database (MySQL, sqlite etc.). The only files required to run MediaWEB are:

* The mediaweb executable
* A configuration file, mediaweb.conf

Optional dependencies are:

* [ffmpeg](https://www.ffmpeg.org/) for video thumbnail support

No additional stuff, such as dockers and similar is required. 

MediaWEB is well suited to run on small platforms such as Raspberry Pi, Banana Pi, ROCK64 and similar. It is still very fast and can be used with advantage on PC:s running Windows, Linux or Mac OS.

## Features

* Simple WEB GUI for viewing your images and videos
* Thumbnail support for images and videos, primary by reading of EXIF thumbnail if it exist, otherwise thumbnails will be created and stored in a thumbnail cache. Video thumbnails requires [ffmpeg](https://www.ffmpeg.org/) to be installed.
* Automatic rotation JPEG images when needed (based on EXIF information)
* Optional authentication with username and password

## Download and install Linux

For PC x64 based Linux write following in a shell:

    export MW_ARCH=x64

For ARM based Linux on for example Raspberry Pi, Banana Pi, ROCK64 etc:

    export MW_ARCH=arm

Then run following for all Linux platforms:

    mkdir ~/mediaweb
    cd ~/mediaweb
    curl -s https://api.github.com/repos/midstar/mediaweb/releases/latest \
    | grep browser_download_url \
    | grep "mediaweb_linux_${MW_ARCH}.tar.gz" \
    | cut -d : -f 2,3 \
    | tr -d \" \
    | wget -qi -
     tar xvzf mediaweb_linux_${MW_ARCH}.tar.gz
     sudo sh service.sh install

Follow the instructions in the service.sh script.

For video thumbnail support, install ffmpeg:

    sudo apt-get install ffmpeg

To perform additional configuration, edit:

    sudo vi /etc/mediaweb.conf

And then restart the MediaWEB service with:

    sudo systemctl restart mediaweb

To uninstall MediaWEB run:

    cd ~/mediaweb
    sudo sh service.sh uninstall


## Download and install on Windows (64bit)

Download mediaweb_windows_x64_setup.exe [here on GitHub](https://github.com/midstar/mediaweb/releases).

Run the installer and follow the instructions.

To modify changes just edit mediaweb.conf in the installation directory and restart the mediaweb
service in task manager.

You need to install [ffmpeg](https://www.ffmpeg.org/) separately and put ffmpeg into your PATH to get video thumbnail support.

## Build from source (any platform)

To build from source on any platform you need to:

* Install Golang 
* Set the GOPATH environment variable

On Windows execute (from cmd.exe):

    go get github.com/midstar/mediaweb
    cd %GOPATH%\src\github.com\midstar\mediaweb
    scripts\install_deps.bat
    scripts\build.bat

On Linux/Mac execute (from a shell):

    go get github.com/midstar/mediaweb
    cd $GOPATH/src/github.com/midstar/mediaweb
    sh scripts/install_deps.sh
    sh scripts/build.sh

The mediaweb executable and an example configuration file will be in 
$GOPATH/src/github.com/midstar/mediaweb. Edit the configuration file
and then run the mediaweb executable.

To install as a Windows service start cmd.exe in administrator mode and run:

    scripts\service.bat install

On Linux platforms execute following to install MediaWEB as a service:

    sudo sh scripts/service.sh install


## Future improvements

* Add support for TLS/SSL


## Author and license

This application is written by Joel Midstjärna and is licensed under the MIT License.