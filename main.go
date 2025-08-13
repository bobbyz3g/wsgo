package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/net/websocket"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wsgo <websocket_url>")
		os.Exit(1)
	}

	url := os.Args[1]
	origin := "http://localhost/" // A default origin is often sufficient

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	fmt.Printf("Connected to %s\n", url)
	fmt.Println("Enter text to send, press Ctrl+C to exit.")

	closeCh := make(chan struct{})
	go receiveMessages(ws, closeCh)
	sendMessages(ws)
}

func receiveMessages(ws *websocket.Conn, closeCh chan struct{}) {
	for {
		select {
		case <-closeCh:
			return
		default:
		}

		var msg string
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			if err == io.EOF {
				close(closeCh)
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

func sendMessages(ws *websocket.Conn) {
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
