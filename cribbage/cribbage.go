package cribbage

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/whereswaldon/cryptage/card"
	"strings"
)

type Card struct {
	Suit string
	Rank string
}

func (c *Card) MarshalText() ([]byte, error) {
	return []byte(c.Rank + " " + c.Suit), nil
}

func (c *Card) UnmarshalText(text []byte) error {
	split := strings.Split(string(text), " ")
	if len(split) < 2 {
		return fmt.Errorf("Invalid card: %v", text)
	}
	c.Rank = split[0]
	c.Suit = split[1]
	return nil
}

type Cribbage struct {
	deck    Deck
	players int
}

type Deck interface {
	Draw() (card.CardFace, error)
	Quit()
	Start() error
}

func NewCribbage(deck Deck) (*Cribbage, error) {
	return &Cribbage{deck: deck, players: 2}, nil
}

func (c *Cribbage) Hand() ([]card.CardFace, error) {
	handSize := getHandSize(c.players)
	hand := make([]card.CardFace, handSize)
	var err error
	for i := range hand {
		hand[i], err = c.deck.Draw()
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to get hand")
		}
	}

	return hand, nil
}

func (c *Cribbage) Quit() error {
	c.deck.Quit()
	return nil
}

func getHandSize(numPlayers int) int {
	switch numPlayers {
	case 2:
		return 6
	case 3:
		fallthrough
	case 4:
		return 5
	default:
		return 0
	}
}
