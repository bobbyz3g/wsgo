package main

import (
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	"golang.org/x/net/websocket"
)

var (
	headers = flag.StringSliceP("header", "H", nil, "Header to add to the request")
	origin  = flag.StringP("origin", "o", "http://localhost/", "Origin to add to the request")
	cookies = flag.StringSliceP("cookie", "b", nil, "Cookie to add to the request")
)

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}
	url := args[0]

	err := wsgoMain(url)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func wsgoMain(url string) error {
	cfg, err := websocket.NewConfig(url, *origin)
	if err != nil {
		return err
	}

	if headers != nil {
		for _, h := range *headers {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				cfg.Header.Add(key, value)
			}
		}
	}
	if cookies != nil {
		for _, c := range *cookies {
			cfg.Header.Add("Cookie", c)
		}
	}
	cli := NewClient(cfg)
	fmt.Printf("Connected to %s\n", url)

	err = cli.Run(nil)
	if err != nil {
		return err
	}
	fmt.Println("Exiting...")
	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: wsgo [flags] <websocket-url>\n")
	flag.PrintDefaults()
}
