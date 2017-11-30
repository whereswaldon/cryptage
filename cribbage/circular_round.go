package cribbage

import (
	"fmt"
)

type CircularState struct {
	Sequences     []*Sequence
	Played        []*Card
	currentPlayer int
	cardPlayed    bool
}

func NewCircularState(firstTurnPlayer int) *CircularState {
	s := make([]*Sequence, 1)
	s[0] = NewSeq()
	return &CircularState{
		Sequences:     s,
		currentPlayer: firstTurnPlayer,
	}
}

func (c *CircularState) PlayCard(player int, card *Card) error {
	if !c.ShouldPlayCard(player) {
		return fmt.Errorf("You cannot play another card this turn")
	}
	if !c.CurrentSequence().CanPlay(card) {
		return fmt.Errorf("This card cannot be played this sequence")
	}
	c.CurrentSequence().Play(player, card)
	c.cardPlayed = true
	return nil
}

func (c *CircularState) ShouldPlayCard(player int) bool {
	if player == c.currentPlayer {
		return !c.cardPlayed
	}
	return false
}

func (c *CircularState) IsCurrent(player int) bool {
	return player == c.currentPlayer
}

func (c *CircularState) EndTurn() {
	//if a player ends their turn without playing a card and they were the last player
	//to play a card, the sequence is over
	if !c.cardPlayed {
		player, _ := c.CurrentSequence().Last()
		if player == c.currentPlayer {
			c.newSeq()
		}
	} else if c.CurrentSequence().Total() == 31 {
		c.newSeq()
	}

	c.rotatePlayer()
	c.cardPlayed = false
}

func (c *CircularState) rotatePlayer() {
	if c.currentPlayer == 1 {
		c.currentPlayer = 2
	} else {
		c.currentPlayer = 1
	}
}

func (c *CircularState) newSeq() {
	c.Sequences = append(c.Sequences, NewSeq())
}

func (c *CircularState) CurrentSequence() *Sequence {
	return c.Sequences[len(c.Sequences)-1]
}

var pointValues map[int]int = map[int]int{CLAIM_FIFTEEN: 2, CLAIM_PAIR: 2}

func (c *CircularState) ClaimPoints(pointType int) int {
	if c.CurrentSequence().WorthPoints(pointType) {
		return pointValues[pointType]
	}
	return 0
}
