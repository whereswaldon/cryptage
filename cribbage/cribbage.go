package cribbage

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/whereswaldon/cryptage/card"
	"os"
	"strconv"
	"strings"
)

type ScoreBoard struct {
	p1current, p1last, p2current, p2last uint
}

type Cribbage struct {
	deck       Deck
	opponent   Opponent
	players    int
	playerNum  int
	hand, crib *Hand
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
		crib:      &Hand{cards: make([]*Card, 0), indicies: make([]uint, 0)},
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

// Crib adds the card at the specified index within the player's hand to the
// crib. This remove it from the player's hand.
func (c *Cribbage) Crib(handIndex uint) error {
	if handIndex >= uint(len(c.hand.cards)) {
		return fmt.Errorf("Index out of bounds %d", handIndex)
	}
	lastIndex := len(c.hand.cards) - 1
	if lastIndex < 4 {
		return fmt.Errorf("Cannot add another card to crib, hand is already minimum size")
	}
	c.crib.cards = append(c.crib.cards, c.hand.cards[handIndex])
	c.crib.indicies = append(c.crib.indicies, c.hand.indicies[handIndex])
	c.hand.cards[handIndex] = c.hand.cards[lastIndex]
	c.hand.indicies[handIndex] = c.hand.indicies[lastIndex]
	c.hand.cards = c.hand.cards[:lastIndex]
	c.hand.indicies = c.hand.indicies[:lastIndex]
	return nil
}

func (c *Cribbage) UI() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		switch input[0] {
		case "quit":
			c.Quit()
			return
		case "hand":
			h, err := c.Hand()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("hand: ", RenderHand(h))
			}
		case "toCrib":
			if len(input) < 2 {
				fmt.Println("Usage: toCrib <card-index>")
				continue
			}
			i, err := strconv.Atoi(input[1])
			if err != nil {
				fmt.Println("Not a valid card index! Use numbers next time")
				continue
			}
			err = c.Crib(uint(i))
			if err != nil {
				fmt.Printf("error adding %s to crib: %v\n", input[1], err)
			}
			fmt.Println("crib: ", RenderHand(c.crib))
		case "crib":
			fmt.Println("crib: ", RenderHand(c.crib))

		default:
			fmt.Println("Uknown command: ", input)
		}
	}
}
