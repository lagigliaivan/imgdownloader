package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	all     = 27
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
	d := mock(img)

	content, err := d.Download(baseURL)
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

func TestThatImagesLinksCanBeFound(t *testing.T) {
	content, err := homePageStub(1)

	assert.NoError(t, err)

	ext := LinkExtractor()
	links := ext.links(content, all)

	assert.NotEmpty(t, links)
	assert.Equal(t, 27, len(links))
}

func TestThatImagesCanBeFoundInAWebContent(t *testing.T) {
	ext := LinkExtractor()

	links, err := ext.FindImages(mock(img), baseURL, all)

	assert.NotEmpty(t, links)
	assert.NoError(t, err)
	assert.Equal(t, 27, len(links))
}

func TestThatImagesCanBeDownloaded(t *testing.T) {
	ext := LinkExtractor()

	links, err := ext.FindImages(mock(img), baseURL, all)
	assert.NoError(t, err)

	linksChannel := ext.GetImagesLinks(links)

	for l := range linksChannel {
		img, err := downloadImage(l.(string), mock(img))
		assert.NoError(t, err)
		assert.NotNil(t, img)
		return
	}

	assert.Fail(t, "shouldn't be here")
}

func TestThatAppCanBeRunIfConfigurationIsOk(t *testing.T) {
	c := Config{
		Extractor:   LinkExtractor(),
		DLoader:     mock(img),
		BaseURL:     "http://icanhas.cheezburger.com/",
		DstDir:      "./images",
		ImgQuantity: 10,
		Goroutines:  2,
	}
	//TODO: Avoid this
	tearDown()

	err := RunApp(c)

	//checking if files were properly created
	for i := 0; i < c.ImgQuantity; i++ {
		fname := fmt.Sprintf("./%s/%d.jpg", c.DstDir, i)
		if _, err := os.Stat(fname); os.IsNotExist(err) {
			assert.Fail(t, err.Error())
		}
	}

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

	if strings.Contains(src, "https://i.chzbgr.com/") {
		return d.b, nil
	}

	switch src {
	case baseURL:
		return homePageStub(1)

	case fmt.Sprintf("%s/page/2", baseURL):
		d.lastPage = 2
		return homePageStub(2)

	case fmt.Sprintf("%s/page/3", baseURL):
		d.lastPage = 3
		return homePageStub(3)

	default:
		return nil, fmt.Errorf("empty content")
	}
}

func homePageStub(page int) ([]byte, error) {
	return ioutil.ReadFile(fmt.Sprintf("../stubs/stub%d.html", page))
}
