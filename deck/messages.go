package deck

import (
	"encoding/gob"
	"fmt"
	"io"
	"math/big"
)

const (
	QUIT       = 0
	START_DECK = 1
	END_DECK = 2
)

// Message is a struct representing a request from one deck to another
type Message struct {
	Type uint64
	Deck []*big.Int
}

// Protocol is an agent implementing the send and recieve sides of the
// cryptage protocol
type Protocol struct {
	r        *gob.Decoder
	w        *gob.Encoder
	recieved chan Message
	done     <-chan struct{}
}

// NewProtocol creates a Protocol instance assuming that the given
// io.ReadWriteCloser is a connection to another Protocol.
func NewProtocol(conn io.ReadWriteCloser, done <-chan struct{}) (*Protocol, error) {
	messages := make(chan Message)
	proto := &Protocol{
		r:        gob.NewDecoder(conn),
		w:        gob.NewEncoder(conn),
		recieved: messages,
		done:     done,
	}
	go func() {
		defer close(proto.recieved)
		for {
			var message Message
			select {
			case <-done:
				return
			default:
				if err := proto.r.Decode(&message); err == nil {
					proto.recieved <- message
				} else if err == io.EOF {
					fmt.Println("Disconnected: EOF")
					return
				} else {
					fmt.Println(err)
				}
			}
		}
	}()
	return proto, nil
}

// SendQuit asks the connected peer to quit
func (p *Protocol) SendQuit() error {
	return p.w.Encode(Message{Type: QUIT})
}

// SendStartDeck ships the first encrypted deck state to the other player
func (p *Protocol) SendStartDeck(encryptedDeck []card) error {
	intArr := make([]*big.Int, len(encryptedDeck))
	for i, c := range encryptedDeck {
    		intArr[i] = c.P1cipher
	}
	return p.w.Encode(Message{Type: START_DECK, Deck: intArr})
}

// SendEndDeck ships the first encrypted deck state to the other player
func (p *Protocol) SendEndDeck(encryptedDeck []card) error {
	intArr := make([]*big.Int, len(encryptedDeck))
	for i, c := range encryptedDeck {
    		intArr[i] = c.BothCipher
	}
	return p.w.Encode(Message{Type: END_DECK, Deck: intArr})
}
// Listen waits for events from the connected peer
func (p *Protocol) Listen() <-chan Message {
	return p.recieved
}
