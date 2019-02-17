# MediaWEB systemd service (un)installation shell script
# 
# For installation:
#  sh service.sh install [mediapath] [port]
# 
#    mediapath: Path to all videos and pictures. Will be
#               promted if not provided
#    port:      Network port number (default 9834)
#

# Global configurations
export CONFIG=/etc/mediaweb.conf
export EXEDESTINATION=/usr/sbin/mediaweb
export SERVICECONFIG=/etc/systemd/system/mediaweb.service


print_usage() {
  echo "Usage:"
  echo
  echo "For MediaWEB installation:"
  echo
  echo "  sudo sh service.sh install [mediapath] [port]"
  echo
  echo "    mediapath: Path to all videos and pictures. Will be"
  echo "               promted if not provided"
  echo "    port:      Network port number (default 9834)"
  echo
  echo
  echo "For MediaWEB uninstallation:"
  echo
  echo "  sudo sh service.sh uninstall"
  echo
  exit 1
}


install_service() {

  echo "--------------------------------------------------"
  echo "Installation of MediaWEB service on Linux         "
  echo "--------------------------------------------------"
  echo

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

  if [ ! -f "mediaweb" ]; then
		echo "ERROR! MediaWEB executable, mediaweb, not found"
		echo
		echo "You need to run this script in the path where you have the"
		echo "MediaWEB executable"
		exit 1
  fi

	# Check if service is already installed
  systemctl is-active --quiet mediaweb && {
	  echo "MediaWEB service is already installed."
		echo 
		echo "Uninstalling previous version..."
		echo
		
		systemctl disable mediaweb || {
			echo "ERROR! Unable to disable MediaWEB service"
			echo
			exit 1
		}
		
		systemctl stop mediaweb || {
			echo "ERROR! Unable to stop MediaWEB service"
			echo
			exit 1
		}
	}


  #########################################################
  # Create mediaweb.conf
  #########################################################

  export KEEP=n
  if [ -f "$CONFIG" ]; then
    echo "There is already a configuration '$CONFIG'"
    echo
    echo "Do you wan't to keep this file?"
    read -p "Keep? (y/n): " KEEP
  fi

  if ! [ "$KEEP" = "y" ]; then

    # Write a new configuration file

    if [ -z $1 ]; then
      echo "Please enter your media path (where your pictures and videos"
      echo "are located)"
      echo
      read -p "Directory: " MEDIAPATH
    else
      echo "Using media path: " $1
      export MEDIAPATH=$1
    fi

    if ! [ -d "$MEDIAPATH" ]; then
      echo "Warning! Directory '$MEDIAPATH' don't exist."
      echo 
      read -p "Continue? (y/n): " CONTINUE
      if ! [ "$CONTINUE" = "y" ]; then
        exit 0
      fi
    fi
    
    if ! [ -z $2 ]; then
      export PORT=$2
    else
      export PORT=9834
    fi

    echo "# MediaWEB configuration file" > $CONFIG || {
      echo "Unable to create $CONFIG"
      echo
      exit 1
    }
    echo >> $CONFIG
    echo "port = $PORT" >> $CONFIG
    echo >> $CONFIG
    echo "mediapath = $MEDIAPATH" >> $CONFIG
    echo >> $CONFIG
    echo "# Thumb cache path is by default your operating systems" >> $CONFIG
    echo "# temp folder + mediaweb. Uncomment below to set to" >> $CONFIG
    echo "# another location. Not used if enablethumbcache = off." >> $CONFIG
    echo "#thumbpath =" >> $CONFIG
    echo >> $CONFIG
    echo "# Thumbnail cache is on by default" >> $CONFIG
    echo "#enablethumbcache = off" >> $CONFIG
    echo >> $CONFIG
    echo "#autorotate = off" >> $CONFIG
    echo >> $CONFIG
    echo "logfile = /var/log/mediaweb.log" >> $CONFIG
    echo >> $CONFIG
    echo "# Log level is 'info' by default" >> $CONFIG
    echo "#loglevel = trace" >> $CONFIG
    echo >> $CONFIG
    echo "# User name and password for authentication. Leave commented" >> $CONFIG
    echo "# for no authentication" >> $CONFIG
    echo "#username = myusername" >> $CONFIG
    echo "#password = mypassword" >> $CONFIG
    
    echo "Wrote  : $CONFIG"

  fi

  #########################################################
  # Copy the mediaweb executable to /usr/sbin
  #########################################################

  cp mediaweb $EXEDESTINATION || {
    echo "ERROR! Unable to copy mediaweb to $EXEDESTINATION"
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


  echo
  echo "MediaWEB service successfully installed :-)"
	echo

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
	
	# Check if service is already installed
  systemctl is-active --quiet mediaweb && {
	  echo "MediaWEB service is running."
		echo 
		echo "Uninstalling MediaWEB service"
		echo
		
		systemctl disable mediaweb || {
			echo "ERROR! Unable to disable MediaWEB service"
			echo
			exit 1
		}
		
		systemctl stop mediaweb || {
			echo "ERROR! Unable to stop MediaWEB service"
			echo
			exit 1
		}
	}
	
	remove_file $SERVICECONFIG
	remove_file $CONFIG
	remove_file $EXEDESTINATION

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
  uninstall_service
else
  echo "ERROR! Unknown command '$1'"
  echo
  print_usage
fi
