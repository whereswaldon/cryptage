package cribbage

import (
	"github.com/pkg/errors"
	"github.com/whereswaldon/cryptage/card"
)

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
