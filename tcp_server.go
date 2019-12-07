package pubsub

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

func (ps *PubSub) server(port string) {
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println(err)
		return
	}
	defer l.Close()

	log.Println("Started TCP server!")

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		ps.handleConnection(c)
	}
}

func (ps *PubSub) handleConnection(c net.Conn) {
	var dataPacket DataPacket
	d := json.NewDecoder(c)
	if err := d.Decode(&dataPacket); err != nil {
		log.Println("Json couldn't decoded.")
	}

	if ps.Config.Debug {
		log.Printf("Got remote data package : %v \n", dataPacket)
	}

	ps.Channel <- dataPacket

	ps.replyMessage(c)
}

func (ps *PubSub) replyMessage(c net.Conn) {
	t := time.Now()
	myTime := t.Format(time.RFC3339) + "\n"
	c.Write([]byte(myTime))
}
