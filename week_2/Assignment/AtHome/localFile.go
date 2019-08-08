package main

import (
	"bufio"
	"os"
	"sync"
)

type LocalFileInfo struct {
	Info FileInfo
}

func (f LocalFileInfo) read(outPut chan Line, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	defer close(outPut)
	file, err := os.Open(f.Info.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		outPut <- Line{scanner.Text()}
	}

	return nil
}
