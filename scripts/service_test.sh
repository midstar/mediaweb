# Test of service.sh - both installation and uninstallation
cd $GOPATH/src/github.com/midstar/mediaweb/
sh scripts/service.sh install $GOPATH/src/github.com/midstar/mediaweb/testmedia
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
echo "Test passed :-)"
sh scripts/service.sh uninstall