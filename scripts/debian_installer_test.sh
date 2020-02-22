#!/bin/sh
set -e

# Test of service.sh - both installation and uninstallation

if [ "$(whoami)" != "root" ]; then
  echo "You need to run this script with root priviledges:"
  echo
  echo "  sudo sh debian_installer_test.sh ..."
  echo
  exit 1
fi

export REL_PATH=`dirname "$0"`
export SCRIPT_PATH=`realpath $REL_PATH`
export MEDIAWEB_PATH=`realpath $REL_PATH/..`

export INSTALLER_NAME=mediaweb_linux_x64.deb

cd $MEDIAWEB_PATH

if [ ! -f $INSTALLER_NAME ]; then
    echo "ERROR! MediaWEB Debian installer $INSTALLER_NAME, found"
    exit 1
fi

systemctl is-enabled --quiet mediaweb 2> /dev/null && {  
    echo "ERROR! MediaWEB is already installed. Uninstall be fore this test!!!"   
    exit 1
}

# Install mediaweb
echo "Installing $INSTALLER_NAME"
dpkg -i $INSTALLER_NAME

echo "Waiting 2 seconds"
sleep 2

systemctl is-active --quiet mediaweb || {
	echo "ERROR! MediaWEB service is not started!"
	exit 1
}
echo "MediaWEB started"

# Edit configuration
echo "Edit configuration"
sh $SCRIPT_PATH/conf_edit.sh /etc/mediaweb.conf mediapath $MEDIAWEB_PATH/testmedia
sh $SCRIPT_PATH/conf_edit.sh /etc/mediaweb.conf logfile /var/log/mediaweb.log

# Restart MediaWEB
echo "Restarting service"
service mediaweb restart

echo "Waiting 2 seconds"
sleep 2

systemctl is-active --quiet mediaweb || {
	echo "ERROR! MediaWEB service is not started after restart!"
	exit 1
}

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

#Unnstall
apt purge -y mediaweb
echo "Test passed :-)"