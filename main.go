package main

func main() {
	/*
		files, err := getFiles("c:/Program Files")
		if err != nil {
			fmt.Print(err)
			return
		}
		fmt.Print("Directories:\n")
		for _, file := range files {
			if file.IsDir {
				fmt.Printf("    %s\n", file.FullPath)
			}
		}
		fmt.Print("Files:\n")
		for _, file := range files {
			if file.IsDir == false {
				fmt.Printf("    %s\n", file.Name)
			}
		}
	*/
	media := createMedia("Y:")
	webAPI := CreateWebAPI(9834, "templates", media)
	httpServerDone := webAPI.Start()
	<-httpServerDone // Block until http server is done
}
