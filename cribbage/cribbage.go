package cribbage

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/whereswaldon/cryptage/card"
	"os"
)

type ScoreBoard struct {
	p1current, p1last, p2current, p2last uint
}

type Cribbage struct {
	deck      Deck
	opponent  Opponent
	players   int
	playerNum int
	hand      *Hand
}

type Deck interface {
	Draw(uint) (card.CardFace, error)
	Quit()
	Start([]card.CardFace) error
}

type Opponent interface {
	Send(message []byte) error
	Recieve() <-chan []byte
}

func NewCribbage(deck Deck, opp Opponent, playerNum int) (*Cribbage, error) {
	if deck == nil {
		return nil, fmt.Errorf("Cannot create Cribbage with nil deck")
	} else if opp == nil {
		return nil, fmt.Errorf("Cannot create Cribbage with nil opponent")
	} else if playerNum < 1 || playerNum > 2 {
		return nil, fmt.Errorf("Illegal playerNum %d", playerNum)
	}
	return &Cribbage{
		deck:      deck,
		players:   2,
		playerNum: playerNum,
		opponent:  opp,
	}, nil
}

func (c *Cribbage) drawHand() (*Hand, error) {
	handSize := getHandSize(c.players)
	c.hand = &Hand{
		cards:    make([]*Card, handSize),
		indicies: make([]uint, handSize),
	}
	var index uint
	for i := range c.hand.cards {
		if c.playerNum == 1 {
			index = 2 * uint(i)
		} else if c.playerNum == 2 {
			index = 2*uint(i) + 1
		} else {
			return nil, fmt.Errorf("Unsupported player number %d", c.playerNum)
		}

		current, err := c.deck.Draw(index)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to get hand")
		}
		c.hand.indicies[i] = index
		c.hand.cards[i] = &Card{}
		c.hand.cards[i].UnmarshalText(current)
	}

	return c.hand, nil
}

// Hand returns the local player's hand
func (c *Cribbage) Hand() (*Hand, error) {
	if c.hand == nil {
		return c.drawHand()
	}
	return c.hand, nil
}

func (c *Cribbage) Quit() error {
	c.deck.Quit()
	return nil
}

func (c *Cribbage) UI() {
	input := ""
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input = scanner.Text()
		switch input {
		case "quit":
			c.Quit()
			return
		case "hand":
			h, err := c.Hand()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(RenderHand(h))
			}
		default:
			fmt.Println("Uknown command: ", input)
		}
	}
}
