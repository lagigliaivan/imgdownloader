package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	all     = -1
	baseURL = "http://icanhas.cheezburger.com/"
)

var img []byte

func setup() {
	img, _ = hex.DecodeString(
		"89504e470d0a1a0a0000000d49484452000000c8000000640806000000c3867f0b0000013a49444154789cecd5311102411404518e42050e908622a4e1001b9ff0ea924e7783f7144cd2358ff7f73737b6f1793d8fd51b38dd570f809d0904824020080482402008048240200804824020080482402008048240200804824020080482402008048240200804824020080482402008048240200804824020080482402008048240200804824020080482402008048240200804824020080482402008048240200804824020080482402008048240201cab07703533b37a03270f02412010040241201004024120100402412010040241201004024120100402412010040241201004024120100402412010040241201004024120100402412010040241201004024120100402412010040241201004024120100402412010040241201004024120100402412010040241201004024120100402e11f0000ffff375608c6eefd06190000000049454e44ae426082")
	os.RemoveAll("./images")
}

func tearDown() {
	os.RemoveAll("./images")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}
func TestThatHTTPContentCanBeDownloaded(t *testing.T) {
	url := ""
	d := mock(img)

	content, err := d.Download(url)
	assert.NoError(t, err)

	assert.NotEmpty(t, content)
	assert.NoError(t, err)
}

func TestAnErrorIsReturnedWhenHTTPContentCannotBeDownloaded(t *testing.T) {
	url := ""
	d := mock([]byte{})
	content, err := d.Download(url)

	assert.Error(t, err)
	assert.Nil(t, content)
}

func TestThatImagesLinksCanBeFoundInAWebContent(t *testing.T) {
	stub, err := homePageStub(1)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	links := ExtractImagesLinks(stub, all)

	for _, l := range links {
		_, err := url.Parse(l)
		assert.NoError(t, err)
	}

	assert.NotEmpty(t, links)
	assert.NoError(t, err)
	assert.Equal(t, 27, len(links))
}

func TestThatImagesCanBeDownloaded(t *testing.T) {
	webContent, err := homePageStub(1)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	links := ExtractImagesLinks(webContent, all)
	images, err := DownloadImages(links, mock(img))

	assert.True(t, len(links) > 0)
	assert.True(t, len(images) > 0)
	assert.NoError(t, err)
}

func TestThat10ImagesCanBeDownloaded(t *testing.T) {
	quantity := 10

	images, err := StartImagesDownload(mock(img), baseURL, quantity)

	assert.True(t, len(images) == quantity)
	assert.NoError(t, err)
}

//TODO: avoid accessing disk for testing purposes
func TestThatImagesAreSavedInADirectory(t *testing.T) {
	quantity := 10
	dir := "./images"

	images, _ := StartImagesDownload(mock(img), baseURL, quantity)

	paths, err := SaveImages(images, dir)

	//checking if files were properly created
	for i, p := range paths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			assert.Fail(t, err.Error())
		}
		name := fmt.Sprintf("./%s/%d.jpg", dir, i+1)
		if p != name {
			assert.Fail(t, fmt.Sprintf("image %s is not well named: %s\n", p, name))
		}
	}

	assert.True(t, len(images) == quantity)
	assert.True(t, len(paths) == quantity)
	assert.NoError(t, err)
}

func TestThatTheQuantityOfImagesToBeDownloadesCanBeSpecified(t *testing.T) {
	cases := []struct {
		desc         string
		imgsQuantity int
	}{
		{"Test downloading 1 image", 1},
		{"Test downloading 2 images", 2},
		{"Test downloading 3 images", 3},
		{"Test downloading 11 images", 11},
	}
	for _, tc := range cases {
		images, err := StartImagesDownload(mock(img), baseURL, tc.imgsQuantity)
		if err != nil {
			assert.Fail(t, err.Error())
		}

		assert.True(t, len(images) == tc.imgsQuantity)
	}
}

func TestThatIfImagesQuantityToDownloadCannotBeFulfilledThenSearchTheNextPage(t *testing.T) {
	quantity := 50
	d := mock(img)

	images, err := StartImagesDownload(d, baseURL, quantity)

	assert.True(t, len(images) == quantity)
	assert.True(t, d.lastPage > 1)
	assert.NoError(t, err)
}

func TestThatAppCanBeRunIfConfigurationIsOk(t *testing.T) {
	c := Config{
		D:        mock(img),
		BaseURL:  "http://icanhas.cheezburger.com/",
		DstDir:   "./images",
		Quantity: 15,
	}

	err := RunApp(c)

	assert.NoError(t, err)
}

func mock(returnValue []byte) *downloaderMock {
	return &downloaderMock{
		b: returnValue,
	}
}

type downloaderMock struct {
	b        []byte
	lastPage int
}

func (d *downloaderMock) Download(src string) ([]byte, error) {
	if src == baseURL {
		return homePageStub(1)
	}

	if src == fmt.Sprintf("%s/page/2", baseURL) {
		d.lastPage = 2
		return homePageStub(2)
	}

	if src == fmt.Sprintf("%s/page/3", baseURL) {
		d.lastPage = 3
		return homePageStub(3)
	}

	if len(d.b) == 0 {
		return nil, fmt.Errorf("empty content")
	}

	return d.b, nil
}

func homePageStub(page int) ([]byte, error) {
	return ioutil.ReadFile(fmt.Sprintf("../stubs/stub%d.html", page))
}
