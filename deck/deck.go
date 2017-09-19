package deck

import (
	"io"
)

type Deck interface {
	Draw() (Card, error)
	Connector() io.ReadWriteCloser
}

type deck struct {
	Cards      []Card
	connection io.ReadWriteCloser
}

var _ Deck = &deck{}

func (d *deck) Draw() (Card, error) {
	return nil, nil
}

func (d *deck) Connector() io.ReadWriteCloser {
	return d.connection
}

type Card interface {
	Face() string
}

// card represents a single card in a deck. The members correspond with:
// p1cipher - face of card encrypted only with "player 1's" key
// p2cipher - face of card encrypted only with "player 2's" key
// bothCipher - face of card encrypted with both keys
// plain - face of card in plaintext
type card struct {
	P1cipher, P2cipher, BothCipher, plain string
}

var _ Card = &card{}

func (c *card) Face() string {
	return ""
}

func NewDeck() (Deck, error) {
	return nil, nil
}

func ConnectToDeck(deckConnection io.ReadWriteCloser) (Deck, error) {
	return nil, nil
}
