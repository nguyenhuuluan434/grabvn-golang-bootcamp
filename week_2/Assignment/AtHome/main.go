package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("missing input")
	}

	//done := make(chan int,100)
	var wg sync.WaitGroup
	result, err := process(args, &wg)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range result {
		fmt.Println("word ", key, "\t occur", value, "time")
	}
	wg.Wait()

	//for {
	//	ok :=true;
	//	select {
	//	case _,ok = <-done:
	//		if !ok {
	//			break
	//		}
	//	}
	//	if !ok {
	//		break
	//	}
	//}

}
