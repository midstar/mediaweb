# Same as build.sh but cross compiles for Linux ARM architecture 
export GOOS=linux 
export GOARCH=arm
export GOARM=5

echo Cross compiling for Linux ARM architecture

BASEDIR=$(dirname $0)
sh $BASEDIR/build.sh $1
