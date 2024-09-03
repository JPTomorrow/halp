package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// represents a message sent
type Message struct {
	Msg string `json:"msg"`
	Id  int    `json:"id"` // customer or support rep id
}

type Chat struct {
	customerID     int
	repID          int
	customerWsConn *websocket.Conn
	repWsConn      *websocket.Conn
	messages       []string
	isComplete     bool
}

type ChatHub struct {
	chats []*Chat
	mu    sync.Mutex
}

func NewChat() *ChatHub {
	return &ChatHub{
		chats: []*Chat{},
	}
}

func (c *ChatHub) RegisterSupportRep(id int, conn *websocket.Conn) *Chat {
	c.mu.Lock()
	for _, chat := range c.chats {
		if chat.repID == id {
			return chat
		}
	}

	newChat := Chat{
		customerID:     -1,
		repID:          id,
		customerWsConn: nil,
		repWsConn:      conn,
		messages:       []string{},
		isComplete:     false,
	}
	c.chats = append(c.chats, &newChat)
	c.mu.Unlock()
	return &newChat
}

func (c *ChatHub) Run(ch *Chat, ec *WsContext, wc *websocket.Conn) {
	ctx, cancel := context.WithTimeout(ec.Request().Context(), time.Minute*15)
	defer cancel()

	var m Message
OUTER:
	for {
		err := wsjson.Read(ctx, wc, &m)
		if err != nil {
			log.Println("no message to read or invalid id.")
			wc.Close(websocket.StatusProtocolError, "no message to read or invalid id.")
			break OUTER
		}

		id := m.Id

		if id != ch.customerID && id != ch.repID {

			log.Printf("invalid id: %v\n", id)
			wc.Close(websocket.StatusProtocolError, "invalid id.")
			break OUTER
		}
		msg := m.Msg

		log.Printf("\n\tid: %v\n\tmsg: %v\n\n", id, msg)
		wsjson.Write(ctx, ch.repWsConn, m)
		wsjson.Write(ctx, ch.customerWsConn, m)
	}
}

func (c *ChatHub) NextAvailableRep(customerId int, conn *websocket.Conn) (*Chat, bool) {
	c.mu.Lock()
	for i, chat := range c.chats {
		if chat.customerID == -1 {
			c.chats[i].customerID = customerId
			c.chats[i].customerWsConn = conn
			return c.chats[i], true
		}
	}
	c.mu.Unlock()

	return nil, false
}

func (c *ChatHub) AppendMsg(ch *Chat, message string) {
	ch.messages = append(ch.messages, message)
}

func (c *ChatHub) SearchAppendMsg(id int, message string) {
	for i, chat := range c.chats {
		if chat.customerID == id {
			c.chats[i].messages = append(c.chats[i].messages, message)
		}
	}
}
