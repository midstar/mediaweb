package main

import (
	"github.com/midstar/llog"
)

func main() {
	llog.SetLevel(llog.LvlTrace)
	media := /* createMedia(".", ".", false) */ createMedia("Y:", "tmpcache/live", true, true)
	webAPI := CreateWebAPI(9834, "templates", media)
	httpServerDone := webAPI.Start()
	<-httpServerDone // Block until http server is done
}
