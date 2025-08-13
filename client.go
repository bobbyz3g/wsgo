package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/net/websocket"
)

type Client struct {
	cfg *websocket.Config
}

func NewClient(cfg *websocket.Config) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Run(stopC chan struct{}) error {
	ws, err := websocket.DialConfig(c.cfg)
	if err != nil {
		return err
	}
	defer ws.Close()
	fmt.Println("Enter text to send, press Ctrl+C to exit.")
	go c.receiveMessages(ws, stopC)
	c.sendMessages(ws)

	return nil
}

func (c *Client) receiveMessages(ws *websocket.Conn, stopC chan struct{}) {
	for {
		var msg string
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			if err == io.EOF {
				log.Println("Connection closed by server.")
				os.Exit(0)
			}
			log.Printf("Receive error: %v\n", err)
			return
		}
		// \r moves cursor to the start of the line.
		// \033[K is an ANSI escape code to clear the line.
		// This prevents received messages from corrupting the input line.
		fmt.Printf("\r\033[K<- %s\n\n", msg)
		fmt.Print("-> ")
	}
}

func (c *Client) sendMessages(ws *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("-> ")
	for {
		if !scanner.Scan() {
			break
		}
		text := scanner.Text()
		err := websocket.Message.Send(ws, text)
		if err != nil {
			log.Printf("Send error: %v\n", err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from stdin: %v\n", err)
	}
	fmt.Println()
}
