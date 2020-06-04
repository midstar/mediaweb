# Same as build.sh but cross compiles for ARM on linux

export GOOS=linux

if [ "$1" = "arm" ]; then
  export GOARCH=arm
  export GOARM=6
  echo Cross compiling for Linux ARMv$GOARM architecture
elif [ "$1" = "arm64" ]; then
  export GOARCH=arm64
  echo Cross compiling for Linux ARM64 architecture
else
  echo No architecture provided 
  echo Valid values are arm or arm64
  exit 1
fi

BASEDIR=$(dirname $0)
sh $BASEDIR/build.sh $2
