# Builds plm using the correct version, build time and git hash
echo building / installing
pushd $GOPATH/src/github.com/midstar/mediaweb
packr2
popd
go install -ldflags="-X 'main.applicationBuildTime=$DATE $TIME' -X main.applicationVersion=$1 -X main.applicationGitHash=$GITHASH" github.com/midstar/mediaweb
