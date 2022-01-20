package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type imageDownloader struct {
}

func (c imageDownloader) Download(u string) ([]byte, error) {
	_, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(u)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
