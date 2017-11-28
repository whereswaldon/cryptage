package main

import (
	"flag"
	"fmt"
	"github.com/whereswaldon/cryptage/cribbage"
	"github.com/whereswaldon/cryptage/deck"
	"log"
	"math/rand"
	"net"
	"os"
)

const (
	MAX_PORT = 65535
	MIN_PORT = 1024
)

// main either listens over TCP or dials depending on whether a command line
// argument was supplied for the port number. If no number was supplied, it
// starts listening on a random port. If you supply a port number as an
// argument, it connects to that port.
func main() {
	stdin := flag.String("stdin", "", "specify an alternate file to read stdin from")
	flag.Parse()
	if *stdin != "" {
		file, err := os.Open(*stdin)
		if err != nil {
			log.Println("Error opening input file: ", err)
			os.Exit(2)
		}
		os.Stdin = file
	}
	if len(flag.Args()) < 1 {
		usage()
	} else if flag.Args()[0] == "join" {
		if len(flag.Args()) < 2 {
			usage()
		}
		dial(flag.Args()[1])
	} else if flag.Args()[0] == "host" {
		port := fmt.Sprintf(":%d", rand.Intn((MAX_PORT)-MIN_PORT)+MIN_PORT)
		if len(flag.Args()) > 2 {
			port = flag.Args()[1]
		}
		listen(port)
	} else {
		fmt.Fprintln(os.Stderr, "Unknown subcommand.")
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: "+os.Args[0]+" (host [<:port>] | join <hostname:port>)")
	os.Exit(1)
}

// dial connects over tcp to the given address.
func dial(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Connected")
	defer conn.Close()
	deck, err := deck.NewDeck(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Starting game...")
	if err := deck.Start(cribbage.Cards()); err != nil {
		fmt.Println(err)
		return
	}
	game, err := cribbage.NewCribbage(deck, deck, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Playing...")
	game.UI()
}

// listen starts listening on the given port.
func listen(address string) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Listening on port ", address)
	conn, err := ln.Accept()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println("Connected")

	defer conn.Close()
	deck, err := deck.NewDeck(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	game, err := cribbage.NewCribbage(deck, deck, 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Playing...")

	game.UI()
}
