package cribbage

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"strconv"
	"strings"
)

const DEALER_PLAYER_NUM = 1

type ScoreBoard struct {
	p1current, p1last, p2current, p2last uint
}

type Cribbage struct {
	deck                *CribbageDeck
	opponent            Opponent
	players             int
	playerNum           int
	dealer              int
	hand, crib          *Hand
	currentState        State
	currentSequence     *Sequence
	myTurn              bool
	cutCard             *Card
	stateChangeRequests chan func()
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
	cDeck, err := NewCribbageDeck(deck)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to create CribbageDeck from supplied deck")
	}
	cribbage := &Cribbage{
		deck:                cDeck,
		players:             2,
		myTurn:              playerNum != DEALER_PLAYER_NUM,
		playerNum:           playerNum,
		opponent:            opp,
		dealer:              DEALER_PLAYER_NUM, //player 1 is always first dealer
		crib:                &Hand{cards: make([]*Card, 0), indicies: make([]uint, 0)},
		currentState:        DRAW_STATE,
		currentSequence:     NewSeq(),
		stateChangeRequests: make(chan func()),
	}
	go func() {
		for req := range cribbage.stateChangeRequests {
			req()
		}
	}()
	go cribbage.listenToMessages()
	return cribbage, nil
}

func (c *Cribbage) listenToMessages() {
	for bytes := range c.opponent.Recieve() {
		m, err := Decode(bytes)
		if err != nil {
			log.Println("Error decoding application message:", err)
		}
		switch m.Type {
		case TO_CRIB:
			log.Println("Recieved TO_CRIB")
			if err = c.addIndexToCrib(m.Val); err != nil {
				log.Println(err)
			}
		case CUT_CARD:
			log.Println("Recieved CUT_CARD")
			if err = c.setCutCard(m.Val); err != nil {
				log.Println(err)
			}
		case PLAYED_CARD:
			log.Println("Recieved PLAYED_CARD")
			if err = c.opponentPlayedCard(m.Val); err != nil {
				log.Println(err)
			}
		default:
			log.Println("Unrecognized message type:", m.Type)
		}
	}
}

// sendToCribMsg sends the opponent a message informing them of which card
// the local player has opted to add to the crib. The value sent is an absolute
// index into the deck, rather than into the local player's hand.
func (c *Cribbage) sentToCribMsg(deckIndex uint) error {
	enc, err := Encode(&Message{Type: TO_CRIB, Val: deckIndex})
	if err != nil {
		return err
	}
	log.Println("Sending TO_CRIB")
	return c.opponent.Send(enc)
}

// sendCutCardMsg sends the opponent a message informing them of which card
// the local player is cutting as the shared card

func (c *Cribbage) sentCutCardMsg(deckIndex uint) error {
	enc, err := Encode(&Message{Type: CUT_CARD, Val: deckIndex})
	if err != nil {
		return err
	}
	log.Println("Sending CUT_CARD")
	return c.opponent.Send(enc)
}

// sendPlayCardMessage sends the opponent a notification that a card has been
// played in the circular count
func (c *Cribbage) sendPlayCardMsg(deckIndex uint) error {
	enc, err := Encode(&Message{Type: PLAYED_CARD, Val: deckIndex})
	if err != nil {
		return err
	}
	log.Println("Sending PLAYED_CARD")
	return c.opponent.Send(enc)
}

func (c *Cribbage) drawHand() (*Hand, error) {
	out := make(chan *Hand, 1) // buffered so that a single send doesn't block
	err := c.requestStateChange(func() error {
		defer close(out)
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
				return fmt.Errorf("Unsupported player number %d", c.playerNum)
			}

			current, err := c.deck.Draw(index)
			if err != nil {
				return errors.Wrapf(err, "Unable to get hand")
			}
			c.hand.indicies[i] = index
			c.hand.cards[i] = current
		}
		out <- c.hand
		return nil
	})
	hand := <-out
	return hand, err
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

// requestStateChange executes the provided function atomically on the
// game's state-managing goroutine. Using this function is the only way
// in which game state should be modified once the game has been initialized.
// Otherwise, you invite race conditions.
func (c *Cribbage) requestStateChange(req func() error) error {
	errs := make(chan error)
	c.stateChangeRequests <- func() {
		errs <- req()
	}
	return <-errs
}

// addIndexToCrib adds the card at the given deck index to the crib.
// This is primarily useful for adding the opponent's selections to
// the local crib, since the local player does not know the faces
// of those cards
func (c *Cribbage) addIndexToCrib(deckIndex uint) error {
	if deckIndex >= c.deck.Size() {
		return fmt.Errorf("Index out of bounds")
	}
	return c.requestStateChange(func() error {
		c.crib.cards = append(c.crib.cards, nil)
		c.crib.indicies = append(c.crib.indicies, deckIndex)
		return nil
	})
}

func (c *Cribbage) setCutCard(deckIndex uint) error {
	if deckIndex > c.deck.Size() {
		return fmt.Errorf("Index out of bounds")
	} else if deckIndex < 12 {
		return fmt.Errorf("Tried to set cut card to a card in a current player's hand")
	}
	return c.requestStateChange(func() error {
		card, err := c.deck.Draw(deckIndex)
		if err != nil {
			return err
		}
		c.cutCard = card
		return nil
	})
}

// opponentPlayedCard indicates that the opponent has played the specified card
// within the circular count.
func (c *Cribbage) opponentPlayedCard(deckIndex uint) error {
	if deckIndex >= c.deck.Size() {
		return fmt.Errorf("Index out of bounds")
	}
	return c.requestStateChange(func() error {
		card, err := c.deck.Draw(deckIndex)
		if err != nil {
			return fmt.Errorf("Unable to draw card played by opponent: %v", err)
		}
		if !c.currentSequence.CanPlay(card) {
			return fmt.Errorf("Cannot play card %v", card)
		}
		opponentNum := 1
		if c.playerNum == 1 {
			opponentNum = 2
		}
		c.currentSequence.Play(opponentNum, card)
		c.myTurn = true
		return nil
	})
}

// Crib adds the card at the specified index within the player's hand to the
// crib. This remove it from the player's hand.
func (c *Cribbage) Crib(handIndex uint) error {
	//ensure hand has been initialized
	c.Hand()
	return c.requestStateChange(func() error {
		if c.currentState != DISCARD_STATE && c.currentState != DRAW_STATE {
			return fmt.Errorf("You can't discard to the crib right now")
		} else if handIndex >= uint(len(c.hand.cards)) {
			return fmt.Errorf("Index out of bounds %d", handIndex)
		}
		lastIndex := len(c.hand.cards) - 1
		if lastIndex < 4 {
			return fmt.Errorf("Cannot add another card to crib, hand is already minimum size")
		}
		c.sentToCribMsg(c.hand.indicies[handIndex])
		c.crib.cards = append(c.crib.cards, c.hand.cards[handIndex])
		c.crib.indicies = append(c.crib.indicies, c.hand.indicies[handIndex])
		c.hand.cards[handIndex] = c.hand.cards[lastIndex]
		c.hand.indicies[handIndex] = c.hand.indicies[lastIndex]
		c.hand.cards = c.hand.cards[:lastIndex]
		c.hand.indicies = c.hand.indicies[:lastIndex]
		return nil
	})
}

// Cut attempts to cut the deck at the specified card
func (c *Cribbage) Cut(deckIndex uint) error {
	return c.requestStateChange(func() error {
		if c.currentState != CUT_STATE {
			return fmt.Errorf("You can't cut the deck right now")
		} else if deckIndex >= uint(c.deck.Size()) {
			return fmt.Errorf("Index out of bounds %d", deckIndex)
		} else if deckIndex < 12 { //cutting into cards that have been dealt
			return fmt.Errorf("Cannot cut at index %d, cards 0-12 are in player hands.", deckIndex)
		}
		c.sentCutCardMsg(deckIndex)
		card, err := c.deck.Draw(deckIndex)
		if err != nil {
			return err
		}
		c.cutCard = card
		return nil
	})
}

func (c *Cribbage) PlayCard(handIndex uint) error {
	if handIndex >= uint(getHandSize(c.players)) {
		return fmt.Errorf("Index out of bounds")
	} else if !c.myTurn {
		return fmt.Errorf("Cannot play cards when it isn't your turn!")
	}

	return c.requestStateChange(func() error {
		err := c.sendPlayCardMsg(c.hand.indicies[handIndex])
		if err != nil {
			return fmt.Errorf("Error sending played card to other player: %v", err)
		}
		card := c.hand.cards[handIndex]
		if !c.currentSequence.CanPlay(card) {
			return fmt.Errorf("Card %s cannot be played!", card)
		}
		c.currentSequence.Play(c.playerNum, card)
		c.myTurn = false
		return nil
	})
}

func (c *Cribbage) updateState() {
	_ = c.requestStateChange(func() error {
		switch c.currentState {
		case DRAW_STATE:
			if c.hand != nil {
				c.currentState = DISCARD_STATE
			}
		case DISCARD_STATE:
			if len(c.hand.cards) == 4 {
				c.currentState = DISCARD_WAIT_STATE
			}
		case DISCARD_WAIT_STATE:
			if len(c.crib.cards) == 4 {
				if c.dealer == c.playerNum {
					c.currentState = CUT_WAIT_STATE
				} else {
					c.currentState = CUT_STATE
				}
			}
		case CUT_STATE:
			if c.cutCard != nil {
				c.currentState = CIRCULAR_STATE
			}
		case CUT_WAIT_STATE:
			if c.cutCard != nil {
				c.currentState = CIRCULAR_WAIT_STATE
			}
		case CIRCULAR_STATE:
			if !c.myTurn {
				c.currentState = CIRCULAR_WAIT_STATE
			}
		case CIRCULAR_WAIT_STATE:
			if c.myTurn {
				c.currentState = CIRCULAR_STATE
			}
		case INTERNAL_STATE:
		case CRIB_STATE:
		case END_STATE:
		}
		return nil
	})
}

func (c *Cribbage) UI() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		c.updateState()
		fmt.Println(instructionsForState(c.currentState))
		fmt.Print("> ")
		scanner.Scan()
		c.updateState()
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
				continue
			}
			fmt.Println("crib: ", RenderHand(c.crib))
		case "crib":
			fmt.Println("crib: ", RenderHand(c.crib))
		case "cut":
			if c.cutCard != nil {
				fmt.Println("cut: ", RenderCard(c.cutCard))
			} else {
				fmt.Println("No cut card yet")
			}
		case "cutAt":
			if len(input) < 2 {
				fmt.Println("Usage: cutAt <card-index>")
				continue
			}
			i, err := strconv.Atoi(input[1])
			if err != nil {
				fmt.Println("Not a valid card index! Use numbers next time")
				continue
			}
			err = c.Cut(uint(i))
			if err != nil {
				fmt.Printf("error cutting card %d: %v\n", i, err)
				continue
			}
			fmt.Println("cut: ", RenderCard(c.cutCard))
		case "seq":
			fmt.Println(RenderSeq(c.currentSequence))
		case "play":
			if len(input) < 2 {
				fmt.Println("Usage: play <hand-index>")
				continue
			}
			i, err := strconv.Atoi(input[1])
			if err != nil {
				fmt.Println("Not a valid card index! Use numbers next time")
				continue
			}
			err = c.PlayCard(uint(i))
			if err != nil {
				fmt.Printf("error playing card %d: %v\n", i, err)
				continue
			}
			fmt.Println("seq: ", RenderSeq(c.currentSequence))
		case "help":
			fmt.Println(STR_HELP)
		default:
			fmt.Println("Uknown command: ", input[0])
			fmt.Println(STR_HELP)
		}
	}
}
