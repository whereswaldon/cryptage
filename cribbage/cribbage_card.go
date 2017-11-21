package cribbage

import (
	"fmt"
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

func (c *Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank, c.Suit)
}

var suits = []string{"Hearts", "Spades", "Clubs", "Diamonds"}
var ranks = []string{"Two", "Three", "Four", "Five", "Six", "Seven",
	"Eight", "Nine", "Ten", "Jack", "Queen", "King", "Ace"}

func Cards() []card.CardFace {
	deck := make([]card.CardFace, len(suits)*len(ranks))
	for i, suit := range suits {
		for j, rank := range ranks {
			c := &Card{Suit: suit, Rank: rank}
			text, _ := c.MarshalText()
			deck[i*len(ranks)+j] = card.CardFace(text)
		}
	}
	return deck
}
