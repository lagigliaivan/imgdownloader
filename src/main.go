package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	appConfig := Config{
		D:       imageDownloader{},
		BaseURL: "http://icanhas.cheezburger.com/",
		DstDir:  "./images",
	}

	flag.IntVar(&appConfig.Quantity, "amount", 10, "Default quantity 10")
	flag.Parse()

	err := RunApp(appConfig)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(-1)
	}
}
