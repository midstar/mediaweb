#!/bin/sh
set -e

# MediaWEB Linux Debian package creator shell script
#
# See print_usage below for instructions

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

copy_w_permission() {
    DIR=`dirname "$2"`
    mkdir -p $DIR
    cp $1 $2
    chmod $3 $2
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

# Create package source create directory
export SCRIPT_PATH=`dirname "$0"`
export PKG_IN_PATH=$SCRIPT_PATH/debian_pkg
export PKG_OUT_PATH=$SCRIPT_PATH/tmp_debian_pkg/$NAME
export PKG_ROOT_PATH=$PKG_OUT_PATH/root
export MEDIAWEB_EXE=$SCRIPT_PATH/../mediaweb
export MEDIAWEB_CFG=$SCRIPT_PATH/../configs/mediaweb.conf

export PACKAGE_DESTINATION=$SCRIPT_PATH/..

# Cleanup - remove previous package path if it exist
if [ -d $PKG_OUT_PATH ]; then
    rm -Rf $PKG_OUT_PATH
fi

# ----------------------------------------------
# COPY IN TO OUT
mkdir -p $PKG_OUT_PATH
cp -r $PKG_IN_PATH/* $PKG_OUT_PATH 
# Set general permissions
find $PKG_OUT_PATH -type d -exec chmod 755 {} +
find $PKG_OUT_PATH  -type f -exec chmod 644 {} +


# ----------------------------------------------
# FILL ROOT DIRECTORY

# Copy mediaweb executable and modify permissions
copy_w_permission $MEDIAWEB_EXE $PKG_ROOT_PATH/usr/sbin/mediaweb 755
strip --strip-unneeded $PKG_ROOT_PATH/usr/sbin/mediaweb

# Copy mediaweb configuration and modify permissions
copy_w_permission $MEDIAWEB_CFG $PKG_ROOT_PATH/etc/mediaweb.conf 644

# Compress the man page and modify permissions
gzip -n --best $PKG_ROOT_PATH/usr/share/man/man1/mediaweb.1

# Create changelog
export CHANGELOG=$PKG_ROOT_PATH/usr/share/doc/mediaweb/changelog
sh $SCRIPT_PATH/generate_changelog.sh mediaweb $VERSION mediaweb-v $CHANGELOG
gzip -n --best $CHANGELOG

# Calculate size of root directory
SIZE=$(du -s ./$PKG_ROOT_PATH | awk '{print $1}')


# ----------------------------------------------
# DEBIAN DIRECTORY

# Modify control file with version, architecture and size
sed -i -e 's/__ARCHITECTURE__/'${ARCHITECTURE}'/g' $PKG_OUT_PATH/DEBIAN/control
sed -i -e 's/__VERSION__/'${VERSION}'/g' $PKG_OUT_PATH/DEBIAN/control
sed -i -e 's/__SIZE__/'${SIZE}'/g' $PKG_OUT_PATH/DEBIAN/control

# Set premissions on scripts
chmod 755 $PKG_OUT_PATH/DEBIAN/post* $PKG_OUT_PATH/DEBIAN/pre*


# ----------------------------------------------
# MOVE ROOT FILES AND REMOVE ROOT DIRECTORY

mv $PKG_ROOT_PATH/* $PKG_OUT_PATH
rm -Rf $PKG_ROOT_PATH

# ----------------------------------------------
# CREATE THE INSTALLER

# Create
dpkg-deb --build $PKG_OUT_PATH

# Move the resulting installer to MediaWEB root folder
mv ${PKG_OUT_PATH}.deb $PACKAGE_DESTINATION/
echo Generated:
realpath $PACKAGE_DESTINATION/${NAME}.deb

# Check with lintian
echo Validating package:
echo
lintian $PACKAGE_DESTINATION/${NAME}.deb

