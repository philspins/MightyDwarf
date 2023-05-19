package websocket

import (
	"fmt"
	"log"
	"github.com/gorilla/websocket"
)

type Client struct {
	scheme string
	host   string
	port   int
	socket websocket.S
}

func (c *Client) Connect() error {
	return nil
}

func (c *Client) GetUriString() string {
	return fmt.Sprintf("%s://%s:%d", c.scheme, c.host, c.port)
}

func (c *Client) Disconnect() error {
	return nil
}

func (c *Client) ReceiveMessage() (interface{}, error) {
	var buffer = make([]byte, 8192)
	_, message, err := c.conn.ReadMessage(buffer)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (c *Client) SendMessage(date string) error {
	var json []byte
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, json)
	if err != nil {
		log.Println(err)
		return
	}
	return nil
}
