#!/bin/sh
set -e

print_usage() {
  echo "Usage:"
  echo
  echo "  sudo sh generate_changelong.sh <modulename> <version> <outfile>"
  echo
  exit 1
}

if [ "$#" -ne 3 ]; then
    echo "Error! Invalid number of parameters"
    print_usage
fi

export MODULE=$1
export VERSION=$2
export OUTFILE=$3

echo "$MODULE ($VERSION) unstable; urgency=low\n" > $OUTFILE;
git log --pretty=format:'  * %s' HEAD^..HEAD >> $OUTFILE;
git log --pretty='format:%n%n -- %aN <%aE>  %aD%n%n' HEAD^..HEAD >> $OUTFILE
