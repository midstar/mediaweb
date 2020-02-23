#!/bin/sh
set -e

print_usage() {
  echo "Usage:"
  echo
  echo "  sudo sh generate_changelong.sh <modulename> <version> <tagprefix> <outfile>"
  echo
  exit 1
}

if [ "$#" -ne 4 ]; then
    echo "Error! Invalid number of parameters"
    print_usage
fi

export MODULE=$1
export VERSION=$2
export TAGPREFIX=$3
export OUTFILE=$4
export LAST_TAG=`git tag -l ${TAGPREFIX}* | sort -V | tail -1`

echo "$MODULE ($VERSION) unstable; urgency=low\n" > $OUTFILE;
git log --pretty=format:'  * %s' $LAST_TAG..HEAD >> $OUTFILE;
git log --pretty='format:%n%n -- %aN <%aE>  %aD%n%n' HEAD^..HEAD >> $OUTFILE
