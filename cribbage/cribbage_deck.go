package cribbage

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/whereswaldon/cryptage/card"
)

type Deck interface {
	Draw(uint) (card.CardFace, error)
	Quit()
	Start([]card.CardFace) error
	Size() uint
}

// CribbageDeck wraps the cryptage deck construct so that its methods
// return the card structures that the cribbage game uses, rather than
// simply returning byte slices. You should only create a CribbageDeck
// from Decks that have already been initialized
type CribbageDeck struct {
	deck Deck
}

func NewCribbageDeck(d Deck) (*CribbageDeck, error) {
	if d == nil {
		return nil, fmt.Errorf("Cannot create CribbageDeck with nil Deck")
	}
	return &CribbageDeck{deck: d}, nil
}

func (c *CribbageDeck) Draw(deckIndex uint) (*Card, error) {
	cardData, err := c.deck.Draw(deckIndex)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to draw card:")
	}
	card := &Card{}
	err = card.UnmarshalText(cardData)
	if err != nil {
		return nil, errors.Wrapf(err, "Error decoding card:")
	}
	return card, nil
}

func (c *CribbageDeck) Quit() {
	c.deck.Quit()
}

func (c *CribbageDeck) Size() uint {
	return c.deck.Size()
}
