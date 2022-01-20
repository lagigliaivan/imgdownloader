package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"regexp"
)

type source string

type Reader interface {
	Read(source) []byte
}

type Downloader struct {
	src Reader
}

func (d Downloader) Download(url string) ([]byte, error) {
	r := d.src.Read(source(url))
	if len(r) <= 0 {
		return nil, errors.New("empty content")
	}

	return r, nil
}

func NewImage(data []byte) (image.Image, error) {
	r := bytes.NewReader(data)

	imageData, err := png.Decode(r)
	if err != nil {
		return nil, err
	}

	return imageData, err
}

func ExtractImagesLinks(content []byte, n int) []string {
	var result []string
	r := regexp.MustCompile(`<img class="resp-media.*" src="data:image.* data-src="(?P<url>https://.*?)" .*`)

	links := r.FindAllString(string(content), n)
	for _, s := range links {
		l := r.FindStringSubmatch(s)
		result = append(result, l[0])
	}

	return result
}

func ReadContent() ([]byte, error) {
	stub, err := ioutil.ReadFile("../stubs/stub1.html")
	return stub, err
}

func DownloadImages(links []string, d Downloader) ([]image.Image, error) {
	var images []image.Image

	for _, l := range links {
		content, err := d.Download(l)
		if err != nil {
			return nil, err
		}

		image, err := NewImage(content)
		if err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	return images, nil
}

func startDownload(d Downloader) ([]image.Image, error) {
	amount := 10

	webContent, err := ReadContent()
	if err != nil {
		return nil, err
	}

	links := ExtractImagesLinks(webContent, amount)
	images, err := DownloadImages(links, d)

	return images, err
}

func save(images []image.Image, dir string) ([]string, error) {
	var filePaths []string

	dstPath := fmt.Sprintf("./%s", dir)

	err := os.Mkdir(dstPath, 0755)
	if err != nil {
		return nil, err
	}

	for i, img := range images {
		filePath := fmt.Sprintf("%s/%d.jpg", dstPath, i)

		f, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}

		err = png.Encode(f, img)
		if err != nil {
			return nil, err
		}

		f.Close()

		filePaths = append(filePaths, filePath)
	}

	return filePaths, nil
}
