#!/bin/sh
set -ex

print_usage() {
  echo "Usage:"
  echo
  echo "  sudo sh conf_edit.sh <conffile> <key> <value>"
  echo
  echo "     conffile: configuration file"
  echo "     key:      configuration key/parameter name"
  echo "     value:    configuration key/parameter value"
  echo
  exit 1
}

if [ "$#" -ne 3 ]; then
    echo "Error! Missing parameters"
    print_usage
fi

export CONFFILE=$1
export KEY=$2
export VALUE=$3

if ! [ -f "$CONFFILE" ]; then
	echo "$CONFFILE does not exist!"
	echo
	print_usage
fi

export MATCH=".*$KEY\s=.*"
export KEY_VALUE="${KEY} = $VALUE"

if grep -q $MATCH "$CONFFILE"; then
  # Key exist - replace
  sed -i -e "s/$MATCH/$KEY_VALUE/" $CONFFILE
else
  # Key don't exist - append
  echo "$KEY_VALUE" >> $CONFFILE
fi
