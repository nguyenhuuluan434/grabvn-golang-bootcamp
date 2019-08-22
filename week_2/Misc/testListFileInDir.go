package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	files, err := ioutil.ReadFile("/etc/hosts")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f)
	}

	ioutil.TempDir()
}
