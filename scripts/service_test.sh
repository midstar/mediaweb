#!/bin/sh
set -e

# Test of service.sh - both installation and uninstallation

export REL_PATH=`dirname "$0"`
export SCRIPT_PATH=`realpath $REL_PATH`
export MEDIAWEB_PATH=`realpath $REL_PATH/..`

cd $MEDIAWEB_PATH

# Move files to a temporary folder to secure that
# embedded resources works as expected
mkdir -p tmpout
mkdir -p tmpout/servicetest
cp mediaweb tmpout/servicetest/
cp configs/mediaweb.conf tmpout/servicetest/
cd tmpout/servicetest 

# Edit configuratino
sh $SCRIPT_PATH/conf_edit.sh mediaweb.conf mediapath $MEDIAWEB_PATH/testmedia
sh $SCRIPT_PATH/conf_edit.sh mediaweb.conf logfile /var/log/mediaweb.log

# Install mediaweb
sh $SCRIPT_PATH/service.sh install

# Thest that it works
echo "Waiting 2 seconds"
sleep 2
echo "Testing connection"
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:9834)
if ! [ "$HTTP_STATUS" = "200" ]; then
	echo "Test Failed! Unable to connect to MediaWEB."
	echo
	echo "Expected status code 200, but got $HTTP_STATUS"
	exit 1
fi
if ! [ -f "/var/log/mediaweb.log" ]; then
	echo "Test Failed! No log file was created in /var/log/mediaweb.log"
	echo
	exit 1
fi
sh $SCRIPT_PATH/service.sh uninstall purge
echo "Test passed :-)"