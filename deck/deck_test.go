package deck

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

var address string = "localhost:8080"

// getDecks creates two Deck instances that are connected over
// a TCP socket.
func getDecks() (Deck, Deck, error) {
	serverChan := make(chan struct{})
	deckChan := make(chan Deck)
	errChan := make(chan error, 10)
	var wg sync.WaitGroup
	go func() {
		ln, err := net.Listen("tcp", address)
		if err != nil {
			errChan <- err
			return
		}
		serverChan <- struct{}{}
		go func() {
			conn, err := ln.Accept()
			if err != nil {
				errChan <- err
				return
			}
			d1, err := NewDeck(conn)
			if err != nil {
				errChan <- err
				return
			}
			if d1 == nil {
    				errChan <- fmt.Errorf("Deck 1 should not be nil")
			}
			wg.Done()
			deckChan <- d1
		}()
	}()
	wg.Add(1)

	// wait for server to be listening
	<-serverChan
	go func() {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			errChan <- err
			return
		}
		d2, err := NewDeck(conn)
		if err != nil {
			errChan <- err
			return
		}
		if d2 == nil {
			errChan <- fmt.Errorf("Deck 2 should not be nil")
		}
		wg.Done()
		deckChan <- d2
	}()
	wg.Add(1)

	wg.Wait()
	close(errChan)
	return <-deckChan, <-deckChan, <-errChan
}

func TestDeckQuit(t *testing.T) {
	d1, d2, err := getDecks()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(d1)
	fmt.Println(d2)
}
