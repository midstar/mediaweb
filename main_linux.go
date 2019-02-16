// Main method for Linux systems
package main

func main() {
	webAPI := mainCommon()
	httpServerDone := webAPI.Start()
	<-httpServerDone // Block until http server is done
}
