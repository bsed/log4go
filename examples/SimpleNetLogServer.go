package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	l4g "github.com/bsed/log4go"
)

var (
	port = flag.String("p", "12124", "Port number to listen on")
)

func e(err error) {
	if err != nil {
		fmt.Printf("Erroring out: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	// Bind to the port
	bind, err := net.ResolveUDPAddr("udp4","0.0.0.0:" + *port)
	e(err)

	// Create listener
	listener, err := net.ListenUDP("udp", bind)
	e(err)

	fmt.Printf("Listening to port %s...\n", *port)

	for {
		l4g.Finest("This will only go to those of you really cool UDP kids!  If you change enabled=true.")
		l4g.Debug("Oh no!  %d + %d = %d!", 2, 2, 2+2)
		l4g.Info("About that time, eh chaps?")
		// read into a new buffer
		buffer := make([]byte, 1024)
		_, _, err := listener.ReadFrom(buffer)
		e(err)

		// log to standard output
		fmt.Println(string(buffer))
	}
}
