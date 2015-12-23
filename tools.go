package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func download(url string) ([]byte, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Respose code != 200")
	}

	return ioutil.ReadAll(response.Body)
}
