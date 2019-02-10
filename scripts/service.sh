# MediaWEB systemd service (un)installation shell script
# 
# For installation:
#  sh service.sh install [mediapath] [port]
# 
#    mediapath: Path to all videos and pictures. Will be
#               promted if not provided
#    port:      Network port number (default 9834)
#

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
    echo TBD exit 1
    echo
  fi

  if ! [ -d "/etc/systemd/system" ]; then
    echo "Unable to access /etc/systemd/system"
    echo
    echo "Try run this script as root, i.e:"
    echo
    echo "sudo sh install_service.sh"
    echo TBD exit 1
    echo
  fi

  if [ ! -f "mediaweb" ]; then
      echo "ERROR! MediaWEB executable, mediaweb, not found"
      echo
      echo "You need to run this script in the path where you have the"
      echo "MediaWEB executable"
      exit 1
  fi

  # Check if running with systemctl is-active --quiet service && echo Service is running


  #########################################################
  # Create mediaweb.conf
  #########################################################

  #TBD export CONFIG=/etc/mediaweb.conf
  export CONFIG=mediaweb.conf

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
      echo "Try run this script as root, i.e:"
      echo
      echo "sudo sh install_service.sh"
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
    
    echo "Wrote  : $CONFIG"

  fi

  #########################################################
  # Copy the mediaweb executable to /usr/sbin
  #########################################################
  echo TBD export EXEDESTINATION=/usr/sbin/mediaweb
  export EXEDESTINATION=mediaweb_new

  cp mediaweb $EXEDESTINATION || {
    echo "Unable to copy mediaweb to $EXEDESTINATION"
    echo
    echo "Try run this script as root, i.e:"
    echo
    echo "sudo sh install_service.sh"
    exit 1
  }
  echo "Copied : $EXEDESTINATION"


  #########################################################
  # Create webmedia.service systemd configuration
  #########################################################

  #TBD export SERVICECONFIG=/etc/systemd/system/mediaweb.service
  export SERVICECONFIG=mediaweb.service

  echo "[Unit]" > $SERVICECONFIG || {
    echo "Unable to create $SERVICECONFIG"
    echo
    echo "Try run this script as root, i.e:"
    echo
    echo "sudo sh install_service.sh"
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

  echo
  echo "MediaWEB service successfully installed :-)"

}




if [ -z $1 ]; then
  print_usage
elif [ "$1" == "install" ]; then
  install_service $2 $3
elif  [ "$1" == "uninstall" ]; then
  echo TBD uninstall_service
else
  echo "ERROR! Unknown command '$1'"
  echo
  print_usage
fi
