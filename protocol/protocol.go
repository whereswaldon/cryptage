package deck

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"math/big"
)

const (
	QUIT            = 0
	START_DECK      = 1
	END_DECK        = 2
	DECRYPT_CARD    = 3
	ONE_CIPHER_CARD = 4
)

type ProtocolHandler interface {
	HandleQuit()
	HandleStartDeck(deck []*big.Int, prime *big.Int)
	HandleEndDeck(deck []*big.Int)
	HandleDecryptCard(index uint64)
	HandleDecryptedCard(index uint64, card *big.Int)
}

// Message is a struct representing a request from one deck to another
type Message struct {
	Type  uint64
	Deck  []*big.Int
	Index uint64
	Value *big.Int
}

// Protocol is an agent implementing the send and recieve sides of the
// cryptage protocol
type Protocol struct {
	r        *gob.Decoder
	w        *gob.Encoder
	recieved chan Message
	done     <-chan struct{}
	handler  ProtocolHandler
}

// NewProtocol creates a Protocol instance assuming that the given
// io.ReadWriteCloser is a connection to another Protocol.
func NewProtocol(conn io.ReadWriteCloser, handler ProtocolHandler, done <-chan struct{}) (*Protocol, error) {
	if conn == nil {
		return nil, fmt.Errorf("Cannot create protocol in nil connection")
	} else if handler == nil {
		return nil, fmt.Errorf("Cannot create protocol with nil handler")
	} else if done == nil {
		return nil, fmt.Errorf("Cannot create protocol with nil done channel")
	}
	messages := make(chan Message)
	proto := &Protocol{
		r:        gob.NewDecoder(conn),
		w:        gob.NewEncoder(conn),
		recieved: messages,
		done:     done,
		handler:  handler,
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
					log.Println("Disconnected: EOF")
					return
				} else {
					log.Println(err)
				}
			}
		}
	}()
	go func() {
		for msg := range proto.recieved {
			switch msg.Type {
			case QUIT:
				proto.handler.HandleQuit()
			case START_DECK:
				proto.handler.HandleStartDeck(msg.Deck, msg.Value)
			case END_DECK:
				proto.handler.HandleEndDeck(msg.Deck)
			case DECRYPT_CARD:
				proto.handler.HandleDecryptCard(msg.Index)
			case ONE_CIPHER_CARD:
				proto.handler.HandleDecryptedCard(msg.Index, msg.Value)
			default:
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
// along with the large prime number used to generate the encryption keys.
func (p *Protocol) SendStartDeck(keyPrime *big.Int, encryptedDeck []*big.Int) error {
	intArr := make([]*big.Int, len(encryptedDeck))
	for i, c := range encryptedDeck {
		intArr[i] = c
	}
	return p.w.Encode(Message{
		Type:  START_DECK,
		Value: keyPrime,
		Deck:  intArr,
	})
}

// SendEndDeck ships the first encrypted deck state to the other player
func (p *Protocol) SendEndDeck(encryptedDeck []*big.Int) error {
	intArr := make([]*big.Int, len(encryptedDeck))
	for i, c := range encryptedDeck {
		intArr[i] = c
	}
	return p.w.Encode(Message{Type: END_DECK, Deck: intArr})
}

// RequestDecryptCard asks the peer to remove their encryption from
// the face of a card.
func (p *Protocol) RequestDecryptCard(cardIndex uint64) error {
	return p.w.Encode(Message{Type: DECRYPT_CARD, Index: cardIndex})
}

func (p *Protocol) SendDecryptedCard(cardIndex uint64, cardCipher *big.Int) error {
	return p.w.Encode(Message{Type: ONE_CIPHER_CARD, Index: cardIndex, Value: cardCipher})
}
