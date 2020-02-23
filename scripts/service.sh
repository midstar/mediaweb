#!/bin/sh
# MediaWEB systemd service (un)installation shell script
# 
# See print_usage below for instructions
#

set -e

# Global configurations
export CONFIGDESTINATION=/etc/mediaweb.conf
export EXEDESTINATION=/usr/sbin/mediaweb
export SERVICECONFIG=/etc/systemd/system/mediaweb.service


print_usage() {
  echo "Usage:"
  echo
  echo "For MediaWEB installation or upgrade (keep config queried):"
  echo
  echo "  sudo sh service.sh install [<exefile>] [<conffile>]"
  echo
  echo "     exefile:  mediaweb executable path (default ./mediaweb)"
  echo "     conffile: mediaweb config file path (default ./mediaweb.conf)"
  echo
  echo "For MediaWEB uninstallation:"
  echo
  echo "  sudo sh service.sh uninstall [purge]"
  echo
  echo "     purge: If purge the mediaweb configuration is also removed"
  echo
  exit 1
}


disable_mediaweb_service() {
  systemctl is-active --quiet mediaweb && {
    echo "Stopping MediaWEB service"
    systemctl stop mediaweb || {
      echo "ERROR! Unable to stop MediaWEB service"
      echo ""
      exit 1
    }
  }

  systemctl is-enabled --quiet mediaweb 2> /dev/null && {  
    echo "Disabling MediaWEB service"   
    systemctl disable mediaweb || {
      echo "ERROR! Unable to disable MediaWEB service"
      echo ""
      exit 1
    }
  }
  return 0
}

install_service() {

  echo "--------------------------------------------------"
  echo "Installation of MediaWEB service on Linux         "
  echo "--------------------------------------------------"
  echo

  #########################################################
  # Check optional arguments
  #########################################################
  if ! [ -z $1 ]; then
    export EXE=$1
  else
    export EXE=mediaweb
  fi

  if ! [ -z $2 ]; then
    export CONFIG=$2
  else
    export CONFIG=mediaweb.conf
  fi

  #########################################################
  # Check for prerequisities
  #########################################################

  if ! [ -x "$(command -v systemctl)" ]; then
    echo "ERROR! systemctl is not installed."
    echo
    echo "To install MediaWEB as a service using this script the system"
    echo "needs to support systemd services. Please manually install as"
    echo "an init/init.d service or similar."
    echo
    exit 1
  fi

  if ! [ -d "/etc/systemd/system" ]; then
    echo "Unable to access /etc/systemd/system"
    echo
    exit 1
  fi

  if [ ! -f "$EXE" ]; then
		echo "ERROR! MediaWEB executable, $EXE, not found"
		echo
		echo "You need to run this script in the path where you have the"
		echo "MediaWEB executable or provide as argument"
		print_usage
  fi

  if [ ! -f "$CONFIG" ]; then
		echo "ERROR! MediaWEB configuration, $CONFIG, not found"
		echo
		echo "You need to run this script in the path where you have the"
		echo "MediaWEB configuration or provide as argument"
		print_usage
  fi

	disable_mediaweb_service

  #########################################################
  # Copy mediaweb configuration to /etc/
  #########################################################

  export KEEP=n
  if [ -f "$CONFIGDESTINATION" ]; then
    echo "There is already a configuration '$CONFIGDESTINATION'"
    echo
    echo "Do you wan't to keep this file?"
    read -p "Keep? (y/n): " KEEP
  fi

  if ! [ "$KEEP" = "y" ]; then

    # Copy configuration file
    cp $CONFIG $CONFIGDESTINATION || {
      echo "ERROR! Unable to copy $CONFIG to $CONFIGDESTINATION"
      echo
      exit 1
    }
    echo "Copied : $CONFIGDESTINATION"

  fi

  #########################################################
  # Copy the mediaweb executable to /usr/sbin
  #########################################################

  cp $EXE $EXEDESTINATION || {
    echo "ERROR! Unable to copy $EXE to $EXEDESTINATION"
    echo
    exit 1
  }
  echo "Copied : $EXEDESTINATION"


  #########################################################
  # Create webmedia.service systemd configuration
  #########################################################

  echo "[Unit]" > $SERVICECONFIG || {
    echo "ERROR! Unable to create $SERVICECONFIG"
    echo
    exit 1
  }
  echo "Description=MediaWEB service" >> $SERVICECONFIG
  echo "After=network.target" >> $SERVICECONFIG
  echo "StartLimitIntervalSec=0" >> $SERVICECONFIG
  echo "" >> $SERVICECONFIG
  echo "[Service]" >> $SERVICECONFIG
  echo "Type=simple" >> $SERVICECONFIG
  echo "Restart=always" >> $SERVICECONFIG
  echo "RestartSec=1" >> $SERVICECONFIG
  echo "User=root" >> $SERVICECONFIG
  echo "ExecStart=/usr/sbin/mediaweb" >> $SERVICECONFIG
  echo "" >> $SERVICECONFIG
  echo "[Install]" >> $SERVICECONFIG
  echo "WantedBy=multi-user.target" >> $SERVICECONFIG

  echo "Wrote  : $SERVICECONFIG"

  #########################################################
  # Start and enable webmedia service
  #########################################################
	
	systemctl enable mediaweb || {
    echo "ERROR! Unable to enable MediaWEB service"
    echo
    exit 1
	}
	
	systemctl start mediaweb || {
    echo "ERROR! Unable to start MediaWEB service"
    echo
    exit 1
	}

  if ! [ -x "$(command -v ffmpeg)" ]; then
    echo
    echo "Warning! ffmpeg is not installed."
    echo ""
    echo "To generate video thumbnails ffmpeg needs to be installed."
    echo ""
  fi

  echo
  echo "-------------------------------------------"
  echo "MediaWEB installation done"
  echo ""
  echo "Edit settings in /etc/mediaweb.conf and"
  echo "restart the service by running:"
  echo "   sudo service mediaweb restart"
  echo "-------------------------------------------"
}

remove_file() {
	if [ -f $1 ]; then 
	  echo "Removing: $1"
	  rm $1 || {
	    echo "WARNING! Unable to remove $1"
			echo
	  }
	fi
}

uninstall_service() {

  echo "--------------------------------------------------"
  echo "Uninstallation of MediaWEB service                "
  echo "--------------------------------------------------"
  echo
	
	disable_mediaweb_service

	remove_file $SERVICECONFIG
	remove_file $EXEDESTINATION

  if [ "$1" = "purge" ]; then
    remove_file $CONFIGDESTINATION
  else
    echo "Keeping $CONFIGDESTINATION (remove with purge argument)"
  fi

	echo
	echo "Uninstallation complete!"
	echo
}



if [ "$(whoami)" != "root" ]; then
  echo "You need to run this script with root priviledges:"
  echo
  echo "  sudo sh service.sh ..."
  echo
  exit 1
fi

if [ -z $1 ]; then
  print_usage
elif [ "$1" = "install" ]; then
  install_service $2 $3
elif  [ "$1" = "uninstall" ]; then
  uninstall_service $2
else
  echo "ERROR! Unknown command '$1'"
  echo
  print_usage
fi
