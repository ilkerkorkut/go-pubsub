package pubsub

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
)

func (ps *PubSub) publishOtherNode(data interface{}, port string) {
	if ps.Config.Debug {
		log.Printf("Sending data %v to node with port: %v", data, port)
	}
	c, err := net.Dial("tcp", ":"+port)
	if err != nil {
		log.Println(err)
		ps.Publish(data)
		return
	}

	e := json.NewEncoder(c)
	if err := e.Encode(data); err != nil {
		log.Println("Json couldn't encoded.")
		return
	}

	if ps.Config.Debug {
		log.Printf("Data published successfully:  %v", data)
	}

	message, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		ps.Publish(data)
	} else {
		log.Printf("Data published successfully to other node, replyMessage : " + message)
	}
}
