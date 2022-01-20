package main

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	all = -1
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

	image, err := NewImage(content)

	assert.NotNil(t, image)
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
	stub, err := ReadContent()
	if err != nil {
		assert.Fail(t, err.Error())
	}

	links := ExtractImagesLinks(stub, all)

	assert.NotEmpty(t, links)
	assert.NoError(t, err)
	assert.Equal(t, 27, len(links))
}

func TestThatImagesCanBeDownloaded(t *testing.T) {
	webContent, err := ReadContent()
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
	amount := 10

	images, err := startDownload(mock(img))

	assert.True(t, len(images) == amount)
	assert.NoError(t, err)
}

//TODO: avoid accessing disk for testing purposes
func TestThatImagesAreSavedInADirectory(t *testing.T) {
	amount := 10
	dir := "./images"

	images, err := startDownload(mock(img))

	paths, err := save(images, dir)

	for _, p := range paths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			assert.Fail(t, err.Error())
		}
	}

	assert.True(t, len(images) == amount)
	assert.True(t, len(paths) == amount)
	assert.NoError(t, err)
}

func mock(returnValue []byte) Downloader {
	return Downloader{
		src: mockReader{
			b: returnValue,
		}}
}

type mockReader struct {
	b []byte
}

func (d mockReader) Read(src source) []byte {
	return d.b
}