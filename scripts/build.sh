# Builds mediaweb using the correct version, build time and git hash
if [ -z $1 ]; then
  echo No version provided. Using unofficial.
  export VERSION=unofficial
else
  export VERSION=$1
fi

DATETIME=`date`
GITHASH=`git rev-parse HEAD`

echo version: $VERSION
echo git hash: $GITHASH
echo date time: $DATETIME

echo building / installing
cd $GOPATH/src/github.com/midstar/mediaweb
rice embed-go
go build -ldflags="-X 'main.applicationBuildTime=$DATETIME' -X main.applicationVersion=$1 -X main.applicationGitHash=$GITHASH" github.com/midstar/mediaweb

