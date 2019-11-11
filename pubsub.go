package pubsub

import (
	"log"
	"sync"
)

type PubSub struct {
	Channel                chan interface{}
	WaitGroup              *sync.WaitGroup
	SubscriberCount        int
	Task                   func(interface{}, func())
	PublishedDataCount     int
	CurrentSubscribedCount int
	Debug                  bool
	Mutex                  sync.Mutex
}

func New(subscriberCount int, publishedDataCount int, taskFunc func(interface{}, func()), debug bool) *PubSub {
	ps := PubSub{
		Channel:            make(chan interface{}, 100),
		WaitGroup:          new(sync.WaitGroup),
		SubscriberCount:    subscriberCount,
		PublishedDataCount: publishedDataCount,
		Task:               taskFunc,
		Debug:              debug,
	}
	return &ps
}

func (ps *PubSub) Publish(data interface{}) {
	if ps.Debug {
		log.Printf("Published : %v", data)
	}
	ps.Channel <- data
}

func (ps *PubSub) StartSubscribers() {
	for i := 0; i <= ps.SubscriberCount; i++ {
		ps.WaitGroup.Add(1)
		go ps.Subscriber(i)
		if ps.Debug {
			log.Printf("Added subscriber with id: %d ", i)
		}
	}
}

func (ps *PubSub) Subscriber(id int) {
	defer ps.WaitGroup.Done()
	for d := range ps.Channel {
		if ps.Debug {
			log.Printf("ID : %d - Task Consuming", id)
		}
		ps.Task(d, ps.taskCallback)

		ps.Mutex.Lock()
		ps.CurrentSubscribedCount++
		ps.Mutex.Unlock()
		// If published data consuming is completely done, close channel
		if ps.CurrentSubscribedCount == ps.PublishedDataCount {
			log.Printf("Subscribed Tasks: %d - Published Tasks: %d", ps.CurrentSubscribedCount, ps.PublishedDataCount)
			ps.CloseChannel()
		}
	}
}

func (ps *PubSub) Wait() {
	ps.WaitGroup.Wait()
}

func (ps *PubSub) CloseChannel() {
	close(ps.Channel)
	log.Println("Channel Closed")
}

func (ps *PubSub) taskCallback() {
	if ps.Debug {
		log.Println("Task Done")
	}
}
