package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Runable interface {
	run(workerId int)
}

type Job Runable

type Worker struct {
	index         int
	job           chan Job
	poolContainer *Pool
}

func (w Worker) start() {
	for {
		if len(w.job) > 0 {
			if job := <-w.job; job != nil {
				job.run(w.index)
				w.poolContainer.workerPool <- w
			}
		}
	}
}

type Pool struct {
	workerPool chan Worker
	job        chan Job
}

func NewWorkerPool(noOfWorkers int, queueSize int64) *Pool {
	p := Pool{
		workerPool: make(chan Worker, noOfWorkers),
		job:        make(chan Job, queueSize),
	}
	for i := 1; i <= noOfWorkers; i++ {
		worker := Worker{
			index:         i,
			job:           make(chan Job, queueSize),
			poolContainer: &p,
		}
		go worker.start()
		p.workerPool <- worker
	}
	go p.start()
	return &p
}

func (p *Pool) dispatch(job Job) {
	p.job <- job
}

func (p *Pool) start() {
	for {
		if len(p.workerPool) > 0 && len(p.job) > 0 {
			job := <-p.job
			worker := <-p.workerPool
			worker.job <- job
		}
	}
}

type fileReaderTask struct {
	fileInfo FileAbstraction
	outPut   chan Line
}

func (task fileReaderTask) run(workerId int) {
	task.fileInfo.read(task.outPut)
}

type counterWordTask struct {
	line   Line
	outPut chan map[string]int
}

func (task counterWordTask) run(workerId int) {
	summary := WordCounter(task.line.value)
	task.outPut <- summary
}

type customError struct {
	message string
}

func (c customError) Error() string {
	return c.message
}

func WordCounter(data string) (result map[string]int) {
	words := strings.Fields(string(data))
	result = make(map[string]int)
	for _, word := range words {
		result[word] += 1
	}
	return
}

const maxLengthLine = 200
const numBufferSize = 10000
const noWorker = 5

type Int int

func (i Int) run(workerId int) {
	//rand.Seed(time.Now().UnixNano())
	//min := 1
	//max := 5
	//sleepTime := rand.Intn(max-min) + min
	//time.Sleep(time.Duration(sleepTime) * time.Second)
	//fmt.Println("worker ", workerId, "get", i, "after sleep ", sleepTime)
	time.Sleep(1 * time.Second)
	fmt.Println("worker ", workerId, "get", i)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("missing input")
	}
	poolRead := NewWorkerPool(noWorker, numBufferSize)

	lineChan := make(chan Line, numBufferSize)
	for _, path := range args {
		fileInfos, err := FileInfoFactory(path)
		if err != nil {
			fmt.Println(err)
		}
		for _, fileInfo := range fileInfos {
			task := fileReaderTask{fileInfo: fileInfo, outPut: lineChan}
			poolRead.dispatch(task)
		}
	}
	result := make(map[string]int)
	for {
		var line Line
		var ok bool
		select {
		case line,ok = <-lineChan:
			if !ok {
				break;
			}
			sumLine := WordCounter(line.value)
			for key, value := range sumLine {
				result[key] = result[key] + value
			}
			fmt.Println(result)
			continue
		}
		if !ok {
			break
		}
	}

	fmt.Scanln()

}

type FileInfo struct {
	FileName string
	Path     string
}

type S3FileInfo struct {
	Info FileInfo
}

type Line struct {
	value string
}

func (f S3FileInfo) read(outPut chan Line) (err error) {
	//defer close(outPut)
	resp, err := http.Get(f.Info.Path)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = customError{fmt.Sprintf("could not get file status code is %d", resp.StatusCode)}
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

type LocalFileInfo struct {
	Info FileInfo
}

func (f LocalFileInfo) read(outPut chan Line) (err error) {
	fmt.Println("read local file")
	//add sleep to test task read data need long time to finish
	time.Sleep(1 * time.Second)
	//defer close(outPut)
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

type FileAbstraction interface {
	read(outPut chan Line) (err error)
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
		return "", customError{message: "invalid path"}
	}
	parts := strings.Split(path, "/")
	return parts[len(parts)-1], nil
}
