package deck

import (
	"fmt"
	"io"
)

type Deck interface {
	Draw() (string, error)
	Start()
	Play()
	String() string
}

// ensure that *deck fulfills Deck interface
var _ Deck = &deck{}

type deck struct {
	cards    []card
	protocol *Protocol
	done     chan struct{}
}

// Draw draws a single card from the deck
func (d *deck) Draw() (string, error) {
	return "", nil
}

func (d *deck) String() string {
	return "Deck"
}

func (d *deck) Quit() {
	close(d.done)
}

// Play runs the game, but does not initiate it.
func (d *deck) Play() {
	d.handleMessages()
}

func (d *deck) handleMessages() {
	defer d.Quit()
	for msg := range d.protocol.Listen() {
		switch msg.Type {
		case QUIT:
			fmt.Println("QUIT")
			return
		default:
			fmt.Println("Unknown message: %v", msg)
		}
	}
}

// Start runs the game, and initiates the first hand
func (d *deck) Start() {
	err := d.protocol.SendQuit()
	if err != nil {
		fmt.Println(err)
	}

	d.handleMessages()
}

// NewDeck creates a deck of cards and assumes that the given
// io.ReadWriteCloser is a connection of some sort to another
// deck.
func NewDeck(deckConnection io.ReadWriteCloser) (Deck, error) {
	d := &deck{}
	done := make(chan struct{})
	p, err := NewProtocol(deckConnection, done)
	if err != nil {
		return nil, err
	}
	d.done = done
	d.protocol = p
	return d, nil
}
