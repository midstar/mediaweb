#!/bin/sh
set -e

# MediaWEB Linux Debian package creator shell script
#
# See print_usage below for instructions
export SCRIPT_PATH=`dirname "$0"`
export PKG_TEMPLATE_PATH=$SCRIPT_PATH/debian_pkg
export PKG_PATH=$SCRIPT_PATH/tmp_debian_pkg
export PKG_APP_PATH=$PKG_PATH/app_tmp
export MEDIAWEB_EXE=$SCRIPT_PATH/../mediaweb
export MEDIAWEB_CFG=$SCRIPT_PATH/../configs/mediaweb.conf
export PACKAGE_DESTINATION=$SCRIPT_PATH/..


print_usage() {
    echo "Usage:"
    echo
    echo "For MediaWEB Debian Package Creator:"
    echo
    echo "  sudo sh $0 <architecture> <version> <name>"
    echo
    echo "    <architecture>: Run 'dpkg-architecture -L' for"
    echo "                    supported architectures"
    echo "    <version>:      Should be in the format 1.1.1.1"
    echo "    <name>:         Name of resulting package (excluding .deb)"
    echo
    echo "MediaWEB exectuable for the correct architecture needs"
    echo "to be built and but directly in the mediaweb path."
    exit 1
}

if [ ! -f $MEDIAWEB_EXE ]; then
    echo "ERROR! MediaWEB executable, mediaweb, not found"
    echo
    echo "You need to run build MediaWEB and the executable"
    echo "shall be located in $MEDIAWEB_EXE"
    exit 1
fi

if [ "$#" -ne 3 ]; then
    echo "Error! Missing parameters"
    print_usage
fi

export ARCHITECTURE=$1
export VERSION=$2
export NAME=$3

# Cleanup - remove previous package path if it exist
if [ -d $PKG_PATH ]; then
    rm -Rf $PKG_PATH
fi


# Create package source create directory
export PKG_SRC_PATH=$PKG_PATH/$NAME

# Copy template files
mkdir -p $PKG_SRC_PATH
cp -r $PKG_TEMPLATE_PATH/DEBIAN $PKG_SRC_PATH
chmod -R 755 $PKG_SRC_PATH/DEBIAN

# Copy files to install to tmp
mkdir -p $PKG_APP_PATH/usr/sbin
cp $MEDIAWEB_EXE $PKG_APP_PATH/usr/sbin/

cp -r $PKG_TEMPLATE_PATH/etc $PKG_APP_PATH
cp $MEDIAWEB_CFG $PKG_APP_PATH/etc/

# Figure out size of application
SIZE=$(du -s ./$PKG_APP_PATH | awk '{print $1}')

# Move files to install to src
mv $PKG_APP_PATH/* $PKG_SRC_PATH
rm -Rf $PKG_APP_PATH

# Modify control file with version and architecture
sed -i -e 's/__ARCHITECTURE__/'${ARCHITECTURE}'/g' $PKG_SRC_PATH/DEBIAN/control
sed -i -e 's/__VERSION__/'${VERSION}'/g' $PKG_SRC_PATH/DEBIAN/control
sed -i -e 's/__SIZE__/'${SIZE}'/g' $PKG_SRC_PATH/DEBIAN/control

# Create changelog
sh $SCRIPT_PATH/generate_changelog.sh mediaweb $VERSION mediaweb-v $PKG_SRC_PATH/DEBIAN/changelog

# Create the installer 
dpkg-deb --build $PKG_SRC_PATH

# Move the resulting installer to MediaWEB root folder
mv $PKG_PATH/${NAME}.deb $PACKAGE_DESTINATION/
echo Generated:
realpath $PACKAGE_DESTINATION/${NAME}.deb
