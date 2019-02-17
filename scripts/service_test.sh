# Test of service.sh - both installation and uninstallation
cd $GOPATH/src/github.com/midstar/mediaweb/

# Move files to a temporary folder to secure that packr 
# (embedded resources) works as expected
mkdir -p tmpout
mkdir -p tmpout/servicetest
cp mediaweb tmpout/servicetest/mediaweb
cd tmpout/servicetest 

export SCRIPTPATH=$GOPATH/src/github.com/midstar/mediaweb/scripts

sh $SCRIPTPATH/service.sh install $GOPATH/src/github.com/midstar/mediaweb/testmedia
echo "Waiting 2 seconds"
sleep 2
echo "Testing connection"
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:9834)
if ! [ "$HTTP_STATUS" = "200" ]; then
	echo "Test Failed! Unable to connect to MediaWEB."
	echo
	echo "Expected status code 200, but got $HTTP_STATUS"
	exit 1
fi
if ! [ -f "/var/log/mediaweb.log" ]; then
	echo "Test Failed! No log file was created in /var/mediaweb.log"
	echo
	exit 1
fi
sh $SCRIPTPATH/service.sh uninstall
echo "Test passed :-)"