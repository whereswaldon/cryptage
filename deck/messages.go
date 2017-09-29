package deck

import (
	"encoding/gob"
	"io"
)

const (
	QUIT = 0
)

type Message struct {
	MessageType uint64
	Payload     map[string]string
}

type Protocol struct {
	r        *gob.Decoder
	w        *gob.Encoder
	recieved chan Message
}

func NewProtocol(conn io.ReadWriteCloser) (*Protocol, error) {
	messages := make(chan Message)
	proto := &Protocol{
		r:        gob.NewDecoder(conn),
		w:        gob.NewEncoder(conn),
		recieved: messages,
	}
	go func() {
		for {
			var message Message
			if err := proto.r.Decode(&message); err == nil {
				proto.recieved <- message
			}

		}
	}()
	return proto, nil
}

func (p *Protocol) SendQuit() error {
	return p.w.Encode(Message{MessageType: QUIT})
}
