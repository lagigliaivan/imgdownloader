package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type (
	Reader interface {
		Read(src string) ([]byte, error)
	}

	Downloader struct {
		src Reader
	}

	Image []byte
)

func AppRun() error {
	d := Downloader{
		src: client{},
	}

	images, err := StartDownload(d)
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

func (d Downloader) Download(url string) ([]byte, error) {
	return d.src.Read(url)
}

func ExtractImagesLinks(content []byte, n int) []string {
	urlGroup := 1
	var result []string
	r := regexp.MustCompile(`<img class="resp-media.*" src="data:image.* data-src="(?P<url>https://.*?)" .*`)

	links := r.FindAllString(string(content), n)
	for _, s := range links {
		l := r.FindStringSubmatch(s)
		result = append(result, l[urlGroup])
	}

	return result
}

func ReadContent() ([]byte, error) {
	return ioutil.ReadFile("../stubs/stub1.html")
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

func StartDownload(d Downloader) ([]Image, error) {
	amount := 10

	webContent, err := ReadContent()
	if err != nil {
		return nil, err
	}

	links := ExtractImagesLinks(webContent, amount)
	images, err := DownloadImages(links, d)

	return images, err
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
