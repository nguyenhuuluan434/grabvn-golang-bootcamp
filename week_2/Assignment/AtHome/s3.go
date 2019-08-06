package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type S3FileInfo struct {
	Info FileInfo
}


func (f S3FileInfo) read(outPut chan Line, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	defer close(outPut)
	resp, err := http.Get(f.Info.Path)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = UnknownError{fmt.Sprintf("could not get file status code is %d", resp.StatusCode)}
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	index := 0
	for len(data) > 0 {
		line := ""
		if len(data) < maxLengthLine {
			line = string(data)
			data = []byte{}
		} else {
			index = bytes.LastIndexByte(data[0:maxLengthLine], ' ')
			//word have len larger than max length line
			if index < 0 {
				line = string(data)
				data = []byte{}
			} else {
				line = string(data[0:index])
				data = data[index:]
			}
		}
		outPut <- Line{value: line}
	}
	return
}