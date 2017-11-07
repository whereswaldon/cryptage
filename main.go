package main

import (
	"bufio"
	"fmt"
	"github.com/whereswaldon/cryptage/cribbage"
	"github.com/whereswaldon/cryptage/deck"
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
	if len(os.Args) < 2 {
		usage()
	} else if os.Args[1] == "join" {
		if len(os.Args) < 3 {
			usage()
		}
		dial(os.Args[2])
	} else if os.Args[1] == "host" {
		port := fmt.Sprintf(":%d", rand.Intn((MAX_PORT)-MIN_PORT)+MIN_PORT)
		if len(os.Args) > 2 {
			port = os.Args[2]
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
	if err := deck.Start(); err != nil {
		fmt.Println(err)
		return
	}
	game, err := cribbage.NewCribbage(deck)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Playing...")
	enterUI(game)
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
	game, err := cribbage.NewCribbage(deck)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Playing...")

	enterUI(game)
}

func enterUI(game *cribbage.Cribbage) {
	input := ""
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input = scanner.Text()
		switch input {
		case "quit":
			game.Quit()
			return
		case "hand":
			h, err := game.Hand()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Hand: ", h)
			}
		default:
			fmt.Println("Uknown command: ", input)
		}
	}
}
