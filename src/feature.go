package main

import (
	"errors"
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

	log.Printf("Finding images links...\n")
	links, err = c.Extractor.GetImgInfo(c.DLoader, c.BaseURL, c.ImgQuantity)
	if err != nil {
		return err
	}

	log.Printf("Downloading...\n")
	linksChannel = c.Extractor.GetImagesLinks(links)

	log.Printf("Saving...\n")
	imagesPaths := StoreImages(linksChannel, c.Goroutines, c.DLoader)

	showImgPaths(imagesPaths)

	return nil
}

func StoreImages(links chan interface{}, qty int, d Downloader) chan interface{} {
	return Process(
		links,
		qty,
		func(link interface{}, out chan interface{}, index int) {
			img, err := downloadImage(link.(string), d)
			if err != nil {
				log.Fatal(err.Error())
				os.Exit(-1)
			}

			err = storeImage(img, out, index)
			if err != nil {
				log.Fatal(err.Error())
				os.Exit(-1)
			}
		},
	)
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

func (e ImgLinksExtractor) GetImagesLinks(links []string) chan interface{} {
	const urlRegexGroup = 1
	outLinks := make(chan interface{}, len(links))

	for _, s := range links {
		link := s
		go func() {
			outLinks <- e.r.FindStringSubmatch(link)[urlRegexGroup]
		}()
	}

	return outLinks
}

func (e ImgLinksExtractor) GetImgInfo(d Downloader, baseURL string, imagesQuantity int) ([]string, error) {
	var (
		page           = baseURL
		links          []string
		linksRemaining = imagesQuantity
	)

	for pageNumber := 2; linksRemaining > 0; pageNumber++ {
		log.Printf("searching in page:%s\n", page)

		webContent, err := downloadWebContent(page, d)
		if err != nil {
			return nil, err
		}

		l := e.links(webContent, linksRemaining)

		if len(l) == 0 {
			return nil, errors.New("no more images")
		}

		links = append(links, l...)

		linksRemaining = imagesQuantity - len(links)

		page = fmt.Sprintf("%spage/%d", baseURL, pageNumber)
	}

	return links, nil
}

func (e ImgLinksExtractor) links(content []byte, quantity int) []string {
	return e.r.FindAllString(string(content), quantity)
}

func downloadImage(link string, d Downloader) (Image, error) {
	return d.Download(link)
}

func downloadWebContent(link string, d Downloader) ([]byte, error) {
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

func showImgPaths(paths chan interface{}) {
	for i := 0; i < cap(paths); i++ {
		p := <-paths
		log.Printf("%s\n", p)
	}
}

func storeImage(v Image, out chan interface{}, index int) error {
	imgPath := fmt.Sprintf("./images/%d.jpg", index)
	err := SaveImage(v, imgPath)
	if err != nil {
		return err
	}

	out <- imgPath
	return nil
}
