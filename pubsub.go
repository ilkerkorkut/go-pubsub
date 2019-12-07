package pubsub

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type PubSub struct {
	Channel                    chan interface{}
	WaitGroup                  *sync.WaitGroup
	SubscriberCount            int
	Task                       func(interface{}, func())
	TotalDataSize              int
	CurrentSubscribedDataCount int
	CurrentPublishedDataCount  int
	LatestPortIndex            int
	RemoteCall                 bool
	Mutex                      sync.Mutex
	Config                     Config
}

type Config struct {
	MultiNode  bool
	ServerPort string
	Name       string
	NodePorts  []string
	Debug      bool
}

type PsError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type DataPacket struct {
	Data interface{} `json:"data"`
	Time time.Time   `json:"time"`
}

func (e *PsError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Error-code: %d %s", e.Code, e.Message)
}

func New(subscriberCount int, totalDataSize int, taskFunc func(interface{}, func()), config *Config) (*PubSub, *PsError) {
	ps := PubSub{
		Channel:         make(chan interface{}, 100),
		WaitGroup:       new(sync.WaitGroup),
		SubscriberCount: subscriberCount,
		TotalDataSize:   totalDataSize,
		Task:            taskFunc,
		Config:          *config,
	}
	if config.MultiNode && config.ServerPort != "" {
		ps.RemoteCall = true
		go ps.server(config.ServerPort)
	} else if config.MultiNode && config.ServerPort == "" {
		return nil, &PsError{
			Code:    0,
			Message: "Please set server port.",
		}
	}
	ps.StartSubscribers()
	return &ps, nil
}

func (ps *PubSub) Publish(data interface{}) {
	if ps.Config.Debug {
		log.Printf("Publishing : %v", data)
	}
	ps.Mutex.Lock()
	ps.CurrentPublishedDataCount++
	ps.Mutex.Unlock()

	if ps.Config.MultiNode {
		if ps.RemoteCall {
			port := ps.getPort()
			go ps.publishOtherNode(data, port)
		} else {
			if ps.Config.Debug {
				log.Printf("Publishing local thread... %d , %d", data, ps.CurrentPublishedDataCount)
			}
			ps.RemoteCall = !ps.RemoteCall
			ps.Channel <- data
		}
	} else {
		ps.Channel <- data
	}
}

func (ps *PubSub) StartSubscribers() {
	for i := 0; i <= ps.SubscriberCount; i++ {
		ps.WaitGroup.Add(1)
		go ps.Subscriber(i)
		if ps.Config.Debug {
			log.Printf("Added subscriber with id: %d ", i)
		}
	}
}

func (ps *PubSub) Subscriber(id int) {
	defer ps.WaitGroup.Done()
	for d := range ps.Channel {
		if ps.Config.Debug {
			log.Printf("ID : %d - Task Consuming", id)
		}
		ps.Task(d, ps.taskCallback)

		ps.Mutex.Lock()
		ps.CurrentSubscribedDataCount++
		ps.Mutex.Unlock()
		// If published data consuming is completely done, close channel
		if ps.CurrentSubscribedDataCount == ps.TotalDataSize && !ps.Config.MultiNode {
			if ps.Config.Debug {
				log.Printf("Subscribed Data Count: %d - Published Tasks: %d", ps.CurrentSubscribedDataCount, ps.TotalDataSize)
			}
			ps.closeChannel()
		}
	}
}

func (ps *PubSub) Wait() {
	ps.WaitGroup.Wait()
}

func (ps *PubSub) closeChannel() {
	close(ps.Channel)
	log.Println("Channel Closed")
}

func (ps *PubSub) taskCallback() {
	if ps.Config.Debug {
		log.Println("Task Done")
	}
}

func (ps *PubSub) getPort() string {
	if ps.LatestPortIndex < len(ps.Config.NodePorts) {
		ps.Mutex.Lock()
		port := ps.Config.NodePorts[ps.LatestPortIndex]
		ps.LatestPortIndex++
		ps.Mutex.Unlock()
		return port
	} else {
		ps.Mutex.Lock()
		ps.LatestPortIndex = 1
		ps.RemoteCall = !ps.RemoteCall
		ps.Mutex.Unlock()
		return ps.Config.NodePorts[0]
	}
}
