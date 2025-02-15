package main

import (
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"log"
	"runtime"
	"time"
)

type Channel struct {
	Channel string `json:"channel"`
}

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func sendJoin(c *gosocketio.Client) {
	log.Println("Acking /join")
	result, err := c.Ack("/join", Channel{"main"}, time.Second*5)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Ack result to /join: ", result)
	}
}

func onConn(h *gosocketio.Channel) {
	log.Println(gosocketio.OnConnection)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

//	c, err := gosocketio.Dial(
	c, err := gosocketio.DialwithConnAndHeader(
		gosocketio.GetUrl("localhost", 3811, false),
		transport.GetDefaultWebsocketTransport(),
		onConn,
		nil)
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("/message", func(h *gosocketio.Channel, args Message) {
		log.Println("--- Got chat message: ", args)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Println(gosocketio.OnDisconnection)
	})
	if err != nil {
		log.Fatal(err)
	}
/*
	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
	})
	if err != nil {
		log.Fatal(err)
	}
*/
	time.Sleep(1 * time.Second)

	go sendJoin(c)
	go sendJoin(c)
	go sendJoin(c)
	go sendJoin(c)
	go sendJoin(c)

	time.Sleep(60 * time.Second)
	c.Close()

	log.Println(" [x] Complete")
}
