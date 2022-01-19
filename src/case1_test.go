package main

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

var img []byte

func init() {
	img, _ = hex.DecodeString(
		"89504e470d0a1a0a0000000d49484452000000c8000000640806000000c3867f0b0000013a49444154789cecd5311102411404518e42050e908622a4e1001b9ff0ea924e7783f7144cd2358ff7f73737b6f1793d8fd51b38dd570f809d0904824020080482402008048240200804824020080482402008048240200804824020080482402008048240200804824020080482402008048240200804824020080482402008048240200804824020080482402008048240200804824020080482402008048240200804824020080482402008048240201cab07703533b37a03270f02412010040241201004024120100402412010040241201004024120100402412010040241201004024120100402412010040241201004024120100402412010040241201004024120100402412010040241201004024120100402412010040241201004024120100402412010040241201004024120100402e11f0000ffff375608c6eefd06190000000049454e44ae426082")
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

	links := ExtractLinks(stub)

	assert.NotEmpty(t, links)
	assert.NoError(t, err)
	assert.Equal(t, 27, len(links))
}

func TestThatImagesCanBeDownloaded(t *testing.T) {
	webContent, err := ReadContent()
	if err != nil {
		assert.Fail(t, err.Error())
	}

	links := ExtractLinks(webContent)
	images, err := DownloadImages(links, mock(img))

	assert.True(t, len(links) > 0)
	assert.True(t, len(images) > 0)
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
