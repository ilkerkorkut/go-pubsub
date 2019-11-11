package main

import (
	"github.com/ilkerkorkut/go-pubsub"
	"log"
	"runtime"
	"time"
)

func main() {

	data := [100000]int{}

	for i := 0; i < len(data); i++ {
		data[i] = i
	}

	go func() {
		ps := pubsub.New(runtime.NumCPU(), len(data), myTask, false)
		ps.StartSubscribers()
		for i := 0; i < len(data); i++ {
			ps.Publish(data[i])
		}
		ps.Wait()
	}()

	c := make(chan string)
	go ticker(c)
	for range c {
	}
	<-c
}

func myTask(data interface{}, done func()) {
	defer done()
	log.Println("My Task is running with data : ", data)
}

func ticker(ch chan<- string) {
	ticker := time.NewTicker(1000 * time.Millisecond)
	for _ = range ticker.C {
		ch <- ""
	}
}
