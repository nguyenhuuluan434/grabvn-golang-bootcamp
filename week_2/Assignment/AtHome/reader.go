package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

const maxLengthLine = 200
const maxBufferSize = 10000

type FileInfo struct {
	FileName string
	Path     string
}

type Line struct {
	value string
}

type FileAbstraction interface {
	read(outPut chan Line, wg *sync.WaitGroup) (err error)
}

func FileInfoFactory(path string) (fileInfos []FileAbstraction, err error) {
	switch {
	case strings.HasPrefix(path, "https://"):
		fileName, err := GetFileName(path)
		if err != nil {
			return []FileAbstraction{}, err
		}
		return []FileAbstraction{S3FileInfo{FileInfo{Path: path, FileName: fileName}}}, nil
	case strings.HasSuffix(path, "/"):
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return []FileAbstraction{}, err
		}
		if len(files) == 0 {
			return []FileAbstraction{}, nil
		}
		for _, file := range files {
			if !file.IsDir() {
				fileInfos = append(fileInfos, LocalFileInfo{FileInfo{Path: path + file.Name(), FileName: file.Name()}})
			}
		}
		return fileInfos, nil
	default:
		fileName, err := GetFileName(path)
		if err != nil {
			return []FileAbstraction{}, nil
		}
		return []FileAbstraction{LocalFileInfo{FileInfo{Path: path, FileName: fileName}}}, nil
	}
}

func GetFileName(path string) (name string, err error) {
	if len(path) == 0 {
		return "", UnknownError{message: "invalid path"}
	}
	parts := strings.Split(path, "/")
	return parts[len(parts)-1], nil
}

func WordCounter(data string, output chan<- map[string]int, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	words := strings.Fields(string(data))
	var result = make(map[string]int)
	for _, word := range words {
		result[word] += 1
	}
	output <- result
	close(output)
}


func process(paths []string, wg *sync.WaitGroup) (result map[string]int, err error) {
	if len(paths) == 0 {
		return make(map[string]int), UnknownError{"missing path "}
	}
	var fileInfoChannel = make(chan FileAbstraction, maxBufferSize)
	var counterChannel = make(chan map[string]int, maxBufferSize)
	var lineChannels = make(chan chan Line, maxBufferSize)
	var wordGroupChannels = make(chan chan map[string]int, maxBufferSize)
	wg.Add(1)
	go readFileToChannel(fileInfoChannel, lineChannels, wg)
	wg.Add(1)
	go wordCounterBaseLine(lineChannels, wordGroupChannels, wg)
	wg.Add(1)
	go mergeCounterLine(wordGroupChannels, counterChannel, wg)
	//init data by add files from parameter
	for _, path := range paths {
		fileInfos, err := FileInfoFactory(path)
		if err != nil {
			//try to learn add log the error here with path
			fmt.Println(err)
		}
		for _, fileInfo := range fileInfos {
			fileInfoChannel <- fileInfo
		}
	}
	close(fileInfoChannel)
	result = sumCounter(counterChannel)
	return
}

func sumCounter(counterChannel <-chan map[string]int) (result map[string]int) {
	var lock sync.Mutex
	result = make(map[string]int)
	for {
		var ok bool
		var wordLineCounter map[string]int
		select {
		case wordLineCounter, ok = <-counterChannel:
			if !ok {
				break
			}
			for key, value := range wordLineCounter {
				lock.Lock()
				result[key] = result[key] + value
				lock.Unlock()
			}
		}
		if !ok {
			break
		}
	}
	return
}

func readFileToChannel(fileInfoChannel <-chan FileAbstraction, lineChannels chan<- chan Line, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		var fileInfo FileAbstraction
		var ok bool
		select {
		case fileInfo, ok = <-fileInfoChannel:
			if !ok {
				break
			}
			lineChannel := make(chan Line, maxBufferSize)
			lineChannels <- lineChannel
			wg.Add(1)
			go fileInfo.read(lineChannel, wg)
		}
		if !ok {
			close(lineChannels)
			break
		}
	}
}

func wordCounterBaseLine(lineChannels <-chan chan Line, wordGroupChannels chan<- chan map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		var ok bool
		var fileLineChannel chan Line
		select {
		case fileLineChannel, ok = <-lineChannels:
			if !ok {
				break
			}
			wordCounterInLine(fileLineChannel, wordGroupChannels, wg)
		}
		if !ok {
			close(wordGroupChannels)
			break
		}
	}
}

func wordCounterInLine(lineChannel <-chan Line, wordGroupChannels chan<- chan map[string]int, wg *sync.WaitGroup) {
	for {
		var ok bool
		var line Line
		select {
		case line, ok = <-lineChannel:
			if !ok {
				break
			}
			wordCountChannel := make(chan map[string]int, maxBufferSize)
			wordGroupChannels <- wordCountChannel
			go WordCounter(line.value, wordCountChannel, wg)
		}
		if !ok {
			break
		}
	}
}

func mergeCounterLine(wordCountChannels <-chan chan map[string]int, outPut chan<- map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		var ok bool
		var wordCountChannel chan map[string]int
		select {
		case wordCountChannel, ok = <-wordCountChannels:
			if !ok {
				break
			}
			for {
				var isClose bool
				var lineCounter map[string]int
				select {
				case lineCounter, isClose = <-wordCountChannel:
					if !isClose {
						break
					}
					outPut <- lineCounter
				}
				if !isClose {
					break
				}
			}
		}
		if !ok {
			close(outPut)
			break
		}
	}
}

//func testFullFlow(fileInfoChannel <-chan FileAbstraction, outPut chan<- map[string]int, wg *sync.WaitGroup) {
//	wg.Add(1)
//	defer wg.Done()
//	lineChannels := make(chan chan Line, maxBufferSize)
//	go func() {
//		for {
//			wg.Add(1)
//			var fileInfo FileAbstraction
//			var ok bool
//			select {
//			case fileInfo, ok = <-fileInfoChannel:
//				if !ok {
//					wg.Done()
//					break
//				}
//				lineChannel := make(chan Line, maxBufferSize)
//				lineChannels <- lineChannel
//				go fileInfo.read(lineChannel, wg)
//			}
//			if !ok {
//				close(lineChannels)
//				break
//			}
//		}
//	}()
//
//	wordCountChannels := make(chan chan map[string]int, maxBufferSize)
//	go func() {
//		for {
//			var isClose bool
//			var fileLineChannel chan Line
//			select {
//			case fileLineChannel, isClose = <-lineChannels:
//				wg.Add(1)
//				if !isClose {
//					wg.Done()
//					break
//				}
//				for {
//					var isSubClose bool
//					var line Line
//					select {
//					case line, isSubClose = <-fileLineChannel:
//						wg.Add(1)
//						if !isSubClose {
//							wg.Done()
//							break
//						}
//						wordCountChannel := make(chan map[string]int, maxBufferSize)
//						wordCountChannels <- wordCountChannel
//						go WordCounter(line.value, wordCountChannel, wg)
//					}
//					if !isSubClose {
//						wg.Done()
//						break
//					}
//				}
//
//			}
//			if !isClose {
//				close(wordCountChannels)
//				break
//			}
//		}
//	}()
//
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		for {
//			var ok bool
//			var wordCountChannel chan map[string]int
//
//			select {
//			case wordCountChannel, ok = <-wordCountChannels:
//				if !ok {
//					break
//				}
//				for {
//					var isClose bool
//					var lineCounter map[string]int
//					select {
//					case lineCounter, isClose = <-wordCountChannel:
//						if !isClose {
//							break
//						}
//						outPut <- lineCounter
//					}
//					if !isClose {
//						break
//					}
//				}
//			}
//			if !ok {
//				close(outPut)
//				break
//			}
//		}
//	}()
//
// }
