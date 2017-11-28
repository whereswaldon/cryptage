package cribbage

import (
	"fmt"
)

const HAND_SIZE = 4

type Hand struct {
	indicies []uint
	cards    []*Card
}

func NewHand() *Hand {
	return &Hand{
		indicies: make([]uint, 0),
		cards:    make([]*Card, 0),
	}
}

func (h *Hand) Remove(handIndex uint) (*Card, uint, error) {
	if handIndex >= h.Size() {
		return nil, 0, fmt.Errorf("Index out of bounds %d", handIndex)
	}
	lastIndex := h.Size() - 1
	if lastIndex < HAND_SIZE {
		return nil, 0, fmt.Errorf("Cannot remove card from hand, hand is already minimum size")
	}
	card := h.cards[handIndex]
	index := h.indicies[handIndex]
	h.cards[handIndex] = h.cards[lastIndex]
	h.indicies[handIndex] = h.indicies[lastIndex]
	h.cards = h.cards[:lastIndex]
	h.indicies = h.indicies[:lastIndex]

	return card, index, nil
}

func (h *Hand) Add(card *Card, deckIndex uint) error {
	h.cards = append(h.cards, card)
	h.indicies = append(h.indicies, deckIndex)
	return nil
}

func (h *Hand) Get(handIndex uint) (*Card, uint, error) {
	if handIndex >= h.Size() {
		return nil, 0, fmt.Errorf("Index out of bounds")
	}
	return h.cards[handIndex], h.indicies[handIndex], nil
}

func (h *Hand) Size() uint {
	return uint(len(h.indicies))
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

func deckIndiciesForPlayer(pi *PlayerInfo) []uint {
	indices := make([]uint, getHandSize(pi.NumPlayers))
	for i := range indices {
		indices[i] = 2 * uint(i)
	}
	if pi.LocalPlayerIsDealer() {
		for i := range indices {
			indices[i] += 1
		}
	}
	return indices
}
