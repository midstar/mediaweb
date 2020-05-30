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

## Content

- [Screenshots](#screenshots)
- [Features](#features)
- [Download and install Linux](#download-and-install-linux)
- [Download and install on Windows](#download-and-install-on-windows)
- [Build from source on any platform](#build-from-source-on-any-platform)
- [Configuration guide](#configuration-guide)
- [Future improvements](#future-improvements)
- [Author and license](#author-and-license)
- [FAQ](FAQ.md)

## Screenshots

Browse your media.

To increase/decrease size of thumbnails do one of following:

* Touch: Use two fingers and pitch
* Keboard: <kbd>+</kbd> (increase) or <kbd>-</kbd> (decrease) 

![browser](testmedia/screenshot_browser.jpg)

View your images and play your videos

![viewer](testmedia/screenshot_viewer.jpg)

To navigate between the images do one of following:

* Touch: Swipe right (previous) or left (next)
* Mouse: Drag right (previous) or left (next)
* Keyboard: <kbd><</KBD> (previous), <kbd>></KBD> (next), <kbd>Page Up</KBD> (first) or <kbd>Page Down</KBD> (last)

![viewer](testmedia/screenshot_viewer_swipe.jpg)

To zoom image:

* Touch: Use two fingers and pitch alternatively double tap to zoom x2
* Mouse: Scroll up (zoom in), down (zoom out) alternatively double click to zoom x2

![viewer](testmedia/screenshot_viewer_zoom.jpg)

Interface suited for mobile devices

![viewer](testmedia/screenshot_mobile.jpg)



## Features

* Simple WEB GUI for viewing your images and videos both on your PC and on your mobile phone
* Thumbnail support for images and videos, primary by reading of EXIF thumbnail if it exist, otherwise thumbnails will be created and stored in a thumbnail cache. Video thumbnails requires [ffmpeg](https://www.ffmpeg.org/) to be installed
* Automatic rotation JPEG images when needed (based on EXIF information)
* Generate thumbnails on the fly, on start-up and/or when new files are added to the media directory
* **NEW!** Automatically resize images to reduce network bandwidth and get a smoother navigation at client
* Optional authentication with username and password

## Download and install Linux

### Debian based distros (Debian, Ubuntu, Raspbian, Armbian etc.)

Debian installer packages are provided for x64 (PC) or arm (Raspberry Pi, Banana Pi, ROCK64 etc.) [here on GitHub](https://github.com/midstar/mediaweb/releases).

To get the latest file from a shell follow the instruction below.

For PC x64 based architectures write following in a shell:

    export MW_ARCH=x64

For ARM 32-bit based architectures (for example Raspberry Pi 1, 2):

    export MW_ARCH=arm

For ARM 64-bit based archtectures (for example Raspberry Pi 3, 4 or ROCK64):

    export MW_ARCH=arm64

Then run following for all Linux platforms:

    mkdir ~/mediaweb
    cd ~/mediaweb
    curl -s https://api.github.com/repos/midstar/mediaweb/releases/latest \
    | grep browser_download_url \
    | grep "mediaweb_linux_${MW_ARCH}.deb" \
    | cut -d : -f 2,3 \
    | tr -d \" \
    | wget -qi -
     sudo dpkg -i mediaweb_linux_${MW_ARCH}.deb

Uninstall by using the normal apt tools. For example:

    sudo apt remove mediaweb


### Other distros (RedHat, CentOS etc.)

For non-Debian based Linux distros, a custom installation / uninstallation script is provided in the linux .gz files [here on GitHub](https://github.com/midstar/mediaweb/releases).

To get the latest file from a shell follow the instruction below.

For PC x64 based Linux write following in a shell:

    export MW_ARCH=x64

For ARM 32-bit based architectures (for example Raspberry Pi 1, 2):

    export MW_ARCH=arm

For ARM 64-bit based archtectures (for example Raspberry Pi 3, 4 or ROCK64):

    export MW_ARCH=arm64

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

The service.sh script can also be used for uninstallation by running:

    sudo sh service.sh uninstall

### Configuration in Linux (any distro)

For video thumbnail support, install ffmpeg:

    sudo apt-get install ffmpeg

To perform additional configuration, edit:

    sudo vi /etc/mediaweb.conf

And then restart the MediaWEB service with:

    sudo service mediaweb restart

Also, checkout the [Configuration guide](#configuration-guide) and [FAQ](FAQ.md).

## Download and install on Windows

Download mediaweb_windows_x64_setup.exe [here on GitHub](https://github.com/midstar/mediaweb/releases).

Run the installer and follow the instructions.

To modify changes just edit mediaweb.conf in the installation directory and restart the mediaweb
service in task manager.

You need to install [ffmpeg](https://www.ffmpeg.org/) separately and put ffmpeg into your PATH to get video thumbnail support.

Also, checkout the [Configuration guide](#configuration-guide) and [FAQ](FAQ.md).

## Build from source on any platform

To build from source on any platform you need to:

* Install Golang (version 1.11 or newer)
* Set the GOPATH environment variable
* Add the bin folder within your GOPATH to your PATH environment variable ($GOPATH/bin)

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


## Configuration guide

See [mediaweb.conf](configs/mediaweb.conf) for the available configuration parameters.

The most important configuration parameter is *mediapath* which points out the 
directory of your media. 

Your might also want to change the *port* configuration parameter.

### Increase performance with preview / resized images

By default the server will provide the full images to the client (WEB browser).
This is generally not an issue if you view your images over your home network, but
over Internet or a mobile network it might take very long time to load the images.
Also, since the images are large, the user interface might feel slow and unresponsive
particulary in mobile web browsers.

To fix these issues, at the cost of more caching storage space, MediaWEB can reduce the
size / resolution and put the resized images in cache. The resizing procedure will
take quite much time, especially for SBCs. Therefore you should configure MediaWEB
to generate the previews at startup and when new files are added to your media
directory:

    enablethumbcache = on
    genthumbsonstartup = on
    genthumbsonadd = on
    enablepreview = on
    genpreviewonstartup = on
    genpreviewonadd = on

Count with ~300 - 800 KB per image, so you need to secure that the cache folder pointed
out with the *cachepath* has enough disk space.

Also, if you have many images and are running on an SBC it might take very, very long
time the first time the images are resized. Count with several days.


## Future improvements

* Add support for TLS/SSL


## Author and license

This application is written by Joel Midstjärna and is licensed under the MIT License.