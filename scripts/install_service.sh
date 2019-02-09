# MediaWEB systemd service installation shell script
# 
# Usage:
#  sh install_service.sh [mediapath] [port]
# 
#    mediapath: Path to all videos and pictures. Will be
#               promted if not provided
#    port:      Network port number (default 9834)
echo "--------------------------------------------------"
echo "Installation of MediaWEB service on Linux         "
echo "--------------------------------------------------"
echo

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

#TBD export CONFIG=/etc/mediaweb.conf
export CONFIG=mediaweb.conf

export KEEP=n
if [ -f "$CONFIG" ]; then
  echo "There is already a configuration in " $CONFIG
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
  
  if ! [ -z $2 ]; then
    export PORT=$2
  else
    export PORT=9834
  fi

  echo "# MediaWEB configuration file" > $CONFIG
  echo "port = " $PORT >> $CONFIG
  echo "mediapath = " $MEDIAPATH >> $CONFIG
  echo "#thumbpath =" >> $CONFIG
  echo "#enablethumbcache = off" >> $CONFIG
  echo "#autorotate = off" >> $CONFIG
  echo "logfile = /var/log/mediaweb.log" >> $CONFIG
  echo "#loglevel = trace" >> $CONFIG
  
fi



echo






