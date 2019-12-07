package main

import (
	"github.com/ilkerkorkut/go-pubsub"
	"log"
	"runtime"
	"time"
)

func main() {

	data := [100]pubsub.DataPacket{}

	for i := 0; i < len(data); i++ {
		data[i] = pubsub.DataPacket{
			Data: i,
			Time: time.Now(),
		}
	}

	go func() {
		ps, err := pubsub.New(runtime.NumCPU(), len(data), myTask, &pubsub.Config{
			MultiNode: false,
			Name:      "my-pubsub-application",
			Debug:     true,
		})
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < len(data); i++ {
			ps.Publish(data[i])
		}

		ps.Wait()
		log.Println("Successfully completed!")
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
