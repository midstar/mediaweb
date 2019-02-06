# MediaWEB - Access your photos and videos over Internet

[![Go Report Card](https://goreportcard.com/badge/github.com/midstar/mediaweb)](https://goreportcard.com/report/github.com/midstar/mediaweb)
[![AppVeyor](https://ci.appveyor.com/api/projects/status/github/midstar/mediaweb?svg=true)](https://ci.appveyor.com/api/projects/status/github/midstar/mediaweb)
[![Coverage Status](https://coveralls.io/repos/github/midstar/mediaweb/badge.svg?branch=master)](https://coveralls.io/github/midstar/mediaweb?branch=master)

MediaWeb is a small self-contained web server software to enable you to access your photos and videos over Internet in your WEB browser.

The main design goal of MediaWEB is that no additional dependencies shall be needed, such as a server (Apache, NGINX, IIS etc.) or database (MySQL, sqlite etc.). The only files required to run MediaWEB are:

* The mediaweb executable
* A configuration file, mediaweb.conf

No additional stuff, such as dockers and similar is required. 

MediaWEB is well suited to run on small platforms such as Raspberry Pi, Banana Pi, ROCK64 and similar. It is still very fast and can be used with advantage on PC:s running Windows, Linux or Mac OS.

## Features

* Simple WEB GUI for viewing your images and videos
* Thumbnail support, primary by reading of EXIF thumbnail if it exist, otherwise thumbnails will be created and stored in a thumbnail cache
* Automatic rotation JPEG images when needed (based on EXIF information)

## Install

TBD

## Configuration

TBD

## Author and license

This application is written by Joel Midstjärna and is licensed under the MIT License.