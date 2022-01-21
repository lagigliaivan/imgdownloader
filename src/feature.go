package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
)

type (
	Downloader interface {
		Download(src string) ([]byte, error)
	}

	Image []byte

	Config struct {
		Goroutines  int
		ImgQuantity int
		BaseURL     string
		DstDir      string
		DLoader     Downloader
		Extractor   ImgLinksExtractor
	}

	ImgLinksExtractor struct {
		r *regexp.Regexp
	}
)

func RunApp(c Config) error {
	var (
		links        []string
		linksChannel chan interface{}
	)

	log.Printf("Creating dir %s ...\n", c.DstDir)
	err := CreateDir(c.DstDir)
	if err != nil {
		return err
	}

	log.Printf("Parsing images links...\n")
	links, err = c.Extractor.ImagesLinks(c.DLoader, c.BaseURL, c.ImgQuantity)
	if err != nil {
		return err
	}

	linksChannel = c.Extractor.ImagesLinksChannel(links)

	log.Printf("Downloading...\n")
	imagesPaths := Process(
		linksChannel,
		c.Goroutines,
		func(link interface{}, out chan interface{}, index int) {
			img, err := DownloadImage(link.(string), c.DLoader)
			if err != nil {
				log.Fatal(err.Error())
				os.Exit(-1)
			}
			storeImage(img, out, index)
		},
	)

	readPaths(imagesPaths)

	return nil
}

func Process(
	in chan interface{},
	goroutines int,
	routine func(interface{}, chan interface{}, int),
) chan interface{} {
	var (
		n   = cap(in)
		wg  = new(sync.WaitGroup)
		out = make(chan interface{}, n)
	)

	go func() { //goroutine to allowed paths being shown while they are downloaded
		for i := 0; i < n; {
			wg.Add(goroutines)

			for t := 0; t < goroutines; t++ {
				index := i
				v := <-in
				go func() {
					routine(v, out, index)
					wg.Done()
				}()

				i++
			}

			wg.Wait()
		}
	}()

	return out
}

func LinkExtractor() ImgLinksExtractor {
	return ImgLinksExtractor{
		r: regexp.MustCompile(`<img class="resp-media.*" src="data:image.* data-src="(?P<url>https://.*?)" .*`),
	}
}

func (ext ImgLinksExtractor) ImagesLinksChannel(links []string) chan interface{} {
	const urlRegexGroup = 1
	outLinks := make(chan interface{}, len(links))

	for _, s := range links {
		link := s
		go func() {
			outLinks <- ext.r.FindStringSubmatch(link)[urlRegexGroup]
		}()
	}

	return outLinks
}

func (ext ImgLinksExtractor) links(content []byte, quantity int) []string {
	return ext.r.FindAllString(string(content), quantity)
}

func (ext ImgLinksExtractor) ImagesLinks(d Downloader, url string, n int) ([]string, error) {
	var (
		linksRemaining = n
		page           = url
		links          []string
	)

	for pageNumber := 2; linksRemaining > 0; pageNumber++ {
		webContent, err := WebContent(page, d)
		if err != nil {
			return nil, err
		}

		l := ext.links(webContent, linksRemaining)
		links = append(links, l...)

		linksRemaining = n - len(links)

		page = fmt.Sprintf("%spage/%d", url, pageNumber)
	}

	return links, nil
}

func DownloadImage(link string, d Downloader) (Image, error) {
	return d.Download(link)
}

func WebContent(link string, d Downloader) ([]byte, error) {
	return d.Download(link)
}

func CreateDir(dirName string) error {
	dstPath := fmt.Sprintf("./%s", dirName)

	err := os.Mkdir(dstPath, 0755)
	if err != nil {
		return err
	}

	return nil
}

func SaveImage(img Image, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(img)
	if err != nil {
		return err
	}

	return nil
}

func readPaths(paths chan interface{}) {
	for i := 0; i < cap(paths); i++ {
		p := <-paths
		log.Printf("%s\n", p)
	}
}

func storeImage(v Image, out chan interface{}, index int) {
	imgPath := fmt.Sprintf("./images/%d.jpg", index)
	SaveImage(v, imgPath)
	out <- imgPath
}
