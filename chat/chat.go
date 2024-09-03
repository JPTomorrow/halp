/*
A Chat is comprised of a
*/

package chat

import (
	"fmt"

	"golang.org/x/net/websocket"
)

type Chat struct {
	CustomerId   int
	RepId        int
	WsConnection *websocket.Conn
}

func (c *Chat) Run() {
	customerMsgPump := make(chan string)
	repMsgPump := make(chan string)

	go c.pollMessage(customerMsgPump)
	go c.pollMessage(repMsgPump)

OUTER:
	for {
		select {
		case customerMsg, open := <-customerMsgPump:
			fmt.Println(customerMsg)
			if !open {
				break OUTER
			}
		case repMsg, open := <-repMsgPump:
			fmt.Println(repMsg)
			if !open {
				break OUTER
			}
		}
	}
}

func (ch *Chat) pollMessage(c chan string) {
	// for {
	// 	// ch.WsConnection.
	// 	_, msg, err := ch.WsConnection.Read()
	// 	if err != nil {
	// 		continue
	// 	}
	// 	fmt.Printf("%s\n", msg)

	// 	c <- string(msg)
	// 	time.Sleep(time.Millisecond * 500)
	// }

	// close(c)
}
