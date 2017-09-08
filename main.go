package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

// main either listens over TCP or dials depending on whether a command line
// argument was supplied for the port number. If no number was supplied, it
// starts listening on a random port. If you supply a port number as an
// argument, it connects to that port.
func main() {
	if len(os.Args) > 1 {
		dial("localhost:" + os.Args[1])
	} else {
		listen(fmt.Sprintf(":%d", rand.Intn((1<<16)-1024)+1024))
	}
}

// dial connects over tcp to the given address.
func dial(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer conn.Close()
}

// listen starts listening on the given port.
func listen(address string) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Listening on port ", address)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		go func(conn net.Conn) {
			fmt.Println("accepted: ", conn)
			time.Sleep(30 * time.Second)
		}(conn)
	}
}
