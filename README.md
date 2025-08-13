# wsgo

`wsgo` is a simple and interactive command-line tool for testing WebSocket servers. It allows you to connect to a WebSocket endpoint, send messages, and view incoming messages in real-time.

## Features

-   Connect to any WebSocket server.
-   Interactive shell for sending messages.
-   Clear distinction between sent (`->`) and received (`<-`) messages.
-   Handles connection closing gracefully.

## Installation

To use `wsgo`, you need to have Go installed on your system.

```sh
go install github.com/bobbyz3g/wsgo@latest
```

## Usage

To connect to a WebSocket server, simply provide the server's URL as a command-line argument.

```sh
./wsgo <websocket_url>
```