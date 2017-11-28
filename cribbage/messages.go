package cribbage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

const (
	TO_CRIB     = 0
	CUT_CARD    = 1
	PLAYED_CARD = 2
	PASSED_TURN = 3
)

type Message struct {
	Type int
	Val  uint
}

// Encode converts a message into a byte slice
func Encode(m *Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(m)
	if err != nil {
		return nil, err
	}
	encoded := buf.Bytes()
	return encoded, nil
}

// Decode converts a byte slice into a Message
func Decode(raw []byte) (*Message, error) {
	m := &Message{}
	buf := bytes.NewBuffer(raw)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type Opponent interface {
	Send(message []byte) error
	Recieve() <-chan []byte
}

type Messenger struct {
	opponent Opponent
}

func NewMessenger(opponent Opponent) (*Messenger, error) {
	if opponent == nil {
		return nil, fmt.Errorf("Cannot create Messenger for nil opponent")
	}
	return &Messenger{opponent}, nil
}

// sendToCribMsg sends the opponent a message informing them of which card
// the local player has opted to add to the crib. The value sent is an absolute
// index into the deck, rather than into the local player's hand.
func (m *Messenger) sendToCribMsg(deckIndex uint) error {
	enc, err := Encode(&Message{Type: TO_CRIB, Val: deckIndex})
	if err != nil {
		return err
	}
	log.Println("Sending TO_CRIB")
	return m.opponent.Send(enc)
}

// sendCutCardMsg sends the opponent a message informing them of which card
// the local player is cutting as the shared card
func (m *Messenger) sendCutCardMsg(deckIndex uint) error {
	enc, err := Encode(&Message{Type: CUT_CARD, Val: deckIndex})
	if err != nil {
		return err
	}
	log.Println("Sending CUT_CARD")
	return m.opponent.Send(enc)
}

// sendPlayCardMessage sends the opponent a notification that a card has been
// played in the circular count
func (m *Messenger) sendPlayCardMsg(deckIndex uint) error {
	enc, err := Encode(&Message{Type: PLAYED_CARD, Val: deckIndex})
	if err != nil {
		return err
	}
	log.Println("Sending PLAYED_CARD")
	return m.opponent.Send(enc)
}

func (m *Messenger) Recieve() <-chan []byte {
	return m.opponent.Recieve()
}
