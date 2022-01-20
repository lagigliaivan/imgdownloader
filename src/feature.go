package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
)

type (
	Downloader interface {
		Download(src string) ([]byte, error)
	}

	Image []byte

	Config struct {
		D        Downloader
		BaseURL  string
		DstDir   string
		Quantity int
	}
)

func RunApp(c Config) error {
	images, err := StartImagesDownload(c.D, c.BaseURL, c.Quantity)
	if err != nil {
		return err
	}

	paths, err := SaveImages(images, c.DstDir)
	if err != nil {
		return err
	}

	for _, p := range paths {
		log.Printf("%s\n", p)
	}

	return nil
}

func StartImagesDownload(d Downloader, baseURL string, imgQuantity int) ([]Image, error) {
	const secondPage = 2

	links, err := ImagesLinks(d, baseURL, imgQuantity)
	if err != nil {
		return nil, err
	}

	for i := secondPage; len(links) < imgQuantity; i++ {
		log.Printf("getting page %d\n", i)

		remainingImages := (imgQuantity - len(links))

		l, err := ImagesLinks(d, fmt.Sprintf("%s/page/%d", baseURL, i), remainingImages)
		if err != nil {
			return nil, err
		}

		links = append(links, l...)
	}

	return DownloadImages(links, d)
}

func ImagesLinks(d Downloader, url string, quantity int) ([]string, error) {
	webContent, err := d.Download(url)
	if err != nil {
		return nil, err
	}

	return ExtractImagesLinks(webContent, quantity), nil

}

func ExtractImagesLinks(content []byte, n int) []string {
	url := 1
	var result []string
	r := regexp.MustCompile(`<img class="resp-media.*" src="data:image.* data-src="(?P<url>https://.*?)" .*`)

	links := r.FindAllString(string(content), n)
	for _, s := range links {
		l := r.FindStringSubmatch(s)
		result = append(result, l[url])
	}

	return result
}

func DownloadImages(links []string, d Downloader) ([]Image, error) {
	var images []Image

	for _, l := range links {
		content, err := d.Download(l)
		if err != nil {
			return nil, err
		}

		images = append(images, Image(content))
	}

	return images, nil
}

func SaveImages(images []Image, dir string) ([]string, error) {
	var filePaths []string

	dstPath := fmt.Sprintf("./%s", dir)

	err := os.Mkdir(dstPath, 0755)
	if err != nil {
		return nil, err
	}

	for i, img := range images {
		filePath := fmt.Sprintf("%s/%d.jpg", dstPath, i+1)

		f, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		_, err = f.Write(img)
		if err != nil {
			return nil, err
		}

		filePaths = append(filePaths, filePath)
	}

	return filePaths, nil
}
