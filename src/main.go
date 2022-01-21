package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	appConfig := Config{
		Extractor: LinkExtractor(),
		DLoader:   imageDownloader{},
		BaseURL:   "http://icanhas.cheezburger.com/",
		DstDir:    "./images",
	}

	flag.IntVar(&appConfig.ImgQuantity, "amount", 10, "Default quantity 10")
	flag.IntVar(&appConfig.Goroutines, "threads", 5, "Default quantity 5")
	flag.Parse()

	if appConfig.Goroutines > 5 || appConfig.Goroutines < 1 {
		log.Fatal("max threads allowed [1 - 5]")
		os.Exit(-1)

	}

	err := RunApp(appConfig)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(-1)
	}
}
