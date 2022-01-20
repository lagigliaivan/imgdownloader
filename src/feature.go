package main

import (
	"fmt"
	"os"
	"regexp"
)

type (
	Downloader interface {
		Download(src string) ([]byte, error)
	}

	Image []byte
)

func AppRun(imgsQuantity int) error {
	d := client{}

	images, err := StartDownload(d, imgsQuantity)
	if err != nil {
		return err
	}

	paths, err := Save(images, "./images")
	if err != nil {
		return err
	}

	for _, p := range paths {
		fmt.Printf("%s\n", p)
	}

	return nil
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

func StartDownload(d Downloader, imgQuantity int) ([]Image, error) {
	webContent, err := d.Download("http://icanhas.cheezburger.com/")
	if err != nil {
		return nil, err
	}

	links := ExtractImagesLinks(webContent, imgQuantity)

	for i := 2; len(links) < imgQuantity; i++ {
		fmt.Printf("getting page %d\n", i)
		webContent, err := d.Download(fmt.Sprintf("http://icanhas.cheezburger.com/page/%d", i))
		if err != nil {
			return nil, err
		}
		l := ExtractImagesLinks(webContent, (imgQuantity - len(links)))

		links = append(links, l...)
	}

	images, err := DownloadImages(links, d)

	return images, err
}

func ImagesLinks() {

}

func Save(images []Image, dir string) ([]string, error) {
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
