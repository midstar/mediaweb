package main

func main() {
	media := createMedia("Y:", ".")
	webAPI := CreateWebAPI(9834, "templates", media)
	httpServerDone := webAPI.Start()
	<-httpServerDone // Block until http server is done
}
