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

type Score struct {
	Old, Current int
}
type ScoreBoard map[int]*Score

type Cribbage struct {
	deck                *CribbageDeck
	opponent            *Messenger
	players             *PlayerInfo
	hand, crib          *Hand
	currentState        State
	circular            *CircularState
	cutCard             *Card
	scores              ScoreBoard
	stateChangeRequests chan func()
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
	opponent, err := NewMessenger(opp)
	if err != nil {
		return nil, errors.Wrapf(err, "Couldn't initialize communication with opponent")
	}
	pi := NewPlayerInfo(2, playerNum, DEALER_PLAYER_NUM)
	cribbage := &Cribbage{
		deck:                cDeck,
		players:             pi,
		opponent:            opponent,
		crib:                NewHand(),
		currentState:        DRAW_STATE,
		circular:            NewCircularState(pi.GetNonDealer()),
		stateChangeRequests: make(chan func()),
		scores:              map[int]*Score{1: {}, 2: {}},
	}
	go func() {
		for req := range cribbage.stateChangeRequests {
			req()
		}
	}()
	go cribbage.listenToMessages()
	cribbage.drawHand()
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
		case PASSED_TURN:
			log.Println("Recieved PASSED_TURN")
			if err = c.opponentEndedTurn(); err != nil {
				log.Println(err)
			}
		default:
			log.Println("Unrecognized message type:", m.Type)
		}
	}
}

func (c *Cribbage) drawHand() (*Hand, error) {
	out := make(chan *Hand, 1) // buffered so that a single send doesn't block
	err := c.requestStateChange(func() error {
		defer close(out)
		hand := NewHand()
		indicies := deckIndiciesForPlayer(c.players)
		for _, index := range indicies {
			current, err := c.deck.Draw(index)
			if err != nil {
				return errors.Wrapf(err, "Unable to get hand")
			}
			hand.Add(current, index)
		}
		c.hand = hand
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
		c.crib.Add(nil, deckIndex)
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
		if !c.circular.ShouldPlayCard(c.players.OpponentNum()) {
			return fmt.Errorf("You cannot play cards right now")
		}
		if err := c.circular.PlayCard(c.players.OpponentNum(), card); err != nil {
			return errors.Wrapf(err, "Unable to play card")
		}
		return nil
	})
}

// opponentEndedTurn indicates that the opponent has ended their turn
// within the circular count.
func (c *Cribbage) opponentEndedTurn() error {
	return c.requestStateChange(func() error {
		c.circular.EndTurn()
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
		}
		card, index, err := c.hand.Remove(handIndex)
		if err != nil {
			return errors.Wrapf(err, "Cannot send card to crib, couldn't remove from hand")
		}
		if err := c.crib.Add(card, index); err != nil {
			return errors.Wrapf(err, "Couldn't add card to crib")
		}
		if err := c.opponent.sendToCribMsg(index); err != nil {
			return errors.Wrapf(err, "Couldn't notify opponent of sending card to Crib")
		}
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
		if err := c.opponent.sendCutCardMsg(deckIndex); err != nil {
			return err
		}
		card, err := c.deck.Draw(deckIndex)
		if err != nil {
			return err
		}
		c.cutCard = card
		return nil
	})
}

func (c *Cribbage) PlayCard(handIndex uint) error {
	if handIndex >= uint(getHandSize(c.players.NumPlayers)) {
		return fmt.Errorf("Index out of bounds")
	} else if !c.circular.IsCurrent(c.players.LocalPlayerNum) {
		return fmt.Errorf("Cannot play cards when it isn't your turn!")
	}

	return c.requestStateChange(func() error {
		card, deckIndex, err := c.hand.Get(handIndex)
		if err != nil {
			return err
		}
		err = c.opponent.sendPlayCardMsg(deckIndex)
		if err != nil {
			return fmt.Errorf("Error sending played card to other player: %v", err)
		}

		if err := c.circular.PlayCard(c.players.LocalPlayerNum, card); err != nil {
			return errors.Wrapf(err, "Unable to play card")
		}
		return nil
	})
}

// ClaimPoints attempts to claim points for the local player based upon the current
// state of the circular count
func (c *Cribbage) ClaimPoints(pointType int) error {
	return c.requestStateChange(func() error {
		if c.currentState != CIRCULAR_STATE {
			return fmt.Errorf("You can't claim points when it isn't your turn")
		}
		gain := c.circular.ClaimPoints(CLAIM_FIFTEEN)
		if gain == 0 {
			return fmt.Errorf("There are no points to claim here")
		}
		old := c.scores[c.players.LocalPlayerNum].Current
		c.scores[c.players.LocalPlayerNum].Old = old
		c.scores[c.players.LocalPlayerNum].Current += gain
		return nil
	})
}

func (c *Cribbage) EndTurn() error {
	return c.requestStateChange(func() error {
		if !c.circular.IsCurrent(c.players.LocalPlayerNum) {
			return fmt.Errorf("Cannot end turn when it is not your turn!")
		}
		c.circular.EndTurn()
		if err := c.opponent.sendEndTurnMsg(); err != nil {
			return errors.Wrapf(err, "unable to notify opponent of ended turn")
		}
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
			if c.hand.Size() == 4 {
				c.currentState = DISCARD_WAIT_STATE
			}
		case DISCARD_WAIT_STATE:
			if c.crib.Size() == 4 {
				if c.players.LocalPlayerIsDealer() {
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
			if !c.circular.IsCurrent(c.players.LocalPlayerNum) {
				c.currentState = CIRCULAR_WAIT_STATE
			}
		case CIRCULAR_WAIT_STATE:
			if c.circular.IsCurrent(c.players.LocalPlayerNum) {
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
	printState := func() {
		fmt.Println("scores: ", RenderScores(c.scores))
		fmt.Println("hand: ", RenderHand(c.hand))
		fmt.Println("crib: ", RenderHand(c.crib))
		fmt.Println("cut: ", RenderCard(c.cutCard))
		fmt.Println("seq: ", RenderSeq(c.circular.CurrentSequence()))
	}
	for {
		c.updateState()
		printState()
		fmt.Println(instructionsForState(c.currentState))
		fmt.Print("> ")
		scanner.Scan()
		c.updateState()
		input := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		switch input[0] {
		case "quit":
			c.Quit()
			return
		case "h":
			fallthrough
		case "hand":
			_, err := c.Hand()
			if err != nil {
				fmt.Println(err)
			}
		case "c":
			fallthrough
		case "crib":
			if len(input) < 2 {
				fmt.Println("Usage: crib <card-index>")
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
		case "cut":
			if len(input) < 2 {
				fmt.Println("Usage: cut <card-index>")
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
		case "pl":
			fallthrough
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
		case "pa":
			fallthrough
		case "pass":
			err := c.EndTurn()
			if err != nil {
				fmt.Printf("Problem ending turn: ", err)
				continue
			}
		case "claim":
			if len(input) < 2 {
				fmt.Println(`Usage: claim <score-type>\nwhere <score-type>s are: 15, pair`)
				continue
			}
			var err error
			switch input[1] {
			case "15":
				err = c.ClaimPoints(CLAIM_FIFTEEN)
			case "pair":
				err = c.ClaimPoints(CLAIM_PAIR)
			default:
				fmt.Println("Unknown score type")
				continue
			}
			if err != nil {
				fmt.Println("Error claiming points", err)
			}

		case "help":
			fmt.Println(STR_HELP)
		default:
			fmt.Println("Unknown command: ", input[0])
		}
	}
}
