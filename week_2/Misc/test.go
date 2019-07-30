package main

var receiveOnly = make(<-chan int)
var sendOnt = make(chan<- int)

//locking and non locking channel

func main() {

	y := <-receiveOnly

	sendOnt <- 1
}

func receiver(msg <-chan string) {

}

func send(msg chan<- string) {

}
