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

func (c *Card) Value() int {
	switch c.Rank {
	case RANK_ACE:
		return 1
	case RANK_TWO:
		return 2
	case RANK_THREE:
		return 3
	case RANK_FOUR:
		return 4
	case RANK_FIVE:
		return 5
	case RANK_SIX:
		return 6
	case RANK_SEVEN:
		return 7
	case RANK_EIGHT:
		return 8
	case RANK_NINE:
		return 9
	case RANK_TEN:
		fallthrough
	case RANK_JACK:
		fallthrough
	case RANK_QUEEN:
		fallthrough
	case RANK_KING:
		return 10
	default:
		return 0
	}
}

func (c *Card) Copy() *Card {
	return &Card{Suit: c.Suit, Rank: c.Rank}
}

const (
	SUIT_HEARTS   = "Hearts"
	SUIT_SPADES   = "Spades"
	SUIT_CLUBS    = "Clubs"
	SUIT_DIAMONDS = "Diamonds"
	RANK_TWO      = "Two"
	RANK_THREE    = "Three"
	RANK_FOUR     = "Four"
	RANK_FIVE     = "Five"
	RANK_SIX      = "Six"
	RANK_SEVEN    = "Seven"
	RANK_EIGHT    = "Eight"
	RANK_NINE     = "Nine"
	RANK_TEN      = "Ten"
	RANK_JACK     = "Jack"
	RANK_QUEEN    = "Queen"
	RANK_KING     = "King"
	RANK_ACE      = "Ace"
)

var suits = []string{SUIT_SPADES, SUIT_HEARTS, SUIT_DIAMONDS, SUIT_CLUBS}
var ranks = []string{RANK_TWO, RANK_THREE, RANK_FOUR, RANK_FIVE, RANK_SIX,
	RANK_SEVEN, RANK_EIGHT, RANK_NINE, RANK_TEN, RANK_JACK, RANK_QUEEN,
	RANK_KING, RANK_ACE}

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
