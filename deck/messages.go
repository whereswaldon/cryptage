package deck

import (
	"encoding/gob"
	"io"
)

const (
	QUIT = 0
)

// Message is a struct representing a request from one deck to another
type Message struct {
	Type uint64
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
		done: done,
	}
	go func() {
        	defer close(proto.recieved)
		for range proto.done {
			var message Message
			if err := proto.r.Decode(&message); err == nil {
				proto.recieved <- message
			}

		}
	}()
	return proto, nil
}

// SendQuit asks the connected peer to quit
func (p *Protocol) SendQuit() error {
	return p.w.Encode(Message{Type: QUIT})
}

// Listen waits for events from the connected peer
func (p *Protocol) Listen() <-chan Message {
	return p.recieved
}
